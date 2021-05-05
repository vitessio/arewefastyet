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
)

var (
	components = []string{
		"vtgate",
		"vttablet",
	}
)

type ExecutionMetrics struct {
	TotalComponentsCPUTime int
	ComponentsCPUTime      map[string]int
}

func GetExecutionMetrics(client influxdb.Client, execUUID string) (ExecutionMetrics, error) {
	execMetrics := ExecutionMetrics{
		ComponentsCPUTime: map[string]int{},
	}

	for _, component := range components {
		cpuTimeForComponent, err := getCPUTimeForComponent(client, "0", "now()", execUUID, component)
		if err != nil {
			return ExecutionMetrics{}, err
		}
		execMetrics.ComponentsCPUTime[component] = int(cpuTimeForComponent)
		execMetrics.TotalComponentsCPUTime += int(cpuTimeForComponent)
	}
	return execMetrics, nil
}

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
