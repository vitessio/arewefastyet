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
	qt "github.com/frankban/quicktest"
	"testing"
)

func Test_skipSteps(t *testing.T) {
	prepare := step{Name: stepPrepare, SysbenchName: stepPrepare}
	warmup := step{Name: stepWarmUp, SysbenchName: stepRun}
	run := step{Name: stepRun, SysbenchName: stepRun}

	type args struct {
		steps []step
		skip  []string
	}
	tests := []struct {
		name         string
		args         args
		wantNewSteps []step
	}{
		{name: "No skip step", args: args{steps: []step{prepare, warmup, run}, skip: []string{}}, wantNewSteps: []step{prepare, warmup, run}},
		{name: "Skip prepare", args: args{steps: []step{prepare, warmup, run}, skip: []string{stepPrepare}}, wantNewSteps: []step{warmup, run}},
		{name: "Skip warm up", args: args{steps: []step{prepare, warmup, run}, skip: []string{stepWarmUp}}, wantNewSteps: []step{prepare, run}},
		{name: "Skip run", args: args{steps: []step{prepare, warmup, run}, skip: []string{stepRun}}, wantNewSteps: []step{prepare, warmup}},
		{name: "Skip prepare and run", args: args{steps: []step{prepare, warmup, run}, skip: []string{stepPrepare, stepRun}}, wantNewSteps: []step{warmup}},
		{name: "Skip all", args: args{steps: []step{prepare, warmup, run}, skip: []string{stepPrepare, stepWarmUp, stepRun}}, wantNewSteps: []step{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			gotNewSteps := skipSteps(tt.args.steps, tt.args.skip)
			c.Assert(gotNewSteps, qt.DeepEquals, tt.wantNewSteps)
		})
	}
}
