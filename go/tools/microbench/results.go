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
		Ops         int
		NSPerOp     float64
		MBPerSec    float64
		BytesPerOp  float64
		AllocsPerOp float64
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

func NewMicroBenchmarkDetails(benchmarkId BenchmarkId, gitRef string, result MicroBenchmarkResult) *MicroBenchmarkDetails {
	return &MicroBenchmarkDetails{
		BenchmarkId: benchmarkId,
		GitRef:      gitRef,
		Result:      result,
	}
}

func NewBenchmarkId(pkgName string, name string) *BenchmarkId {
	return &BenchmarkId{
		PkgName: pkgName,
		Name:    name,
	}
}

func NewMicroBenchmarkResult(ops int, NSPerOp, MBPerSec, BytesPerOp, AllocsPerOp float64) *MicroBenchmarkResult {
	return &MicroBenchmarkResult{
		Ops:         ops,
		NSPerOp:     NSPerOp,
		MBPerSec:    MBPerSec,
		BytesPerOp:  BytesPerOp,
		AllocsPerOp: AllocsPerOp,
	}
}

func MergeMicroBenchmarkDetails(currentMbd, lastReleaseMbd MicroBenchmarkDetailsArray) (compareMbs MicroBenchmarkComparisonArray) {
	for _, details := range currentMbd {
		var compareMb MicroBenchmarkComparison
		compareMb.BenchmarkId = details.BenchmarkId
		compareMb.Current = details.Result
		for j := 0; j < len(lastReleaseMbd); j++ {
			if lastReleaseMbd[j].BenchmarkId == details.BenchmarkId {
				compareMb.Last = lastReleaseMbd[j].Result
				break
			}
		}
		compareMbs = append(compareMbs, compareMb)
	}
	return compareMbs
}

func (mbd MicroBenchmarkDetailsArray) ReduceSimpleMedian() (reduceMbd MicroBenchmarkDetailsArray) {
	sort.SliceStable(mbd, func(i, j int) bool {
		return mbd[i].Name < mbd[j].Name
	})
	sort.SliceStable(mbd, func(i, j int) bool {
		return mbd[i].PkgName < mbd[j].PkgName
	})
	for i := 0; i < len(mbd); {
		var j int
		var interOps []int
		var interNSPerOp []float64
		var interMBPerSec []float64
		var interBytesPerOp []float64
		var interAllocsPerOp []float64
		for j = i; j < len(mbd) && mbd[i].Name == mbd[j].Name; j++ {
			interOps = append(interOps, mbd[j].Result.Ops)
			interNSPerOp = append(interNSPerOp, mbd[j].Result.NSPerOp)
			interMBPerSec = append(interMBPerSec, mbd[j].Result.MBPerSec)
			interBytesPerOp = append(interBytesPerOp, mbd[j].Result.BytesPerOp)
			interAllocsPerOp = append(interAllocsPerOp, mbd[j].Result.AllocsPerOp)
		}

		sort.Ints(interOps)
		sort.Float64s(interNSPerOp)
		sort.Float64s(interMBPerSec)
		sort.Float64s(interBytesPerOp)
		sort.Float64s(interAllocsPerOp)
		interOpsResult := medianInt(interOps)
		interNSPerOpResult := medianFloat(interNSPerOp)
		interMBPerSecResult := medianFloat(interMBPerSec)
		interBytesPerOpResult := medianFloat(interBytesPerOp)
		interAllocsPerOpResult := medianFloat(interAllocsPerOp)
		reduceMbd = append(reduceMbd, *NewMicroBenchmarkDetails(
			*NewBenchmarkId(mbd[i].PkgName, mbd[i].Name),
			mbd[i].GitRef,
			*NewMicroBenchmarkResult(interOpsResult, interNSPerOpResult, interMBPerSecResult, interBytesPerOpResult, interAllocsPerOpResult),
		))
		i = j
	}
	return reduceMbd
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
	rows, err := client.Select("select m.pkg_name, m.name, md.n, md.ns_per_op, md.bytes_per_op,"+
		" md.allocs_per_op, md.mb_per_sec FROM microbenchmark m, microbenchmark_details md where m.git_ref = ? AND "+
		"md.microbenchmark_no = m.microbenchmark_no order by m.microbenchmark_no desc", ref)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var res MicroBenchmarkDetails
		res.GitRef = ref
		err = rows.Scan(&res.PkgName, &res.Name, &res.Result.Ops, &res.Result.NSPerOp, &res.Result.BytesPerOp,
			&res.Result.AllocsPerOp, &res.Result.MBPerSec)
		if err != nil {
			return nil, err
		}
		mrs = append(mrs, res)
	}
	return mrs, nil
}
