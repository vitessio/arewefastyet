/*
Copyright 2024 The Vitess Authors.

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

package macrobench

import (
	"github.com/vitessio/arewefastyet/go/storage"
	"golang.org/x/perf/benchmath"
)

type (
	statisticalSummary struct {
		Center     float64 `json:"center"`
		Confidence float64 `json:"confidence"`
	}

	statisticalResult struct {
		Insignificant bool               `json:"insignificant"`
		Delta         float64            `json:"delta"`
		P             float64            `json:"p"`
		N1            int                `json:"n1"`
		N2            int                `json:"n2"`
		Old           statisticalSummary `json:"old"`
		New           statisticalSummary `json:"new"`
	}

	// StatisticalCompareResults is the full representation of the results
	// obtained by comparing two samples using the Mann Whitney U Test.
	StatisticalCompareResults struct {
		TotalQPS  statisticalResult `json:"total_qps"`
		ReadsQPS  statisticalResult `json:"reads_qps"`
		WritesQPS statisticalResult `json:"writes_qps"`
		OtherQPS  statisticalResult `json:"other_qps"`

		TPS     statisticalResult `json:"tps"`
		Latency statisticalResult `json:"latency"`
		Errors  statisticalResult `json:"errors"`

		TotalComponentsCPUTime statisticalResult            `json:"total_components_cpu_time"`
		ComponentsCPUTime      map[string]statisticalResult `json:"components_cpu_time"`

		TotalComponentsMemStatsAllocBytes statisticalResult            `json:"total_components_mem_stats_alloc_bytes"`
		ComponentsMemStatsAllocBytes      map[string]statisticalResult `json:"components_mem_stats_alloc_bytes"`
	}

	StatisticalCompare struct {
		RightSHA   string
		LeftSHA    string
		Planner    PlannerVersion
		MacroTypes []string
	}
)

var (
	// defaultThresholds sets the thresholds below which the null hypothesis that two samples
	// come from the same distribution can be rejected.
	defaultThresholds = benchmath.DefaultThresholds

	// defaultConfidence sets the desired confidence interval when doing a summary of a sample.
	defaultConfidence = 0.95
)

func (sc StatisticalCompare) Compare(client storage.SQLClient) (map[string]StatisticalCompareResults, error) {
	results := make(map[string]StatisticalCompareResults, len(sc.MacroTypes))
	for _, macroType := range sc.MacroTypes {
		leftResult, err := GetBenchmarkResults(client, macroType, sc.LeftSHA, sc.Planner)
		if err != nil {
			return nil, err
		}

		rightResult, err := GetBenchmarkResults(client, macroType, sc.RightSHA, sc.Planner)
		if err != nil {
			return nil, err
		}

		leftResultsAsSlice := leftResult.asSlice()
		rightResultsAsSlice := rightResult.asSlice()

		scr := performAnalysis(leftResultsAsSlice, rightResultsAsSlice)
		results[macroType] = scr
	}
	return results, nil
}

func compare(old, new []float64) statisticalResult {
	var sr statisticalResult

	s1 := benchmath.NewSample(old, &defaultThresholds)
	s2 := benchmath.NewSample(new, &defaultThresholds)

	s1Summary := benchmath.AssumeNothing.Summary(s1, defaultConfidence)
	s2Summary := benchmath.AssumeNothing.Summary(s2, defaultConfidence)

	sr.Old = statisticalSummary{
		Center:     s1Summary.Center,
		Confidence: s1Summary.Confidence,
	}
	sr.New = statisticalSummary{
		Center:     s2Summary.Center,
		Confidence: s2Summary.Confidence,
	}

	c := benchmath.AssumeNothing.Compare(s1, s2)
	sr.P = c.P
	sr.N1 = c.N1
	sr.N2 = c.N2

	if c.P > c.Alpha {
		sr.Insignificant = true
	}

	switch {
	case sr.Old.Center == sr.New.Center || sr.Old.Center == 0:
		sr.Delta = 0
	default:
		sr.Delta = ((sr.New.Center / sr.Old.Center) - 1.0) * 100.0
	}
	return sr
}

func performAnalysis(old, new resultAsSlice) StatisticalCompareResults {
	scr := StatisticalCompareResults{
		ComponentsCPUTime:            map[string]statisticalResult{},
		ComponentsMemStatsAllocBytes: map[string]statisticalResult{},
	}

	scr.TotalQPS = compare(old.qps.total, new.qps.total)
	scr.ReadsQPS = compare(old.qps.reads, new.qps.reads)
	scr.WritesQPS = compare(old.qps.writes, new.qps.writes)
	scr.OtherQPS = compare(old.qps.other, new.qps.other)

	scr.TPS = compare(old.tps, new.tps)
	scr.Latency = compare(old.latency, new.latency)
	scr.Errors = compare(old.errors, new.errors)

	scr.TotalComponentsCPUTime = compare(old.metrics.totalComponentsCPUTime, new.metrics.totalComponentsCPUTime)
	for name, values := range old.metrics.componentsCPUTime {
		scr.ComponentsCPUTime[name] = compare(values, new.metrics.componentsCPUTime[name])
	}

	scr.TotalComponentsMemStatsAllocBytes = compare(old.metrics.totalComponentsMemStatsAllocBytes, new.metrics.totalComponentsMemStatsAllocBytes)
	for name, values := range old.metrics.componentsMemStatsAllocBytes {
		scr.ComponentsMemStatsAllocBytes[name] = compare(values, new.metrics.componentsMemStatsAllocBytes[name])
	}
	return scr
}
