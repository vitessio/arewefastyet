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

func TestMicroBenchmarkComparisonArray_Regression(t *testing.T) {
	tests := []struct {
		name         string
		microsMatrix ComparisonArray
		wantReason   string
	}{
		{name: "No regression", microsMatrix: ComparisonArray{
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench3", SubBenchmarkName: "bench3-pkg1"}, CurrLastDiff: 1},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1", SubBenchmarkName: "bench1-pkg1"}, CurrLastDiff: 1},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench2", SubBenchmarkName: "bench2-pkg1"}, CurrLastDiff: 1},
			{BenchmarkId: BenchmarkId{PkgName: "pkg2", Name: "bench2", SubBenchmarkName: "bench2-pkg2"}, CurrLastDiff: 1},
			{BenchmarkId: BenchmarkId{PkgName: "pkg2", Name: "bench1", SubBenchmarkName: "bench1-pkg2"}, CurrLastDiff: 1},
			{BenchmarkId: BenchmarkId{PkgName: "pkg3", Name: "bench1", SubBenchmarkName: "bench1-pkg3"}, CurrLastDiff: 1},
		}, wantReason: ""},

		{name: "Few regressions", microsMatrix: ComparisonArray{
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench3", SubBenchmarkName: "bench3-pkg1"}, CurrLastDiff: 0.5},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1", SubBenchmarkName: "bench1-pkg1"}, CurrLastDiff: 0.89},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench2", SubBenchmarkName: "bench2-pkg1"}, CurrLastDiff: 0.25},
		}, wantReason: "- pkg1/bench3-pkg1 decreased by 50.00%\n- pkg1/bench1-pkg1 decreased by 11.00%\n- pkg1/bench2-pkg1 decreased by 75.00%\n"},

		{name: "Close call regressions", microsMatrix: ComparisonArray{
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench3", SubBenchmarkName: "bench3-pkg1"}, CurrLastDiff: 0.90},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench1", SubBenchmarkName: "bench1-pkg1"}, CurrLastDiff: 0.91},
			{BenchmarkId: BenchmarkId{PkgName: "pkg1", Name: "bench2", SubBenchmarkName: "bench2-pkg1"}, CurrLastDiff: 0.899},
		}, wantReason: "- pkg1/bench2-pkg1 decreased by 10.10%\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			reason := tt.microsMatrix.Regression()
			c.Assert(reason, qt.Contains, tt.wantReason)
		})
	}
}
