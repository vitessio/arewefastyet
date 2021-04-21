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

package stats

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	statsRemoteDBHost     = "stats-remote-db-host"
	statsRemoteDBDatabase = "stats-remote-db-database"
	statsRemoteDBPort     = "stats-remote-db-port"
	statsRemoteDBUser     = "stats-remote-db-user"
	statsRemoteDBPassword = "stats-remote-db-password"
)

type RemoteDBConfig struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
}

func (rdbcfg *RemoteDBConfig) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&rdbcfg.host, statsRemoteDBHost, "", "Hostname of the stats remote database.")
	cmd.Flags().StringVar(&rdbcfg.port, statsRemoteDBPort, "", "Port of the stats remote database.")
	cmd.Flags().StringVar(&rdbcfg.dbName, statsRemoteDBDatabase, "", "Name of the stats remote database.")
	cmd.Flags().StringVar(&rdbcfg.user, statsRemoteDBUser, "", "User used to connect to the stats remote database")
	cmd.Flags().StringVar(&rdbcfg.password, statsRemoteDBPassword, "", "Password to authenticate the stats remote database.")

	_ = viper.BindPFlag(statsRemoteDBHost, cmd.Flags().Lookup(statsRemoteDBHost))
	_ = viper.BindPFlag(statsRemoteDBPort, cmd.Flags().Lookup(statsRemoteDBPort))
	_ = viper.BindPFlag(statsRemoteDBDatabase, cmd.Flags().Lookup(statsRemoteDBDatabase))
	_ = viper.BindPFlag(statsRemoteDBUser, cmd.Flags().Lookup(statsRemoteDBUser))
	_ = viper.BindPFlag(statsRemoteDBPassword, cmd.Flags().Lookup(statsRemoteDBPassword))
}
