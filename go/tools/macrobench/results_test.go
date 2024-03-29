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
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/vitessio/arewefastyet/go/exec/metrics"
)

func TestMacroBenchmarkResultsArray_mergeMedian(t *testing.T) {
	qpsOfOne := *newQPS(1.0, 1.0, 1.0, 1.0)
	qpsOfTwo := *newQPS(2.0, 2.0, 2.0, 2.0)
	resultOfOne := *newResult(qpsOfOne, 1.0, 1.0, 1.0, 1.0, 1, 1.0)
	resultOfTwo := *newResult(qpsOfTwo, 2.0, 2.0, 2.0, 2.0, 2, 2.0)

	var tests = []struct {
		name             string
		mrs              ResultsArray
		wantMergedResult Result
	}{
		{name: "Single result in array", mrs: ResultsArray{resultOfOne}, wantMergedResult: resultOfOne},
		{name: "Even number of results in array", mrs: ResultsArray{resultOfOne, resultOfTwo}, wantMergedResult: *newResult(*newQPS(1.5, 1.5, 1.5, 1.5), 1.5, 1.5, 1.5, 1.5, 1, 1.5)},
		{name: "Multiple results in array", mrs: ResultsArray{resultOfOne, resultOfOne, resultOfOne}, wantMergedResult: resultOfOne},
		{name: "Multiple and different results in array", mrs: ResultsArray{
			*newResult(*newQPS(1.0, 1.0, 1.0, 3), 1.0, 1.0, 1.0, 1.5, 1, 10.0),
			*newResult(*newQPS(2.0, 5.0, 1.5, 6), 2.0, 5.0, 3.0, 2.5, 1000, 20.0),
			*newResult(*newQPS(3.0, 10.0, 2.0, 9), 3.0, 10.0, 2.0, 3.5, 500, 30.0),
		}, wantMergedResult: *newResult(*newQPS(2.0, 5.0, 1.5, 6), 2.0, 5.0, 2.0, 2.5, 500, 20.0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			gotMergedResult := tt.mrs.mergeMedian()

			c.Assert(gotMergedResult, qt.DeepEquals, tt.wantMergedResult)
		})
	}
}

func TestMacroBenchmarkDetailsArray_ReduceSimpleMedian(t *testing.T) {
	qpsOfOne := *newQPS(1.0, 1.0, 1.0, 1.0)
	qpsOfTwo := *newQPS(2.0, 2.0, 2.0, 2.0)
	qpsOfOneHalf := *newQPS(1.5, 1.5, 1.5, 1.5)

	resultOfOne := *newResult(qpsOfOne, 1.0, 1.0, 1.0, 1.0, 1, 1.0)
	resultOfTwo := *newResult(qpsOfTwo, 2.0, 2.0, 2.0, 2.0, 2, 2.0)
	resultOfOneHalf := *newResult(qpsOfOneHalf, 1.5, 1.5, 1.5, 1.5, 1, 1.5)

	tests := []struct {
		name           string
		mabd           DetailsArray
		wantReduceMabd DetailsArray
	}{
		{name: "Few elements with same git ref", mabd: []Details{
			*newDetails(*newBenchmarkID(1, "webhook", nil), "11bbAAA", resultOfOne, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
			*newDetails(*newBenchmarkID(2, "webhook", nil), "11bbAAA", resultOfTwo, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
		}, wantReduceMabd: []Details{
			*newDetails(BenchmarkID{}, "11bbAAA", resultOfOneHalf, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
		}},

		{name: "Few elements with different git refs", mabd: []Details{
			*newDetails(*newBenchmarkID(1, "webhook", nil), "11bbAAA", resultOfOne, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
			*newDetails(*newBenchmarkID(2, "webhook", nil), "11bbAAA", resultOfTwo, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
			*newDetails(*newBenchmarkID(3, "api_call", nil), "f78gh1p", resultOfOne, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
			*newDetails(*newBenchmarkID(4, "webhook", nil), "f78gh1p", resultOfTwo, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
		}, wantReduceMabd: []Details{
			*newDetails(BenchmarkID{}, "11bbAAA", resultOfOneHalf, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
			*newDetails(BenchmarkID{}, "f78gh1p", resultOfOneHalf, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			gotReduceMabd := tt.mabd.ReduceSimpleMedian()
			c.Assert(gotReduceMabd, qt.DeepEquals, tt.wantReduceMabd)
		})
	}
}

func BenchmarkReduceSimpleMedian(b *testing.B) {
	qpsOfOne := *newQPS(1.0, 1.0, 1.0, 1.0)
	qpsOfTwo := *newQPS(2.0, 2.0, 2.0, 2.0)
	resultOfOne := *newResult(qpsOfOne, 1.0, 1.0, 1.0, 1.0, 1, 1.0)
	resultOfTwo := *newResult(qpsOfTwo, 2.0, 2.0, 2.0, 2.0, 2, 2.0)

	mabd := DetailsArray{
		*newDetails(*newBenchmarkID(1, "webhook", nil), "11bbAAA", resultOfOne, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}}),
		*newDetails(*newBenchmarkID(2, "webhook", nil), "11bbAAA", resultOfTwo, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}}),
		*newDetails(*newBenchmarkID(3, "api_call", nil), "f78gh1p", resultOfOne, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}}),
		*newDetails(*newBenchmarkID(4, "webhook", nil), "f78gh1p", resultOfTwo, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}}),
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		got := mabd.ReduceSimpleMedian()
		if len(got) != 2 {
			b.Error("benchmark results failed: result must contain 2 elements")
		}
	}
}
