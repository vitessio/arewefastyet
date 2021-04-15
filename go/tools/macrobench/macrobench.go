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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vitessio/arewefastyet/go/mysql"
	"os"
	"os/exec"
	"strings"
)

const (
	ErrorNoSysBenchResult          = "no sysbench results were found"

	prefixMacroBenchSysbenchConfig = "macrobench_"
)

func buildSysbenchArgString(m map[string]string, step string) []string {
	output := map[string]string{}
	for k, v := range m {
		idx := strings.Index(k, "_")
		if idx < 0 {
			continue
		}
		head := k[:idx]
		tail := k[idx+1:]
		if head == step || head == "all" {
			if _, exists := output[tail]; exists {
				if head == "all" {
					continue
				}
			}
			output[tail] = v
		}
	}

	var results []string
	for k, v := range output {
		results = append(results, fmt.Sprintf("--%s=%s", k, v))
	}
	return results
}

// Run executes a macro benchmark by using sysbench.
// Based on the given MacroBenchConfig, the function will
// parse the configuration to send down to sysbench (size of tables
// duration of benchmark, mysql targets, etc...).
// After the execution, the output of the last step (stepRun) is
// converted to a slice of MacroBenchmarkResult, which is then
// uploaded to MySQL using the mysql.ConfigDB in MacroBenchConfig.
//
// We use two forks of sysbench, one for oltp workloads
// and the other for tpcc workload. We use these forks because
// they implement a custom method to print results in JSON.
//
// Regular Sysbench: https://github.com/planetscale/sysbench
// Sysbench-TPCC: https://github.com/planetscale/sysbench-tpcc
func Run(mabcfg Config) error {
	var err error
	var sqlClient *mysql.Client

	if mabcfg.DatabaseConfig != nil && mabcfg.DatabaseConfig.IsValid() {
		sqlClient, err = mysql.New(*mabcfg.DatabaseConfig)
		if err != nil {
			return err
		}
		defer sqlClient.Close()
	}

	// Create new macro benchmark in MySQL
	var macrobenchID int
	if sqlClient != nil {
		macrobenchID, err = mabcfg.insertBenchmarkToSQL(sqlClient)
		if err != nil {
			return err
		}
	}

	// Prepare
	if mabcfg.WorkingDirectory == "" {
		mabcfg.WorkingDirectory, _ = os.Getwd()
	}
	mabcfg.parseIntoMap(prefixMacroBenchSysbenchConfig)
	newSteps := skipSteps(steps, mabcfg.SkipSteps)

	// Execution
	var resStr []byte
	for _, step := range newSteps {
		args := buildSysbenchArgString(mabcfg.M, step.Name)
		args = append(args, mabcfg.WorkloadPath, step.SysbenchName)
		command := exec.Command(mabcfg.SysbenchExec, args...)
		command.Dir = mabcfg.WorkingDirectory
		out, err := command.Output()
		if err != nil {
			return fmt.Errorf("%s:\n%s", err.Error(), string(out))
		}
		if step.Name == stepRun {
			resStr = out
		}
	}

	// Parse results
	var results []MacroBenchmarkResult
	err = json.Unmarshal(resStr, &results)
	if err != nil {
		return fmt.Errorf("unmarshal results: %+v\n", err)
	} else if len(results) == 0 {
		return errors.New(ErrorNoSysBenchResult)
	}

	// Save results
	if sqlClient != nil {
		err = results[0].insertToMySQL(mabcfg.Type, macrobenchID, sqlClient)
		if err != nil {
			return err
		}
	}
	return nil
}
