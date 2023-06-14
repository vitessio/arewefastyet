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
import { formatByteForGB } from '../../utils/utils';

import '../Macrobench/macrobench.css'

const Macrobench = ({data, dropDownLeft, dropDownRight, swiperSlide}) => {
    return (
        <div className='macrobench__component'>
            <h3>{data.type}</h3>
            <div className='macrobench__component__container flex'>
                <div className='macrobench__data flex--column'>
                    <h4>{dropDownLeft}</h4>
                    <span>{data.diff.Left.Result.qps.total}</span>
                    <span>{data.diff.Left.Result.qps.reads}</span>
                    <span>{data.diff.Left.Result.qps.writes}</span>
                    <span>{data.diff.Left.Result.qps.other}</span>
                    <span>{data.diff.Left.Result.tps}</span>
                    <span>{data.diff.Left.Result.latency}</span>
                    <span>{data.diff.Left.Result.errors}</span>
                    <span>{data.diff.Left.Result.reconnects}</span>
                    <span>{data.diff.Left.Result.time}</span>
                    <span>{data.diff.Left.Result.threads}</span>
                    <span>{data.diff.Left.Metrics.TotalComponentsCPUTime}</span>
                    {/* <span>{data.diff.Left.Metrics.ComponentsCPUTime.vtgate}</span>
                    <span>{data.diff.Left.Metrics.ComponentsCPUTime.vttablet}</span> */}
                    <span>{formatByteForGB(data.diff.Left.Metrics.TotalComponentsMemStatsAllocBytes)}</span>
                    {/* <span>{data.diff.Left.Metrics.ComponentsMemStatsAllocBytes.vtgate}</span>
                    <span>{data.diff.Left.Metrics.ComponentsMemStatsAllocBytes.vttablet}</span> */}
                    
                </div>
                <div className='macrobench__data flex--column' >
                    <h4>{dropDownRight}</h4>
                    <span>{data.diff.Right.Result.qps.total}</span>
                    <span>{data.diff.Right.Result.qps.reads}</span>
                    <span>{data.diff.Right.Result.qps.writes}</span>
                    <span>{data.diff.Right.Result.qps.other}</span>
                    <span>{data.diff.Right.Result.tps}</span>
                    <span>{data.diff.Right.Result.latency}</span>
                    <span>{data.diff.Right.Result.errors}</span>
                    <span>{data.diff.Right.Result.reconnects}</span>
                    <span>{data.diff.Right.Result.time}</span>
                    <span>{data.diff.Right.Result.threads}</span>
                    <span>{data.diff.Right.Metrics.TotalComponentsCPUTime}</span>
                    {/* <span>{data.diff.Right.Metrics.ComponentsCPUTime.vtgate}</span>
                    <span>{data.diff.Right.Metrics.ComponentsCPUTime.vttablet}</span> */}
                    <span>{data.diff.Right.Metrics.TotalComponentsMemStatsAllocBytes}</span>
                    {/* <span>{data.diff.Right.Metrics.ComponentsMemStatsAllocBytes.vtgate}</span>
                    <span>{data.diff.Right.Metrics.ComponentsMemStatsAllocBytes.vttablet}</span> */}
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