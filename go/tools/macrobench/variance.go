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
	"github.com/montanaflynn/stats"
	"sort"
)

func GetVarianceForMacroBenchmarks(macrobenchmarks DetailsArray) (Details, Details) {
	var qps, r, w, o, tps, lat, cpu, cpuGate, cpuTablet, mem, memGate, memTablet stats.Float64Data
	res := newEmptyDetails()
	resPercentage := newEmptyDetails()

	for _, m := range macrobenchmarks {
		qps = append(qps, m.Result.QPS.Total)
		r = append(r, m.Result.QPS.Reads)
		w = append(w, m.Result.QPS.Writes)
		o = append(o, m.Result.QPS.Other)
		tps = append(tps, m.Result.TPS)
		lat = append(lat, m.Result.Latency)
		cpu = append(cpu, m.Metrics.TotalComponentsCPUTime)
		cpuGate = append(cpuGate, m.Metrics.ComponentsCPUTime["vtgate"])
		cpuTablet = append(cpuTablet, m.Metrics.ComponentsCPUTime["vttablet"])
		mem = append(mem, m.Metrics.TotalComponentsMemStatsAllocBytes)
		memGate = append(memGate, m.Metrics.ComponentsMemStatsAllocBytes["vtgate"])
		memTablet = append(memTablet, m.Metrics.ComponentsMemStatsAllocBytes["vttablet"])
	}

	sort.Sort(qps)
	sort.Sort(r)
	sort.Sort(w)
	sort.Sort(o)
	sort.Sort(tps)
	sort.Sort(lat)
	sort.Sort(cpu)
	sort.Sort(cpuGate)
	sort.Sort(cpuTablet)
	sort.Sort(mem)
	sort.Sort(memGate)
	sort.Sort(memTablet)

	res.Result.QPS.Total, _ = stats.StandardDeviationSample(qps)
	res.Result.QPS.Reads, _ = stats.StandardDeviationSample(r)
	res.Result.QPS.Writes, _ = stats.StandardDeviationSample(w)
	res.Result.QPS.Other, _ = stats.StandardDeviationSample(o)
	res.Result.TPS, _ = stats.StandardDeviationSample(tps)
	res.Result.Latency, _ = stats.StandardDeviationSample(lat)
	res.Metrics.TotalComponentsCPUTime, _ = stats.StandardDeviationSample(cpu)
	res.Metrics.ComponentsCPUTime["vtgate"], _ = stats.StandardDeviationSample(cpuGate)
	res.Metrics.ComponentsCPUTime["vttablet"], _ = stats.StandardDeviationSample(cpuTablet)
	res.Metrics.TotalComponentsMemStatsAllocBytes, _ = stats.StandardDeviationSample(mem)
	res.Metrics.ComponentsMemStatsAllocBytes["vtgate"], _ = stats.StandardDeviationSample(memGate)
	res.Metrics.ComponentsMemStatsAllocBytes["vttablet"], _ = stats.StandardDeviationSample(memTablet)

	resPercentage.Result.QPS.Total = res.Result.QPS.Total * 100 / mean(qps)
	resPercentage.Result.QPS.Reads = res.Result.QPS.Reads * 100 / mean(r)
	resPercentage.Result.QPS.Writes = res.Result.QPS.Writes * 100 / mean(w)
	resPercentage.Result.QPS.Other = res.Result.QPS.Other * 100 / mean(o)
	resPercentage.Result.TPS = res.Result.TPS * 100 / mean(tps)
	resPercentage.Result.Latency = res.Result.Latency * 100 / mean(lat)
	resPercentage.Metrics.TotalComponentsCPUTime = res.Metrics.TotalComponentsCPUTime * 100 / mean(cpu)
	resPercentage.Metrics.ComponentsCPUTime["vtgate"] = res.Metrics.ComponentsCPUTime["vtgate"] * 100 / mean(cpuGate)
	resPercentage.Metrics.ComponentsCPUTime["vttablet"] = res.Metrics.ComponentsCPUTime["vttablet"] * 100 / mean(cpuTablet)
	resPercentage.Metrics.TotalComponentsMemStatsAllocBytes = res.Metrics.TotalComponentsMemStatsAllocBytes * 100 / mean(mem)
	resPercentage.Metrics.ComponentsMemStatsAllocBytes["vtgate"] = res.Metrics.ComponentsMemStatsAllocBytes["vtgate"] * 100 / mean(memGate)
	resPercentage.Metrics.ComponentsMemStatsAllocBytes["vttablet"] = res.Metrics.ComponentsMemStatsAllocBytes["vttablet"] * 100 / mean(memTablet)
	return res, resPercentage
}

func mean(f stats.Float64Data) float64 {
	m, _ := stats.Mean(f)
	return m
}

func newEmptyDetails() Details {
	res := Details{}
	res.Metrics.ComponentsCPUTime = map[string]float64{}
	res.Metrics.ComponentsMemStatsAllocBytes = map[string]float64{}
	return res
}
