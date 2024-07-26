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
	"sort"
	"sync"
	"time"

	"github.com/vitessio/arewefastyet/go/exec/metrics"
	"github.com/vitessio/arewefastyet/go/storage"
)

type (
	// sysbenchQPS is a subset of the results we get from sysbench, it only contains QPS information
	sysbenchQPS struct {
		ID     int
		RefID  int
		Total  float64 `json:"total"`
		Reads  float64 `json:"reads"`
		Writes float64 `json:"writes"`
		Other  float64 `json:"other"`
	}

	// sysbenchResult is the full representation of the results we get after executing sysbench once.
	sysbenchResult struct {
		ID         int
		Queries    int         `json:"queries"`
		QPS        sysbenchQPS `json:"qps"`
		TPS        float64     `json:"tps"`
		Latency    float64     `json:"latency"`
		Errors     float64     `json:"errors"`
		Reconnects float64     `json:"reconnects"`
		Time       int         `json:"time"`
		Threads    float64     `json:"threads"`
	}

	sysbenchResultArray []sysbenchResult

	qpsAsSlice struct {
		total  []float64
		reads  []float64
		writes []float64
		other  []float64
	}

	metricsAsSlice struct {
		totalComponentsCPUTime []float64
		componentsCPUTime      map[string][]float64

		totalComponentsMemStatsAllocBytes []float64
		componentsMemStatsAllocBytes      map[string][]float64
	}

	executionGroupResultsAsSlice struct {
		qps qpsAsSlice

		tps        []float64
		latency    []float64
		errors     []float64
		reconnects []float64
		time       []int
		threads    []float64

		metrics metricsAsSlice
	}

	executionGroupResults struct {
		GitRef  string
		Results sysbenchResultArray
		Metrics metrics.ExecutionMetricsArray
	}

	// executionID is used to identify a macro benchmark using its database's ID, the
	// source from which the benchmark was triggered and its creation date.
	executionID struct {
		ID        int
		Source    string
		CreatedAt *time.Time
		ExecUUID  string
	}

	// executionResults is the full representation of a single execution,
	// it contains the ID and the results (both sysbench and metrics results)
	executionResults struct {
		executionID

		GitRef  string
		Result  sysbenchResult
		Metrics metrics.ExecutionMetrics
	}
)

func (br executionGroupResults) asSlice() executionGroupResultsAsSlice {
	s := br.Results.resultsArrayToSlice()
	s.metrics = metricsToSlice(br.Metrics)
	return s
}

func (br executionGroupResults) toStatisticalSingleResult() StatisticalSingleResult {
	ssr := StatisticalSingleResult{
		GitRef:                       br.GitRef,
		ComponentsCPUTime:            map[string]StatisticalSummary{},
		ComponentsMemStatsAllocBytes: map[string]StatisticalSummary{},
	}

	resultSlice := br.asSlice()

	ssr.TotalQPS, _ = getSummary(resultSlice.qps.total)
	ssr.ReadsQPS, _ = getSummary(resultSlice.qps.reads)
	ssr.WritesQPS, _ = getSummary(resultSlice.qps.writes)
	ssr.OtherQPS, _ = getSummary(resultSlice.qps.other)

	ssr.TPS, _ = getSummary(resultSlice.tps)
	ssr.Latency, _ = getSummary(resultSlice.latency)
	ssr.Errors, _ = getSummary(resultSlice.errors)

	ssr.TotalComponentsCPUTime, _ = getSummary(resultSlice.metrics.totalComponentsCPUTime)
	for name, value := range resultSlice.metrics.componentsCPUTime {
		ssr.ComponentsCPUTime[name], _ = getSummary(value)
	}

	ssr.TotalComponentsMemStatsAllocBytes, _ = getSummary(resultSlice.metrics.totalComponentsMemStatsAllocBytes)
	for name, value := range resultSlice.metrics.componentsMemStatsAllocBytes {
		ssr.ComponentsMemStatsAllocBytes[name], _ = getSummary(value)
	}
	return ssr
}

func (br executionGroupResults) toShortStatisticalSingleResult() ShortStatisticalSingleResult {
	var sssr ShortStatisticalSingleResult

	resultSlice := br.asSlice()

	sssr.TotalQPS, _ = getSummary(resultSlice.qps.total)
	return sssr
}

func metricsToSlice(metrics metrics.ExecutionMetricsArray) metricsAsSlice {
	var s metricsAsSlice
	s.componentsCPUTime = make(map[string][]float64)
	s.componentsMemStatsAllocBytes = make(map[string][]float64)
	for _, metricRow := range metrics {
		s.totalComponentsCPUTime = append(s.totalComponentsCPUTime, metricRow.TotalComponentsCPUTime)
		for name, value := range metricRow.ComponentsCPUTime {
			s.componentsCPUTime[name] = append(s.componentsCPUTime[name], value)
		}

		s.totalComponentsMemStatsAllocBytes = append(s.totalComponentsMemStatsAllocBytes, metricRow.TotalComponentsMemStatsAllocBytes)
		for name, value := range metricRow.ComponentsMemStatsAllocBytes {
			s.componentsMemStatsAllocBytes[name] = append(s.componentsMemStatsAllocBytes[name], value)
		}
	}
	return s
}

