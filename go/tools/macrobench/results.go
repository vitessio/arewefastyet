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
	"github.com/dustin/go-humanize"
	"fmt"
	"github.com/vitessio/arewefastyet/go/mysql"
	"github.com/vitessio/arewefastyet/go/tools/math"
	"sort"
	"strings"
	"time"
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

	// MacroBenchmarkResult represents both OLTP and TPCC tables.
	// The two tables share the same schema and can thus be grouped
	// under an unique go struct.
	MacroBenchmarkResult struct {
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
	}

	// MacroBenchmarkDetails represents the entire macro benchmark and its sub
	// components. It has a BenchmarkID (ID, creation date, source of the benchmark),
	// the git reference that was used, and its results represented by a MacroBenchmarkResult.
	// This struct encapsulates the "benchmark", "qps" and ("OLTP" or "TPCC") database tables.
	MacroBenchmarkDetails struct {
		BenchmarkID

		// refers to commit
		GitRef string
		Result MacroBenchmarkResult
	}

	// Comparison contains two MacroBenchmarkDetails and their difference in a
	// MacroBenchmarkResult field.
	Comparison struct {
		GitRef string
		Reference, Compare MacroBenchmarkResult
		Diff MacroBenchmarkResult
	}

	MacroBenchmarkResultsArray []MacroBenchmarkResult
	MacroBenchmarkDetailsArray []MacroBenchmarkDetails

	ComparisonArray []Comparison
)

func newBenchmarkID(ID int, source string, createdAt *time.Time) *BenchmarkID {
	return &BenchmarkID{ID: ID, Source: source, CreatedAt: createdAt}
}

func newMacroBenchmarkDetails(benchmarkID BenchmarkID, gitRef string, result MacroBenchmarkResult) *MacroBenchmarkDetails {
	return &MacroBenchmarkDetails{BenchmarkID: benchmarkID, GitRef: gitRef, Result: result}
}

func newQPS(total float64, reads float64, writes float64, other float64) *QPS {
	return &QPS{Total: total, Reads: reads, Writes: writes, Other: other}
}

func newMacroBenchmarkResult(QPS QPS, TPS float64, latency float64, errors float64, reconnects float64, time int, threads float64) *MacroBenchmarkResult {
	return &MacroBenchmarkResult{QPS: QPS, TPS: TPS, Latency: latency, Errors: errors, Reconnects: reconnects, Time: time, Threads: threads}
}

func CompareDetailsArrays(references, compares MacroBenchmarkDetailsArray) (compared ComparisonArray) {
	for _, ref := range references {
		var cmp Comparison
		cmp.GitRef = ref.GitRef
		cmp.Reference = ref.Result
		for j := 0; j < len(compares); j++ {
			if compares[j].GitRef == ref.GitRef {
				cmp.Compare = compares[j].Result
				cmp.Diff.QPS.Total = cmp.Reference.QPS.Total / cmp.Compare.QPS.Total
				cmp.Diff.QPS.Reads = cmp.Reference.QPS.Reads / cmp.Compare.QPS.Reads
				cmp.Diff.QPS.Writes = cmp.Reference.QPS.Writes / cmp.Compare.QPS.Writes
				cmp.Diff.QPS.Other = cmp.Reference.QPS.Other / cmp.Compare.QPS.Other
				cmp.Diff.TPS = cmp.Reference.TPS / cmp.Compare.TPS
				cmp.Diff.Latency = cmp.Reference.Latency / cmp.Compare.Latency
				cmp.Diff.Reconnects = cmp.Reference.Reconnects / cmp.Compare.Reconnects
				cmp.Diff.Errors = cmp.Reference.Errors / cmp.Compare.Errors
				cmp.Diff.Time = cmp.Reference.Time / cmp.Compare.Time
				cmp.Diff.Threads = cmp.Reference.Threads / cmp.Compare.Threads
				break
			}
		}
		compared = append(compared, cmp)
	}
	return compared
}

func (mrs MacroBenchmarkResultsArray) mergeMedian() (mergedResult MacroBenchmarkResult) {
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

	mergedResult.QPS.Total = math.MedianFloat(inter.total)
	mergedResult.QPS.Reads = math.MedianFloat(inter.reads)
	mergedResult.QPS.Writes = math.MedianFloat(inter.writes)
	mergedResult.QPS.Other = math.MedianFloat(inter.other)
	mergedResult.TPS = math.MedianFloat(inter.tps)
	mergedResult.Latency = math.MedianFloat(inter.latency)
	mergedResult.Errors = math.MedianFloat(inter.errors)
	mergedResult.Reconnects = math.MedianFloat(inter.reconnects)
	mergedResult.Time = int(math.MedianInt(inter.time))
	mergedResult.Threads = math.MedianFloat(inter.threads)
	return mergedResult
}

