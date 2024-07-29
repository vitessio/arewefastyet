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
	"database/sql"
	"errors"
	"strings"

	"github.com/vitessio/arewefastyet/go/exec/metrics"
	"github.com/vitessio/arewefastyet/go/storage"
	"github.com/vitessio/arewefastyet/go/storage/mysql"
)

// getExecutionGroupResults the results of an execution group
func getExecutionGroupResults(workload string, ref string, planner PlannerVersion, client storage.SQLClient) (executionGroupResults, error) {
	query := `
        SELECT 
            IFNULL(e.uuid, '') AS exec_uuid, 
            results.tps, 
            results.latency, 
            results.errors, 
            results.reconnects, 
            results.time, 
            results.threads, 
            results.total_qps, 
            results.reads_qps, 
            results.writes_qps, 
            results.other_qps, 
            m.name AS metric_name, 
            m.value AS metric_value
        FROM 
            execution AS e
        JOIN 
            macrobenchmark AS info ON e.uuid = info.exec_uuid
        JOIN 
            macrobenchmark_results AS results ON info.macrobenchmark_id = results.macrobenchmark_id
        LEFT JOIN 
            metrics AS m ON e.uuid = m.exec_uuid
        WHERE 
            e.status = 'finished'
            AND e.git_ref = ? 
            AND info.vtgate_planner_version = ? 
            AND info.workload = ?
        ORDER BY 
            e.uuid, m.name
    `

	rows, err := client.Read(query, ref, planner, strings.ToUpper(workload))
	if err != nil {
		return executionGroupResults{}, err
	}
	defer rows.Close()

	results := executionGroupResults{GitRef: ref}
	var execRes *executionResults
	var currentExecUUID string

	for rows.Next() {
		var (
			execUUID    string
			sr          sysbenchResult
			metricName  sql.NullString
			metricValue sql.NullFloat64
		)

		err := rows.Scan(
			&execUUID, &sr.TPS, &sr.Latency, &sr.Errors, &sr.Reconnects, &sr.Time, &sr.Threads, &sr.QPS.Total,
			&sr.QPS.Reads, &sr.QPS.Writes, &sr.QPS.Other, &metricName, &metricValue,
		)
		if err != nil {
			return executionGroupResults{}, err
		}

		// If execUUID is different it means we are looking at another set of results
		// we then add the current execRes to our results, and start again with a new executionResults
		if currentExecUUID != execUUID {
			if execRes != nil {
				results.Results = append(results.Results, execRes.Result)
				results.Metrics = append(results.Metrics, execRes.Metrics)
			}
			execRes = &executionResults{
				Result: sr,
				Metrics: metrics.ExecutionMetrics{
					ComponentsCPUTime:            make(map[string]float64),
					ComponentsMemStatsAllocBytes: make(map[string]float64),
				},
			}
			currentExecUUID = execUUID
		}

		// For each execution we will have multiple rows since we are doing a LEFT JOIN on metrics
		// here we just all the metrics value to the executionResults, later when we are done consuming
		// all the metrics for our current execUUID we will create a new executionResults
		if metricName.Valid {
			switch {
			case metricName.String == "TotalComponentsCPUTime":
				execRes.Metrics.TotalComponentsCPUTime = metricValue.Float64
			case metricName.String == "TotalComponentsMemStatsAllocBytes":
				execRes.Metrics.TotalComponentsMemStatsAllocBytes = metricValue.Float64
			case strings.HasPrefix(metricName.String, "ComponentsCPUTime."):
				execRes.Metrics.ComponentsCPUTime[strings.Split(metricName.String, ".")[1]] = metricValue.Float64
			case strings.HasPrefix(metricName.String, "ComponentsMemStatsAllocBytes."):
				execRes.Metrics.ComponentsMemStatsAllocBytes[strings.Split(metricName.String, ".")[1]] = metricValue.Float64
			}
		}
	}

	// This is used to add the results of our very last execUUID since we will exit the
	// previous loop before adding the results.
	if execRes != nil {
		results.Results = append(results.Results, execRes.Result)
		results.Metrics = append(results.Metrics, execRes.Metrics)
	}
	return results, nil
}

