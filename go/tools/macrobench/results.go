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
	"errors"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/vitessio/arewefastyet/go/storage"

	"github.com/dustin/go-humanize"
	"github.com/vitessio/arewefastyet/go/exec/metrics"
	"github.com/vitessio/arewefastyet/go/storage/mysql"
	awftmath "github.com/vitessio/arewefastyet/go/tools/math"
)

type (
	// QPS represents the QPS table. This table contains the raw
	// results of a macro benchmark.
	QPS struct {
		ID     int
		RefID  int
		Total  float64 `json:"total"`
		Reads  float64 `json:"reads"`
		Writes float64 `json:"writes"`
		Other  float64 `json:"other"`
	}

	// Result represents both OLTP and TPCC tables.
	// The two tables share the same schema and can thus be grouped
	// under an unique go struct.
	Result struct {
		ID         int
		QPS        QPS     `json:"qps"`
		TPS        float64 `json:"tps"`
		Latency    float64 `json:"latency"`
		Errors     float64 `json:"errors"`
		Reconnects float64 `json:"reconnects"`
		Time       int     `json:"time"`
		Threads    float64 `json:"threads"`
	}

	// BenchmarkID is used to identify a macro benchmark using its database's ID, the
	// source from which the benchmark was triggered and its creation date.
	BenchmarkID struct {
		ID        int
		Source    string
		CreatedAt *time.Time
		ExecUUID  string
	}

	// Details represents the entire macro benchmark and its sub
	// components. It has a BenchmarkID (ID, creation date, source of the benchmark),
	// the git reference that was used, and its results represented by a Result.
	// This struct encapsulates the "benchmark", "qps" and ("OLTP" or "TPCC") database tables.
	Details struct {
		BenchmarkID

		// refers to commit
		GitRef  string
		Result  Result
		Metrics metrics.ExecutionMetrics
	}

	// Comparison contains two Details and their difference in a
	// Result field.
	Comparison struct {
		Right, Left Details
		Diff        Result
		DiffMetrics metrics.ExecutionMetrics
	}

	ResultsArray []Result
	DetailsArray []Details

	ComparisonArray []Comparison
)

func newBenchmarkID(ID int, source string, createdAt *time.Time) *BenchmarkID {
	return &BenchmarkID{ID: ID, Source: source, CreatedAt: createdAt}
}

func newDetails(benchmarkID BenchmarkID, gitRef string, result Result, metrics metrics.ExecutionMetrics) *Details {
	return &Details{BenchmarkID: benchmarkID, GitRef: gitRef, Result: result, Metrics: metrics}
}

func newQPS(total float64, reads float64, writes float64, other float64) *QPS {
	return &QPS{Total: total, Reads: reads, Writes: writes, Other: other}
}

func newResult(QPS QPS, TPS float64, latency float64, errors float64, reconnects float64, time int, threads float64) *Result {
	return &Result{QPS: QPS, TPS: TPS, Latency: latency, Errors: errors, Reconnects: reconnects, Time: time, Threads: threads}
}

// CompareDetailsArrays compare two DetailsArray and return
// their comparison in a ComparisonArray.
func CompareDetailsArrays(references, compares DetailsArray) (compared ComparisonArray) {
	emptyCmp := Comparison{
		Right: Details{
			Metrics: metrics.NewExecMetrics(),
		},
		Left: Details{
			Metrics: metrics.NewExecMetrics(),
		},
		DiffMetrics: metrics.NewExecMetrics(),
	}
	for i := 0; i < int(math.Max(float64(len(references)), float64(len(compares)))); i++ {
		cmp := emptyCmp
		if i < len(references) {
			cmp.Right = references[i]
		}
		if i < len(compares) {
			cmp.Left = compares[i]
		}
		if cmp.Left.GitRef != "" && cmp.Right.GitRef != "" {
			compareResult := cmp.Left.Result
			referenceResult := cmp.Right.Result
			cmp.Diff.QPS.Total = (referenceResult.QPS.Total - compareResult.QPS.Total) / compareResult.QPS.Total * 100
			cmp.Diff.QPS.Reads = (referenceResult.QPS.Reads - compareResult.QPS.Reads) / compareResult.QPS.Reads * 100
			cmp.Diff.QPS.Writes = (referenceResult.QPS.Writes - compareResult.QPS.Writes) / compareResult.QPS.Writes * 100
			cmp.Diff.QPS.Other = (referenceResult.QPS.Other - compareResult.QPS.Other) / compareResult.QPS.Other * 100
			cmp.Diff.TPS = (referenceResult.TPS - compareResult.TPS) / compareResult.TPS * 100
			cmp.Diff.Latency = (compareResult.Latency - referenceResult.Latency) / referenceResult.Latency * 100
			cmp.Diff.Reconnects = (compareResult.Reconnects - referenceResult.Reconnects) / referenceResult.Reconnects * 100
			cmp.Diff.Errors = (compareResult.Errors - referenceResult.Errors) / referenceResult.Errors * 100
			cmp.Diff.Time = int(float64(compareResult.Time) - (float64(referenceResult.Time))/float64(referenceResult.Time)*100)
			cmp.Diff.Threads = (compareResult.Threads - referenceResult.Threads) / referenceResult.Threads * 100
			awftmath.CheckForNaN(&cmp.Diff, 0)
			awftmath.CheckForNaN(&cmp.Diff.QPS, 0)
			cmp.DiffMetrics = metrics.CompareTwo(cmp.Left.Metrics, cmp.Right.Metrics)
		}
		compared = append(compared, cmp)
	}
	if len(compared) == 0 {
		compared = append(compared, emptyCmp)
	}
	return compared
}

