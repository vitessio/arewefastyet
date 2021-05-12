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
	"github.com/vitessio/arewefastyet/go/exec/metrics"
	"testing"
)

func TestComparison_Regression(t *testing.T) {
	tests := []struct {
		name       string
		cmp        Comparison
		wantReason string
	}{
		{name: "No regression", cmp: Comparison{}, wantReason: ""},
		{name: "Total CPU time increase (1)", cmp: Comparison{DiffMetrics: metrics.ExecutionMetrics{TotalComponentsCPUTime: -100}}, wantReason: "- Total CPU time increased by 100.00% \n"},
		{name: "Total CPU time increase (2)", cmp: Comparison{DiffMetrics: metrics.ExecutionMetrics{TotalComponentsCPUTime: -5}}, wantReason: "- Total CPU time increased by 5.00% \n"},
		{name: "Total CPU time no increase", cmp: Comparison{DiffMetrics: metrics.ExecutionMetrics{TotalComponentsCPUTime: -4}}, wantReason: ""},
		{name: "VTTablet time increase (1)", cmp: Comparison{DiffMetrics: metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{"vttablet": -35}}}, wantReason: "- vttablet CPU time increased by 35.00% \n"},
		{name: "VTTablet time increase (2)", cmp: Comparison{DiffMetrics: metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{"vttablet": -5}}}, wantReason: "- vttablet CPU time increased by 5.00% \n"},
		{name: "VTTablet time no increase", cmp: Comparison{DiffMetrics: metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{"vttablet": -3.59}}}, wantReason: ""},
		{name: "VTGate time increase (1)", cmp: Comparison{DiffMetrics: metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{"vtgate": -11.98}}}, wantReason: "- vtgate CPU time increased by 11.98% \n"},
		{name: "VTGate time increase (2)", cmp: Comparison{DiffMetrics: metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{"vtgate": -5}}}, wantReason: "- vtgate CPU time increased by 5.00% \n"},
		{name: "VTGate time no increase", cmp: Comparison{DiffMetrics: metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{"vtgate": -4.99}}}, wantReason: ""},
		{name: "VTGate and VTTablet times increase", cmp: Comparison{DiffMetrics: metrics.ExecutionMetrics{ComponentsCPUTime: map[string]float64{"vtgate": -5, "vttablet": -5}}}, wantReason: "- vtgate CPU time increased by 5.00% \n- vttablet CPU time increased by 5.00% \n"},
		{name: "TPS decrease", cmp: Comparison{Diff: Result{TPS: -50}}, wantReason: "- TPS decreased by 50.00% \n"},
		{name: "Latency increase", cmp: Comparison{Diff: Result{Latency: -15}}, wantReason: "- Latency increased by 15.00% \n"},
		{name: "QPS decrease", cmp: Comparison{Diff: Result{QPS: QPS{Total: -10}}}, wantReason: "- QPS decreased by 10.00% \n"},
		{name: "TPS and QPS decrease", cmp: Comparison{Diff: Result{TPS: -32.5, QPS: QPS{Total: -27.7}}}, wantReason: "- TPS decreased by 32.50% \n- QPS decreased by 27.70% \n"},
		{name: "TPS, QPS decrease and Latency increase", cmp: Comparison{Diff: Result{Latency: -10, TPS: -32.5, QPS: QPS{Total: -27.7}}}, wantReason: "- TPS decreased by 32.50% \n- QPS decreased by 27.70% \n- Latency increased by 10.00% \n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			reason := tt.cmp.Regression()
			c.Assert(reason, qt.Contains, tt.wantReason)
		})
	}
}
