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
	"github.com/planetscale/planetscale-go/planetscale"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagPsdbOrg              = "planetscale-db-org"
	flagPsdbServiceTokenName = "planetscale-db-service-token-name"
	flagPsdbServiceToken     = "planetscale-db-service-token"
)

type (
	ServiceToken struct {
		Name  string
		Token string
	}

	Config struct {
		ServiceToken
		Org string
	}

	Client struct {
		Config *Config
		client *planetscale.Client
	}
)

func (cfg *Config) AddToViper(v *viper.Viper) {
	_ = v.UnmarshalKey(flagPsdbOrg, &cfg.Org)
	_ = v.UnmarshalKey(flagPsdbServiceTokenName, &cfg.Name)
	_ = v.UnmarshalKey(flagPsdbServiceToken, &cfg.Token)
}

func (cfg *Config) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&cfg.Org, flagPsdbOrg, "", "Name of the PlanetscaleDB organization.")
	cmd.Flags().StringVar(&cfg.Name, flagPsdbServiceTokenName, "", "PlanetscaleDB service token name.")
	cmd.Flags().StringVar(&cfg.Token, flagPsdbServiceToken, "", "PlanetscaleDB service token value.")

	_ = viper.BindPFlag(flagPsdbOrg, cmd.Flags().Lookup(flagPsdbOrg))
	_ = viper.BindPFlag(flagPsdbServiceTokenName, cmd.Flags().Lookup(flagPsdbServiceTokenName))
	_ = viper.BindPFlag(flagPsdbServiceToken, cmd.Flags().Lookup(flagPsdbServiceToken))
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
	return client, nil
}
