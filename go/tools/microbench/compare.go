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
	"github.com/vitessio/arewefastyet/go/mysql"
	"sort"
)

func CompareMicroBenchmarks(dbClient *mysql.Client, reference string, compare string) (MicroBenchmarkComparisonArray, error) {
	// compare micro benchmarks
	SHAs := []string{reference, compare}
	micros := map[string]MicroBenchmarkDetailsArray{}
	for _, sha := range SHAs {
		micro, err := GetResultsForGitRef(sha, dbClient)
		if err != nil {
			return nil, err
		}
		micros[sha] = micro.ReduceSimpleMedian()
	}
	microsMatrix := MergeMicroBenchmarkDetails(micros[reference], micros[compare])
	sort.SliceStable(microsMatrix, func(i, j int) bool {
		return !(microsMatrix[i].Current.NSPerOp < microsMatrix[j].Current.NSPerOp)
	})
	return microsMatrix, nil
}
