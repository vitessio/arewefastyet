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
	"github.com/vitessio/arewefastyet/go/infra/ansible"
	"strings"
)

const (
	statsRemoteDBHost     = "stats-remote-db-host"
	statsRemoteDBDatabase = "stats-remote-db-database"
	statsRemoteDBPort     = "stats-remote-db-port"
	statsRemoteDBUser     = "stats-remote-db-user"
	statsRemoteDBPassword = "stats-remote-db-password"
)

type RemoteDBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

func (rdbcfg *RemoteDBConfig) AddToViper(v *viper.Viper) {
	_ = v.UnmarshalKey(statsRemoteDBHost, &rdbcfg.Host)
	_ = v.UnmarshalKey(statsRemoteDBPort, &rdbcfg.Port)
	_ = v.UnmarshalKey(statsRemoteDBDatabase, &rdbcfg.DbName)
	_ = v.UnmarshalKey(statsRemoteDBUser, &rdbcfg.User)
	_ = v.UnmarshalKey(statsRemoteDBPassword, &rdbcfg.Password)
}

func (rdbcfg *RemoteDBConfig) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&rdbcfg.Host, statsRemoteDBHost, "", "Hostname of the stats remote database.")
	cmd.Flags().StringVar(&rdbcfg.Port, statsRemoteDBPort, "", "Port of the stats remote database.")
	cmd.Flags().StringVar(&rdbcfg.DbName, statsRemoteDBDatabase, "", "Name of the stats remote database.")
	cmd.Flags().StringVar(&rdbcfg.User, statsRemoteDBUser, "", "User used to connect to the stats remote database")
	cmd.Flags().StringVar(&rdbcfg.Password, statsRemoteDBPassword, "", "Password to authenticate the stats remote database.")

	_ = viper.BindPFlag(statsRemoteDBHost, cmd.Flags().Lookup(statsRemoteDBHost))
	_ = viper.BindPFlag(statsRemoteDBPort, cmd.Flags().Lookup(statsRemoteDBPort))
	_ = viper.BindPFlag(statsRemoteDBDatabase, cmd.Flags().Lookup(statsRemoteDBDatabase))
	_ = viper.BindPFlag(statsRemoteDBUser, cmd.Flags().Lookup(statsRemoteDBUser))
	_ = viper.BindPFlag(statsRemoteDBPassword, cmd.Flags().Lookup(statsRemoteDBPassword))
}

func (rdbcfg RemoteDBConfig) valid() bool {
	return rdbcfg.Host != "" && rdbcfg.Port != "" && rdbcfg.DbName != ""
}

// AddToAnsible will add the stats remote database configuration
// to the list of Ansible ExtraVars.
func (rdbcfg RemoteDBConfig) AddToAnsible(ansibleCfg *ansible.Config) {
	if !rdbcfg.valid() {
		return
	}
	ansibleCfg.ExtraVars[strings.ReplaceAll(statsRemoteDBHost, "-", "_")] = rdbcfg.Host
	ansibleCfg.ExtraVars[strings.ReplaceAll(statsRemoteDBDatabase, "-", "_")] = rdbcfg.DbName
	ansibleCfg.ExtraVars[strings.ReplaceAll(statsRemoteDBPort, "-", "_")] = rdbcfg.Port
	ansibleCfg.ExtraVars[strings.ReplaceAll(statsRemoteDBUser, "-", "_")] = rdbcfg.User
	ansibleCfg.ExtraVars[strings.ReplaceAll(statsRemoteDBPassword, "-", "_")] = rdbcfg.Password
}
