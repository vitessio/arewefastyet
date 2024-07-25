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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/vitessio/arewefastyet/go/storage"
	"github.com/vitessio/arewefastyet/go/storage/mysql"
	"vitess.io/vitess/go/vt/sqlparser"
)

type VTGateQueryPlanValue struct {
	QueryType    string      `json:"query_type"`
	Original     string      `json:"original"`
	Instructions interface{} `json:"instructions"`

	ExecCount    int `json:"exec_count"`    // Count of times this plan was executed
	ExecTime     int `json:"exec_time"`     // Average execution time per query
	ShardQueries int `json:"shard_queries"` // Total number of shard queries
	RowsReturned int `json:"rows_returned"` // Total number of rows
	RowsAffected int `json:"rows_affected"` // Total number of rows
	Errors       int `json:"errors"`        // Total number of errors

	TablesUsed interface{} `json:"tables_used"`
}

type VTGateQueryPlan struct {
	Key   string               `json:"key"`
	Value VTGateQueryPlanValue `json:"value"`
}

type VTGateQueryPlanComparer struct {
	Left             *VTGateQueryPlan `json:"left"`
	Right            *VTGateQueryPlan `json:"right"`
	SamePlan         bool             `json:"same_plan"`
	Key              string           `json:"key"`
	ExecCountDiff    int              `json:"exec_count_diff"`
	ExecTimeDiff     int              `json:"exec_time_diff"`
	RowsReturnedDiff int              `json:"rows_returned_diff"`
	ErrorsDiff       int              `json:"errors_diff"`
}

type VTGateQueryPlanMap map[string]VTGateQueryPlanValue

func CompareVTGateQueryPlans(left, right []VTGateQueryPlan) []VTGateQueryPlanComparer {
	res := []VTGateQueryPlanComparer{}
	for i, plan := range left {
		newCompare := VTGateQueryPlanComparer{
			Key:  plan.Key,
			Left: &left[i],
		}
		for j, rightPlan := range right {
			if rightPlan.Key == plan.Key {
				switch instructionRight := rightPlan.Value.Instructions.(type) {
				case string:
					switch instructionLeft := plan.Value.Instructions.(type) {
					case string:
						newCompare.SamePlan = instructionLeft == instructionRight
					}
				}
				newCompare.Right = &right[j]
				if plan.Value.ExecCount != 0 {
					newCompare.ExecCountDiff = int(float64(rightPlan.Value.ExecCount-plan.Value.ExecCount) / float64(plan.Value.ExecCount) * 100)
				}
				if plan.Value.ExecTime != 0 {
					newCompare.ExecTimeDiff = int(float64(rightPlan.Value.ExecTime-plan.Value.ExecTime) / float64(plan.Value.ExecTime) * 100)
				}
				if plan.Value.RowsReturned != 0 {
					newCompare.RowsReturnedDiff = int(float64(rightPlan.Value.RowsReturned-plan.Value.RowsReturned) / float64(plan.Value.RowsReturned) * 100)
				}
				if plan.Value.Errors != 0 {
					newCompare.ErrorsDiff = int(float64(rightPlan.Value.Errors-plan.Value.Errors) / float64(plan.Value.Errors) * 100)
				}
				break
			}
		}
		res = append(res, newCompare)
	}
	for _, plan := range right {
		found := false
		for _, resPlan := range res {
			if plan.Key == resPlan.Key {
				found = true
				break
			}
		}
		if !found {
			res = append(res, VTGateQueryPlanComparer{
				Right: &plan,
				Key:   plan.Key,
			})
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].ExecTimeDiff > res[j].ExecTimeDiff
	})
	return res
}

func normalizeVTGateQueryPlan(plan *VTGateQueryPlanValue) {
	if plan.ExecCount == 0 {
		return
	}
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

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]VTGateQueryPlanValue
	err = json.NewDecoder(bytes.NewReader(respBytes)).Decode(&response)
	if err != nil {
		return getOldVTGateQueryPlans(respBytes)
	}
	res := make(map[string]VTGateQueryPlanValue)
	for key, plan := range response {
		jsonPlan, err := json.MarshalIndent(plan.Instructions, "", "\t")
		if err != nil {
			return nil, err
		}
		plan.Instructions = string(jsonPlan)
		res[key] = plan
	}
	return res, nil
}

func getOldVTGateQueryPlans(respBytes []byte) (VTGateQueryPlanMap, error) {
	var response []VTGateQueryPlan
	err := json.NewDecoder(bytes.NewReader(respBytes)).Decode(&response)
	if err != nil {
		return nil, err
	}
	planMap := VTGateQueryPlanMap{}
	for _, plan := range response {
		jsonPlan, err := json.MarshalIndent(plan.Value.Instructions, "", "\t")
		if err != nil {
			return nil, err
		}
		plan.Value.Instructions = string(jsonPlan)
		planMap[plan.Key] = plan.Value
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

func GetVTGateSelectQueryPlansWithFilter(gitRef string, macroType Type, planner PlannerVersion, client storage.SQLClient) ([]VTGateQueryPlan, error) {
	query := "select " +
		"qp.`key` as `key`, " +
		"qp.plan as plan, " +
		"convert(avg(qp.exec_time), signed) as time, " +
		"convert(avg(qp.exec_count), signed) as count, " +
		"convert(avg(qp.rows), signed) as r, " +
		"convert(avg(qp.errors), signed) as errors " +
		"from " +
		"query_plans qp, macrobenchmark ma, execution ex " +
		"where " +
		"qp.macrobenchmark_id = ma.macrobenchmark_id " +
		"and ex.uuid = qp.exec_uuid " +
		"and ex.uuid = ma.exec_uuid " +
		"and ex.type = ? " +
		"and ma.commit = ? " +
		"and ma.vtgate_planner_version = ? " +
		"group by " +
		"qp.`key`, qp.plan " +
		"order by qp.`key` " +
		"limit 1500;"

	result, err := client.Read(query, macroType.String(), gitRef, string(planner))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	parser, err := sqlparser.New(sqlparser.Options{})
	if err != nil {
		return nil, err
	}
	res := []VTGateQueryPlan{}
	for result.Next() {
		var plan VTGateQueryPlan
		err = result.Scan(&plan.Key, &plan.Value.Instructions, &plan.Value.ExecTime, &plan.Value.ExecCount, &plan.Value.RowsReturned, &plan.Value.Errors)
		if err != nil {
			return nil, err
		}
		switch p := plan.Value.Instructions.(type) {
		case []byte:
			plan.Value.Instructions = string(p)
		}

		// Remove all comments from the query
		// This prevents the query from not match across two versions
		// of Vitess where we changed query hints and added comments
		stmt, err := parser.Parse(plan.Key)
		if err != nil {
			return nil, err
		}
		cmmtd, isCmmtd := stmt.(sqlparser.Commented)
		if isCmmtd && cmmtd != nil {
			cmmtd.SetComments(nil)
			plan.Key = sqlparser.String(stmt)
		}

		res = append(res, plan)
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
		_, err := client.Write(query, execUUID, macrobenchmarkID, key, fmt.Sprintf("%v", value.Instructions), value.ExecCount, value.ExecTime, value.RowsReturned, value.Errors)
		if err != nil {
			return err
		}
	}
	return nil
}
