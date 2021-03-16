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
	"github.com/vitessio/arewefastyet/go/mysql"
	"github.com/vitessio/arewefastyet/go/tools/microbench"
)

func run() *cobra.Command {
	var mbcfg microbench.MicroBenchConfig
	mbcfg.DatabaseConfig = &mysql.ConfigDB{}

	cmd := &cobra.Command{
		Use:   "run [root dir] <pkg> <output path>",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Run micro benchmarks from the <root dir> on <pkg>, and outputs to <output path>.",
		RunE: func(cmd *cobra.Command, args []string) error {
			idx := 0
			if len(args) == 3 {
				mbcfg.RootDir = args[idx]
				idx++
			} else {
				mbcfg.RootDir = "."
			}
			mbcfg.Package = args[idx]
			mbcfg.Output = args[idx + 1]
			microbench.MicroBenchmark(mbcfg)
			return nil
		},
	}

	mbcfg.DatabaseConfig.AddToCommand(cmd)

	return cmd
}
