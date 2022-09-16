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

package exec

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagRootExec             = "exec-root-dir"
	flagGitRefExec           = "exec-git-ref"
	flagSourceExec           = "exec-source"
	flagExecType             = "exec-type"
	flagVtgatePlannerVersion = "exec-vtgate-planner-version"
	flagExecPullNB           = "exec-pull-nb"
	flagGolangVersion        = "exec-go-version"
	flagServerAddress        = "exec-server-address"
	flagVitessConfig         = "exec-vitess-config"
)

func (e *Exec) AddToViper(v *viper.Viper) (err error) {
	_ = v.UnmarshalKey(flagRootExec, &e.rootDir)
	_ = v.UnmarshalKey(flagGitRefExec, &e.GitRef)
	_ = v.UnmarshalKey(flagSourceExec, &e.Source)
	_ = v.UnmarshalKey(flagExecType, &e.TypeOf)
	_ = v.UnmarshalKey(flagVtgatePlannerVersion, &e.VtgatePlannerVersion)
	_ = v.UnmarshalKey(flagExecPullNB, &e.PullNB)
	_ = v.UnmarshalKey(flagGolangVersion, &e.GolangVersion)
	_ = v.UnmarshalKey(flagServerAddress, &e.ServerAddress)
	_ = v.UnmarshalKey(flagVitessConfig, &e.rawVitessConfig)

	e.AnsibleConfig.AddToViper(v)
	e.configDB.AddToViper(v)
	e.statsRemoteDBConfig.AddToViper(v)
	return nil
}

func (e *Exec) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&e.rootDir, flagRootExec, "", "Path to the root directory of exec.")
	cmd.Flags().StringVar(&e.GitRef, flagGitRefExec, "", "Git reference on which the benchmarks will run.")
	cmd.Flags().StringVar(&e.Source, flagSourceExec, "", "Name of the source that triggered the execution.")
	cmd.Flags().StringVar(&e.TypeOf, flagExecType, "", "Defines the execution type (oltp, tpcc, micro).")
	cmd.Flags().StringVar(&e.VtgatePlannerVersion, flagVtgatePlannerVersion, "V3", "Defines the vtgate planner version to use. Valid values are: V3, Gen4, Gen4Greedy and Gen4Fallback.")
	cmd.Flags().IntVar(&e.PullNB, flagExecPullNB, 0, "Defines the number of the pull request against which to execute.")
	cmd.Flags().StringVar(&e.GolangVersion, flagGolangVersion, "1.17", "Defines the golang version that will be used by this execution.")
	cmd.Flags().StringVar(&e.ServerAddress, flagServerAddress, "", "The IP address of the server on which the benchmark will be executed.")

	_ = viper.BindPFlag(flagRootExec, cmd.Flags().Lookup(flagRootExec))
	_ = viper.BindPFlag(flagGitRefExec, cmd.Flags().Lookup(flagGitRefExec))
	_ = viper.BindPFlag(flagSourceExec, cmd.Flags().Lookup(flagSourceExec))
	_ = viper.BindPFlag(flagExecType, cmd.Flags().Lookup(flagExecType))
	_ = viper.BindPFlag(flagExecPullNB, cmd.Flags().Lookup(flagExecPullNB))
	_ = viper.BindPFlag(flagGolangVersion, cmd.Flags().Lookup(flagGolangVersion))
	_ = viper.BindPFlag(flagServerAddress, cmd.Flags().Lookup(flagServerAddress))

	e.AnsibleConfig.AddToPersistentCommand(cmd)
	e.statsRemoteDBConfig.AddToCommand(cmd)
	e.configDB.AddToCommand(cmd)
}