// mergeMedian will merge a ResultsArray into a single Result
// by calculating the median of all elements in the array.
func (mrs ResultsArray) mergeMedian() (mergedResult Result) {
	inter := struct {
		total      []float64
		reads      []float64
		writes     []float64
		other      []float64
		tps        []float64
		latency    []float64
		errors     []float64
		reconnects []float64
		time       []int
		threads    []float64
	}{}

	for _, mr := range mrs {
		inter.total = append(inter.total, mr.QPS.Total)
		inter.reads = append(inter.reads, mr.QPS.Reads)
		inter.writes = append(inter.writes, mr.QPS.Writes)
		inter.other = append(inter.other, mr.QPS.Other)
		inter.tps = append(inter.tps, mr.TPS)
		inter.latency = append(inter.latency, mr.Latency)
		inter.errors = append(inter.errors, mr.Errors)
		inter.reconnects = append(inter.reconnects, mr.Reconnects)
		inter.time = append(inter.time, mr.Time)
		inter.threads = append(inter.threads, mr.Threads)
	}

	mergedResult.QPS.Total = awftmath.MedianFloat(inter.total)
	mergedResult.QPS.Reads = awftmath.MedianFloat(inter.reads)
	mergedResult.QPS.Writes = awftmath.MedianFloat(inter.writes)
	mergedResult.QPS.Other = awftmath.MedianFloat(inter.other)
	mergedResult.TPS = awftmath.MedianFloat(inter.tps)
	mergedResult.Latency = awftmath.MedianFloat(inter.latency)
	mergedResult.Errors = awftmath.MedianFloat(inter.errors)
	mergedResult.Reconnects = awftmath.MedianFloat(inter.reconnects)
	mergedResult.Time = int(awftmath.MedianInt(inter.time))
	mergedResult.Threads = awftmath.MedianFloat(inter.threads)
	return mergedResult
}

// ReduceSimpleMedian reduces the given DetailsArray by
// merging altogether the elements that share the same GitRef.
// During the reduce, the math.MedianFloat and math.MedianInt methods
// are applied on the different Result.
func (mabd DetailsArray) ReduceSimpleMedian() (reduceMabd DetailsArray) {
	sort.SliceStable(mabd, func(i, j int) bool {
		return mabd[i].GitRef < mabd[j].GitRef
	})
	for i := 0; i < len(mabd); {
		var j int
		interResults := ResultsArray{}
		interMetrics := metrics.ExecutionMetricsArray{}
		for j = i; j < len(mabd) && mabd[i].GitRef == mabd[j].GitRef; j++ {
			interResults = append(interResults, mabd[j].Result)
			interMetrics = append(interMetrics, mabd[j].Metrics)
		}

		reducedResult := interResults.mergeMedian()
		reduceMabd = append(reduceMabd, Details{
			GitRef:  mabd[i].GitRef,
			Result:  reducedResult,
			Metrics: interMetrics.Median(),
		})
		i = j
	}
	return reduceMabd
}

func (mbr Result) TPSStr() string {
	return humanize.FormatFloat("#,###.#", mbr.TPS)
}

func (mbr Result) LatencyStr() string {
	return humanize.FormatFloat("#,###.#", mbr.Latency)
}

func (mbr Result) ErrorsStr() string {
	return humanize.FormatFloat("#,###.#", mbr.Errors)
}

func (mbr Result) ReconnectsStr() string {
	return humanize.FormatFloat("#,###.#", mbr.Reconnects)
}

func (mbr Result) TimeStr() string {
	return humanize.Comma(int64(mbr.Time))
}

func (mbr Result) ThreadsStr() string {
	return humanize.FormatFloat("#,###.#", mbr.Threads)
}

func (qps QPS) TotalStr() string {
	return humanize.FormatFloat("#,###.#", qps.Total)
}

func (qps QPS) ReadsStr() string {
	return humanize.FormatFloat("#,###.#", qps.Reads)
}

func (qps QPS) WritesStr() string {
	return humanize.FormatFloat("#,###.#", qps.Writes)
}

func (qps QPS) OtherStr() string {
	return humanize.FormatFloat("#,###.#", qps.Other)
}

