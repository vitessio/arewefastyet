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
	"fmt"

	"github.com/vitessio/arewefastyet/go/storage/influxdb"
	"github.com/vitessio/arewefastyet/go/storage/mysql"
)

// CompareMacroBenchmarks takes in 3 arguments, the database, and 2 SHAs. It reads from the database, the macrobenchmark
// results for the 2 SHAs and compares them. The result is a map with the key being the macrobenchmark name.
func CompareMacroBenchmarks(dbClient *mysql.Client, metricsClient *influxdb.Client, reference, compare string, planner PlannerVersion) (map[Type]interface{}, error) {
	// Get macro benchmarks from all the different types
	SHAs := []string{reference, compare}
	var err error
	macros := map[string]map[Type]DetailsArray{}
	for _, sha := range SHAs {
		macros[sha], err = GetDetailsArraysFromAllTypes(sha, planner, dbClient, metricsClient)
		if err != nil {
			return nil, err
		}
		for mtype := range macros[sha] {
			macros[sha][mtype] = macros[sha][mtype].ReduceSimpleMedian()
		}
	}
	macrosMatrixes := map[Type]interface{}{}
	for _, mtype := range Types {
		macrosMatrixes[mtype] = CompareDetailsArrays(macros[reference][mtype], macros[compare][mtype])
	}
	return macrosMatrixes, nil
}

// ComparePlanners takes in 2 arguments, the database, and a SHA. It reads from the database, the macrobenchmark
// results for the 2 planners corresponding to the sha and compares them. The result is a map with the key being the macrobenchmark name.
func ComparePlanners(dbClient *mysql.Client, metricsClient *influxdb.Client, sha string) (map[Type]interface{}, error) {
	// Get macro benchmarks from all the different types
	var err error
	macros := map[string]map[Type]DetailsArray{}
	for _, planner := range PlannerVersions {
		macros[string(planner)], err = GetDetailsArraysFromAllTypes(sha, planner, dbClient, metricsClient)
		if err != nil {
			return nil, err
		}
		for mtype := range macros[string(planner)] {
			macros[string(planner)][mtype] = macros[string(planner)][mtype].ReduceSimpleMedian()
		}
	}
	macrosMatrixes := map[Type]interface{}{}
	for _, mtype := range Types {
		macrosMatrixes[mtype] = CompareDetailsArrays(macros[string(Gen4FallbackPlanner)][mtype], macros[string(V3Planner)][mtype])
	}
	return macrosMatrixes, nil
}

// Regression returns a string containing the reason of the regression, if no regression is found, the string
// will be returned empty.
func (c Comparison) Regression() (reason string) {
	if c.DiffMetrics.TotalComponentsCPUTime <= -5.00 {
		reason += fmt.Sprintf("- Total CPU time increased by %.2f%% \n", c.DiffMetrics.TotalComponentsCPUTime*-1)
	}
	for key, value := range c.DiffMetrics.ComponentsCPUTime {
		if value <= -5.00 {
			reason += fmt.Sprintf("- %s CPU time increased by %.2f%% \n", key, value*-1)
		}
	}
	if c.Diff.TPS <= -10 {
		reason += fmt.Sprintf("- TPS decreased by %.2f%% \n", c.Diff.TPS*-1)
	}
	if c.Diff.QPS.Total <= -10 {
		reason += fmt.Sprintf("- QPS decreased by %.2f%% \n", c.Diff.QPS.Total*-1)
	}
	if c.Diff.Latency <= -10 {
		reason += fmt.Sprintf("- Latency increased by %.2f%% \n", c.Diff.Latency*-1)
	}
	return
}
