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
	"math"

	"github.com/aclements/go-moremath/mathx"
	"golang.org/x/perf/benchmath"
)

type (
	Range struct {
		Infinite bool    `json:"infinite"`
		Unknown  bool    `json:"unknown"`
		Value    float64 `json:"value"`
	}

	StatisticalSummary struct {
		Center     float64 `json:"center"`
		Confidence float64 `json:"confidence"`
		Range      Range   `json:"range"`
	}

	StatisticalResult struct {
		Insignificant bool               `json:"insignificant"`
		Delta         float64            `json:"delta"`
		P             float64            `json:"p"`
		N1            int                `json:"n1"`
		N2            int                `json:"n2"`
		Old           StatisticalSummary `json:"old"`
		New           StatisticalSummary `json:"new"`
	}

	ShortStatisticalSingleResult struct {
		TotalQPS StatisticalSummary `json:"total_qps"`
	}

	// StatisticalSingleResult represents a single benchmark's statistical summary.
	StatisticalSingleResult struct {
		GitRef string `json:"git_ref"`

		TotalQPS  StatisticalSummary `json:"total_qps"`
		ReadsQPS  StatisticalSummary `json:"reads_qps"`
		WritesQPS StatisticalSummary `json:"writes_qps"`
		OtherQPS  StatisticalSummary `json:"other_qps"`

		TPS     StatisticalSummary `json:"tps"`
		Latency StatisticalSummary `json:"latency"`
		Errors  StatisticalSummary `json:"errors"`

		TotalComponentsCPUTime StatisticalSummary            `json:"total_components_cpu_time"`
		ComponentsCPUTime      map[string]StatisticalSummary `json:"components_cpu_time"`

		TotalComponentsMemStatsAllocBytes StatisticalSummary            `json:"total_components_mem_stats_alloc_bytes"`
		ComponentsMemStatsAllocBytes      map[string]StatisticalSummary `json:"components_mem_stats_alloc_bytes"`
	}

	// StatisticalCompareResults is the full representation of the results
	// obtained by comparing two samples using the Mann Whitney U Test.
	StatisticalCompareResults struct {
		TotalQPS  StatisticalResult `json:"total_qps"`
		ReadsQPS  StatisticalResult `json:"reads_qps"`
		WritesQPS StatisticalResult `json:"writes_qps"`
		OtherQPS  StatisticalResult `json:"other_qps"`

		TPS     StatisticalResult `json:"tps"`
		Latency StatisticalResult `json:"latency"`
		Errors  StatisticalResult `json:"errors"`

		TotalComponentsCPUTime StatisticalResult            `json:"total_components_cpu_time"`
		ComponentsCPUTime      map[string]StatisticalResult `json:"components_cpu_time"`

		TotalComponentsMemStatsAllocBytes StatisticalResult            `json:"total_components_mem_stats_alloc_bytes"`
		ComponentsMemStatsAllocBytes      map[string]StatisticalResult `json:"components_mem_stats_alloc_bytes"`
	}
)

var (
	// defaultThresholds sets the thresholds below which the null hypothesis that two samples
	// come from the same distribution can be rejected.
	defaultThresholds = benchmath.DefaultThresholds

	// defaultConfidence sets the desired confidence interval when doing a summary of a sample.
	defaultConfidence = 0.95
)

func getRangeFromSummary(s benchmath.Summary) Range {
	if math.IsInf(s.Lo, 0) || math.IsInf(s.Hi, 0) {
		return Range{Infinite: true}
	}

	// If the signs of the bounds differ from the center, we can't
	// render it as a percent.
	var csign = mathx.Sign(s.Center)
	if csign != mathx.Sign(s.Lo) || csign != mathx.Sign(s.Hi) {
		return Range{Unknown: true}
	}

	// If center is 0, avoid dividing by zero. But we can only get
	// here if lo and hi are also 0, in which case is seems
	// reasonable to call this 0%.
	if s.Center == 0 {
		return Range{Value: 0.00}
	}

	// Phew. Compute the range percent.
	v := math.Max(s.Hi/s.Center-1, 1-s.Lo/s.Center)
	return Range{Value: v * 100}
}

func getSummary(values []float64) (StatisticalSummary, *benchmath.Sample) {
	sample := benchmath.NewSample(values, &defaultThresholds)
	summary := benchmath.AssumeNothing.Summary(sample, defaultConfidence)
	return StatisticalSummary{
		Center:     summary.Center,
		Confidence: summary.Confidence,
		Range:      getRangeFromSummary(summary),
	}, sample
}

func compare(old, new []float64) StatisticalResult {
	var sr StatisticalResult

	ssOld, s1 := getSummary(old)
	ssNew, s2 := getSummary(new)

	sr.Old = ssOld
	sr.New = ssNew

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
		ComponentsCPUTime:            map[string]StatisticalResult{},
		ComponentsMemStatsAllocBytes: map[string]StatisticalResult{},
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
