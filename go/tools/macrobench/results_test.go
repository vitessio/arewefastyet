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
	qt "github.com/frankban/quicktest"
	"github.com/vitessio/arewefastyet/go/exec/metrics"
	"testing"
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
			*newDetails(BenchmarkID{}, "11bbAAA", resultOfOneHalf, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes:      map[string]float64{}}),
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

func TestHumanReadableStrings(t *testing.T) {
	c := qt.New(t)
	r := Result{
		QPS: QPS{
			Total:  72029.0,
			Reads:  45018.1,
			Writes: 18007.2,
			Other:  9003.6,
		},
		TPS:        4501.8,
		Latency:    189848.3,
		Errors:     8999.6,
		Reconnects: 7984.3,
		Time:       11064,
		Threads:    1107794.12,
	}
	c.Assert(r.QPS.TotalStr(), qt.Equals, "72,029.0")
	c.Assert(r.QPS.ReadsStr(), qt.Equals, "45,018.1")
	c.Assert(r.QPS.WritesStr(), qt.Equals, "18,007.2")
	c.Assert(r.QPS.OtherStr(), qt.Equals, "9,003.6")
	c.Assert(r.TPSStr(), qt.Equals, "4,501.8")
	c.Assert(r.LatencyStr(), qt.Equals, "189,848.3")
	c.Assert(r.ErrorsStr(), qt.Equals, "8,999.6")
	c.Assert(r.ReconnectsStr(), qt.Equals, "7,984.3")
	c.Assert(r.TimeStr(), qt.Equals, "11,064")
	c.Assert(r.ThreadsStr(), qt.Equals, "1,107,794.1")
}

func TestCompareDetailsArrays(t *testing.T) {
	type args struct {
		references DetailsArray
		compares   DetailsArray
	}
	qpsOfOne := *newQPS(1.0, 1.0, 1.0, 1.0)
	qpsOfTwo := *newQPS(2.0, 2.0, 2.0, 2.0)
	resultOfOne := *newResult(qpsOfOne, 1.0, 1.0, 1.0, 1.0, 1, 1.0)
	resultOfTwo := *newResult(qpsOfTwo, 2.0, 2.0, 2.0, 2.0, 2, 2.0)
	qpsOfFifty := *newQPS(-100, -100, -100, -100)
	resultOfFifty := *newResult(qpsOfFifty, -100, 50, -100, -100, -100, -100)

	tests := []struct {
		name         string
		args         args
		wantCompared ComparisonArray
	}{
		{name: "Simple comparison array", args: args{
			references: DetailsArray{
				*newDetails(*newBenchmarkID(1, "webhook", nil), "11bbAAA", resultOfOne, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
			},
			compares: DetailsArray{
				*newDetails(*newBenchmarkID(2, "webhook", nil), "11bbAAA", resultOfTwo, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
			},
		}, wantCompared: ComparisonArray{
			Comparison{
				Reference:   *newDetails(*newBenchmarkID(1, "webhook", nil), "11bbAAA", resultOfOne, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
				Compare:     *newDetails(*newBenchmarkID(2, "webhook", nil), "11bbAAA", resultOfTwo, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
				Diff:        resultOfFifty,
				DiffMetrics: metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}},
			},
		}},

		{name: "Simple comparison array with multiple sources", args: args{
			references: DetailsArray{
				*newDetails(*newBenchmarkID(1, "webhook", nil), "11bbAAA", resultOfOne, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
				*newDetails(*newBenchmarkID(4, "webhook", nil), "f78gh1p", resultOfTwo, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
			},
			compares: DetailsArray{
				*newDetails(*newBenchmarkID(2, "webhook", nil), "11bbAAA", resultOfTwo, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
				*newDetails(*newBenchmarkID(3, "api_call", nil), "f78gh1p", resultOfOne, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
			},
		}, wantCompared: ComparisonArray{
			Comparison{
				Reference:   *newDetails(*newBenchmarkID(1, "webhook", nil), "11bbAAA", resultOfOne, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
				Compare:     *newDetails(*newBenchmarkID(2, "webhook", nil), "11bbAAA", resultOfTwo, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
				Diff:        resultOfFifty,
				DiffMetrics: metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}},
			},
			Comparison{
				Reference:   *newDetails(*newBenchmarkID(4, "webhook", nil), "f78gh1p", resultOfTwo, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
				Compare:     *newDetails(*newBenchmarkID(3, "api_call", nil), "f78gh1p", resultOfOne, metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}}),
				Diff:        *newResult(*newQPS(50, 50, 50, 50), 50, -100, 50, 50, 50, 50),
				DiffMetrics: metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{}, ComponentsMemStatsAllocBytes: map[string]float64{}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			gotCompared := CompareDetailsArrays(tt.args.references, tt.args.compares)
			c.Assert(gotCompared, qt.DeepEquals, tt.wantCompared)
		})
	}
}
