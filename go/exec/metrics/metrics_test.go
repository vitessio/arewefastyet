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

package metrics

import (
	qt "github.com/frankban/quicktest"
	"testing"
)

func emptyExecMetrics() ExecutionMetrics {
	return ExecutionMetrics{
		ComponentsCPUTime: map[string]float64{},
		ComponentsMemStatsAllocBytes: map[string]float64{},
	}
}

func simpleExecMetrics(vtgate, vttablet float64) ExecutionMetrics {
	return ExecutionMetrics{
		TotalComponentsCPUTime: vtgate + vttablet,
		ComponentsCPUTime: map[string]float64{
			"vtgate":   vtgate,
			"vttablet": vttablet,
		},
		TotalComponentsMemStatsAllocBytes: vtgate + vttablet,
		ComponentsMemStatsAllocBytes: map[string]float64{
			"vtgate":   vtgate,
			"vttablet": vttablet,
		},
	}
}

func complexExecMetrics(vtgate, vttablet, all float64) ExecutionMetrics {
	return ExecutionMetrics{
		TotalComponentsCPUTime: all,
		ComponentsCPUTime: map[string]float64{
			"vtgate":   vtgate,
			"vttablet": vttablet,
		},
		TotalComponentsMemStatsAllocBytes: all,
		ComponentsMemStatsAllocBytes: map[string]float64{
			"vtgate":   vtgate,
			"vttablet": vttablet,
		},
	}
}

func TestCompareTwo(t *testing.T) {
	tests := []struct {
		name  string
		left  ExecutionMetrics
		right ExecutionMetrics
		want  ExecutionMetrics
	}{
		{name: "Nil metrics", want: emptyExecMetrics()},
		{name: "Empty metrics", left: emptyExecMetrics(), right: emptyExecMetrics(), want: emptyExecMetrics()},
		{name: "Same metrics", left: simpleExecMetrics(9816.56, 15789.36), right: simpleExecMetrics(9816.56, 15789.36), want: simpleExecMetrics(0, 0)},
		{name: "50% improvement (less CPU time)", left: simpleExecMetrics(20, 20), right: simpleExecMetrics(10, 10), want: complexExecMetrics(50, 50, 50)},
		{name: "100% improvement (zero CPU time)", left: simpleExecMetrics(20, 20), right: simpleExecMetrics(0, 0), want: complexExecMetrics(100, 100, 100)},
		{name: "50% regression (more CPU time)", left: simpleExecMetrics(20, 20), right: simpleExecMetrics(30, 30), want: complexExecMetrics(-50, -50, -50)},
		{name: "100% regression (twice as much CPU time)", left: simpleExecMetrics(20, 20), right: simpleExecMetrics(40, 40), want: complexExecMetrics(-100, -100, -100)},
		{name: "50% vtgate improvement (less CPU time)", left: simpleExecMetrics(20, 20), right: simpleExecMetrics(10, 20), want: complexExecMetrics(50, 0, 25)},
		{name: "50% vttablet improvement (less CPU time)", left: simpleExecMetrics(20, 20), right: simpleExecMetrics(20, 10), want: complexExecMetrics(0, 50, 25)},
		{name: "50% vtgate regression (more CPU time)", left: simpleExecMetrics(20, 20), right: simpleExecMetrics(30, 20), want: complexExecMetrics(-50, 0, -25)},
		{name: "50% vttablet regression (more CPU time)", left: simpleExecMetrics(20, 20), right: simpleExecMetrics(20, 30), want: complexExecMetrics(0, -50, -25)},
		{name: "Left with zero values", left: simpleExecMetrics(0, 0), right: simpleExecMetrics(20, 20), want: complexExecMetrics(-100, -100, -100)},
		{name: "Right with zero values", left: simpleExecMetrics(10, 10), right: simpleExecMetrics(0, 0), want: complexExecMetrics(100, 100, 100)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			got := CompareTwo(tt.left, tt.right)
			c.Assert(got, qt.DeepEquals, tt.want)
		})
	}
}

func Test_compareSafe(t *testing.T) {
	tests := []struct {
		name       string
		left       float64
		right      float64
		wantResult float64
	}{
		{name: "50% performance increase", left: 100, right: 50, wantResult: 50},
		{name: "100% performance decrease", left: 50, right: 100, wantResult: -100},
		{name: "200% performance decrease", left: 50, right: 150, wantResult: -200},
		{name: "left zero", left: 0, right: 100, wantResult: -100},
		{name: "right zero", left: 100, right: 0, wantResult: 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			got := compareSafe(tt.left, tt.right)
			c.Assert(got, qt.Equals, tt.wantResult)
		})
	}
}

func Test_compareSafeComponentMap(t *testing.T) {
	tests := []struct {
		name  string
		left  map[string]float64
		right map[string]float64
		want  map[string]float64
	}{
		{name: "Same component", left: map[string]float64{"vtgate": 10}, right: map[string]float64{"vtgate": 10}, want: map[string]float64{"vtgate": 0}},
		{name: "Same component with performance increase", left: map[string]float64{"vtgate": 10}, right: map[string]float64{"vtgate": 5}, want: map[string]float64{"vtgate": 50}},
		{name: "Compare multiple components", left: map[string]float64{"vtgate": 10, "vttablet": 10}, right: map[string]float64{"vtgate": 5, "vttablet": 20}, want: map[string]float64{"vtgate": 50, "vttablet": -100}},
		{name: "Left empty", left: map[string]float64{}, right: map[string]float64{"vtgate": 10, "vttablet": 10}, want: map[string]float64{"vtgate": -100, "vttablet": -100}},
		{name: "Right empty", right: map[string]float64{}, left: map[string]float64{"vtgate": 10, "vttablet": 10}, want: map[string]float64{"vtgate": 0, "vttablet": 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			got := compareSafeComponentMap(tt.left, tt.right)
			c.Assert(got, qt.DeepEquals, tt.want)
		})
	}
}
