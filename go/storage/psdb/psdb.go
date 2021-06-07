/*
 *
 * Copyright 2021 The Vitess Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * /
 */

package psdb

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/planetscale/planetscale-go/planetscale"
	"github.com/planetscale/planetscale-go/planetscale/dbutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagPsdbOrg              = "planetscale-db-org"
	flagPsdbServiceTokenName = "planetscale-db-service-token-name"
	flagPsdbServiceToken     = "planetscale-db-service-token"
	flagPsdbDatabase         = "planetscale-db-database"
	flagPsdbBranch           = "planetscale-db-branch"

	ErrorClientConnectionNotInitialized = "the client connection to the database is not initialized"
)

type (
	ServiceToken struct {
		Name  string
		Token string
	}

	Config struct {
		ServiceToken
		Org      string
		Database string
		Branch   string
	}

	Client struct {
		Config *Config
		client *planetscale.Client
		dial   *sql.DB
	}
)

func (cfg *Config) AddToViper(v *viper.Viper) {
	_ = v.UnmarshalKey(flagPsdbOrg, &cfg.Org)
	_ = v.UnmarshalKey(flagPsdbServiceTokenName, &cfg.Name)
	_ = v.UnmarshalKey(flagPsdbServiceToken, &cfg.Token)
	_ = v.UnmarshalKey(flagPsdbDatabase, &cfg.Database)
	_ = v.UnmarshalKey(flagPsdbBranch, &cfg.Branch)
}

func (cfg *Config) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&cfg.Org, flagPsdbOrg, "", "Name of the PlanetscaleDB organization.")
	cmd.Flags().StringVar(&cfg.Name, flagPsdbServiceTokenName, "", "PlanetscaleDB service token name.")
	cmd.Flags().StringVar(&cfg.Token, flagPsdbServiceToken, "", "PlanetscaleDB service token value.")
	cmd.Flags().StringVar(&cfg.Database, flagPsdbDatabase, "", "PlanetscaleDB database name.")
	cmd.Flags().StringVar(&cfg.Branch, flagPsdbBranch, "main", "PlanetscaleDB branch to use.")

	_ = viper.BindPFlag(flagPsdbOrg, cmd.Flags().Lookup(flagPsdbOrg))
	_ = viper.BindPFlag(flagPsdbServiceTokenName, cmd.Flags().Lookup(flagPsdbServiceTokenName))
	_ = viper.BindPFlag(flagPsdbServiceToken, cmd.Flags().Lookup(flagPsdbServiceToken))
	_ = viper.BindPFlag(flagPsdbDatabase, cmd.Flags().Lookup(flagPsdbDatabase))
	_ = viper.BindPFlag(flagPsdbBranch, cmd.Flags().Lookup(flagPsdbBranch))
}

func (cfg Config) NewClient() (*Client, error) {
	client := &Client{
		Config: &cfg,
	}
	psdbClient, err := planetscale.NewClient(
		planetscale.WithServiceToken(cfg.Name, cfg.Token),
	)
	if err != nil {
		return nil, err
	}
	client.client = psdbClient

	dialCfg := &dbutil.DialConfig{
		Organization: cfg.Org,
		Database:     cfg.Database,
		Branch:       cfg.Branch,
		Client:       psdbClient,
		MySQLConfig:  mysql.NewConfig(),
	}
	db, err := dbutil.Dial(context.Background(), dialCfg)
	if err != nil {
		return nil, err
	}
	client.dial = db
	return client, nil
}

func (c *Client) Close() error {
	if c.dial == nil {
		return errors.New(ErrorClientConnectionNotInitialized)
	}
	return c.dial.Close()
}

func (c *Client) Insert(query string, args ...interface{}) (int64, error) {
	if c.dial == nil {
		return 0, errors.New(ErrorClientConnectionNotInitialized)
	}
	stms, err := c.dial.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stms.Close()

	res, err := stms.Exec(args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (c *Client) Select(query string, args ...interface{}) (*sql.Rows, error) {
	if c.dial == nil {
		return nil, errors.New(ErrorClientConnectionNotInitialized)
	}
	rows, err := c.dial.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