// ReduceSimpleMedian reduces the given MacroBenchmarkDetailsArray by
// merging altogether the elements that share the same GitRef.
// During the reduce, the math.MedianFloat and math.MedianInt methods
// are applied on the different MacroBenchmarkResult.
func (mabd MacroBenchmarkDetailsArray) ReduceSimpleMedian() (reduceMabd MacroBenchmarkDetailsArray) {
	sort.SliceStable(mabd, func(i, j int) bool {
		return mabd[i].GitRef < mabd[j].GitRef
	})
	for i := 0; i < len(mabd); {
		var j int
		interResults := MacroBenchmarkResultsArray{}
		for j = i; j < len(mabd) && mabd[i].GitRef == mabd[j].GitRef; j++ {
			interResults = append(interResults, mabd[j].Result)
		}

		reducedResult := interResults.mergeMedian()
		reduceMabd = append(reduceMabd, MacroBenchmarkDetails{
			GitRef: mabd[i].GitRef,
			Result: reducedResult,
		})
		i = j
	}
	return reduceMabd
}

func (mbr MacroBenchmarkResult) TPSStr() string {
	return humanize.FormatFloat("#,###.#", mbr.TPS)
}

func (mbr MacroBenchmarkResult) LatencyStr() string {
	return humanize.FormatFloat("#,###.#", mbr.Latency)
}

func (mbr MacroBenchmarkResult) ErrorsStr() string {
	return humanize.FormatFloat("#,###.#", mbr.Errors)
}

func (mbr MacroBenchmarkResult) ReconnectsStr() string {
	return humanize.FormatFloat("#,###.#", mbr.Reconnects)
}

func (mbr MacroBenchmarkResult) TimeStr() string {
	return humanize.Comma(int64(mbr.Time))
}

func (mbr MacroBenchmarkResult) ThreadsStr() string {
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

// GetResultsForLastDays returns a slice MacroBenchmarkDetails based on a given macro benchmark type.
// The type can either be OLTP or TPCC. Using that type, the function will generate a query using
// the *mysql.Client. The query will select only results that were added between now and lastDays.
func GetResultsForLastDays(macroType Type, source string, lastDays int, client *mysql.Client) (macrodetails MacroBenchmarkDetailsArray, err error) {
	if macroType != OLTP && macroType != TPCC {
		return nil, errors.New(IncorrectMacroBenchmarkType)
	}

	upperMacroType := macroType.ToUpper().String()
	query := "SELECT b.macrobenchmark_id, b.commit, b.source, b.DateTime, " +
		"macrotype.tps, macrotype.latency, macrotype.errors, macrotype.reconnects, macrotype.time, macrotype.threads, " +
		"qps.qps_no, qps.total_qps, qps.reads_qps, qps.writes_qps, qps.other_qps " +
		"FROM macrobenchmark AS b, $(MBTYPE) AS macrotype, qps AS qps " +
		"WHERE b.DateTime BETWEEN DATE(NOW()) - INTERVAL ? DAY AND DATE(NOW()) " +
		"AND b.source = ? AND b.macrobenchmark_id = macrotype.macrobenchmark_id AND macrotype.$(MBTYPE)_no = qps.$(MBTYPE)_no"

	query = strings.ReplaceAll(query, "$(MBTYPE)", upperMacroType)

	rows, err := client.Select(query, lastDays, source)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var res MacroBenchmarkDetails
		err = rows.Scan(&res.ID, &res.GitRef, &res.Source, &res.CreatedAt, &res.Result.TPS, &res.Result.Latency,
			&res.Result.Errors, &res.Result.Reconnects, &res.Result.Time, &res.Result.Threads, &res.Result.QPS.ID,
			&res.Result.QPS.Total, &res.Result.QPS.Reads, &res.Result.QPS.Writes, &res.Result.QPS.Other)
		if err != nil {
			return nil, err
		}
		macrodetails = append(macrodetails, res)
	}
	return macrodetails, nil
}

func GetResultsForGitRef(macroType Type, ref string, client *mysql.Client) (macrodetails MacroBenchmarkDetailsArray, err error) {
	if macroType != OLTP && macroType != TPCC {
		return nil, errors.New(IncorrectMacroBenchmarkType)
	}
	upperMacroType := macroType.ToUpper().String()
	query := "SELECT b.test_no, b.commit, b.source, b.DateTime, " +
		"macrotype.tps, macrotype.latency, macrotype.errors, macrotype.reconnects, macrotype.time, macrotype.threads, " +
		"qps.qps_no, qps.total_qps, qps.reads_qps, qps.writes_qps, qps.other_qps " +
		"FROM benchmark AS b, $(MBTYPE) AS macrotype, qps AS qps " +
		"WHERE b.commit = ? AND b.test_no = macrotype.test_no AND macrotype.$(MBTYPE)_no = qps.$(MBTYPE)_no"

	query = strings.ReplaceAll(query, "$(MBTYPE)", upperMacroType)

	rows, err := client.Select(query, ref)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var res MacroBenchmarkDetails
		err = rows.Scan(&res.ID, &res.GitRef, &res.Source, &res.CreatedAt, &res.Result.TPS, &res.Result.Latency,
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
func (mbr *MacroBenchmarkResult) insertToMySQL(benchmarkType MacroBenchmarkType, macrobenchmarkID int, client *mysql.Client) error {
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
	queryQPS := fmt.Sprintf("INSERT INTO qps(%s, total_qps, reads_qps, writes_qps, other_qps) VALUES(?, ?, ?, ?, ?)", benchmarkType.ToUpper().String() + "_no")
	_, err = client.Insert(queryQPS, resultID, mbr.QPS.Total, mbr.QPS.Reads, mbr.QPS.Writes, mbr.QPS.Other)
	if err != nil {
		return err
	}
	return nil
}