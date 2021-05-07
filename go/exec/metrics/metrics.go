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
	}

	for _, component := range components {
		cpuTimeForComponent, err := getCPUTimeForComponent(client, "0", "now()", execUUID, component)
		if err != nil {
			return ExecutionMetrics{}, err
		}
		execMetrics.ComponentsCPUTime[component] = cpuTimeForComponent
		execMetrics.TotalComponentsCPUTime += cpuTimeForComponent
	}
	return execMetrics, nil
}

// getCPUTimeForComponent return the CPU Time taken by a component (vtgate, vttablet, ...)
// for the given execUUID. The range of the select is defined through the start and end
// arguments, to select the whole span one can use: "start:0, end:now()".
func getCPUTimeForComponent(client influxdb.Client, start, end, execUUID, component string) (float64, error) {
	result, err := client.Select(fmt.Sprintf(`from(bucket:"%s")
			|> range(start: %s, stop: %s)
			|> filter(fn:(r) => r._measurement == "process_cpu_seconds_total" and r.exec_uuid == "%s" and r.component == "%s")
			|> max()`,
		client.Config.Database, start, end, execUUID, component))
	if err != nil {
		return 0, err
	}

	time := 0.0
	for _, value := range result {
		time += value["_value"].(float64)
	}
	return time, nil
}

// Median computes the median of the ExecutionMetricsArray.
// It returns an ExecutionMetrics struct containing the medians.
func (metricsArray ExecutionMetricsArray) Median() ExecutionMetrics {
	interResults := struct {
		totalComponentsCPUTime []float64
		componentsCPUTime      map[string][]float64
	}{
		totalComponentsCPUTime: []float64{},
		componentsCPUTime:      map[string][]float64{},
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
	}
	result := ExecutionMetrics{
		ComponentsCPUTime: map[string]float64{},
	}
	result.TotalComponentsCPUTime = awftmath.MedianFloat(interResults.totalComponentsCPUTime)
	for component, value := range interResults.componentsCPUTime {
		result.ComponentsCPUTime[component] = awftmath.MedianFloat(value)
	}
	return result
}

// CompareTwo computes the percentage decrease between left and right.
// If left is equal to 20 and right is equal to 10, then the decrease will be 50%.
// The percentage are returned through ExecutionMetrics.
func CompareTwo(left, right ExecutionMetrics) ExecutionMetrics {
	result := ExecutionMetrics{
		ComponentsCPUTime: map[string]float64{},
	}
	if left.TotalComponentsCPUTime != 0 {
		result.TotalComponentsCPUTime = (left.TotalComponentsCPUTime - right.TotalComponentsCPUTime) / left.TotalComponentsCPUTime * 100
	} else if right.TotalComponentsCPUTime > 0 {
		result.TotalComponentsCPUTime = -100
	}
	for component, value := range left.ComponentsCPUTime {
		result.ComponentsCPUTime[component] = 0
		if _, ok := right.ComponentsCPUTime[component]; !ok {
			continue
		}
		if value != 0 {
			result.ComponentsCPUTime[component] = (value - right.ComponentsCPUTime[component]) / value * 100
		} else if right.ComponentsCPUTime[component] > 0 {
			result.ComponentsCPUTime[component] = -100
		}
	}
	awftmath.CheckForNaN(&result, 0)
	awftmath.CheckForInf(&result, 0)
	return result
}
