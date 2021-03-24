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

type MicroBenchmarkResult struct {
	PkgName string
	Name string
	NSPerOp float64
}

type MicroBenchmarkResults []MicroBenchmarkResult

func (mrs MicroBenchmarkResults) ReduceSimpleMedian() MicroBenchmarkResults {
	var results MicroBenchmarkResults

	sort.Slice(mrs, func(i, j int) bool {
		return mrs[i].PkgName < mrs[i].PkgName && mrs[i].Name < mrs[j].Name
	})
	for i := 0; i < len(mrs); {
		var j int
		var interNSPerOp []float64
		for j = i; j < len(mrs) && mrs[i].Name == mrs[j].Name; j++ {
			interNSPerOp = append(interNSPerOp, mrs[j].NSPerOp)
		}

		sort.Float64s(interNSPerOp)
		var interNSPerOpResult float64
		middle := len(interNSPerOp) / 2
		if len(interNSPerOp) % 2 == 1 {
			interNSPerOpResult = interNSPerOp[middle]
		} else {
			interNSPerOpResult = (interNSPerOp[middle - 1] + interNSPerOp[middle]) / 2
		}

		results = append(results, MicroBenchmarkResult{
			PkgName: mrs[i].PkgName,
			Name:    mrs[i].Name,
			NSPerOp: interNSPerOpResult,
		})
		i = j
	}
	return results
}

func GetResultsForGitRef(ref string, client *mysql.Client) (mrs MicroBenchmarkResults, err error) {
	rows, err := client.Select("select m.pkg_name, m.name, md.ns_per_op FROM " +
		"microbenchmark m, microbenchmark_details  md where m.git_ref = ? AND " +
		"md.microbenchmark_no = m.microbenchmark_no order by m.microbenchmark_no desc;", ref)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var res MicroBenchmarkResult
		err = rows.Scan(&res.PkgName, &res.Name, &res.NSPerOp)
		if err != nil {
			return nil, err
		}
		mrs = append(mrs, res)
	}
	return mrs, nil
}