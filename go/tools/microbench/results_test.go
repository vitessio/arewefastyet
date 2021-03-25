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
		mbd  MicroBenchmarkDetailsArray
		want MicroBenchmarkDetailsArray
	}{
		// tc1
		{name: "Few simple values in same package but different benchmark names", mbd: MicroBenchmarkDetailsArray{
			// input bench 1
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 1.00, 0, 0, 0)),
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 1.00, 0, 0, 0)),
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 1.00, 0, 0, 0)),

			// input bench 2
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench2-pkg1"), "", *NewMicroBenchmarkResult(0, 2.00, 0, 0, 0)),
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench2-pkg1"), "", *NewMicroBenchmarkResult(0, 2.00, 0, 0, 0)),
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench2-pkg1"), "", *NewMicroBenchmarkResult(0, 2.00, 0, 0, 0)),
		}, want: MicroBenchmarkDetailsArray{
			// want bench 1
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 1.00, 0, 0, 0)),

			// want bench 2
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench2-pkg1"), "", *NewMicroBenchmarkResult(0, 2.00, 0, 0, 0)),
		}},

		// tc2
		{name: "Few values in different packages and different benchmark names", mbd: MicroBenchmarkDetailsArray{
			// input bench 1 in pkg 1
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 1.00, 0, 0, 0)),
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 5.00, 0, 0, 0)),
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 10.00, 0, 0, 0)),

			// input bench 1 in pkg 2
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg2", "bench1-pkg2"), "", *NewMicroBenchmarkResult(0, 2.00, 0, 0, 0)),
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg2", "bench1-pkg2"), "", *NewMicroBenchmarkResult(0, 2.50, 0, 0, 0)),
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg2", "bench1-pkg2"), "", *NewMicroBenchmarkResult(0, 3.00, 0, 0, 0)),
		}, want: MicroBenchmarkDetailsArray{
			// want bench 1 from pkg1
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 5.00, 0, 0, 0)),

			// want bench 1 from pkg2
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg2", "bench1-pkg2"), "", *NewMicroBenchmarkResult(0, 2.50, 0, 0, 0)),
		}},

		// tc3
		{name: "More unordered values with single package and benchmark name", mbd: MicroBenchmarkDetailsArray{
			// input bench 1
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 30.00, 0, 0, 0)),
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 5.00, 0, 0, 0)),
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 15.00, 0, 0, 0)),
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 10.00, 0, 0, 0)),
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 40.00, 0, 0, 0)),
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 25.00, 0, 0, 0)),
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 20.00, 0, 0, 0)),
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 0.00, 0, 0, 0)),
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 35.00, 0, 0, 0)),
		}, want: MicroBenchmarkDetailsArray{
			// want bench 1
			*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 20.00, 0, 0, 0)),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			got := tt.mbd.ReduceSimpleMedian()
			c.Assert(got, qt.DeepEquals, tt.want)
		})
	}
}

func BenchmarkReduceSimpleMedian(b *testing.B) {
	mbd := MicroBenchmarkDetailsArray{
		*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 1.00, 0, 0, 0)),
		*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 1.00, 0, 0, 0)),
		*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 1.00, 0, 0, 0)),
		*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench2-pkg1"), "", *NewMicroBenchmarkResult(0, 2.00, 0, 0, 0)),
		*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench2-pkg1"), "", *NewMicroBenchmarkResult(0, 2.00, 0, 0, 0)),
		*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench2-pkg1"), "", *NewMicroBenchmarkResult(0, 2.00, 0, 0, 0)),
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		got := mbd.ReduceSimpleMedian()

		// Must be reduced to two indexes (bench1-pkg1 AND bench2-pkg1).
		if len(got) != 2 {
			b.Error("must be reduced to two elements")
		}
		if got[0].Result.NSPerOp != 1.00 || got[1].Result.NSPerOp != 2.00 {
			b.Error("wrong output from ReduceSimpleMedian")
		}
	}
}