func (mrs sysbenchResultArray) resultsArrayToSlice() executionGroupResultsAsSlice {
	var ras executionGroupResultsAsSlice
	for _, mr := range mrs {
		ras.qps.total = append(ras.qps.total, mr.QPS.Total)
		ras.qps.reads = append(ras.qps.reads, mr.QPS.Reads)
		ras.qps.writes = append(ras.qps.writes, mr.QPS.Writes)
		ras.qps.other = append(ras.qps.other, mr.QPS.Other)
		ras.tps = append(ras.tps, mr.TPS)
		ras.latency = append(ras.latency, mr.Latency)
		ras.errors = append(ras.errors, mr.Errors)
		ras.reconnects = append(ras.reconnects, mr.Reconnects)
		ras.time = append(ras.time, mr.Time)
		ras.threads = append(ras.threads, mr.Threads)
	}
	sort.Float64s(ras.qps.total)
	sort.Float64s(ras.qps.reads)
	sort.Float64s(ras.qps.writes)
	sort.Float64s(ras.qps.other)
	sort.Float64s(ras.tps)
	sort.Float64s(ras.latency)
	sort.Float64s(ras.reconnects)
	sort.Ints(ras.time)
	sort.Float64s(ras.threads)
	return ras
}

func Compare(client storage.SQLClient, old, new string, types []string, planner PlannerVersion) (map[string]StatisticalCompareResults, error) {
	results := make(map[string]StatisticalCompareResults, len(types))
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	var err error
	for _, macroType := range types {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var oldResult, newResult executionGroupResults
			oldResult, err = getExecutionGroupResults(macroType, old, planner, client)
			if err != nil {
				return
			}

			newResult, err = getExecutionGroupResults(macroType, new, planner, client)
			if err != nil {
				return
			}

			if len(oldResult.Results) == 0 && len(newResult.Results) == 0 {
				mu.Lock()
				defer mu.Unlock()
				results[macroType] = StatisticalCompareResults{
					ComponentsCPUTime: map[string]StatisticalResult{
						"vtgate":   {},
						"vttablet": {},
					},
					ComponentsMemStatsAllocBytes: map[string]StatisticalResult{
						"vtgate":   {},
						"vttablet": {},
					},
				}
				return
			}

			oldResultsAsSlice := oldResult.asSlice()
			newResultsAsSlice := newResult.asSlice()

			scr := performAnalysis(oldResultsAsSlice, newResultsAsSlice)

			mu.Lock()
			defer mu.Unlock()
			results[macroType] = scr
		}()
	}
	wg.Wait()
	return results, err
}

func CompareFKs(client storage.SQLClient, oldWorkload, newWorkload string, sha string, planner PlannerVersion) (StatisticalCompareResults, error) {
	oldResult, err := getExecutionGroupResults(oldWorkload, sha, planner, client)
	if err != nil {
		return StatisticalCompareResults{}, err
	}

	newResult, err := getExecutionGroupResults(newWorkload, sha, planner, client)
	if err != nil {
		return StatisticalCompareResults{}, err
	}

	if len(oldResult.Results) == 0 && len(newResult.Results) == 0 {
		return StatisticalCompareResults{
			ComponentsCPUTime: map[string]StatisticalResult{
				"vtgate":   {},
				"vttablet": {},
			},
			ComponentsMemStatsAllocBytes: map[string]StatisticalResult{
				"vtgate":   {},
				"vttablet": {},
			},
		}, nil
	}

	oldResultsAsSlice := oldResult.asSlice()
	newResultsAsSlice := newResult.asSlice()

	scr := performAnalysis(oldResultsAsSlice, newResultsAsSlice)

	return scr, nil
}

func Search(client storage.SQLClient, sha string, types []string, planner PlannerVersion) (map[string]StatisticalSingleResult, error) {
	results := make(map[string]StatisticalSingleResult, len(types))
	for _, macroType := range types {
		result, err := getExecutionGroupResults(macroType, sha, planner, client)
		if err != nil {
			return nil, err
		}
		if len(result.Results) == 0 {
			results[macroType] = StatisticalSingleResult{
				ComponentsCPUTime: map[string]StatisticalSummary{
					"vtgate":   {},
					"vttablet": {},
				},
				ComponentsMemStatsAllocBytes: map[string]StatisticalSummary{
					"vtgate":   {},
					"vttablet": {},
				},
			}
			continue
		}
		results[macroType] = result.toStatisticalSingleResult()
	}
	return results, nil
}

func SearchForLast30Days(client storage.SQLClient, macroType string, planner PlannerVersion) ([]StatisticalSingleResult, error) {
	var ssrs []StatisticalSingleResult
	results, err := getExecutionGroupResultsFromLast30Days(macroType, planner, client)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		ssrs = append(ssrs, result.toStatisticalSingleResult())
	}
	return ssrs, nil
}

func SearchForLast30DaysQPSOnly(client storage.SQLClient, types []string, planner PlannerVersion, days int) (map[string][]ShortStatisticalSingleResult, error) {
	results := make(map[string][]ShortStatisticalSingleResult)
	for _, macroType := range types {
		resultsForType, err := getSummaryLast30Days(macroType, planner, client)
		if err != nil {
			return nil, err
		}

		for _, result := range resultsForType {
			// If we do not have a decent number of results in the set of benchmark, let's skip the result.
			if len(result.Results) < 6 {
				continue
			}
			results[macroType] = append(results[macroType], result.toShortStatisticalSingleResult())
		}
	}
	return results, nil
}
