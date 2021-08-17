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
	"database/sql"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagPsdbOrg      = "planetscale-db-org"
	flagPsdbPassword = "planetscale-db-password"
	flagPsdbUser     = "planetscale-db-user"
	flagPsdbHost     = "planetscale-db-host"
	flagPsdbDatabase = "planetscale-db-database"
	flagPsdbBranch   = "planetscale-db-branch"

	ErrorClientConnectionNotInitialized = "the client connection to the database is not initialized"
)

type (
	Config struct {
		Org      string
		Database string
		Branch   string
		User string
		Password string
		Host string
	}

	Client struct {
		Config *Config
		dial   *sql.DB
	}
)

func (cfg Config) IsValid() bool {
	return !(cfg.User == "" || cfg.Password == "" || cfg.Host == "" || cfg.Org == "" || cfg.Database == "" || cfg.Branch == "")
}

func (cfg *Config) AddToViper(v *viper.Viper) {
	_ = v.UnmarshalKey(flagPsdbOrg, &cfg.Org)
	_ = v.UnmarshalKey(flagPsdbHost, &cfg.Host)
	_ = v.UnmarshalKey(flagPsdbPassword, &cfg.Password)
	_ = v.UnmarshalKey(flagPsdbUser, &cfg.User)
	_ = v.UnmarshalKey(flagPsdbDatabase, &cfg.Database)
	_ = v.UnmarshalKey(flagPsdbBranch, &cfg.Branch)
}

func (cfg *Config) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&cfg.Org, flagPsdbOrg, "", "Name of the PlanetscaleDB organization.")
	cmd.Flags().StringVar(&cfg.User, flagPsdbUser, "", "Username used to authenticate to PlanetscaleDB.")
	cmd.Flags().StringVar(&cfg.Password, flagPsdbPassword, "", "Password used to authenticate to PlanetscaleDB.")
	cmd.Flags().StringVar(&cfg.Host, flagPsdbHost, "", "Hostname of the PlanetscaleDB database.")
	cmd.Flags().StringVar(&cfg.Database, flagPsdbDatabase, "", "PlanetscaleDB database name.")
	cmd.Flags().StringVar(&cfg.Branch, flagPsdbBranch, "main", "PlanetscaleDB branch to use.")

	_ = viper.BindPFlag(flagPsdbOrg, cmd.Flags().Lookup(flagPsdbOrg))
	_ = viper.BindPFlag(flagPsdbHost, cmd.Flags().Lookup(flagPsdbHost))
	_ = viper.BindPFlag(flagPsdbUser, cmd.Flags().Lookup(flagPsdbUser))
	_ = viper.BindPFlag(flagPsdbPassword, cmd.Flags().Lookup(flagPsdbPassword))
	_ = viper.BindPFlag(flagPsdbDatabase, cmd.Flags().Lookup(flagPsdbDatabase))
	_ = viper.BindPFlag(flagPsdbBranch, cmd.Flags().Lookup(flagPsdbBranch))
}

func (cfg Config) NewClient() (*Client, error) {
	client := &Client{
		Config: &cfg,
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&tls=true", cfg.User, cfg.Password, cfg.Host, cfg.Database))
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