func TestMergeMicroBenchmarkDetails(t *testing.T) {
	type args struct {
		currentMbd     MicroBenchmarkDetailsArray
		lastReleaseMbd MicroBenchmarkDetailsArray
	}
	tests := []struct {
		name string
		args args
		want MicroBenchmarkComparisonArray
	}{
		// tc1
		{name: "Simple compare with ordered array of one elements", args: args{
			currentMbd: MicroBenchmarkDetailsArray{
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 1.00, 0, 0, 0)),
			},
			lastReleaseMbd: MicroBenchmarkDetailsArray{
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 5.00, 0, 0, 0)),
			},
		}, want: MicroBenchmarkComparisonArray{
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"}, Current: MicroBenchmarkResult{NSPerOp: 1.00}, Last: MicroBenchmarkResult{NSPerOp: 5.00}},
		}},

		// tc2
		{name: "Simple compare with ordered array of two elements", args: args{
			currentMbd: MicroBenchmarkDetailsArray{
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 1.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench2-pkg1"), "", *NewMicroBenchmarkResult(0, 98.00, 0, 0, 0)),
			},
			lastReleaseMbd: MicroBenchmarkDetailsArray{
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "", *NewMicroBenchmarkResult(0, 5.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench2-pkg1"), "", *NewMicroBenchmarkResult(0, 89.00, 0, 0, 0)),
			},
		}, want: MicroBenchmarkComparisonArray{
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"}, Current: *NewMicroBenchmarkResult(0, 1.00, 0, 0, 0), Last: *NewMicroBenchmarkResult(0, 5.00, 0, 0, 0)},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench2-pkg1"}, Current: *NewMicroBenchmarkResult(0, 98.00, 0, 0, 0), Last: *NewMicroBenchmarkResult(0, 89.00, 0, 0, 0)},
		}},

		// tc3
		{name: "Compare with unordered array", args: args{
			currentMbd: MicroBenchmarkDetailsArray{
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench3-pkg1"), "aabb", *NewMicroBenchmarkResult(0, 58.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "aabb", *NewMicroBenchmarkResult(0, 1.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench2-pkg1"), "aabb", *NewMicroBenchmarkResult(0, 98.00, 0, 0, 0)),
			},
			lastReleaseMbd: MicroBenchmarkDetailsArray{
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench2-pkg1"), "ppbb", *NewMicroBenchmarkResult(0, 89.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "ppbb", *NewMicroBenchmarkResult(0, 5.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench3-pkg1"), "ppbb", *NewMicroBenchmarkResult(0, 56.00, 0, 0, 0)),
			},
		}, want: MicroBenchmarkComparisonArray{
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench3-pkg1"}, Current: *NewMicroBenchmarkResult(0, 58.00, 0, 0, 0), Last: *NewMicroBenchmarkResult(0, 56.00, 0, 0, 0)},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"}, Current: *NewMicroBenchmarkResult(0, 1.00, 0, 0, 0), Last: *NewMicroBenchmarkResult(0, 5.00, 0, 0, 0)},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench2-pkg1"}, Current: *NewMicroBenchmarkResult(0, 98.00, 0, 0, 0), Last: *NewMicroBenchmarkResult(0, 89.00, 0, 0, 0)},
		}},

		// tc4
		{name: "Compare with unordered array from multiple package", args: args{
			currentMbd: MicroBenchmarkDetailsArray{
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench3-pkg1"), "aabb", *NewMicroBenchmarkResult(0, 58.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "aabb", *NewMicroBenchmarkResult(0, 1.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench2-pkg1"), "aabb", *NewMicroBenchmarkResult(0, 98.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg2", "bench2-pkg2"), "ppbb", *NewMicroBenchmarkResult(0, 3.50, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg2", "bench1-pkg2"), "ppbb", *NewMicroBenchmarkResult(0, 5.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg3", "bench1-pkg3"), "ppbb", *NewMicroBenchmarkResult(0, 2385.00, 0, 0, 0)),
			},
			lastReleaseMbd: MicroBenchmarkDetailsArray{
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench2-pkg1"), "ppbb", *NewMicroBenchmarkResult(0, 89.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg3", "bench1-pkg3"), "ppbb", *NewMicroBenchmarkResult(0, 2560.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench3-pkg1"), "ppbb", *NewMicroBenchmarkResult(0, 56.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg2", "bench2-pkg2"), "ppbb", *NewMicroBenchmarkResult(0, 6.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "ppbb", *NewMicroBenchmarkResult(0, 5.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg2", "bench1-pkg2"), "ppbb", *NewMicroBenchmarkResult(0, 4.20, 0, 0, 0)),
			},
		}, want: MicroBenchmarkComparisonArray{
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench3-pkg1"}, Current: *NewMicroBenchmarkResult(0, 58.00, 0, 0, 0), Last: *NewMicroBenchmarkResult(0, 56.00, 0, 0, 0)},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"}, Current: *NewMicroBenchmarkResult(0, 1.00, 0, 0, 0), Last: *NewMicroBenchmarkResult(0, 5.00, 0, 0, 0)},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench2-pkg1"}, Current: *NewMicroBenchmarkResult(0, 98.00, 0, 0, 0), Last: *NewMicroBenchmarkResult(0, 89.00, 0, 0, 0)},
			{BenchmarkId: BenchmarkId{PkgName: "pkg2", Name: "bench2-pkg2"}, Current: *NewMicroBenchmarkResult(0, 3.50, 0, 0, 0), Last: *NewMicroBenchmarkResult(0, 6.00, 0, 0, 0)},
			{BenchmarkId: BenchmarkId{PkgName: "pkg2", Name: "bench1-pkg2"}, Current: *NewMicroBenchmarkResult(0, 5.00, 0, 0, 0), Last: *NewMicroBenchmarkResult(0, 4.20, 0, 0, 0)},
			{BenchmarkId: BenchmarkId{PkgName: "pkg3", Name: "bench1-pkg3"}, Current: *NewMicroBenchmarkResult(0, 2385.00, 0, 0, 0), Last: *NewMicroBenchmarkResult(0, 2560.00, 0, 0, 0)},
		}},

		// tc5
		{name: "Compare with unordered and different size array from multiple package", args: args{
			currentMbd: MicroBenchmarkDetailsArray{
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench3-pkg1"), "aabb", *NewMicroBenchmarkResult(0, 58.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "aabb", *NewMicroBenchmarkResult(0, 1.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench2-pkg1"), "aabb", *NewMicroBenchmarkResult(0, 98.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg2", "bench2-pkg2"), "ppbb", *NewMicroBenchmarkResult(0, 3.50, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg2", "bench1-pkg2"), "ppbb", *NewMicroBenchmarkResult(0, 5.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg3", "bench1-pkg3"), "ppbb", *NewMicroBenchmarkResult(0, 2385.00, 0, 0, 0)),
			},
			lastReleaseMbd: MicroBenchmarkDetailsArray{
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench2-pkg1"), "ppbb", *NewMicroBenchmarkResult(0, 89.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench3-pkg1"), "ppbb", *NewMicroBenchmarkResult(0, 56.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg1", "bench1-pkg1"), "ppbb", *NewMicroBenchmarkResult(0, 5.00, 0, 0, 0)),
				*NewMicroBenchmarkDetails(*NewBenchmarkId("pkg2", "bench1-pkg2"), "ppbb", *NewMicroBenchmarkResult(0, 4.20, 0, 0, 0)),
			},
		}, want: MicroBenchmarkComparisonArray{
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench3-pkg1"}, Current: *NewMicroBenchmarkResult(0, 58.00, 0, 0, 0), Last: *NewMicroBenchmarkResult(0, 56.00, 0, 0, 0)},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1-pkg1"}, Current: *NewMicroBenchmarkResult(0, 1.00, 0, 0, 0), Last: *NewMicroBenchmarkResult(0, 5.00, 0, 0, 0)},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench2-pkg1"}, Current: *NewMicroBenchmarkResult(0, 98.00, 0, 0, 0), Last: *NewMicroBenchmarkResult(0, 89.00, 0, 0, 0)},
			{BenchmarkId: BenchmarkId{PkgName: "pkg2", Name: "bench2-pkg2"}, Current: *NewMicroBenchmarkResult(0, 3.50, 0, 0, 0), Last: *NewMicroBenchmarkResult(0, 0, 0, 0, 0)},
			{BenchmarkId: BenchmarkId{PkgName: "pkg2", Name: "bench1-pkg2"}, Current: *NewMicroBenchmarkResult(0, 5.00, 0, 0, 0), Last: *NewMicroBenchmarkResult(0, 4.20, 0, 0, 0)},
			{BenchmarkId: BenchmarkId{PkgName: "pkg3", Name: "bench1-pkg3"}, Current: *NewMicroBenchmarkResult(0, 2385.00, 0, 0, 0), Last: *NewMicroBenchmarkResult(0, 0.00, 0, 0, 0)},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			got := MergeMicroBenchmarkDetails(tt.args.currentMbd, tt.args.lastReleaseMbd)
			c.Assert(got, qt.HasLen, len(tt.want))
			c.Assert(got, qt.DeepEquals, tt.want)
		})
	}
}