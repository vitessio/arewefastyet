/*
Copyright 2024 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package macrobench

import (
	"errors"
	"strings"

	"github.com/vitessio/arewefastyet/go/storage"
	"github.com/vitessio/arewefastyet/go/storage/mysql"
)

// getResultsForGitRefAndPlanner returns a slice of details based on the given git ref
// and macro benchmark Type.
func getResultsForGitRefAndPlanner(macroType string, ref string, planner PlannerVersion, client storage.SQLClient) (macrodetails detailsArray, err error) {
	upperMacroType := strings.ToUpper(macroType)
	query := "SELECT " +
		"info.macrobenchmark_id, e.git_ref, e.source, e.finished_at, IFNULL(info.exec_uuid, ''), " +
		"results.tps, results.latency, results.errors, results.reconnects, results.time, results.threads, " +
		"results.total_qps, results.reads_qps, results.writes_qps, results.other_qps " +
		"FROM execution AS e, macrobenchmark AS info, macrobenchmark_results AS results " +
		"WHERE e.uuid = info.exec_uuid " +
		"AND e.status = \"finished\" " +
		"AND e.git_ref = ? " +
		"AND info.vtgate_planner_version = ? " +
		"AND info.macrobenchmark_id = results.macrobenchmark_id " +
		"AND info.type = ?"

	result, err := client.Read(query, ref, planner, upperMacroType)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		var res details
		err = result.Scan(
			&res.ID,
			&res.GitRef,
			&res.Source,
			&res.CreatedAt,
			&res.ExecUUID,
			&res.Result.TPS,
			&res.Result.Latency,
			&res.Result.Errors,
			&res.Result.Reconnects,
			&res.Result.Time,
			&res.Result.Threads,
			&res.Result.QPS.Total,
			&res.Result.QPS.Reads,
			&res.Result.QPS.Writes,
			&res.Result.QPS.Other,
		)
		if err != nil {
			return nil, err
		}
		macrodetails = append(macrodetails, res)
	}
	return macrodetails, nil
}

func getResultsLastXDays(macroType string, source string, planner PlannerVersion, lastDays int, client storage.SQLClient) (macrodetails detailsArray, err error) {
	macrodetails = []details{}
	upperMacroType := strings.ToUpper(macroType)
	query := "SELECT " +
		"info.macrobenchmark_id, e.git_ref, e.source, e.finished_at, IFNULL(e.uuid, ''), " +
		"results.tps, results.latency, results.errors, results.reconnects, results.time, results.threads, " +
		"results.total_qps, results.reads_qps, results.writes_qps, results.other_qps " +
		"FROM execution AS e, macrobenchmark AS info, macrobenchmark_results AS results " +
		"WHERE e.uuid = info.exec_uuid AND e.status = \"finished\" " +
		"AND e.finished_at BETWEEN DATE(NOW()) - INTERVAL ? DAY " +
		"AND DATE(NOW() + INTERVAL 1 DAY) " +
		"AND e.source = ? " +
		"AND info.vtgate_planner_version = ? " +
		"AND info.macrobenchmark_id = results.macrobenchmark_id " +
		"AND info.type = ? " +
		"ORDER BY e.finished_at "

	result, err := client.Read(query, lastDays, source, planner, upperMacroType)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		var res details
		err = result.Scan(
			&res.ID,
			&res.GitRef,
			&res.Source,
			&res.CreatedAt,
			&res.ExecUUID,
			&res.Result.TPS,
			&res.Result.Latency,
			&res.Result.Errors,
			&res.Result.Reconnects,
			&res.Result.Time,
			&res.Result.Threads,
			&res.Result.QPS.Total,
			&res.Result.QPS.Reads,
			&res.Result.QPS.Writes,
			&res.Result.QPS.Other,
		)
		if err != nil {
			return nil, err
		}
		macrodetails = append(macrodetails, res)
	}
	return macrodetails, nil
}

func getSummaryLastXDays(macroType string, source string, planner PlannerVersion, lastDays int, client storage.SQLClient) (results detailsArray, err error) {
	upperMacroType := strings.ToUpper(macroType)
	query := "SELECT " +
		"info.macrobenchmark_id, e.git_ref, results.total_qps, IFNULL(e.uuid, '') " +
		"FROM execution AS e, macrobenchmark AS info, macrobenchmark_results AS results " +
		"WHERE e.uuid = info.exec_uuid " +
		"AND e.status = \"finished\" " +
		"AND e.finished_at BETWEEN DATE(NOW()) - INTERVAL ? DAY " +
		"AND DATE(NOW() + INTERVAL 1 DAY) " +
		"AND e.source = ? " +
		"AND info.vtgate_planner_version = ? " +
		"AND info.macrobenchmark_id = results.macrobenchmark_id " +
		"AND info.type = ? " +
		"ORDER BY e.finished_at "

	result, err := client.Read(query, lastDays, source, planner, upperMacroType)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		var res details
		err = result.Scan(&res.ID, &res.GitRef, &res.Result.QPS.Total, &res.ExecUUID)
		if err != nil {
			return nil, err
		}
		results = append(results, res)
	}
	return
}

// insertToMySQL inserts the given MacroBenchmarkResult to MySQL using a *mysql.Client.
// The MacroBenchmarkResults gets added in one of macrobenchmark's children tables.
// Depending on the MacroBenchmarkType, the insert will be routed to a specific children table.
// The children table qps is also inserted.
func (mbr *result) insertToMySQL(macrobenchmarkID int, client storage.SQLClient) error {
	if client == nil {
		return errors.New(mysql.ErrorClientConnectionNotInitialized)
	}

	// insert result
	queryResult := "INSERT INTO macrobenchmark_results(macrobenchmark_id, queries, tps, latency, errors, reconnects, time, threads, total_qps, reads_qps, writes_qps, other_qps) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := client.Write(queryResult, macrobenchmarkID, mbr.Queries, mbr.TPS, mbr.Latency, mbr.Errors, mbr.Reconnects, mbr.Time, mbr.Threads, mbr.QPS.Total, mbr.QPS.Reads, mbr.QPS.Writes, mbr.QPS.Other)
	if err != nil {
		return err
	}
	return nil
}
