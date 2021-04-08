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

package macrobench

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vitessio/arewefastyet/go/mysql"
	"strings"
)

type MacroBenchConfig struct {
	SysbenchExec   string
	WorkloadPath   string
	DatabaseConfig *mysql.ConfigDB
	M              map[string]string
}

const (
	flagSysbenchExecutable = "macrobench-sysbench-executable"
	flagSysbenchPath       = "macrobench-workload-path"
)

func (mabcfg *MacroBenchConfig) AddToCommand(cmd *cobra.Command) {
	mabcfg.DatabaseConfig.AddToCommand(cmd)

	cmd.Flags().StringVar(&mabcfg.WorkloadPath, flagSysbenchPath, "", "")
	cmd.Flags().StringVar(&mabcfg.SysbenchExec, flagSysbenchExecutable, "", "")

	_ = viper.BindPFlag(flagSysbenchPath, cmd.Flags().Lookup(flagSysbenchPath))
	_ = viper.BindPFlag(flagSysbenchExecutable, cmd.Flags().Lookup(flagSysbenchExecutable))
}

func (mabcfg *MacroBenchConfig) parseIntoMap(prefix string) {
	mabcfg.M = map[string]string{}
	keys := viper.AllKeys()
	for _, key := range keys {
		if strings.Index(key, prefix) == 0 {
			mabcfg.M[key[len(prefix):]] = viper.GetString(key)
		}
	}
}
