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
	"strings"

	"github.com/vitessio/arewefastyet/go/storage"
	"github.com/vitessio/arewefastyet/go/storage/influxdb"
	"github.com/vitessio/arewefastyet/go/storage/psdb"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vitessio/arewefastyet/go/storage/mysql"
)

// Config defines a configuration used to execute macro benchmark.
// For instance, the Run method uses MacroBenchConfig.
type Config struct {
	// SysbenchExec defines the path to sysbench binary
	SysbenchExec string

	// WorkloadPath defines the path to the lua file used by sysbench.
	WorkloadPath string

	// DatabaseConfig points to the configuration used to create
	// a *psdb.Client. If no configuration, results and reports will
	// not be saved to a database, though the program won't fail.
	DatabaseConfig *psdb.Config

	// MetricsDatabaseConfig points to the required configuration to create
	// a *influxdb.Client. If no configuration is provided results will not be
	// saved to the database and the program will not fail.
	MetricsDatabaseConfig *influxdb.Config

	// M contains all metadata used to parameter sysbench execution.
	// This key value map stores the value of each CLI parameters.
	M map[string]string

	// SkipSteps is a slice of string that is used to skip some of
	// sysbench steps.
	SkipSteps []string

	// Type will be used to differentiate macro benchmarks.
	Type Type

	// Source defines from where the macro benchmark is triggered.
	// This field is used to distinguish runs triggered by webhooks,
	// local, nightly build, and so on.
	Source string

	// GitRef refers to the commit SHA pointing to the version
	// of Vitess that we are currently macro benchmarking.
	GitRef string

	// VtgatePlannerVersion refers to the planner version that the vtgate is using
	// in Vitess that we are currently macro benchmarking.
	VtgatePlannerVersion string

	// WorkingDirectory defines from where sysbench commands will be executed.
	// This parameter
	WorkingDirectory string

	// execUUID refers to the parent execution for this macro benchmark.
	// If this field if empty, the corresponding column in SQL will be set
	// to NULL.
	execUUID string

	// vtgateWebPorts lists web endpoint of each VTGate
	vtgateWebPorts []string
}

const (
	flagSysbenchExecutable   = "macrobench-sysbench-executable"
	flagSysbenchPath         = "macrobench-workload-path"
	flagSkipSteps            = "macrobench-skip-steps"
	flagType                 = "macrobench-type"
	flagGitRef               = "macrobench-git-ref"
	flagWorkingDirectory     = "macrobench-working-directory"
	flagExecUUID             = "macrobench-exec-uuid"
	flagVtgatePlannerVersion = "macrobench-vtgate-planner-version"
	flagVtgateWebPorts       = "macrobench-vtgate-web-ports"
)

// AddToCommand will add the different CLI flags used by MacroBenchConfig into
// the given *cobra.Command.
func (mabcfg *Config) AddToCommand(cmd *cobra.Command) {
	mabcfg.DatabaseConfig.AddToCommand(cmd)
	mabcfg.MetricsDatabaseConfig.AddToCommand(cmd)

	cmd.Flags().StringVar(&mabcfg.WorkloadPath, flagSysbenchPath, "", "Path to the workload used by sysbench.")
	cmd.Flags().StringVar(&mabcfg.SysbenchExec, flagSysbenchExecutable, "", "Path to the sysbench binary.")
	cmd.Flags().StringSliceVar(&mabcfg.SkipSteps, flagSkipSteps, []string{}, "Slice of sysbench steps to skip.")
	cmd.Flags().Var(&mabcfg.Type, flagType, "Type of macro benchmark.")
	cmd.Flags().StringVar(&mabcfg.VtgatePlannerVersion, flagVtgatePlannerVersion, "", "Vtgate planner version running on Vitess")
	cmd.Flags().StringVar(&mabcfg.GitRef, flagGitRef, "", "Git SHA referring to the macro benchmark.")
	cmd.Flags().StringVar(&mabcfg.WorkingDirectory, flagWorkingDirectory, "", "Directory on which to execute sysbench.")
	cmd.Flags().StringVar(&mabcfg.execUUID, flagExecUUID, "", "UUID of the parent execution, an empty string will set to NULL.")
	cmd.Flags().StringSliceVar(&mabcfg.vtgateWebPorts, flagVtgateWebPorts, nil, "List of the web port for each VTGate.")

	_ = viper.BindPFlag(flagSysbenchPath, cmd.Flags().Lookup(flagSysbenchPath))
	_ = viper.BindPFlag(flagSysbenchExecutable, cmd.Flags().Lookup(flagSysbenchExecutable))
	_ = viper.BindPFlag(flagSkipSteps, cmd.Flags().Lookup(flagSkipSteps))
	_ = viper.BindPFlag(flagType, cmd.Flags().Lookup(flagType))
	_ = viper.BindPFlag(flagGitRef, cmd.Flags().Lookup(flagGitRef))
	_ = viper.BindPFlag(flagVtgatePlannerVersion, cmd.Flags().Lookup(flagVtgatePlannerVersion))
	_ = viper.BindPFlag(flagWorkingDirectory, cmd.Flags().Lookup(flagWorkingDirectory))
	_ = viper.BindPFlag(flagExecUUID, cmd.Flags().Lookup(flagExecUUID))
	_ = viper.BindPFlag(flagVtgateWebPorts, cmd.Flags().Lookup(flagVtgateWebPorts))
}

func (mabcfg *Config) parseIntoMap(prefix string) {
	mabcfg.M = map[string]string{}
	keys := viper.AllKeys()
	for _, key := range keys {
		if strings.Index(key, prefix) == 0 {
			mabcfg.M[key[len(prefix):]] = viper.GetString(key)
		}
	}
}

// insertBenchmarkToSQL will insert a new row in the benchmark table based on
// the given MacroBenchConfig. The newly created row's unique ID is returned.
func (mabcfg Config) insertBenchmarkToSQL(client storage.SQLClient) (newMacroBenchmarkID int, err error) {
	if client == nil {
		return 0, errors.New(mysql.ErrorClientConnectionNotInitialized)
	}
	query := "INSERT INTO macrobenchmark(exec_uuid, commit, vtgate_planner_version, type) VALUES(NULLIF(?, ''), ?, ?, ?, ?)"
	res, err := client.Insert(query, mabcfg.execUUID, mabcfg.GitRef, mabcfg.VtgatePlannerVersion, mabcfg.Type.ToUpper().String())
	if err != nil {
		return 0, err
	}
	newMacroBenchmarkID = int(res)
	return newMacroBenchmarkID, nil
}
