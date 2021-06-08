/*
Copyright 2021 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package microbench

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vitessio/arewefastyet/go/storage/psdb"
)

const (
	flagExecUUID = "microbench-exec-uuid"
)

type Config struct {
	// RootDir is the root path from where micro benchmarks will
	// be executed.
	RootDir        string

	// Package we want to microbenchmark.
	Package        string

	// Output file on which to print the intermediate results.
	Output         string

	// DatabaseConfig used to save results to SQL. If this field
	// is nil, saving results will be skipped and no error will
	// be returned.
	DatabaseConfig *psdb.Config

	// execUUID refers to parent execution of the microbenchmark.
	// If this field is empty, the corresponding column in SQL
	// will be set to NULL.
	execUUID string
}

func (mbc *Config) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&mbc.execUUID, flagExecUUID, "", "UUID of the parent execution, an empty string will set to NULL.")

	_ = viper.BindPFlag(flagExecUUID, cmd.Flags().Lookup(flagExecUUID))

	mbc.DatabaseConfig.AddToCommand(cmd)
}
