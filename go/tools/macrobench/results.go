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
	"sort"
	"strings"
	"time"

	"github.com/vitessio/arewefastyet/go/exec/metrics"
	"github.com/vitessio/arewefastyet/go/storage"
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
		Queries    int     `json:"queries"`
		QPS        QPS     `json:"qps"`
		TPS        float64 `json:"tps"`
		Latency    float64 `json:"latency"`
		Errors     float64 `json:"errors"`
		Reconnects float64 `json:"reconnects"`
		Time       int     `json:"time"`
		Threads    float64 `json:"threads"`
	}

	qpsAsSlice struct {
		total  []float64
		reads  []float64
		writes []float64
		other  []float64
	}

	metricsAsSlice struct {
		totalComponentsCPUTime []float64
		componentsCPUTime      map[string][]float64

		totalComponentsMemStatsAllocBytes []float64
		componentsMemStatsAllocBytes      map[string][]float64
	}

	resultAsSlice struct {
		qps qpsAsSlice

		tps        []float64
		latency    []float64
		errors     []float64
		reconnects []float64
		time       []int
		threads    []float64

		metrics metricsAsSlice
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

	BenchmarkResults struct {
		Results ResultsArray
		Metrics metrics.ExecutionMetricsArray
	}

	ResultsArray []Result
	DetailsArray []Details

	ComparisonArray []Comparison

	DailySummary struct {
		CreatedAt *time.Time
		QPSTotal  float64
	}
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

func (br BenchmarkResults) asSlice() resultAsSlice {
	s := br.Results.resultsArrayToSlice()
	s.metrics = metricsToSlice(br.Metrics)
	return s
}

func metricsToSlice(metrics metrics.ExecutionMetricsArray) metricsAsSlice {
	var s metricsAsSlice
	s.componentsCPUTime = make(map[string][]float64)
	s.componentsMemStatsAllocBytes = make(map[string][]float64)
	for _, metricRow := range metrics {
		s.totalComponentsCPUTime = append(s.totalComponentsCPUTime, metricRow.TotalComponentsCPUTime)
		for name, value := range metricRow.ComponentsCPUTime {
			s.componentsCPUTime[name] = append(s.componentsCPUTime[name], value)
		}

		s.totalComponentsMemStatsAllocBytes = append(s.totalComponentsMemStatsAllocBytes, metricRow.TotalComponentsMemStatsAllocBytes)
		for name, value := range metricRow.ComponentsMemStatsAllocBytes {
			s.componentsMemStatsAllocBytes[name] = append(s.componentsMemStatsAllocBytes[name], value)
		}
	}
	return s
}

func (mrs ResultsArray) resultsArrayToSlice() resultAsSlice {
	var ras resultAsSlice
	for _, mr := range mrs {
		ras.qps.total = append(ras.qps.total, mr.QPS.Total)
		ras.qps.reads = append(ras.qps.reads, mr.QPS.Reads)
		ras.qps.writes = append(ras.qps.writes, mr.QPS.Writes)
		ras.qps.other = append(ras.qps.other, mr.QPS.Other)
		ras.tps = append(ras.tps, mr.TPS)
		ras.latency = append(ras.latency, mr.Latency)
		ras.errors = append(ras.errors, mr.Errors)
		ras.reconnects = append(ras.reconnects, mr.Reconnects)
		ras.time = append(ras.time, mr.Time)
		ras.threads = append(ras.threads, mr.Threads)
	}
	return ras
}

// mergeMedian will merge a ResultsArray into a single Result
// by calculating the median of all elements in the array.
func (mrs ResultsArray) mergeMedian() (mergedResult Result) {
	ras := mrs.resultsArrayToSlice()
	mergedResult.QPS.Total = awftmath.MedianFloat(ras.qps.total)
	mergedResult.QPS.Reads = awftmath.MedianFloat(ras.qps.reads)
	mergedResult.QPS.Writes = awftmath.MedianFloat(ras.qps.writes)
	mergedResult.QPS.Other = awftmath.MedianFloat(ras.qps.other)
	mergedResult.TPS = awftmath.MedianFloat(ras.tps)
	mergedResult.Latency = awftmath.MedianFloat(ras.latency)
	mergedResult.Errors = awftmath.MedianFloat(ras.errors)
	mergedResult.Reconnects = awftmath.MedianFloat(ras.reconnects)
	mergedResult.Time = int(awftmath.MedianInt(ras.time))
	mergedResult.Threads = awftmath.MedianFloat(ras.threads)
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

func GetDetailsFromAllTypes(sha string, planner PlannerVersion, dbclient storage.SQLClient, types []string) (map[string]Details, error) {
	details, err := GetDetailsArraysFromAllTypes(sha, planner, dbclient, types)
	if err != nil {
		return nil, err
	}
	result := make(map[string]Details, len(details))
	for s, array := range details {
		var d Details
		d.Metrics = metrics.NewExecMetrics()
		if len(array) == 1 {
			d = array[0]
		}
		result[s] = d
	}
	return result, nil
}

// GetDetailsArraysFromAllTypes returns a slice of Details based on the given git ref and Types.
func GetDetailsArraysFromAllTypes(sha string, planner PlannerVersion, dbclient storage.SQLClient, types []string) (map[string]DetailsArray, error) {
	macros := map[string]DetailsArray{}
	for _, mtype := range types {
		macro, err := getResultsForGitRefAndPlanner(mtype, sha, planner, dbclient)
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
	macrodetails = []Details{}
	upperMacroType := strings.ToUpper(macroType)
	query := "SELECT info.macrobenchmark_id, e.git_ref, e.source, e.finished_at, IFNULL(e.uuid, ''), " +
		"results.tps, results.latency, results.errors, results.reconnects, results.time, results.threads, " +
		"results.total_qps, results.reads_qps, results.writes_qps, results.other_qps " +
		"FROM execution AS e, macrobenchmark AS info, macrobenchmark_results AS results " +
		"WHERE e.uuid = info.exec_uuid AND e.status = \"finished\" AND e.finished_at BETWEEN DATE(NOW()) - INTERVAL ? DAY AND DATE(NOW() + INTERVAL 1 DAY) " +
		"AND e.source = ? AND info.vtgate_planner_version = ? AND info.macrobenchmark_id = results.macrobenchmark_id AND info.type = ? " +
		"ORDER BY e.finished_at "

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

func GetSummaryForLastDays(macroType string, source string, planner PlannerVersion, lastDays int, client storage.SQLClient) (dailySummary []DailySummary, err error) {
	upperMacroType := strings.ToUpper(macroType)
	query := "SELECT e.finished_at, results.total_qps " +
		"FROM execution AS e, macrobenchmark AS info, macrobenchmark_results AS results " +
		"WHERE e.uuid = info.exec_uuid AND e.status = \"finished\" AND e.finished_at BETWEEN DATE(NOW()) - INTERVAL ? DAY AND DATE(NOW() + INTERVAL 1 DAY) " +
		"AND e.source = ? AND info.vtgate_planner_version = ? AND info.macrobenchmark_id = results.macrobenchmark_id AND info.type = ? " +
		"ORDER BY e.finished_at "

	result, err := client.Select(query, lastDays, source, planner, upperMacroType)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		var res DailySummary
		err = result.Scan(&res.CreatedAt, &res.QPSTotal)
		if err != nil {
			return nil, err
		}
		dailySummary = append(dailySummary, res)
	}
	return
}

func getBenchmarkResults(client storage.SQLClient, macroType, gitSHA string, planner PlannerVersion) (BenchmarkResults, error) {
	results, err := getResultsForGitRefAndPlanner(macroType, gitSHA, planner, client)
	if err != nil {
		return BenchmarkResults{}, err
	}

	if len(results) == 0 {
		return BenchmarkResults{}, nil
	}

	var br BenchmarkResults
	for _, result := range results {
		br.Results = append(br.Results, result.Result)

		metricsResult, err := metrics.GetExecutionMetricsSQL(client, result.ExecUUID)
		if err != nil {
			return BenchmarkResults{}, err
		}
		br.Metrics = append(br.Metrics, metricsResult)
	}
	return br, nil
}

// getResultsForGitRefAndPlanner returns a slice of Details based on the given git ref
// and macro benchmark Type.
func getResultsForGitRefAndPlanner(macroType string, ref string, planner PlannerVersion, client storage.SQLClient) (macrodetails DetailsArray, err error) {
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
	queryResult := "INSERT INTO macrobenchmark_results(macrobenchmark_id, queries, tps, latency, errors, reconnects, time, threads, total_qps, reads_qps, writes_qps, other_qps) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := client.Insert(queryResult, macrobenchmarkID, mbr.Queries, mbr.TPS, mbr.Latency, mbr.Errors, mbr.Reconnects, mbr.Time, mbr.Threads, mbr.QPS.Total, mbr.QPS.Reads, mbr.QPS.Writes, mbr.QPS.Other)
	if err != nil {
		return err
	}
	return nil
}
