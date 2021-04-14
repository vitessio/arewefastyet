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
	"testing"
)

func TestMacroBenchmarkResultsArray_mergeMedian(t *testing.T) {
	tests := []struct {
		name             string
		mrs              MacroBenchmarkResultsArray
		wantMergedResult MacroBenchmarkResult
	}{
		{name:"Single result in array", mrs: MacroBenchmarkResultsArray{
			*NewMacroBenchmarkResult(*NewQPS(1.0, 1.0, 1.0, 1.0), 1.0, 1.0, 1.0, 1.0, 1, 1.0),
		}, wantMergedResult: *NewMacroBenchmarkResult(*NewQPS(1.0, 1.0, 1.0, 1.0), 1.0, 1.0, 1.0, 1.0, 1, 1.0)},

		{name:"Multiple results in array", mrs: MacroBenchmarkResultsArray{
			*NewMacroBenchmarkResult(*NewQPS(1.0, 1.0, 1.0, 1.0), 1.0, 1.0, 1.0, 1.0, 1, 1.0),
			*NewMacroBenchmarkResult(*NewQPS(1.0, 1.0, 1.0, 1.0), 1.0, 1.0, 1.0, 1.0, 1, 1.0),
			*NewMacroBenchmarkResult(*NewQPS(1.0, 1.0, 1.0, 1.0), 1.0, 1.0, 1.0, 1.0, 1, 1.0),
		}, wantMergedResult: *NewMacroBenchmarkResult(*NewQPS(1.0, 1.0, 1.0, 1.0), 1.0, 1.0, 1.0, 1.0, 1, 1.0)},

		{name:"Multiple and different results in array", mrs: MacroBenchmarkResultsArray{
			*NewMacroBenchmarkResult(*NewQPS(1.0, 1.0, 1.0, 3), 1.0, 1.0, 1.0, 1.5, 1, 10.0),
			*NewMacroBenchmarkResult(*NewQPS(2.0, 5.0, 1.5, 6), 2.0, 5.0, 3.0, 2.5, 1000, 20.0),
			*NewMacroBenchmarkResult(*NewQPS(3.0, 10.0, 2.0, 9), 3.0, 10.0, 2.0, 3.5, 500, 30.0),
		}, wantMergedResult: *NewMacroBenchmarkResult(*NewQPS(2.0, 5.0, 1.5, 6), 2.0, 5.0, 2.0, 2.5, 500, 20.0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			gotMergedResult := tt.mrs.mergeMedian()

			c.Assert(gotMergedResult, qt.DeepEquals, tt.wantMergedResult)
		})
	}
}
