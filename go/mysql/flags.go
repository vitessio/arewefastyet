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

func (cfg *ConfigDB) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&cfg.Database, "db-database", "", "Database to use.")
	cmd.Flags().StringVar(&cfg.Host, "db-host", "", "Hostname of the database")
	cmd.Flags().StringVar(&cfg.Password, "db-password", "", "Password to authenticate the database.")
	cmd.Flags().StringVar(&cfg.User, "db-user", "", "User used to connect to the database")

	_ = viper.BindPFlag("db-database", cmd.Flags().Lookup("db-database"))
	_ = viper.BindPFlag("db-host", cmd.Flags().Lookup("db-host"))
	_ = viper.BindPFlag("db-password", cmd.Flags().Lookup("db-password"))
	_ = viper.BindPFlag("db-user", cmd.Flags().Lookup("db-user"))
}
