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

package mysql

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagDatabaseName     = "db-database"
	flagDatabaseHost     = "db-host"
	flagDatabasePassword = "db-password"
	flagDatabaseUser     = "db-user"
)

func (cfg *ConfigDB) AddToViper(v *viper.Viper) {
	_ = v.UnmarshalKey(flagDatabaseName, &cfg.Database)
	_ = v.UnmarshalKey(flagDatabaseHost, &cfg.Host)
	_ = v.UnmarshalKey(flagDatabasePassword, &cfg.Password)
	_ = v.UnmarshalKey(flagDatabaseUser, &cfg.User)
}

func (cfg *ConfigDB) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&cfg.Database, flagDatabaseName, "", "Database to use.")
	cmd.Flags().StringVar(&cfg.Host, flagDatabaseHost, "", "Hostname of the database")
	cmd.Flags().StringVar(&cfg.Password, flagDatabasePassword, "", "Password to authenticate the database.")
	cmd.Flags().StringVar(&cfg.User, flagDatabaseUser, "", "User used to connect to the database")

	_ = viper.BindPFlag(flagDatabaseName, cmd.Flags().Lookup(flagDatabaseName))
	_ = viper.BindPFlag(flagDatabaseHost, cmd.Flags().Lookup(flagDatabaseHost))
	_ = viper.BindPFlag(flagDatabasePassword, cmd.Flags().Lookup(flagDatabasePassword))
	_ = viper.BindPFlag(flagDatabaseUser, cmd.Flags().Lookup(flagDatabaseUser))
}
