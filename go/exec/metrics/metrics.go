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
	"github.com/vitessio/arewefastyet/go/storage/influxdb"
	awftmath "github.com/vitessio/arewefastyet/go/tools/math"
)

const (
	cpuSecondsPerComponent = `from(bucket:"%s")
			|> range(start: %s, stop: %s)
			|> filter(fn:(r) => r._measurement == "process_cpu_seconds_total" and r.exec_uuid == "%s" and r.component == "%s")
			|> max()`

	memAllocBytesPerComponent = `from(bucket:"%s")
			|> range(start: %s, stop: %s)
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
		// The sum of the time taken by every component.
		TotalComponentsCPUTime float64

		// Map of string/float that contains the name of the component as a key
		// and the time taken by that component as a value.
		ComponentsCPUTime map[string]float64

		// TotalComponentsMemStatsAllocBytes represents the total number of bytes
		// allocated even if freed, by all the components of the execution.
		// The underlying go metrics used is go_memstats_alloc_bytes_total.
		TotalComponentsMemStatsAllocBytes float64

		// ComponentsMemStatsAllocBytes represents the number of bytes allocated
		// and freed that each component used. The go metrics used is go_memstats_alloc_bytes_total.
		ComponentsMemStatsAllocBytes map[string]float64
	}

	// ExecutionMetricsArray is a slice of ExecutionMetrics, it has a Median method
	// to compute the overall median of the slice.
	ExecutionMetricsArray []ExecutionMetrics
)

// GetExecutionMetrics fetches and computes a single execution's metrics.
// Metrics are fetched using the given influxdb.Client and execUUID.
func GetExecutionMetrics(client influxdb.Client, execUUID string) (ExecutionMetrics, error) {
	execMetrics := ExecutionMetrics{
		ComponentsCPUTime: map[string]float64{},
		ComponentsMemStatsAllocBytes: map[string]float64{},
	}

	var err error
	for _, component := range components {
		execMetrics.ComponentsCPUTime[component], err = getSumFloatValueForQuery(client, fmt.Sprintf(cpuSecondsPerComponent, client.Config.Database, "0", "now()", execUUID, component))
		if err != nil {
			return ExecutionMetrics{}, err
		}
		execMetrics.TotalComponentsCPUTime += execMetrics.ComponentsCPUTime[component]

		execMetrics.ComponentsMemStatsAllocBytes[component], err = getSumFloatValueForQuery(client, fmt.Sprintf(memAllocBytesPerComponent, client.Config.Database, "0", "now()", execUUID, component))
		if err != nil {
			return ExecutionMetrics{}, err
		}
		execMetrics.TotalComponentsMemStatsAllocBytes += execMetrics.ComponentsMemStatsAllocBytes[component]
	}
	return execMetrics, nil
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

// Median computes the median of the ExecutionMetricsArray.
// It returns an ExecutionMetrics struct containing the medians.
func (metricsArray ExecutionMetricsArray) Median() ExecutionMetrics {
	interResults := struct {
		totalComponentsCPUTime []float64
		componentsCPUTime      map[string][]float64

		totalComponentsMemStatsAllocBytes []float64
		componentsMemStatsAllocBytes      map[string][]float64
	}{
		totalComponentsCPUTime: []float64{},
		componentsCPUTime:      map[string][]float64{},

		totalComponentsMemStatsAllocBytes: []float64{},
		componentsMemStatsAllocBytes:      map[string][]float64{},
	}

	// Append all the metrics into interResults
	for _, metrics := range metricsArray {
		// If an execution is missing metrics, we do not count it toward
		// the median of all execution.
		if metrics.TotalComponentsCPUTime == 0 {
			continue
		}
		interResults.totalComponentsCPUTime = append(interResults.totalComponentsCPUTime, metrics.TotalComponentsCPUTime)
		for component, value := range metrics.ComponentsCPUTime {
			interResults.componentsCPUTime[component] = append(interResults.componentsCPUTime[component], value)
		}

		interResults.totalComponentsMemStatsAllocBytes = append(interResults.totalComponentsMemStatsAllocBytes, metrics.TotalComponentsMemStatsAllocBytes)
		for component, value := range metrics.ComponentsMemStatsAllocBytes {
			interResults.componentsMemStatsAllocBytes[component] = append(interResults.componentsMemStatsAllocBytes[component], value)
		}
	}
	result := ExecutionMetrics{
		ComponentsCPUTime: map[string]float64{},
		ComponentsMemStatsAllocBytes: map[string]float64{},
	}
	result.TotalComponentsCPUTime = awftmath.MedianFloat(interResults.totalComponentsCPUTime)
	for component, value := range interResults.componentsCPUTime {
		result.ComponentsCPUTime[component] = awftmath.MedianFloat(value)
	}

	result.TotalComponentsMemStatsAllocBytes = awftmath.MedianFloat(interResults.totalComponentsMemStatsAllocBytes)
	for component, value := range interResults.componentsMemStatsAllocBytes {
		result.ComponentsMemStatsAllocBytes[component] = awftmath.MedianFloat(value)
	}
	return result
}

// CompareTwo computes the percentage decrease between left and right.
// If left is equal to 20 and right is equal to 10, then the decrease will be 50%.
// The percentage are returned through ExecutionMetrics.
func CompareTwo(left, right ExecutionMetrics) ExecutionMetrics {
	result := ExecutionMetrics{
		ComponentsCPUTime: map[string]float64{},
		ComponentsMemStatsAllocBytes: map[string]float64{},
	}
	result.TotalComponentsCPUTime = compareSafe(left.TotalComponentsCPUTime, right.TotalComponentsCPUTime)
	result.TotalComponentsMemStatsAllocBytes = compareSafe(left.TotalComponentsMemStatsAllocBytes, right.TotalComponentsMemStatsAllocBytes)
	result.ComponentsCPUTime = compareSafeComponentMap(left.ComponentsCPUTime, right.ComponentsCPUTime)
	result.ComponentsMemStatsAllocBytes = compareSafeComponentMap(left.ComponentsMemStatsAllocBytes, right.ComponentsMemStatsAllocBytes)
	return result
}

func compareSafeComponentMap(left, right map[string]float64) map[string]float64 {
	result := map[string]float64{}
	for component := range left {
		result[component] = 0
		if _, ok := right[component]; !ok {
			continue
		}
		result[component] = compareSafe(left[component], right[component])
	}
	for component := range right {
		if _, ok := left[component]; ok {
			continue
		}
		result[component] = compareSafe(0, right[component])
	}
	return result
}

// Compare the decrease between left and right.
// The more decrease, the higher result will be.
// Ex: left=100, right=50, we decreased by 1/2, thus result=50
func compareSafe(left, right float64) (result float64) {
	if left != 0 {
		result = (left - right) / left * 100
	} else if right > 0 {
		result = -100
	}
	return
}
