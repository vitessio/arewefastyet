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
	"github.com/vitessio/arewefastyet/go/storage"
)

type StatisticalComparison struct {
	RightSHA   string
	LeftSHA    string
	Planner    PlannerVersion
	MacroTypes []string
}

func (sc StatisticalComparison) Compare(client storage.SQLClient) (map[string]StatisticalCompareResults, error) {
	results := make(map[string]StatisticalCompareResults, len(sc.MacroTypes))
	for _, macroType := range sc.MacroTypes {
		leftResult, err := GetBenchmarkResults(client, macroType, sc.LeftSHA, sc.Planner)
		if err != nil {
			return nil, err
		}

		rightResult, err := GetBenchmarkResults(client, macroType, sc.RightSHA, sc.Planner)
		if err != nil {
			return nil, err
		}

		leftResultsAsSlice := leftResult.asSlice()
		rightResultsAsSlice := rightResult.asSlice()

		scr := performAnalysis(leftResultsAsSlice, rightResultsAsSlice)
		results[macroType] = scr
	}
	return results, nil
}
