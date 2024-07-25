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

package metrics

import (
	"fmt"
	"strings"

	"github.com/vitessio/arewefastyet/go/storage"
	"github.com/vitessio/arewefastyet/go/storage/influxdb"
)

const (
	cpuSecondsPerComponentStart = `from(bucket:"%s")
			|> range(start: 0, stop: now())
			|> filter(fn:(r) => r._measurement == "process_cpu_seconds_total" and r.exec_uuid == "%s" and r.component == "%s")
			|> filter(fn: (r) => r._value > 0)
			|> min()`

	cpuSecondsPerComponentEnd = `from(bucket:"%s")
			|> range(start: 0, stop: now())
			|> filter(fn:(r) => r._measurement == "process_cpu_seconds_total" and r.exec_uuid == "%s" and r.component == "%s")
			|> max()`

	memAllocBytesPerComponentStart = `from(bucket:"%s")
			|> range(start: 0, stop: now())
			|> filter(fn:(r) => r._measurement == "go_memstats_alloc_bytes_total" and r.exec_uuid == "%s" and r.component == "%s")
			|> filter(fn: (r) => r._value > 0)
			|> min()`

	memAllocBytesPerComponentEnd = `from(bucket:"%s")
			|> range(start: 0, stop: now())
			|> filter(fn:(r) => r._measurement == "go_memstats_alloc_bytes_total" and r.exec_uuid == "%s" and r.component == "%s")
			|> max()`
)

var (
	components = []string{
		"vtgate",
		"vttablet",
	}
)

type (
	// ExecutionMetrics contains all the different system and service metrics
	// that were gathered during the execution of a benchmark.
	ExecutionMetrics struct {
		// The sum of the time taken by every component to run one query on average.
		TotalComponentsCPUTime float64

		// Map of string/float that contains the name of the component as a key
		// and the time taken by that component on average per query as a value.
		ComponentsCPUTime map[string]float64

		// TotalComponentsMemStatsAllocBytes represents the total number of bytes
		// allocated even if freed, by all the components of the execution and on average per query.
		// The underlying go metrics used is go_memstats_alloc_bytes_total.
		TotalComponentsMemStatsAllocBytes float64

		// ComponentsMemStatsAllocBytes represents the number of bytes allocated
		// and freed that each component used on average per query. The go metrics used is go_memstats_alloc_bytes_total.
		ComponentsMemStatsAllocBytes map[string]float64
	}

	// ExecutionMetricsArray is a slice of ExecutionMetrics, it has a Median method
	// to compute the overall median of the slice.
	ExecutionMetricsArray []ExecutionMetrics
)

// GetExecutionMetrics fetches and computes a single execution's metrics.
// Metrics are fetched using the given influxdb.Client and execUUID.
func GetExecutionMetrics(client influxdb.Client, execUUID string, queries int) (ExecutionMetrics, error) {
	execMetrics := NewExecMetrics()

	for _, component := range components {
		// CPU time
		endValue, err := getSumFloatValueForQuery(client, fmt.Sprintf(cpuSecondsPerComponentEnd, client.Config.Database, execUUID, component))
		if err != nil {
			return ExecutionMetrics{}, err
		}
		startValue, err := getSumFloatValueForQuery(client, fmt.Sprintf(cpuSecondsPerComponentStart, client.Config.Database, execUUID, component))
		if err != nil {
			return ExecutionMetrics{}, err
		}
		execMetrics.ComponentsCPUTime[component] = endValue - startValue
		execMetrics.TotalComponentsCPUTime += execMetrics.ComponentsCPUTime[component]

		// Memory
		endValue, err = getSumFloatValueForQuery(client, fmt.Sprintf(memAllocBytesPerComponentEnd, client.Config.Database, execUUID, component))
		if err != nil {
			return ExecutionMetrics{}, err
		}

		startValue, err = getSumFloatValueForQuery(client, fmt.Sprintf(memAllocBytesPerComponentStart, client.Config.Database, execUUID, component))
		if err != nil {
			return ExecutionMetrics{}, err
		}

		execMetrics.ComponentsMemStatsAllocBytes[component] = endValue - startValue
		execMetrics.TotalComponentsMemStatsAllocBytes += execMetrics.ComponentsMemStatsAllocBytes[component]
	}

	// Divide all metrics by the number of queries that were executed
	if queries > 0 {
		execMetrics.TotalComponentsCPUTime = execMetrics.TotalComponentsCPUTime / float64(queries)
		execMetrics.TotalComponentsMemStatsAllocBytes = execMetrics.TotalComponentsMemStatsAllocBytes / float64(queries)
		for key, val := range execMetrics.ComponentsCPUTime {
			execMetrics.ComponentsCPUTime[key] = val / float64(queries)
		}
		for key, val := range execMetrics.ComponentsMemStatsAllocBytes {
			execMetrics.ComponentsMemStatsAllocBytes[key] = val / float64(queries)
		}
	}
	return execMetrics, nil
}

