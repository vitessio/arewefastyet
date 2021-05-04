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

package influxdb

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagInfluxHostname = "influx-hostname"
	flagInfluxPort     = "influx-port"
	flagInfluxUsername = "influx-username"
	flagInfluxPassword = "influx-password"
	flagInfluxDatabase = "influx-database"
)

// Config defines the required configuration used to authenticate
// to an InfluxDB database.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// IsValid return true if Config is ready to be used, and false otherwise.
func (cfg Config) IsValid() bool {
	return cfg.Host != ""
}

func (cfg *Config) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&cfg.Host, flagInfluxHostname, "", "Hostname of InfluxDB.")
	cmd.Flags().StringVar(&cfg.Port, flagInfluxPort, "8086", "Port on which to InfluxDB listens.")
	cmd.Flags().StringVar(&cfg.User, flagInfluxUsername, "", "Username used to connect to InfluxDB.")
	cmd.Flags().StringVar(&cfg.Password, flagInfluxPassword, "", "Password used to connect to InfluxDB.")
	cmd.Flags().StringVar(&cfg.Database, flagInfluxDatabase, "", "Name of the database to use in InfluxDB.")

	_ = cmd.MarkFlagRequired(flagInfluxHostname)

	_ = viper.BindPFlag(flagInfluxHostname, cmd.Flags().Lookup(flagInfluxHostname))
	_ = viper.BindPFlag(flagInfluxPort, cmd.Flags().Lookup(flagInfluxPort))
	_ = viper.BindPFlag(flagInfluxUsername, cmd.Flags().Lookup(flagInfluxUsername))
	_ = viper.BindPFlag(flagInfluxPassword, cmd.Flags().Lookup(flagInfluxPassword))
	_ = viper.BindPFlag(flagInfluxDatabase, cmd.Flags().Lookup(flagInfluxDatabase))
}
