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
	"time"

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

	errorClientConnectionNotInitialized = "the client connection to the database is not initialized"

	// These two values are used to configure our connection pools.
	// The values are subject to change, but are based off the recommendations on:
	// https://github.com/go-sql-driver/mysql#important-settings
	connMaxLifetime = 3 * time.Minute
	maxOpenedConns  = 20
)

type (
	auth struct {
		username string
		password string
	}

	Config struct {
		organisation string
		database     string
		branch       string
		hostname     string
		authWrite    auth
		authRead     auth
	}

	Client struct {
		config  *Config
		writeDB *sql.DB
		readDB  *sql.DB
	}
)

func (au auth) isValid() bool {
	return au.username != "" && au.password != ""
}

func (cfg *Config) IsValid() bool {
	return cfg.organisation != "" && cfg.database != "" && cfg.branch != "" && cfg.hostname != "" && cfg.authWrite.isValid() && cfg.authRead.isValid()
}

func (cfg *Config) AddToViper(v *viper.Viper) {
	// General settings
	_ = v.UnmarshalKey(flagPsdbOrg, &cfg.organisation)
	_ = v.UnmarshalKey(flagPsdbHost, &cfg.hostname)
	_ = v.UnmarshalKey(flagPsdbDatabase, &cfg.database)
	_ = v.UnmarshalKey(flagPsdbBranch, &cfg.branch)

	// Write authentication
	_ = v.UnmarshalKey(flagPsdbPasswordWrite, &cfg.authWrite.password)
	_ = v.UnmarshalKey(flagPsdbUserWrite, &cfg.authWrite.username)

	// Read authentication
	_ = v.UnmarshalKey(flagPsdbPasswordWrite, &cfg.authRead.password)
	_ = v.UnmarshalKey(flagPsdbUserWrite, &cfg.authRead.username)
}

func (cfg *Config) AddToCommand(cmd *cobra.Command) {
	// General settings
	cmd.Flags().StringVar(&cfg.organisation, flagPsdbOrg, "", "Name of the PlanetScaleDB organization.")
	cmd.Flags().StringVar(&cfg.hostname, flagPsdbHost, "", "Hostname of the PlanetScaleDB database.")
	cmd.Flags().StringVar(&cfg.database, flagPsdbDatabase, "", "PlanetScaleDB database name.")
	cmd.Flags().StringVar(&cfg.branch, flagPsdbBranch, "main", "PlanetScaleDB branch to use.")
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
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&tls=true&interpolateParams=true", a.username, a.password, cfg.hostname, cfg.database)
}

func (cfg *Config) NewClient() (*Client, error) {
	// Open a connection pool to the write servers
	writedb, err := sql.Open("mysql", cfg.connectionString(cfg.authWrite))
	if err != nil {
		return nil, err
	}
	setConnDefault(writedb)

	// Open a connection pool to the read-only servers
	readdb, err := sql.Open("mysql", cfg.connectionString(cfg.authRead))
	if err != nil {
		return nil, err
	}
	setConnDefault(readdb)

	return &Client{
		config:  cfg,
		writeDB: writedb,
		readDB:  readdb,
	}, nil
}

func setConnDefault(db *sql.DB) {
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetMaxOpenConns(maxOpenedConns)
	db.SetMaxIdleConns(maxOpenedConns)
}

func (c *Client) Close() error {
	if c.writeDB == nil || c.readDB == nil {
		return errors.New(errorClientConnectionNotInitialized)
	}
	err := c.writeDB.Close()
	if err != nil {
		return err
	}
	return c.readDB.Close()
}

func (c *Client) Write(query string, args ...interface{}) (int64, error) {
	if c.writeDB == nil {
		return 0, errors.New(errorClientConnectionNotInitialized)
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

func (c *Client) Read(query string, args ...interface{}) (*sql.Rows, error) {
	if c.readDB == nil {
		return nil, errors.New(errorClientConnectionNotInitialized)
	}
	rows, err := c.readDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