func getExecutionGroupResultsFromLast30Days(workload string, planner PlannerVersion, client storage.SQLClient) ([]executionGroupResults, error) {
	query := `
        SELECT 
            IFNULL(e.uuid, '') AS exec_uuid, 
            e.git_ref,
            results.tps, 
            results.latency, 
            results.errors, 
            results.reconnects, 
            results.time, 
            results.threads, 
            results.total_qps, 
            results.reads_qps, 
            results.writes_qps, 
            results.other_qps, 
            m.name AS metric_name, 
            m.value AS metric_value
        FROM 
            execution AS e
        JOIN 
            macrobenchmark AS info ON e.uuid = info.exec_uuid
        JOIN 
            macrobenchmark_results AS results ON info.macrobenchmark_id = results.macrobenchmark_id
        LEFT JOIN 
            metrics AS m ON e.uuid = m.exec_uuid
        WHERE 
            e.finished_at BETWEEN DATE(NOW()) - INTERVAL 30 DAY AND DATE(NOW() + INTERVAL 1 DAY)
            AND e.source = 'cron'
            AND e.status = 'finished'
            AND info.vtgate_planner_version = ? 
            AND info.workload = ?
        ORDER BY 
            e.finished_at ASC, e.uuid, m.name
    `

	rows, err := client.Read(query, planner, strings.ToUpper(workload))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allResults []executionGroupResults
	var results executionGroupResults
	var execRes *executionResults
	var currentExecUUID string
	var currentGitRef string

	for rows.Next() {
		var (
			execUUID    string
			gitRef      string
			sr          sysbenchResult
			metricName  sql.NullString
			metricValue sql.NullFloat64
		)

		err := rows.Scan(
			&execUUID, &gitRef, &sr.TPS, &sr.Latency, &sr.Errors, &sr.Reconnects, &sr.Time, &sr.Threads, &sr.QPS.Total,
			&sr.QPS.Reads, &sr.QPS.Writes, &sr.QPS.Other, &metricName, &metricValue,
		)
		if err != nil {
			return nil, err
		}

		// If gitRef is different it means we are looking at another group of results
		if currentGitRef != gitRef {
			if execRes != nil {
				results.Results = append(results.Results, execRes.Result)
				results.Metrics = append(results.Metrics, execRes.Metrics)
				execRes = nil
			}
			if len(results.Results) > 0 {
				allResults = append(allResults, results)
			}
			results = executionGroupResults{GitRef: gitRef}
			currentGitRef = gitRef
		}

		// If execUUID is different it means we are looking at another set of results
		if currentExecUUID != execUUID {
			if execRes != nil {
				results.Results = append(results.Results, execRes.Result)
				results.Metrics = append(results.Metrics, execRes.Metrics)
			}
			execRes = &executionResults{
				Result: sr,
				Metrics: metrics.ExecutionMetrics{
					ComponentsCPUTime:            make(map[string]float64),
					ComponentsMemStatsAllocBytes: make(map[string]float64),
				},
			}
			currentExecUUID = execUUID
		}

		// Add the metrics values to the current executionResults
		if metricName.Valid {
			switch {
			case metricName.String == "TotalComponentsCPUTime":
				execRes.Metrics.TotalComponentsCPUTime = metricValue.Float64
			case metricName.String == "TotalComponentsMemStatsAllocBytes":
				execRes.Metrics.TotalComponentsMemStatsAllocBytes = metricValue.Float64
			case strings.HasPrefix(metricName.String, "ComponentsCPUTime."):
				execRes.Metrics.ComponentsCPUTime[strings.Split(metricName.String, ".")[1]] = metricValue.Float64
			case strings.HasPrefix(metricName.String, "ComponentsMemStatsAllocBytes."):
				execRes.Metrics.ComponentsMemStatsAllocBytes[strings.Split(metricName.String, ".")[1]] = metricValue.Float64
			}
		}
	}

	// Add the last set of results
	if execRes != nil {
		results.Results = append(results.Results, execRes.Result)
		results.Metrics = append(results.Metrics, execRes.Metrics)
	}
	if len(results.Results) > 0 {
		allResults = append(allResults, results)
	}

	return allResults, nil
}

func getSummaryLast30Days(workload string, planner PlannerVersion, client storage.SQLClient) ([]executionGroupResults, error) {
	query := `
        SELECT 
            e.git_ref, 
            results.total_qps 
        FROM 
            execution AS e
        JOIN 
            macrobenchmark AS info ON e.uuid = info.exec_uuid
        JOIN 
            macrobenchmark_results AS results ON info.macrobenchmark_id = results.macrobenchmark_id
        WHERE 
            e.finished_at BETWEEN DATE(NOW()) - INTERVAL 30 DAY AND DATE(NOW() + INTERVAL 1 DAY) 
            AND e.status = "finished" 
            AND e.source = "cron" 
            AND info.vtgate_planner_version = ? 
            AND info.workload = ? 
        ORDER BY 
            e.finished_at ASC
    `

	rows, err := client.Read(query, planner, strings.ToUpper(workload))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allResults []executionGroupResults
	var results *executionGroupResults
	var currentGitRef string

	for rows.Next() {
		var (
			gitRef string
			sr     sysbenchResult
		)

		if err := rows.Scan(&gitRef, &sr.QPS.Total); err != nil {
			return nil, err
		}

		if currentGitRef != gitRef {
			if results != nil {
				allResults = append(allResults, *results)
			}
			results = &executionGroupResults{
				GitRef:  gitRef,
				Results: []sysbenchResult{sr},
			}
			currentGitRef = gitRef
		} else if results != nil {
			results.Results = append(results.Results, sr)
		}
	}

	if results != nil {
		allResults = append(allResults, *results)
	}

	return allResults, nil
}

// insertToMySQL inserts the given sysbenchResult to MySQL.
func (mbr *sysbenchResult) insertToMySQL(macrobenchmarkID int, client storage.SQLClient) error {
	if client == nil {
		return errors.New(mysql.ErrorClientConnectionNotInitialized)
	}

	// insert sysbenchResult
	queryResult := "INSERT INTO macrobenchmark_results(macrobenchmark_id, queries, tps, latency, errors, reconnects, time, threads, total_qps, reads_qps, writes_qps, other_qps) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := client.Write(queryResult, macrobenchmarkID, mbr.Queries, mbr.TPS, mbr.Latency, mbr.Errors, mbr.Reconnects, mbr.Time, mbr.Threads, mbr.QPS.Total, mbr.QPS.Reads, mbr.QPS.Writes, mbr.QPS.Other)
	if err != nil {
		return err
	}
	return nil
}
