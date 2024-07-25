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
	flagPsdbOrg           = "planetscale-db-org"
	flagPsdbPasswordWrite = "planetscale-db-password-write"
	flagPsdbUserWrite     = "planetscale-db-user-write"
	flagPsdbPasswordRead  = "planetscale-db-password-read"
	flagPsdbUserRead      = "planetscale-db-user-read"
	flagPsdbHost          = "planetscale-db-host"
	flagPsdbDatabase      = "planetscale-db-database"
	flagPsdbBranch        = "planetscale-db-branch"

	ErrorClientConnectionNotInitialized = "the client connection to the database is not initialized"
)

type (
	auth struct {
		username string
		password string
	}

	Config struct {
		Org       string
		Database  string
		Branch    string
		Host      string
		authWrite auth
		authRead  auth
	}

	Client struct {
		Config  *Config
		writeDB *sql.DB
		readDB  *sql.DB
	}
)

func (au auth) isValid() bool {
	return au.username != "" && au.password != ""
}

func (cfg *Config) IsValid() bool {
	return cfg.Org != "" && cfg.Database != "" && cfg.Branch != "" && cfg.Host != "" && cfg.authWrite.isValid() && cfg.authRead.isValid()
}

func (cfg *Config) AddToViper(v *viper.Viper) {
	// General settings
	_ = v.UnmarshalKey(flagPsdbOrg, &cfg.Org)
	_ = v.UnmarshalKey(flagPsdbHost, &cfg.Host)
	_ = v.UnmarshalKey(flagPsdbDatabase, &cfg.Database)
	_ = v.UnmarshalKey(flagPsdbBranch, &cfg.Branch)

	// Write authentication
	_ = v.UnmarshalKey(flagPsdbPasswordWrite, &cfg.authWrite.password)
	_ = v.UnmarshalKey(flagPsdbUserWrite, &cfg.authWrite.username)

	// Read authentication
	_ = v.UnmarshalKey(flagPsdbPasswordWrite, &cfg.authWrite.password)
	_ = v.UnmarshalKey(flagPsdbUserWrite, &cfg.authWrite.username)
}

func (cfg *Config) AddToCommand(cmd *cobra.Command) {
	// General settings
	cmd.Flags().StringVar(&cfg.Org, flagPsdbOrg, "", "Name of the PlanetscaleDB organization.")
	cmd.Flags().StringVar(&cfg.Host, flagPsdbHost, "", "Hostname of the PlanetscaleDB database.")
	cmd.Flags().StringVar(&cfg.Database, flagPsdbDatabase, "", "PlanetscaleDB database name.")
	cmd.Flags().StringVar(&cfg.Branch, flagPsdbBranch, "main", "PlanetscaleDB branch to use.")
	_ = viper.BindPFlag(flagPsdbOrg, cmd.Flags().Lookup(flagPsdbOrg))
	_ = viper.BindPFlag(flagPsdbHost, cmd.Flags().Lookup(flagPsdbHost))
	_ = viper.BindPFlag(flagPsdbDatabase, cmd.Flags().Lookup(flagPsdbDatabase))
	_ = viper.BindPFlag(flagPsdbBranch, cmd.Flags().Lookup(flagPsdbBranch))

	// Write authentication
	cmd.Flags().StringVar(&cfg.authWrite.username, flagPsdbUserWrite, "", "Username used to authenticate to the write servers of PlanetScaleDB.")
	cmd.Flags().StringVar(&cfg.authWrite.password, flagPsdbPasswordWrite, "", "Password used to authenticate to the write servers of PlanetScaleDB.")
	_ = viper.BindPFlag(flagPsdbUserWrite, cmd.Flags().Lookup(flagPsdbUserWrite))
	_ = viper.BindPFlag(flagPsdbPasswordWrite, cmd.Flags().Lookup(flagPsdbPasswordWrite))

	// Read authentication
	cmd.Flags().StringVar(&cfg.authRead.username, flagPsdbUserRead, "", "Username used to authenticate to the read-only servers of PlanetScaleDB.")
	cmd.Flags().StringVar(&cfg.authRead.password, flagPsdbPasswordRead, "", "Password used to authenticate to the read-only servers of PlanetScaleDB.")
	_ = viper.BindPFlag(flagPsdbUserRead, cmd.Flags().Lookup(flagPsdbUserRead))
	_ = viper.BindPFlag(flagPsdbPasswordRead, cmd.Flags().Lookup(flagPsdbPasswordRead))
}

func (cfg *Config) connectionString(a auth) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&tls=true&interpolateParams=true", a.username, a.password, cfg.Host, cfg.Database)
}

func (cfg *Config) NewClient() (*Client, error) {
	// Open a connection pool to the write servers
	writedb, err := sql.Open("mysql", cfg.connectionString(cfg.authWrite))
	if err != nil {
		return nil, err
	}

	// Open a connection pool to the read-only servers
	readdb, err := sql.Open("mysql", cfg.connectionString(cfg.authRead))
	if err != nil {
		return nil, err
	}

	return &Client{
		Config:  cfg,
		writeDB: writedb,
		readDB:  readdb,
	}, nil
}

func (c *Client) Close() error {
	if c.writeDB == nil {
		return errors.New(ErrorClientConnectionNotInitialized)
	}
	return c.writeDB.Close()
}

func (c *Client) Insert(query string, args ...interface{}) (int64, error) {
	if c.writeDB == nil {
		return 0, errors.New(ErrorClientConnectionNotInitialized)
	}
	stms, err := c.writeDB.Prepare(query)
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
	if c.readDB == nil {
		return nil, errors.New(ErrorClientConnectionNotInitialized)
	}
	rows, err := c.readDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
