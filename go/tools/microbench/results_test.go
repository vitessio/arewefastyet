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
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestMicroBenchmarkResults_ReduceSimpleMedianByName(t *testing.T) {
	tests := []struct {
		name string
		mbd  DetailsArray
		want DetailsArray
	}{
		// tc1
		{name: "Few simple values in same package but different benchmark names", mbd: DetailsArray{
			// input bench 1
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(100, 1.00, 1.00, 1, 9)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(100, 1.00, 2.00, 50, 18)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(100, 1.00, 3.00, 100, 27)),

			// input bench 2
			*NewDetails(*NewBenchmarkId("pkg1", "bench2", "bench2-pkg1"), "", "", *NewResult(150, 2.00, 3.00, 55.00, 42)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench2", "bench2-pkg1"), "", "", *NewResult(300, 2.00, 4.00, 55.00, 84)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench2", "bench2-pkg1"), "", "", *NewResult(450, 2.00, 5.00, 55.00, 126)),
		}, want: DetailsArray{
			// want bench 1
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(100, 1.00, 2, 50.00, 18)),

			// want bench 2
			*NewDetails(*NewBenchmarkId("pkg1", "bench2", "bench2-pkg1"), "", "", *NewResult(300, 2.00, 4, 55.00, 84)),
		}},

		// tc2
		{name: "Few values in different packages and different benchmark names", mbd: DetailsArray{
			// input bench 1 in pkg 1
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "2021-05-10T13:20:11Z", *NewResult(0, 1.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "2021-05-10T13:20:11Z", *NewResult(0, 5.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "2021-05-10T13:20:11Z", *NewResult(0, 10.00, 0, 0, 0)),

			// input bench 1 in pkg 2
			*NewDetails(*NewBenchmarkId("pkg2", "bench1", "bench1-pkg2"), "", "2021-05-10T13:20:11Z", *NewResult(0, 2.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg2", "bench1", "bench1-pkg2"), "", "2021-05-10T13:20:11Z", *NewResult(0, 2.50, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg2", "bench1", "bench1-pkg2"), "", "2021-05-10T13:20:11Z", *NewResult(0, 3.00, 0, 0, 0)),
		}, want: DetailsArray{
			// want bench 1 from pkg1
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "2021-05-10T13:20:11Z", *NewResult(0, 5.00, 0, 0, 0)),

			// want bench 1 from pkg2
			*NewDetails(*NewBenchmarkId("pkg2", "bench1", "bench1-pkg2"), "", "2021-05-10T13:20:11Z", *NewResult(0, 2.50, 0, 0, 0)),
		}},

		// tc3
		{name: "More unordered values with single package and benchmark name", mbd: DetailsArray{
			// input bench 1
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(0, 30.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(0, 5.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(0, 15.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(0, 10.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(0, 40.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(0, 25.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(0, 20.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(0, 0.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(0, 35.00, 0, 0, 0)),
		}, want: DetailsArray{
			// want bench 1
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(0, 20.00, 0, 0, 0)),
		}},

		// tc4
		{name: "Few values in different packages and different benchmark names", mbd: DetailsArray{
			// input bench 1 in pkg 2
			*NewDetails(*NewBenchmarkId("pkg2", ", ", "bench1"), "", "", *NewResult(0, 2.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg2", ", ", "bench1"), "", "", *NewResult(0, 2.50, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg2", ", ", "bench1"), "", "", *NewResult(0, 3.00, 0, 0, 0)),

			// input bench 1 in pkg 1
			*NewDetails(*NewBenchmarkId("pkg1", ", ", "bench2"), "", "", *NewResult(0, 1.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", ", ", "bench2"), "", "", *NewResult(0, 5.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", ", ", "bench2"), "", "", *NewResult(0, 10.00, 0, 0, 0)),
		}, want: DetailsArray{
			// want bench 1 from pkg1
			*NewDetails(*NewBenchmarkId("pkg1", ", ", "bench2"), "", "", *NewResult(0, 5.00, 0, 0, 0)),

			// want bench 1 from pkg2
			*NewDetails(*NewBenchmarkId("pkg2", ", ", "bench1"), "", "", *NewResult(0, 2.50, 0, 0, 0)),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			got := tt.mbd.ReduceSimpleMedianByName()
			c.Assert(got, qt.DeepEquals, tt.want)
		})
	}
}

func TestMicroBenchmarkResults_ReduceSimpleMedianByGitRef(t *testing.T) {
	tests := []struct {
		name string
		mbd  DetailsArray
		want DetailsArray
	}{
		// tc1
		{name: "Few simple values in same package with different git refs", mbd: DetailsArray{
			// input for git ref `abcd`
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "abcd", "", *NewResult(100, 1.00, 1.00, 1, 9)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "abcd", "", *NewResult(100, 1.00, 2.00, 50, 18)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "abcd", "", *NewResult(100, 1.00, 3.00, 100, 27)),

			// input for git ref `efgh`
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "efgh", "", *NewResult(150, 2.00, 3.00, 55.00, 42)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "efgh", "", *NewResult(300, 2.00, 4.00, 55.00, 84)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "efgh", "", *NewResult(450, 2.00, 5.00, 55.00, 126)),
		}, want: DetailsArray{
			// want git ref `abcd`
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "abcd", "", *NewResult(100, 1.00, 2, 50.00, 18)),

			// want git ref `efgh`
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "efgh", "", *NewResult(300, 2.00, 4, 55.00, 84)),
		}},

		// tc2
		{name: "More unordered values with same git ref", mbd: DetailsArray{
			// input bench 1
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "abcd", "", *NewResult(0, 30.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "abcd", "", *NewResult(0, 5.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "abcd", "", *NewResult(0, 15.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "abcd", "", *NewResult(0, 10.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "abcd", "", *NewResult(0, 40.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "abcd", "", *NewResult(0, 25.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "abcd", "", *NewResult(0, 20.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "abcd", "", *NewResult(0, 0.00, 0, 0, 0)),
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "abcd", "", *NewResult(0, 35.00, 0, 0, 0)),
		}, want: DetailsArray{
			// want bench 1
			*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "abcd", "", *NewResult(0, 20.00, 0, 0, 0)),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			got := tt.mbd.ReduceSimpleMedianByGitRef()
			c.Assert(got, qt.DeepEquals, tt.want)
		})
	}
}

func BenchmarkReduceSimpleMedianByName(b *testing.B) {
	mbd := DetailsArray{
		*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(0, 1.00, 0, 0, 0)),
		*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(0, 1.00, 0, 0, 0)),
		*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(0, 1.00, 0, 0, 0)),
		*NewDetails(*NewBenchmarkId("pkg1", "bench2", "bench2-pkg1"), "", "", *NewResult(0, 2.00, 0, 0, 0)),
		*NewDetails(*NewBenchmarkId("pkg1", "bench2", "bench2-pkg1"), "", "", *NewResult(0, 2.00, 0, 0, 0)),
		*NewDetails(*NewBenchmarkId("pkg1", "bench2", "bench2-pkg1"), "", "", *NewResult(0, 2.00, 0, 0, 0)),
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		got := mbd.ReduceSimpleMedianByName()

		// Must be reduced to two indexes (bench1-pkg1 AND bench2-pkg1).
		if len(got) != 2 {
			b.Error("must be reduced to two elements")
		}
		if got[0].Result.NSPerOp != 1.00 || got[1].Result.NSPerOp != 2.00 {
			b.Error("wrong output from ReduceSimpleMedianByName")
		}
	}
}

func TestMergeMicroBenchmarkDetails(t *testing.T) {
	type args struct {
		currentMbd     DetailsArray
		lastReleaseMbd DetailsArray
	}
	tests := []struct {
		name string
		args args
		want ComparisonArray
	}{
		// tc1
		{name: "Simple compare with ordered array of one elements", args: args{
			currentMbd: DetailsArray{
				*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(0, 1.00, 0, 0, 0)),
			},
			lastReleaseMbd: DetailsArray{
				*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(0, 5.00, 0, 0, 0)),
			},
		}, want: ComparisonArray{
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1", SubBenchmarkName: "bench1-pkg1"}, Current: Result{NSPerOp: 1.00}, Last: Result{NSPerOp: 5.00}, Diff: Result{NSPerOp: -400}},
		}},

		// tc2
		{name: "Simple compare with ordered array of two elements", args: args{
			currentMbd: DetailsArray{
				*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(0, 1.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg1", "bench2", "bench2-pkg1"), "", "", *NewResult(0, 98.00, 0, 0, 0)),
			},
			lastReleaseMbd: DetailsArray{
				*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "", "", *NewResult(0, 5.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg1", "bench2", "bench2-pkg1"), "", "", *NewResult(0, 89.00, 0, 0, 0)),
			},
		}, want: ComparisonArray{
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1", SubBenchmarkName: "bench1-pkg1"}, Current: *NewResult(0, 1.00, 0, 0, 0), Last: *NewResult(0, 5.00, 0, 0, 0), Diff: Result{NSPerOp: -400}},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench2", SubBenchmarkName: "bench2-pkg1"}, Current: *NewResult(0, 98.00, 0, 0, 0), Last: *NewResult(0, 89.00, 0, 0, 0), Diff: Result{NSPerOp: 9.183673469387756}},
		}},

		// tc3
		{name: "Compare with unordered array", args: args{
			currentMbd: DetailsArray{
				*NewDetails(*NewBenchmarkId("pkg1", "bench3", "bench3-pkg1"), "aabb", "", *NewResult(0, 58.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "aabb", "", *NewResult(0, 1.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg1", "bench2", "bench2-pkg1"), "aabb", "", *NewResult(0, 98.00, 0, 0, 0)),
			},
			lastReleaseMbd: DetailsArray{
				*NewDetails(*NewBenchmarkId("pkg1", "bench2", "bench2-pkg1"), "ppbb", "", *NewResult(0, 89.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "ppbb", "", *NewResult(0, 5.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg1", "bench3", "bench3-pkg1"), "ppbb", "", *NewResult(0, 56.00, 0, 0, 0)),
			},
		}, want: ComparisonArray{
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench3", SubBenchmarkName: "bench3-pkg1"}, Current: *NewResult(0, 58.00, 0, 0, 0), Last: *NewResult(0, 56.00, 0, 0, 0), Diff: Result{NSPerOp: 3.4482758620689653}},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1", SubBenchmarkName: "bench1-pkg1"}, Current: *NewResult(0, 1.00, 0, 0, 0), Last: *NewResult(0, 5.00, 0, 0, 0), Diff: Result{NSPerOp:  -400}},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench2", SubBenchmarkName: "bench2-pkg1"}, Current: *NewResult(0, 98.00, 0, 0, 0), Last: *NewResult(0, 89.00, 0, 0, 0), Diff: Result{NSPerOp:  9.183673469387756}},
		}},

		// tc4
		{name: "Compare with unordered array from multiple package", args: args{
			currentMbd: DetailsArray{
				*NewDetails(*NewBenchmarkId("pkg1", "bench3", "bench3-pkg1"), "aabb", "", *NewResult(0, 58.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "aabb", "", *NewResult(0, 1.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg1", "bench2", "bench2-pkg1"), "aabb", "", *NewResult(0, 98.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg2", "bench2", "bench2-pkg2"), "ppbb", "", *NewResult(0, 3.50, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg2", "bench1", "bench1-pkg2"), "ppbb", "", *NewResult(0, 5.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg3", "bench1", "bench1-pkg3"), "ppbb", "", *NewResult(0, 2385.00, 0, 0, 0)),
			},
			lastReleaseMbd: DetailsArray{
				*NewDetails(*NewBenchmarkId("pkg1", "bench2", "bench2-pkg1"), "ppbb", "", *NewResult(0, 89.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg3", "bench1", "bench1-pkg3"), "ppbb", "", *NewResult(0, 2560.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg1", "bench3", "bench3-pkg1"), "ppbb", "", *NewResult(0, 56.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg2", "bench2", "bench2-pkg2"), "ppbb", "", *NewResult(0, 6.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "ppbb", "", *NewResult(0, 5.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg2", "bench1", "bench1-pkg2"), "ppbb", "", *NewResult(0, 4.20, 0, 0, 0)),
			},
		}, want: ComparisonArray{
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench3", SubBenchmarkName: "bench3-pkg1"}, Current: *NewResult(0, 58.00, 0, 0, 0), Last: *NewResult(0, 56.00, 0, 0, 0), Diff: Result{NSPerOp:  3.4482758620689653}},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1", SubBenchmarkName: "bench1-pkg1"}, Current: *NewResult(0, 1.00, 0, 0, 0), Last: *NewResult(0, 5.00, 0, 0, 0), Diff: Result{NSPerOp:  -400}},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench2", SubBenchmarkName: "bench2-pkg1"}, Current: *NewResult(0, 98.00, 0, 0, 0), Last: *NewResult(0, 89.00, 0, 0, 0), Diff: Result{NSPerOp: 9.183673469387756}},
			{BenchmarkId: BenchmarkId{PkgName: "pkg2", Name: "bench2", SubBenchmarkName: "bench2-pkg2"}, Current: *NewResult(0, 3.50, 0, 0, 0), Last: *NewResult(0, 6.00, 0, 0, 0), Diff: Result{NSPerOp:  -71.42857142857143}},
			{BenchmarkId: BenchmarkId{PkgName: "pkg2", Name: "bench1", SubBenchmarkName: "bench1-pkg2"}, Current: *NewResult(0, 5.00, 0, 0, 0), Last: *NewResult(0, 4.20, 0, 0, 0), Diff: Result{NSPerOp:  15.999999999999998}},
			{BenchmarkId: BenchmarkId{PkgName: "pkg3", Name: "bench1", SubBenchmarkName: "bench1-pkg3"}, Current: *NewResult(0, 2385.00, 0, 0, 0), Last: *NewResult(0, 2560.00, 0, 0, 0), Diff: Result{NSPerOp:  -7.337526205450734}},
		}},

		// tc5
		{name: "Compare with unordered and different size array from multiple package", args: args{
			currentMbd: DetailsArray{
				*NewDetails(*NewBenchmarkId("pkg1", "bench3", "bench3-pkg1"), "aabb", "", *NewResult(0, 58.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "aabb", "", *NewResult(0, 1.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg1", "bench2", "bench2-pkg1"), "aabb", "", *NewResult(0, 98.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg2", "bench2", "bench2-pkg2"), "ppbb", "", *NewResult(0, 3.50, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg2", "bench1", "bench1-pkg2"), "ppbb", "", *NewResult(0, 5.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg3", "bench1", "bench1-pkg3"), "ppbb", "", *NewResult(0, 2385.00, 0, 0, 0)),
			},
			lastReleaseMbd: DetailsArray{
				*NewDetails(*NewBenchmarkId("pkg1", "bench2", "bench2-pkg1"), "ppbb", "", *NewResult(0, 89.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg1", "bench3", "bench3-pkg1"), "ppbb", "", *NewResult(0, 56.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg1", "bench1", "bench1-pkg1"), "ppbb", "", *NewResult(0, 5.00, 0, 0, 0)),
				*NewDetails(*NewBenchmarkId("pkg2", "bench1", "bench1-pkg2"), "ppbb", "", *NewResult(0, 4.20, 0, 0, 0)),
			},
		}, want: ComparisonArray{
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench3", SubBenchmarkName: "bench3-pkg1"}, Current: *NewResult(0, 58.00, 0, 0, 0), Last: *NewResult(0, 56.00, 0, 0, 0), Diff: Result{NSPerOp:  3.4482758620689653}},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1", SubBenchmarkName: "bench1-pkg1"}, Current: *NewResult(0, 1.00, 0, 0, 0), Last: *NewResult(0, 5.00, 0, 0, 0), Diff: Result{NSPerOp:  -400}},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench2", SubBenchmarkName: "bench2-pkg1"}, Current: *NewResult(0, 98.00, 0, 0, 0), Last: *NewResult(0, 89.00, 0, 0, 0), Diff: Result{NSPerOp:  9.183673469387756}},
			{BenchmarkId: BenchmarkId{PkgName: "pkg2", Name: "bench2", SubBenchmarkName: "bench2-pkg2"}, Current: *NewResult(0, 3.50, 0, 0, 0), Last: *NewResult(0, 0, 0, 0, 0), Diff: Result{NSPerOp:  0}},
			{BenchmarkId: BenchmarkId{PkgName: "pkg2", Name: "bench1", SubBenchmarkName: "bench1-pkg2"}, Current: *NewResult(0, 5.00, 0, 0, 0), Last: *NewResult(0, 4.20, 0, 0, 0), Diff: Result{NSPerOp:  15.999999999999998}},
			{BenchmarkId: BenchmarkId{PkgName: "pkg3", Name: "bench1", SubBenchmarkName: "bench1-pkg3"}, Current: *NewResult(0, 2385.00, 0, 0, 0), Last: *NewResult(0, 0.00, 0, 0, 0), Diff: Result{NSPerOp:  0}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			got := MergeDetails(tt.args.currentMbd, tt.args.lastReleaseMbd)
			c.Assert(got, qt.HasLen, len(tt.want))
			c.Assert(got, qt.DeepEquals, tt.want)
		})
	}
}

func TestHumanReadableStrings(t *testing.T) {
	c := qt.New(t)
	r := Result{
		Ops:         876543,
		NSPerOp:     141650883.50,
		MBPerSec:    45030859.00,
		BytesPerOp:  4528.14,
		AllocsPerOp: 1106456.00,
	}
	c.Assert(r.AllocsPerOpStr(), qt.Equals, "1.106456 M")
	c.Assert(r.MBPerSecStr(), qt.Equals, "45 MB/s")
	c.Assert(r.NSPerOpStr(), qt.Equals, "141,650,883.5")
	c.Assert(r.NSPerOpToDurationStr(), qt.Equals, "141.65 ms")
	c.Assert(r.OpsStr(), qt.Equals, "876,543")
	c.Assert(r.BytesPerOpStr(), qt.Equals, "4.5 kB/op")

	r = Result{
		NSPerOp: 2.5149999999999997,
	}
	c.Assert(r.NSPerOpStr(), qt.Equals, "2.5")
	c.Assert(r.NSPerOpToDurationStr(), qt.Equals, "2.00 ns")
}
