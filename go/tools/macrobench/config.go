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
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vitessio/arewefastyet/go/mysql"
	"strings"
)

// MacroBenchConfig defines a configuration used to execute macro benchmark.
// For instance, the MacroBench method uses MacroBenchConfig.
type MacroBenchConfig struct {
	SysbenchExec   string
	WorkloadPath   string
	DatabaseConfig *mysql.ConfigDB
	M              map[string]string
	SkipSteps      []string
	Type           MacroBenchmarkType

	// Source defines from where the macro benchmark is triggered.
	// This field is used to distinguish runs triggered by webhooks,
	// local, nightly build, and so on.
	Source string

	// GitRef refers to the commit SHA pointing to the version
	// of Vitess that we are currently macro benchmarking.
	GitRef string
}

const (
	flagSysbenchExecutable = "macrobench-sysbench-executable"
	flagSysbenchPath       = "macrobench-workload-path"
	flagSkipSteps          = "macrobench-skip-steps"
	flagType               = "macrobench-type"
	flagSource             = "macrobench-source"
	flagGitRef             = "macrobench-git-ref"
)

// AddToCommand will add the different CLI flags used by MacroBenchConfig into
// the given *cobra.Command.
func (mabcfg *MacroBenchConfig) AddToCommand(cmd *cobra.Command) {
	mabcfg.DatabaseConfig.AddToCommand(cmd)

	cmd.Flags().StringVar(&mabcfg.WorkloadPath, flagSysbenchPath, "", "")
	cmd.Flags().StringVar(&mabcfg.SysbenchExec, flagSysbenchExecutable, "", "")
	cmd.Flags().StringSliceVar(&mabcfg.SkipSteps, flagSkipSteps, []string{}, "")
	cmd.Flags().Var(&mabcfg.Type, flagType, "")
	cmd.Flags().StringVar(&mabcfg.Source, flagSource, "", "")
	cmd.Flags().StringVar(&mabcfg.GitRef, flagGitRef, "", "")

	_ = viper.BindPFlag(flagSysbenchPath, cmd.Flags().Lookup(flagSysbenchPath))
	_ = viper.BindPFlag(flagSysbenchExecutable, cmd.Flags().Lookup(flagSysbenchExecutable))
	_ = viper.BindPFlag(flagSkipSteps, cmd.Flags().Lookup(flagSkipSteps))
	_ = viper.BindPFlag(flagType, cmd.Flags().Lookup(flagType))
	_ = viper.BindPFlag(flagSource, cmd.Flags().Lookup(flagSource))
	_ = viper.BindPFlag(flagGitRef, cmd.Flags().Lookup(flagGitRef))
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

func (mabcfg MacroBenchConfig) RegisterNewBenchmarkToMySQL(client *mysql.Client) (newMacroBenchmarkID int, err error) {
	if client == nil {
		return 0, errors.New(mysql.ErrorClientConnectionNotInitialized)
	}
	query := "INSERT INTO benchmark(commit, source) VALUES(?, ?)"
	id, err := client.Insert(query, mabcfg.GitRef, mabcfg.Source)
	if err != nil {
		return 0, err
	}
	newMacroBenchmarkID = int(id)
	return newMacroBenchmarkID, nil
}
