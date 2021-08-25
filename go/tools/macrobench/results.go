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
	"fmt"
	"github.com/vitessio/arewefastyet/go/storage"
	"math"
	"sort"
	"strings"
	"time"

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
		Reference, Compare Details
		Diff               Result
		DiffMetrics        metrics.ExecutionMetrics
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
	for i := 0; i < int(math.Max(float64(len(references)), float64(len(compares)))); i++ {
		var cmp Comparison
		if i < len(references) {
			cmp.Reference = references[i]
		}
		if i < len(compares) {
			cmp.Compare = compares[i]
		}
		if cmp.Compare.GitRef != "" && cmp.Reference.GitRef != "" {
			cmp.Diff.QPS.Total = (cmp.Reference.Result.QPS.Total - cmp.Compare.Result.QPS.Total) / cmp.Reference.Result.QPS.Total * 100
			cmp.Diff.QPS.Reads = (cmp.Reference.Result.QPS.Reads - cmp.Compare.Result.QPS.Reads) / cmp.Reference.Result.QPS.Reads * 100
			cmp.Diff.QPS.Writes = (cmp.Reference.Result.QPS.Writes - cmp.Compare.Result.QPS.Writes) / cmp.Reference.Result.QPS.Writes * 100
			cmp.Diff.QPS.Other = (cmp.Reference.Result.QPS.Other - cmp.Compare.Result.QPS.Other) / cmp.Reference.Result.QPS.Other * 100
			cmp.Diff.TPS = (cmp.Reference.Result.TPS - cmp.Compare.Result.TPS) / cmp.Reference.Result.TPS * 100
			cmp.Diff.Latency = (cmp.Compare.Result.Latency - cmp.Reference.Result.Latency) / cmp.Compare.Result.Latency * 100
			cmp.Diff.Reconnects = (cmp.Reference.Result.Reconnects - cmp.Compare.Result.Reconnects) / cmp.Reference.Result.Reconnects * 100
			cmp.Diff.Errors = (cmp.Compare.Result.Errors - cmp.Reference.Result.Errors) / cmp.Compare.Result.Errors * 100
			cmp.Diff.Time = int((float64(cmp.Reference.Result.Time) - float64(cmp.Compare.Result.Time)) / float64(cmp.Reference.Result.Time) * 100)
			cmp.Diff.Threads = (cmp.Reference.Result.Threads - cmp.Compare.Result.Threads) / cmp.Reference.Result.Threads * 100
			awftmath.CheckForNaN(&cmp.Diff, 0)
			awftmath.CheckForNaN(&cmp.Diff.QPS, 0)
			cmp.DiffMetrics = metrics.CompareTwo(cmp.Compare.Metrics, cmp.Reference.Metrics)
		}
		compared = append(compared, cmp)
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
func GetDetailsArraysFromAllTypes(sha string, planner PlannerVersion, dbclient storage.SQLClient) (map[Type]DetailsArray, error) {
	macros := map[Type]DetailsArray{}
	for _, mtype := range Types {
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
func GetResultsForLastDays(macroType Type, source string, planner PlannerVersion, lastDays int, client storage.SQLClient) (macrodetails DetailsArray, err error) {
	if macroType != OLTP && macroType != TPCC {
		return nil, errors.New(IncorrectMacroBenchmarkType)
	}

	upperMacroType := macroType.ToUpper().String()
	query := "SELECT b.macrobenchmark_id, b.commit, b.source, b.DateTime, IFNULL(b.exec_uuid, ''), " +
		"macrotype.tps, macrotype.latency, macrotype.errors, macrotype.reconnects, macrotype.time, macrotype.threads, " +
		"qps.qps_no, qps.total_qps, qps.reads_qps, qps.writes_qps, qps.other_qps " +
		"FROM execution AS e, macrobenchmark AS b, $(MBTYPE) AS macrotype, qps AS qps " +
		"WHERE e.uuid = b.exec_uuid AND e.status = \"finished\" AND b.DateTime BETWEEN DATE(NOW()) - INTERVAL ? DAY AND DATE(NOW()) " +
		"AND b.source = ? AND b.vtgate_planner_version = ? AND b.macrobenchmark_id = macrotype.macrobenchmark_id AND macrotype.$(MBTYPE)_no = qps.$(MBTYPE)_no"

	query = strings.ReplaceAll(query, "$(MBTYPE)", upperMacroType)

	result, err := client.Select(query, lastDays, source, planner)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		var res Details
		err = result.Scan(&res.ID, &res.GitRef, &res.Source, &res.CreatedAt, &res.ExecUUID, &res.Result.TPS, &res.Result.Latency,
			&res.Result.Errors, &res.Result.Reconnects, &res.Result.Time, &res.Result.Threads, &res.Result.QPS.ID,
			&res.Result.QPS.Total, &res.Result.QPS.Reads, &res.Result.QPS.Writes, &res.Result.QPS.Other)
		if err != nil {
			return nil, err
		}
		macrodetails = append(macrodetails, res)
	}
	return macrodetails, nil
}

// GetResultsForGitRefAndPlanner returns a slice of Details based on the given git ref
// and macro benchmark Type. The type must be either OLTP or TPCC.
func GetResultsForGitRefAndPlanner(macroType Type, ref string, planner PlannerVersion, client storage.SQLClient) (macrodetails DetailsArray, err error) {
	if macroType != OLTP && macroType != TPCC {
		return nil, errors.New(IncorrectMacroBenchmarkType)
	}
	upperMacroType := macroType.ToUpper().String()
	query := "SELECT b.macrobenchmark_id, b.commit, b.source, b.DateTime, IFNULL(b.exec_uuid, ''), " +
		"macrotype.tps, macrotype.latency, macrotype.errors, macrotype.reconnects, macrotype.time, macrotype.threads, " +
		"qps.qps_no, qps.total_qps, qps.reads_qps, qps.writes_qps, qps.other_qps " +
		"FROM execution AS e, macrobenchmark AS b, $(MBTYPE) AS macrotype, qps AS qps " +
		"WHERE e.uuid = b.exec_uuid AND e.status = \"finished\" AND b.commit = ? AND b.vtgate_planner_version = ? AND b.macrobenchmark_id = macrotype.macrobenchmark_id AND macrotype.$(MBTYPE)_no = qps.$(MBTYPE)_no"

	query = strings.ReplaceAll(query, "$(MBTYPE)", upperMacroType)

	result, err := client.Select(query, ref, planner)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		var res Details
		err = result.Scan(&res.ID, &res.GitRef, &res.Source, &res.CreatedAt, &res.ExecUUID, &res.Result.TPS, &res.Result.Latency,
			&res.Result.Errors, &res.Result.Reconnects, &res.Result.Time, &res.Result.Threads, &res.Result.QPS.ID,
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
func (mbr *Result) insertToMySQL(benchmarkType Type, macrobenchmarkID int, client storage.SQLClient) error {
	if client == nil {
		return errors.New(mysql.ErrorClientConnectionNotInitialized)
	}
	if benchmarkType == "" {
		return errors.New(IncorrectMacroBenchmarkType)
	}

	// insert Result
	queryResult := fmt.Sprintf("INSERT INTO %s(macrobenchmark_id, tps, latency, errors, reconnects, time, threads) VALUES(?, ?, ?, ?, ?, ?, ?)", benchmarkType.ToUpper().String())
	resultID, err := client.Insert(queryResult, macrobenchmarkID, mbr.TPS, mbr.Latency, mbr.Errors, mbr.Reconnects, mbr.Time, mbr.Threads)
	if err != nil {
		return err
	}

	// insert QPS
	queryQPS := fmt.Sprintf("INSERT INTO qps(%s, total_qps, reads_qps, writes_qps, other_qps) VALUES(?, ?, ?, ?, ?)", benchmarkType.ToUpper().String()+"_no")
	_, err = client.Insert(queryQPS, resultID, mbr.QPS.Total, mbr.QPS.Reads, mbr.QPS.Writes, mbr.QPS.Other)
	if err != nil {
		return err
	}
	return nil
}