// GetDetailsArraysFromAllTypes returns a slice of Details based on the given git ref and Types.
func GetDetailsArraysFromAllTypes(sha string, planner PlannerVersion, dbclient storage.SQLClient, types []string) (map[string]DetailsArray, error) {
	macros := map[string]DetailsArray{}
	for _, mtype := range types {
		macro, err := GetResultsForGitRefAndPlanner(mtype, sha, planner, dbclient)
		if err != nil {
			return nil, err
		}

		// Get the execution metrics of each macrobenchmark details
		for i, details := range macro {
			macro[i].Metrics, err = metrics.GetExecutionMetricsSQL(dbclient, details.ExecUUID)
			if err != nil {
				return nil, err
			}
		}

		macros[mtype] = macro.ReduceSimpleMedian()
	}
	return macros, nil
}

// GetResultsForLastDays returns a slice Details based on a given macro benchmark type.
// The type can either be OLTP or TPCC. Using that type, the function will generate a query using
// the *mysql.Client. The query will select only results that were added between now and lastDays.
func GetResultsForLastDays(macroType string, source string, planner PlannerVersion, lastDays int, client storage.SQLClient) (macrodetails DetailsArray, err error) {
	upperMacroType := strings.ToUpper(macroType)
	query := "SELECT info.macrobenchmark_id, e.git_ref, e.source, e.finished_at, IFNULL(e.uuid, ''), " +
		"results.tps, results.latency, results.errors, results.reconnects, results.time, results.threads, " +
		"results.total_qps, results.reads_qps, results.writes_qps, results.other_qps " +
		"FROM execution AS e, macrobenchmark AS info, macrobenchmark_results AS results " +
		"WHERE e.uuid = info.exec_uuid AND e.status = \"finished\" AND e.finished_at BETWEEN DATE(NOW()) - INTERVAL ? DAY AND DATE(NOW() + INTERVAL 1 DAY) " +
		"AND e.source = ? AND info.vtgate_planner_version = ? AND info.macrobenchmark_id = results.macrobenchmark_id AND info.type = ?"

	result, err := client.Select(query, lastDays, source, planner, upperMacroType)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		var res Details
		err = result.Scan(&res.ID, &res.GitRef, &res.Source, &res.CreatedAt, &res.ExecUUID, &res.Result.TPS, &res.Result.Latency,
			&res.Result.Errors, &res.Result.Reconnects, &res.Result.Time, &res.Result.Threads,
			&res.Result.QPS.Total, &res.Result.QPS.Reads, &res.Result.QPS.Writes, &res.Result.QPS.Other)
		if err != nil {
			return nil, err
		}
		macrodetails = append(macrodetails, res)
	}
	return macrodetails, nil
}

// GetResultsForGitRefAndPlanner returns a slice of Details based on the given git ref
// and macro benchmark Type.
func GetResultsForGitRefAndPlanner(macroType string, ref string, planner PlannerVersion, client storage.SQLClient) (macrodetails DetailsArray, err error) {
	upperMacroType := strings.ToUpper(macroType)
	query := "SELECT info.macrobenchmark_id, e.git_ref, e.source, e.finished_at, IFNULL(info.exec_uuid, ''), " +
		"results.tps, results.latency, results.errors, results.reconnects, results.time, results.threads, " +
		"results.total_qps, results.reads_qps, results.writes_qps, results.other_qps " +
		"FROM execution AS e, macrobenchmark AS info, macrobenchmark_results AS results " +
		"WHERE e.uuid = info.exec_uuid AND e.status = \"finished\" AND e.git_ref = ? AND info.vtgate_planner_version = ? AND info.macrobenchmark_id = results.macrobenchmark_id AND info.type = ?"

	result, err := client.Select(query, ref, planner, upperMacroType)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		var res Details
		err = result.Scan(&res.ID, &res.GitRef, &res.Source, &res.CreatedAt, &res.ExecUUID, &res.Result.TPS, &res.Result.Latency,
			&res.Result.Errors, &res.Result.Reconnects, &res.Result.Time, &res.Result.Threads,
			&res.Result.QPS.Total, &res.Result.QPS.Reads, &res.Result.QPS.Writes, &res.Result.QPS.Other)
		if err != nil {
			return nil, err
		}
		macrodetails = append(macrodetails, res)
	}
	return macrodetails, nil
}

// insertToMySQL inserts the given MacroBenchmarkResult to MySQL using a *mysql.Client.
// The MacroBenchmarkResults gets added in one of macrobenchmark's children tables.
// Depending on the MacroBenchmarkType, the insert will be routed to a specific children table.
// The children table QPS is also inserted.
func (mbr *Result) insertToMySQL(macrobenchmarkID int, client storage.SQLClient) error {
	if client == nil {
		return errors.New(mysql.ErrorClientConnectionNotInitialized)
	}

	// insert Result
	queryResult := "INSERT INTO macrobenchmark_results(macrobenchmark_id, tps, latency, errors, reconnects, time, threads, total_qps, reads_qps, writes_qps, other_qps) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := client.Insert(queryResult, macrobenchmarkID, mbr.TPS, mbr.Latency, mbr.Errors, mbr.Reconnects, mbr.Time, mbr.Threads, mbr.QPS.Total, mbr.QPS.Reads, mbr.QPS.Writes, mbr.QPS.Other)
	if err != nil {
		return err
	}
	return nil
}
