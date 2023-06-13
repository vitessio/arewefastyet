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
	"fmt"

	"github.com/vitessio/arewefastyet/go/storage"
)

// Compare takes in 3 arguments, the database, and 2 SHAs. It reads from the database, the microbenchmark
// results for the 2 SHAs and compares them. The result is a comparison array.
func Compare(client storage.SQLClient, right string, left string) (ComparisonArray, error) {
	// compare micro benchmarks
	SHAs := []string{right, left}
	micros := map[string]DetailsArray{}
	for _, sha := range SHAs {
		micro, err := GetResultsForGitRef(sha, client)
		if err != nil {
			return nil, err
		}
		micros[sha] = micro.ReduceSimpleMedianByName()
	}
	microsMatrix := MergeDetails(micros[right], micros[left])
	// The result of the merge will be sorted by the package name and then the benchmark name
	return microsMatrix, nil
}

// Regression returns a string containing the reason of the regression of the given ComparisonArray,
// if no regression was evaluated, the reason will be an empty string.
// The format of a single benchmark regression's reason is like this:
//
// "- {pkg name}/{benchmark name} decreased by {decrease percentage}%\n"
//
func (microsMatrix ComparisonArray) Regression() (reason string) {
	for _, micro := range microsMatrix {
		m := []struct{
			value float64
			name string
		}{
			{name: "total operation", value: micro.Diff.Ops},
			{name: "nanosecond per operation", value: micro.Diff.NSPerOp},
			{name: "bytes per operation", value: micro.Diff.BytesPerOp},
			{name: "MB per second", value: micro.Diff.MBPerSec},
			{name: "allocations per operation", value: micro.Diff.AllocsPerOp},
		}

		for _, s := range m {
			if s.value < -10 {
				reason += fmt.Sprintf("- %s/%s: metric: %s, decreased by %.2f%%\n", micro.PkgName, micro.SubBenchmarkName, s.name, -1*s.value)
			}
		}
	}
	return
}
