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
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/dustin/go-humanize"
	"github.com/vitessio/arewefastyet/go/storage/mysql"
	"github.com/vitessio/arewefastyet/go/tools/math"
)

type (
	// Result contains all the metrics measured by a microbenchmark.
	Result struct {
		Ops         float64
		NSPerOp     float64
		MBPerSec    float64
		BytesPerOp  float64
		AllocsPerOp float64
	}

	// BenchmarkId represents the identification of a microbenchmark.
	BenchmarkId struct {
		PkgName          string
		Name             string
		SubBenchmarkName string
	}

	// Details refers to a single microbenchmark.
	Details struct {
		BenchmarkId
		GitRef    string
		StartedAt string
		Result    Result
	}

	// Comparison allows comparison of two Result
	// that share the same BenchmarkId.
	Comparison struct {
		BenchmarkId
		Current, Last Result

		// Difference between Current and Last.
		Diff Result
	}

	DetailsArray    []Details
	ComparisonArray []Comparison
)

// NewDetails creates a new Details.
func NewDetails(benchmarkId BenchmarkId, gitRef string, startedAt string, result Result) *Details {
	return &Details{
		BenchmarkId: benchmarkId,
		GitRef:      gitRef,
		Result:      result,
		StartedAt:   startedAt,
	}
}

// NewBenchmarkId creates a new BenchmarkId.
func NewBenchmarkId(pkgName string, name string, subBenchmarkName string) *BenchmarkId {
	return &BenchmarkId{
		PkgName:          pkgName,
		Name:             name,
		SubBenchmarkName: subBenchmarkName,
	}
}

// NewResult creates a new Result.
func NewResult(ops, NSPerOp, MBPerSec, BytesPerOp, AllocsPerOp float64) *Result {
	return &Result{
		Ops:         ops,
		NSPerOp:     NSPerOp,
		MBPerSec:    MBPerSec,
		BytesPerOp:  BytesPerOp,
		AllocsPerOp: AllocsPerOp,
	}
}

// MergeDetails merges two DetailsArray into a single
// ComparisonArray.
func MergeDetails(currentMbd, lastReleaseMbd DetailsArray) (compareMbs ComparisonArray) {
	for _, details := range currentMbd {
		compareMb := Comparison{
			BenchmarkId: details.BenchmarkId,
			Current:     details.Result,
		}
		for j := 0; j < len(lastReleaseMbd); j++ {
			if lastReleaseMbd[j].BenchmarkId == details.BenchmarkId {
				compareMb.Last = lastReleaseMbd[j].Result
				compareMb.Diff.NSPerOp = (compareMb.Current.NSPerOp - compareMb.Last.NSPerOp) / compareMb.Current.NSPerOp * 100 * -1
				compareMb.Diff.Ops = (compareMb.Current.Ops - compareMb.Last.Ops) / compareMb.Current.Ops * 100 * -1
				compareMb.Diff.BytesPerOp = (compareMb.Current.BytesPerOp - compareMb.Last.BytesPerOp) / compareMb.Current.BytesPerOp * 100 * -1
				compareMb.Diff.MBPerSec = (compareMb.Current.MBPerSec - compareMb.Last.MBPerSec) / compareMb.Current.MBPerSec * 100 * -1
				compareMb.Diff.AllocsPerOp = (compareMb.Current.AllocsPerOp - compareMb.Last.AllocsPerOp) / compareMb.Current.AllocsPerOp * 100 * -1
				math.CheckForNaN(&compareMb.Diff, 0)
				break
			}
		}
		compareMbs = append(compareMbs, compareMb)
	}
	return compareMbs
}

// ReduceSimpleMedianByName reduces a DetailsArray by merging
// all Details with the same benchmark name into a single
// one. The results of each Details correspond to the median
// of the merged elements.
func (mbd DetailsArray) ReduceSimpleMedianByName() (reduceMbd DetailsArray) {
	sort.SliceStable(mbd, func(i, j int) bool {
		if mbd[i].PkgName == mbd[j].PkgName {
			if mbd[i].Name == mbd[j].Name {
				return mbd[i].SubBenchmarkName < mbd[j].SubBenchmarkName
			}
			return mbd[i].Name < mbd[j].Name
		}
		return mbd[i].PkgName < mbd[j].PkgName
	})

	reduceMbd = mbd.mergeUsingCondition(func(i, j int) bool {
		return mbd[i].Name == mbd[j].Name && mbd[i].PkgName == mbd[j].PkgName && mbd[i].SubBenchmarkName == mbd[j].SubBenchmarkName
	})
	return reduceMbd
}

// ReduceSimpleMedianByGitRef reduces a DetailsArray by merging
// all Details with the same git ref into a single
// one. The results of each Details correspond to the median
// of the merged elements.
func (mbd DetailsArray) ReduceSimpleMedianByGitRef() (reduceMbd DetailsArray) {
	sort.SliceStable(mbd, func(i, j int) bool {
		return mbd[i].GitRef < mbd[j].GitRef
	})
	reduceMbd = mbd.mergeUsingCondition(func(i, j int) bool {
		return mbd[i].GitRef == mbd[j].GitRef
	})
	return reduceMbd
}

