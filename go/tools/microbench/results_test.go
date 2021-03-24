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

package microbench

import (
	qt "github.com/frankban/quicktest"
	"testing"
)

func TestMicroBenchmarkResults_ReduceSimpleMedian(t *testing.T) {
	tests := []struct {
		name string
		mrs  MicroBenchmarkDetailsArray
		want MicroBenchmarkDetailsArray
	}{
		// tc1
		{name: "Few simple values in same package but different benchmark names", mrs: MicroBenchmarkDetailsArray{
			// input bench 1
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 1.00 }},
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 1.00 }},
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 1.00 }},

			// input bench 2
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench2-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 2.00 }},
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench2-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 2.00 }},
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench2-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 2.00 }},
		}, want: MicroBenchmarkDetailsArray{
			// want bench 1
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 1.00 }},

			// want bench 2
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench2-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 2.00 }},
		}},

		// tc2
		{name: "Few values in different packages and different benchmark names", mrs: MicroBenchmarkDetailsArray{
			// input bench 1 in pkg 1
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 1.00 }},
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 5.00 }},
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 10.00 }},

			// input bench 1 in pkg 2
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg2", Name: "bench1-pkg2"} , Result: MicroBenchmarkResult{NSPerOp: 2.00 }},
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg2", Name: "bench1-pkg2"} , Result: MicroBenchmarkResult{NSPerOp: 2.50 }},
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg2", Name: "bench1-pkg2"} , Result: MicroBenchmarkResult{NSPerOp: 3.00 }},
		}, want: MicroBenchmarkDetailsArray{
			// want bench 1 from pkg1
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 5.00 }},

			// want bench 1 from pkg2
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg2", Name: "bench1-pkg2"} , Result: MicroBenchmarkResult{NSPerOp: 2.50 }},
		}},

		// tc3
		{name: "More unordered values with single package and benchmark name", mrs: MicroBenchmarkDetailsArray{
			// input bench 1
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 30.00 }},
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 5.00 }},
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 15.00 }},
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 10.00 }},
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 40.00 }},
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 25.00 }},
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 20.00 }},
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 0.00 }},
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 35.00 }},

		}, want: MicroBenchmarkDetailsArray{
			// want bench 1
			MicroBenchmarkDetails{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"} , Result: MicroBenchmarkResult{NSPerOp: 20.00 }},
		}},

		// tc4
		// {name: "Verify ordering by names", mrs: MicroBenchmarkDetailsArray{
		// 	// input bench 2
		// 	MicroBenchmarkDetails{PkgName: "pkg1", NSPerOp: 2.00, Name: "bench2-pkg1"},
		// 	MicroBenchmarkDetails{PkgName: "pkg1", NSPerOp: 2.00, Name: "bench2-pkg1"},
		// 	MicroBenchmarkDetails{PkgName: "pkg1", NSPerOp: 2.00, Name: "bench2-pkg1"},
		//
		// 	// input bench 1
		// 	MicroBenchmarkDetails{PkgName: "pkg1", NSPerOp: 1.00, Name: "bench1-pkg1"},
		// 	MicroBenchmarkDetails{PkgName: "pkg1", NSPerOp: 1.00, Name: "bench1-pkg1"},
		// 	MicroBenchmarkDetails{PkgName: "pkg1", NSPerOp: 1.00, Name: "bench1-pkg1"},
		//
		// }, want: MicroBenchmarkResults{
		// 	// want bench 1
		// 	MicroBenchmarkDetails{PkgName: "pkg1", NSPerOp: 1.00, Name: "bench1-pkg1"},
		//
		// 	// want bench 2
		// 	MicroBenchmarkDetails{PkgName: "pkg1", NSPerOp: 2.00, Name: "bench2-pkg1"},
		// }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			got := tt.mrs.ReduceSimpleMedian()
			c.Assert(got, qt.DeepEquals, tt.want)
		})
	}
}
