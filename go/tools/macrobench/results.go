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
	"github.com/vitessio/arewefastyet/go/mysql"
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

	MacroBenchmarkDetailsArray []MacroBenchmarkDetails
)

// GetResultsForLastDays returns a slice MacroBenchmarkDetails based on a given macro benchmark type.
// The type can either be OLTP or TPCC. Using that type, the function will generate a query using
// the *mysql.Client. The query will select only results that were added between now and lastDays.
func GetResultsForLastDays(macroType MacroBenchmarkType, source string, lastDays int, client *mysql.Client) (macrodetails MacroBenchmarkDetailsArray, err error) {
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
