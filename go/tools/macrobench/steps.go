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

import "strings"

type step struct {
	Name         string
	SysbenchName string
}

const (
	stepPrepare = "prepare"
	stepRun     = "run"
)

var (
	steps = []step{
		{Name: stepPrepare, SysbenchName: stepPrepare},
		{Name: stepRun, SysbenchName: stepRun},
	}
)

func skipSteps(steps []step, skip string) (newSteps []step) {
	skips := strings.Split(skip, ",")
	newSteps = []step{}
	for _, step := range steps {
		add := true
		for _, skipStep := range skips {
			if step.Name == skipStep {
				add = false
				break
			}
		}
		if add {
			newSteps = append(newSteps, step)
		}
	}
	return newSteps
}
