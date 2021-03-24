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

type (
	MicroBenchmarkResult struct {
		Ops     int
		NSPerOp float64
	}
	BenchmarkId struct {
		PkgName string
		Name    string
	}
	MicroBenchmarkDetails struct {
		BenchmarkId
		GitRef string
		Result MicroBenchmarkResult
	}

	MicroBenchmarkComparison struct {
		BenchmarkId
		Current, Last MicroBenchmarkResult
	}

	MicroBenchmarkDetailsArray    []MicroBenchmarkDetails
	MicroBenchmarkComparisonArray []MicroBenchmarkComparison
)

func MergeMicroBenchmarkDetails(currDetails, lastDetails MicroBenchmarkDetailsArray) MicroBenchmarkComparisonArray {
	var comparisons MicroBenchmarkComparisonArray

	for _, details := range currDetails {
		var comparison MicroBenchmarkComparison
		comparison.BenchmarkId = details.BenchmarkId
		comparison.Current = details.Result
		for j := 0; j < len(lastDetails); j++ {
			if lastDetails[j].BenchmarkId == details.BenchmarkId {
				comparison.Last = lastDetails[j].Result
				break
			}
		}
		comparisons = append(comparisons, comparison)
	}
	return comparisons
}

func (mrs MicroBenchmarkDetailsArray) ReduceSimpleMedian() MicroBenchmarkDetailsArray {
	var details MicroBenchmarkDetailsArray

	sort.Slice(mrs, func(i, j int) bool {
		return mrs[i].PkgName < mrs[j].PkgName && mrs[i].Name < mrs[j].Name
	})
	for i := 0; i < len(mrs); {
		var j int
		var interOps []int
		var interNSPerOp []float64
		for j = i; j < len(mrs) && mrs[i].Name == mrs[j].Name; j++ {
			interOps = append(interOps, mrs[j].Result.Ops)
			interNSPerOp = append(interNSPerOp, mrs[j].Result.NSPerOp)
		}

		sort.Ints(interOps)
		sort.Float64s(interNSPerOp)
		interOpsResult := medianInt(interOps)
		interNSPerOpResult := medianFloat(interNSPerOp)

		details = append(details, MicroBenchmarkDetails{
			BenchmarkId: BenchmarkId{
				PkgName: mrs[i].PkgName,
				Name:    mrs[i].Name,
			},
			GitRef: mrs[i].GitRef,
			Result: MicroBenchmarkResult{
				Ops:     interOpsResult,
				NSPerOp: interNSPerOpResult,
			},
		})
		i = j
	}
	return details
}

func medianInt(values []int) int {
	middle := len(values) / 2
	if len(values)%2 == 1 {
		return values[middle]
	}
	return (values[middle-1] + values[middle]) / 2
}

func medianFloat(values []float64) float64 {
	middle := len(values) / 2
	if len(values)%2 == 1 {
		return values[middle]
	}
	return (values[middle-1] + values[middle]) / 2
}

func GetResultsForGitRef(ref string, client *mysql.Client) (mrs MicroBenchmarkDetailsArray, err error) {
	rows, err := client.Select("select m.pkg_name, m.name, md.n, md.ns_per_op FROM "+
		"microbenchmark m, microbenchmark_details md where m.git_ref = ? AND "+
		"md.microbenchmark_no = m.microbenchmark_no order by m.microbenchmark_no desc", ref)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var res MicroBenchmarkDetails
		res.GitRef = ref
		err = rows.Scan(&res.PkgName, &res.Name, &res.Result.Ops, &res.Result.NSPerOp)
		if err != nil {
			return nil, err
		}
		mrs = append(mrs, res)
	}
	return mrs, nil
}