// mergeUsingCondition is used to merge the DetailsArray based on the compare condition provided
func (mbd DetailsArray) mergeUsingCondition(compareCondition func(i, j int) bool) (reduceMbd DetailsArray) {
	for i := 0; i < len(mbd); {
		var j int
		var interOps []float64
		var interNSPerOp []float64
		var interMBPerSec []float64
		var interBytesPerOp []float64
		var interAllocsPerOp []float64
		for j = i; j < len(mbd) && compareCondition(i, j); j++ {
			interOps = append(interOps, mbd[j].Result.Ops)
			interNSPerOp = append(interNSPerOp, mbd[j].Result.NSPerOp)
			interMBPerSec = append(interMBPerSec, mbd[j].Result.MBPerSec)
			interBytesPerOp = append(interBytesPerOp, mbd[j].Result.BytesPerOp)
			interAllocsPerOp = append(interAllocsPerOp, mbd[j].Result.AllocsPerOp)
		}

		interOpsResult := math.MedianFloat(interOps)
		interNSPerOpResult := math.MedianFloat(interNSPerOp)
		interMBPerSecResult := math.MedianFloat(interMBPerSec)
		interBytesPerOpResult := math.MedianFloat(interBytesPerOp)
		interAllocsPerOpResult := math.MedianFloat(interAllocsPerOp)
		reduceMbd = append(reduceMbd, *NewDetails(*NewBenchmarkId(mbd[i].PkgName, mbd[i].Name, mbd[i].SubBenchmarkName), mbd[i].GitRef, mbd[i].StartedAt, *NewResult(interOpsResult, interNSPerOpResult, interMBPerSecResult, interBytesPerOpResult, interAllocsPerOpResult)))
		i = j
	}
	return reduceMbd
}

// GetResultsForGitRef will fetch and return a DetailsArray
// containing all the Details linked to the given git commit SHA.
func GetResultsForGitRef(ref string, client *mysql.Client) (mrs DetailsArray, err error) {
	result, err := client.Select("select m.pkg_name, m.name, md.name, md.n, md.ns_per_op, md.bytes_per_op,"+
		" md.allocs_per_op, md.mb_per_sec FROM microbenchmark m, microbenchmark_details md where m.git_ref = ? AND "+
		"md.microbenchmark_no = m.microbenchmark_no order by m.microbenchmark_no desc", ref)
	if err != nil {
		return nil, err
	}

	for result.Next() {
		var res Details
		res.GitRef = ref
		err = result.Scan(&res.PkgName, &res.Name, &res.SubBenchmarkName, &res.Result.Ops, &res.Result.NSPerOp, &res.Result.BytesPerOp,
			&res.Result.AllocsPerOp, &res.Result.MBPerSec)
		if err != nil {
			return nil, err
		}
		mrs = append(mrs, res)
	}
	return mrs, nil
}

// GetLatestResultsFor will fetch and return a DetailsArray
// containing all the Details linked to latest runs of the given benchmark name.
func GetLatestResultsFor(name, subBenchmarkName string, count int, client *mysql.Client) (mrs DetailsArray, err error) {
	query := "select m.pkg_name, m.name, md.name, m.git_ref , md.n, md.ns_per_op, md.bytes_per_op," +
		" md.allocs_per_op, md.mb_per_sec, m.started_at  from (select microbenchmark_no, pkg_name, name, microbenchmark.git_ref, started_at" +
		" from microbenchmark join execution on exec_uuid = uuid where name = ? and source = \"cron\" and status = \"finished\" order by started_at desc limit ?) m, " +
		"microbenchmark_details md where md.microbenchmark_no = m.microbenchmark_no and md.name = ?"
	rows, err := client.Select(query, name, count, subBenchmarkName)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var res Details
		err = rows.Scan(&res.PkgName, &res.Name, &res.SubBenchmarkName, &res.GitRef, &res.Result.Ops, &res.Result.NSPerOp, &res.Result.BytesPerOp,
			&res.Result.AllocsPerOp, &res.Result.MBPerSec, &res.StartedAt)
		if err != nil {
			return nil, err
		}
		mrs = append(mrs, res)
	}
	return mrs, nil
}

func (r Result) OpsStr() string {
	if r.Ops == 0 {
		return "N/A"
	}
	return humanize.Comma(int64(r.Ops))
}

func (r Result) NSPerOpStr() string {
	if r.NSPerOp == 0 {
		return "N/A"
	}

	return humanize.FormatFloat("#,###.#", r.NSPerOp)
}

func (r Result) NSPerOpToDurationStr() string {
	if r.NSPerOp == 0 {
		return "N/A"
	}

	dur, _ := time.ParseDuration(fmt.Sprintf("%fns", r.NSPerOp))
	str := dur.String()
	i := strings.IndexFunc(str, func(r rune) bool {
		return !unicode.IsNumber(r) && r != '.'
	})
	durStr := str[:i]
	durUnit := str[i:]
	durFloat, _ := strconv.ParseFloat(durStr, 64)
	return fmt.Sprintf("%.2f %s", durFloat, durUnit)
}

func (r Result) MBPerSecStr() string {
	if r.MBPerSec == 0 {
		return "N/A"
	}

	return humanize.Bytes(uint64(r.MBPerSec)) + "/s"
}
func (r Result) BytesPerOpStr() string {
	if r.BytesPerOp == 0 {
		return "N/A"
	}

	return humanize.Bytes(uint64(r.BytesPerOp)) + "/op"
}
func (r Result) AllocsPerOpStr() string {
	if r.AllocsPerOp == 0 {
		return "N/A"
	}

	return humanize.SI(r.AllocsPerOp, "")
}
