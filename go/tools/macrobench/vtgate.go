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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vitessio/arewefastyet/go/storage"
	"github.com/vitessio/arewefastyet/go/storage/mysql"
	"net/http"
	"strings"
)

type VTGateQueryPlanValue struct {
	QueryType    string
	Original     string
	Instructions interface{}
	ExecCount    int
	ExecTime     int
	ShardQueries int
	RowsReturned int
	Errors       int
}

type VTGateQueryPlan struct {
	Key   string
	Value VTGateQueryPlanValue
}

type VTGateQueryPlanMap map[string]VTGateQueryPlanValue

func normalizeVTGateQueryPlan(plan *VTGateQueryPlanValue) {
	plan.ExecTime = plan.ExecTime / plan.ExecCount
}

func mergeVTGateQueryPlans(plansLeft, plansRight VTGateQueryPlanMap) VTGateQueryPlanMap {
	res := VTGateQueryPlanMap{}

	for keyLeft, planLeft := range plansLeft {
		if planRight, found := plansRight[keyLeft]; found {
			newPlanValue := VTGateQueryPlanValue{
				QueryType:    planLeft.QueryType,
				Original:     planLeft.Original,
				Instructions: planLeft.Instructions,
				ExecCount:    planLeft.ExecCount + planRight.ExecCount,
				ExecTime:     planLeft.ExecTime + planRight.ExecTime,
				ShardQueries: planLeft.ShardQueries + planRight.ShardQueries,
				RowsReturned: planLeft.RowsReturned + planRight.RowsReturned,
				Errors:       planLeft.Errors + planRight.Errors,
			}
			res[keyLeft] = newPlanValue
		} else {
			res[keyLeft] = planLeft
		}
	}
	for keyRight, planRight := range plansRight {
		if _, found := res[keyRight]; !found {
			res[keyRight] = planRight
		}
	}
	return res
}

func getVTGateQueryPlans(port string) (VTGateQueryPlanMap, error) {
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/debug/query_plans", port))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response []VTGateQueryPlan
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	planMap := VTGateQueryPlanMap{}
	for _, plan := range response {
		// keeping only select statements
		if strings.HasPrefix(plan.Key, "select") {
			planMap[plan.Key] = plan.Value
		}
	}
	return planMap, nil
}

func getVTGatesQueryPlans(ports []string) (VTGateQueryPlanMap, error) {
	res := VTGateQueryPlanMap{}

	for _, port := range ports {
		plans, err := getVTGateQueryPlans(port)
		if err != nil {
			return nil, err
		}
		res = mergeVTGateQueryPlans(res, plans)
	}
	return res, nil
}

func GetVTGateSelectQueryPlansWithFilter(gitRef string, macroType Type, planner PlannerVersion, client storage.SQLClient) (VTGateQueryPlanMap, error) {
	if macroType != OLTP && macroType != TPCC {
		return nil, errors.New(IncorrectMacroBenchmarkType)
	}
	query := "select qp.`key` as `key`, qp.plan as plan, avg(qp.exec_time) as time, avg(qp.exec_count) as count, avg(qp.rows) as rows, avg(qp.errors) as errors " +
		"from query_plans qp, macrobenchmark ma, execution ex " +
		"where qp.macrobenchmark_id = ma.macrobenchmark_id " +
		"and ex.uuid = qp.exec_uuid " +
		"and qp.`key` like \"select%\" " +
		"and ex.type = ?" +
		"and ma.commit = ? " +
		"and ma.vtgate_planner_version = ? " +
		"group by qp.`key`, qp.plan order by time desc limit 100;"

	result, err := client.Select(query, macroType.String(), gitRef, string(planner))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	res := VTGateQueryPlanMap{}
	for result.Next() {
		var plan VTGateQueryPlan
		err = result.Scan(&plan.Key, &plan.Value.Instructions, &plan.Value.ExecTime, &plan.Value.ExecCount, &plan.Value.RowsReturned, &plan.Value.Errors)
		if err != nil {
			return nil, err
		}
		res[plan.Key] = plan.Value
	}
	return res, nil
}

func insertVTGateQueryMapToMySQL(client storage.SQLClient, execUUID string, result VTGateQueryPlanMap, macrobenchmarkID int) error {
	if client == nil {
		return errors.New(mysql.ErrorClientConnectionNotInitialized)
	}

	query := "INSERT INTO query_plans(`exec_uuid`, `macrobenchmark_id`, `key`, `plan`, `exec_count`, `exec_time`, `rows`, `errors`) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"
	for key, value := range result {
		normalizeVTGateQueryPlan(&value)
		_, err := client.Insert(query, execUUID, macrobenchmarkID, key, fmt.Sprintf("%v", value.Instructions), value.ExecCount, value.ExecTime, value.RowsReturned, value.Errors)
		if err != nil {
			return err
		}
	}
	return nil
}