func NewExecMetrics() ExecutionMetrics {
	return ExecutionMetrics{
		ComponentsCPUTime: map[string]float64{
			"vtgate":   0,
			"vttablet": 0,
		},
		ComponentsMemStatsAllocBytes: map[string]float64{
			"vtgate":   0,
			"vttablet": 0,
		},
	}
}

func InsertExecutionMetrics(client storage.SQLClient, execUUID string, execMetrics ExecutionMetrics) error {
	query := "INSERT INTO metrics(exec_uuid, `name`, `value`) VALUES (?, ?, ?), (?, ?, ?)"
	args := []interface{}{
		execUUID, "TotalComponentsCPUTime", execMetrics.TotalComponentsCPUTime,
		execUUID, "TotalComponentsMemStatsAllocBytes", execMetrics.TotalComponentsMemStatsAllocBytes,
	}
	for k, v := range execMetrics.ComponentsCPUTime {
		query += ", (?,?,?)"
		args = append(args, []interface{}{
			execUUID, "ComponentsCPUTime." + k, v,
		}...)
	}
	for k, v := range execMetrics.ComponentsMemStatsAllocBytes {
		query += ", (?,?,?)"
		args = append(args, []interface{}{
			execUUID, "ComponentsMemStatsAllocBytes." + k, v,
		}...)
	}
	_, err := client.Write(query, args...)
	return err
}

func GetExecutionMetricsSQL(client storage.SQLClient, execUUID string) (ExecutionMetrics, error) {
	query := "select `name`, value from metrics where exec_uuid = ?"
	rows, err := client.Read(query, execUUID)
	if err != nil {
		return ExecutionMetrics{}, err
	}
	defer rows.Close()

	result := NewExecMetrics()
	for rows.Next() {
		var name string
		var value float64
		err = rows.Scan(&name, &value)
		if err != nil {
			return ExecutionMetrics{}, err
		}
		switch {
		case name == "TotalComponentsCPUTime":
			result.TotalComponentsCPUTime = value
		case name == "TotalComponentsMemStatsAllocBytes":
			result.TotalComponentsMemStatsAllocBytes = value
		case strings.HasPrefix(name, "ComponentsCPUTime."):
			key := strings.Split(name, ".")[1]
			result.ComponentsCPUTime[key] = value
		case strings.HasPrefix(name, "ComponentsMemStatsAllocBytes."):
			key := strings.Split(name, ".")[1]
			result.ComponentsMemStatsAllocBytes[key] = value
		}
	}
	return result, nil
}

// getSumFloatValueForQuery return the sum of a float value based on the given query, for
// each row.
func getSumFloatValueForQuery(client influxdb.Client, query string) (float64, error) {
	result, err := client.Select(query)
	if err != nil {
		return 0, err
	}

	res := 0.0
	for _, value := range result {
		res += value["_value"].(float64)
	}
	return res, nil
}
