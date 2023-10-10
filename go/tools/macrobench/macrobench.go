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
	"os"
	"os/exec"
	"strings"

	"github.com/vitessio/arewefastyet/go/exec/metrics"
	"github.com/vitessio/arewefastyet/go/storage/influxdb"
	"github.com/vitessio/arewefastyet/go/storage/psdb"
)

type PlannerVersion string

const (
	ErrorNoSysBenchResult = "no sysbench results were found"

	prefixMacroBenchSysbenchConfig = "macrobench_"

	V3Planner           PlannerVersion = "V3"
	Gen4FallbackPlanner PlannerVersion = "Gen4Fallback"
	Gen4Planner         PlannerVersion = "Gen4"
)

var (
	LegacyPlannerVersions = []PlannerVersion{
		V3Planner,
		Gen4FallbackPlanner,
	}
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
	// get sql database client
	sqlClient, err := createSQLClient(mabcfg.DatabaseConfig)
	if err != nil {
		return err
	}
	defer sqlClient.Close()

	// get metrics database client
	metricsClient, err := createMetricsDatabaseClient(mabcfg.MetricsDatabaseConfig)
	if err != nil {
		return err
	}

	// Create new macro benchmark in MySQL
	var macrobenchID int
	if sqlClient != nil {
		macrobenchID, err = mabcfg.insertBenchmarkToSQL(sqlClient)
		if err != nil {
			return err
		}
	}

	fmt.Println("Step insert is done")

	// Prepare
	if mabcfg.WorkingDirectory == "" {
		mabcfg.WorkingDirectory, _ = os.Getwd()
	}
	mabcfg.parseIntoMap(prefixMacroBenchSysbenchConfig)
	newSteps := skipSteps(steps, mabcfg.SkipSteps)

	fmt.Println("Step prepare is done")

	// Execution
	var resStr []byte
	for _, step := range newSteps {
		fmt.Printf("Step %s begins\n", step)
		args := buildSysbenchArgString(mabcfg.M, step.Name)
		args = append(args, mabcfg.WorkloadPath, step.SysbenchName)
		fmt.Println("Execute:", mabcfg.SysbenchExec, strings.Join(args, " "))
		command := exec.Command(mabcfg.SysbenchExec, args...)
		command.Dir = mabcfg.WorkingDirectory
		out, err := command.Output()
		if err != nil {
			return fmt.Errorf("%s:\n%s", err.Error(), string(out))
		}
		fmt.Printf("Step %s is done\n", step)
		if step.Name == stepRun {
			resStr = out
		}
	}

	err = handleResults(mabcfg, resStr, sqlClient, metricsClient, macrobenchID)
	if err != nil {
		return err
	}
	return nil
}

func handleResults(mabcfg Config, resStr []byte, sqlClient *psdb.Client, metricsClient *influxdb.Client, macrobenchID int) error {
	err := handleSysBenchResults(resStr, sqlClient, macrobenchID)
	if err != nil {
		return err
	}
	err = handleMetricsResults(metricsClient, sqlClient, mabcfg.execUUID)
	if err != nil {
		return err
	}
	err = handleVTGateResults(mabcfg.vtgateWebPorts, sqlClient, mabcfg.execUUID, macrobenchID)
	if err != nil {
		return err
	}
	return nil
}

func createSQLClient(dbConfig *psdb.Config) (client *psdb.Client, err error) {
	if dbConfig != nil && dbConfig.IsValid() {
		client, err = dbConfig.NewClient()
		if err != nil {
			return
		}
	}
	return
}

func createMetricsDatabaseClient(dbConfig *influxdb.Config) (client *influxdb.Client, err error) {
	if dbConfig != nil && dbConfig.IsValid() {
		client, err = dbConfig.NewClient()
		if err != nil {
			return
		}
	}
	return
}

func handleVTGateResults(ports []string, sqlClient *psdb.Client, execUUID string, macrobenchID int) error {
	plans, err := getVTGatesQueryPlans(ports)
	if err != nil {
		return err
	}
	return insertVTGateQueryMapToMySQL(sqlClient, execUUID, plans, macrobenchID)
}

func handleMetricsResults(client *influxdb.Client, sqlClient *psdb.Client, execUUID string) error {
	execMetrics, err := metrics.GetExecutionMetrics(*client, execUUID)
	if err != nil {
		return err
	}
	err = metrics.InsertExecutionMetrics(sqlClient, execUUID, execMetrics)
	if err != nil {
		return err
	}
	return nil
}

func handleSysBenchResults(resStr []byte, sqlClient *psdb.Client, macrobenchID int) error {
	// Parse results
	var results []Result
	err := json.Unmarshal(resStr, &results)
	if err != nil {
		return fmt.Errorf("unmarshal results: %+v\n", err)
	}
	if len(results) == 0 {
		return errors.New(ErrorNoSysBenchResult)
	}

	// Save results
	if sqlClient != nil {
		err = results[0].insertToMySQL(macrobenchID, sqlClient)
		if err != nil {
			return err
		}
	}
	return nil
}
