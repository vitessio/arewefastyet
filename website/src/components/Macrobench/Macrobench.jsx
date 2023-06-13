/*
Copyright 2023 The Vitess Authors.

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

import React from 'react';

import '../Macrobench/macrobench.css'

const Macrobench = ({data, dropDownLeft, dropDownRight}) => {
    return (
        <div className='macrobench__component'>
            <h3>{data.type}</h3>
            <div className='macrobench__component__container flex'>
                <div className='macrobench__data flex--column'>
                    <h4>{dropDownLeft}</h4>
                    <span>{data.diff.Compare.Result.qps.total}</span>
                    <span>{data.diff.Compare.Result.qps.reads}</span>
                    <span>{data.diff.Compare.Result.qps.writes}</span>
                    <span>{data.diff.Compare.Result.qps.other}</span>
                    <span>{data.diff.Compare.Result.tps}</span>
                    <span>{data.diff.Compare.Result.latency}</span>
                    <span>{data.diff.Compare.Result.errors}</span>
                    <span>{data.diff.Compare.Result.reconnects}</span>
                    <span>{data.diff.Compare.Result.time}</span>
                    <span>{data.diff.Compare.Result.threads}</span>
                    <span>{data.diff.Compare.Metrics.TotalComponentsCPUTime}</span>
                    {/* <span>{data.diff.Compare.Metrics.ComponentsCPUTime.vtgate}</span>
                    <span>{data.diff.Compare.Metrics.ComponentsCPUTime.vttablet}</span> */}
                    <span>{data.diff.Compare.Metrics.TotalComponentsMemStatsAllocBytes}</span>
                    {/* <span>{data.diff.Compare.Metrics.ComponentsMemStatsAllocBytes.vtgate}</span>
                    <span>{data.diff.Compare.Metrics.ComponentsMemStatsAllocBytes.vttablet}</span> */}
                    
                </div>
                <div className='macrobench__data flex--column' >
                    <h4>{dropDownRight}</h4>
                    <span>{data.diff.Reference.Result.qps.total}</span>
                    <span>{data.diff.Reference.Result.qps.reads}</span>
                    <span>{data.diff.Reference.Result.qps.writes}</span>
                    <span>{data.diff.Reference.Result.qps.other}</span>
                    <span>{data.diff.Reference.Result.tps}</span>
                    <span>{data.diff.Reference.Result.latency}</span>
                    <span>{data.diff.Reference.Result.errors}</span>
                    <span>{data.diff.Reference.Result.reconnects}</span>
                    <span>{data.diff.Reference.Result.time}</span>
                    <span>{data.diff.Reference.Result.threads}</span>
                    <span>{data.diff.Reference.Metrics.TotalComponentsCPUTime}</span>
                    {/* <span>{data.diff.Reference.Metrics.ComponentsCPUTime.vtgate}</span>
                    <span>{data.diff.Reference.Metrics.ComponentsCPUTime.vttablet}</span> */}
                    <span>{data.diff.Reference.Metrics.TotalComponentsMemStatsAllocBytes}</span>
                    {/* <span>{data.diff.Reference.Metrics.ComponentsMemStatsAllocBytes.vtgate}</span>
                    <span>{data.diff.Reference.Metrics.ComponentsMemStatsAllocBytes.vttablet}</span> */}
                </div>
                <div className='macrobench__data flex--column'>
                    <h4>Impoved by %</h4>
                    <span>{data.diff.Diff.qps.total}</span>
                    <span>{data.diff.Diff.qps.reads}</span>
                    <span>{data.diff.Diff.qps.writes}</span>
                    <span>{data.diff.Diff.qps.other}</span>
                    <span>{data.diff.Diff.tps}</span>
                    <span>{data.diff.Diff.latency}</span>
                    <span>{data.diff.Diff.errors}</span>
                    <span>{data.diff.Diff.reconnects}</span>
                    <span>{data.diff.Diff.time}</span>
                    <span>{data.diff.Diff.threads}</span>
                    <span>{data.diff.DiffMetrics.TotalComponentsCPUTime}</span>
                    {/* <span>{data.diff.DiffMetrics.ComponentsCPUTime.vtgate}</span>
                    <span>{data.diff.DiffMetrics.ComponentsCPUTime.vttablet}</span> */}
                    <span>{data.diff.DiffMetrics.TotalComponentsMemStatsAllocBytes}</span>
                    {/* <span>{data.diff.DiffMetrics.ComponentsMemStatsAllocBytes.vtgate}</span>
                    <span>{data.diff.DiffMetrics.ComponentsMemStatsAllocBytes.vttablet}</span> */}
                </div>
            </div>
        </div>
    );
};

export default Macrobench;