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
import { Link } from 'react-router-dom';
import { formatByteForGB, fixed } from '../../../utils/Utils';

import './macrobench.css'

const Macrobench = React.memo(({data, gitRefLeft, gitRefRight, showSpan, commitHashLeft, commitHashRight}) => {
    return (
        <div className='macrobench__component'>
            <div className='macrobench__component__header'>
            <h3>{data.type}</h3>
            <span className='linkQueryT'>Click <Link to="/macrobench/queries/compare" state={{data : data, commitHashLeft : commitHashLeft, commitHashRight : commitHashRight }}>here</Link> to see the query plans comparison for this benchmark.</span>
            </div>
            
            <table>
                <thead>
                    <tr>
                        <th><h4>{gitRefLeft ? gitRefLeft : 'Left'}</h4></th>
                        <th><h4>{gitRefRight ? gitRefRight : 'Right'}</h4></th>
                        <th><h4>Impoved by %</h4></th>
                    </tr>
                </thead>
                <tbody>
                    <tr className='border--bottom'>
                        <td><span>{data.diff.Left.Result.qps.total}</span></td>
                        <td><span>{data.diff.Right.Result.qps.total}</span></td>
                        <td><span>{fixed(data.diff.Diff.qps.total, 2)}</span></td>
                    </tr>
                    <tr className='border--bottom'>
                        <td><span>{data.diff.Left.Result.qps.reads}</span></td>
                        <td><span>{data.diff.Right.Result.qps.reads}</span></td>
                        <td><span>{fixed(data.diff.Diff.qps.reads, 2)}</span></td>
                    </tr>
                    <tr className='border--bottom'>
                        <td><span>{data.diff.Left.Result.qps.writes}</span></td>
                        <td><span>{data.diff.Right.Result.qps.writes}</span></td>
                        <td><span>{fixed(data.diff.Diff.qps.writes, 2)}</span></td>
                    </tr>
                    <tr className='border--bottom'>
                        <td><span>{data.diff.Left.Result.qps.other}</span></td>
                        <td><span>{data.diff.Right.Result.qps.other}</span></td>
                        <td><span>{fixed(data.diff.Diff.qps.other, 2)}</span></td>
                    </tr>
                    <tr className='border--bottom'>
                        <td><span>{data.diff.Left.Result.tps}</span></td>
                        <td><span>{data.diff.Right.Result.tps}</span></td>
                        <td><span>{fixed(data.diff.Diff.tps, 2)}</span></td>
                    </tr>
                    <tr className='border--bottom'>
                        <td><span>{data.diff.Left.Result.latency}</span></td>
                        <td><span>{data.diff.Right.Result.latency}</span></td>
                        <td><span>{fixed(data.diff.Diff.latency, 2)}</span></td>
                    </tr>
                    <tr className='border--bottom'>
                        <td><span>{data.diff.Left.Result.errors}</span></td>
                        <td><span>{data.diff.Right.Result.errors}</span></td>
                        <td><span>{fixed(data.diff.Diff.errors, 2)}</span></td>
                    </tr>
                    <tr className='border--bottom'>
                        <td><span>{data.diff.Left.Result.reconnects}</span></td>
                        <td><span>{data.diff.Right.Result.reconnects}</span></td>
                        <td><span>{fixed(data.diff.Diff.reconnects, 2)}</span></td>
                    </tr>
                    <tr className='border--bottom'>
                        <td><span>{data.diff.Left.Result.time}</span></td>
                        <td><span>{data.diff.Right.Result.time}</span></td>
                        <td><span>{fixed(data.diff.Diff.time, 2)}</span></td>
                    </tr>
                    <tr className='border--bottom'>
                        <td><span>{data.diff.Left.Result.threads}</span></td>
                        <td><span>{data.diff.Right.Result.threads}</span></td>
                        <td><span>{fixed(data.diff.Diff.threads, 2)}</span></td>
                    </tr>
                    <tr className='border--bottom'>
                        <td><span>{fixed(data.diff.Left.Metrics.TotalComponentsCPUTime, 0)}</span></td>
                        <td><span>{fixed(data.diff.Right.Metrics.TotalComponentsCPUTime, 0)}</span></td>
                        <td><span>{fixed(data.diff.DiffMetrics.TotalComponentsCPUTime, 2)}</span></td>
                    </tr>
                    <tr className='border--bottom'>
                        <td><span>{data.diff.Left.Metrics.ComponentsCPUTime.vtgate}</span></td>
                        <td><span>{data.diff.Right.Metrics.ComponentsCPUTime.vtgate}</span></td>
                        <td><span>{fixed(data.diff.DiffMetrics.ComponentsCPUTime.vtgate, 2)}</span></td>
                    </tr>
                    <tr className='border--bottom'>
                        <td><span>{data.diff.Left.Metrics.ComponentsCPUTime.vttablet}</span></td>
                        <td><span>{data.diff.Right.Metrics.ComponentsCPUTime.vttablet}</span></td>
                        <td><span>{fixed(data.diff.DiffMetrics.ComponentsCPUTime.vttablet, 2)}</span></td>
                    </tr>
                    <tr className='border--bottom'>
                        <td><span>{formatByteForGB(data.diff.Left.Metrics.TotalComponentsMemStatsAllocBytes)}</span></td>
                        <td><span>{formatByteForGB(data.diff.Right.Metrics.TotalComponentsMemStatsAllocBytes)}</span></td>
                        <td><span>{fixed(data.diff.DiffMetrics.TotalComponentsMemStatsAllocBytes, 2)}</span></td>
                    </tr>
                    <tr className='border--bottom'>
                        <td><span>{formatByteForGB(data.diff.Left.Metrics.ComponentsMemStatsAllocBytes.vtgate)}</span></td>
                        <td><span>{formatByteForGB(data.diff.Right.Metrics.ComponentsMemStatsAllocBytes.vtgate)}</span></td>
                        <td><span>{fixed(data.diff.DiffMetrics.ComponentsMemStatsAllocBytes.vtgate, 2)}</span></td>
                    </tr>
                    <tr>
                        <td><span>{formatByteForGB(data.diff.Left.Metrics.ComponentsMemStatsAllocBytes.vttablet)}</span></td>
                        <td><span>{formatByteForGB(data.diff.Right.Metrics.ComponentsMemStatsAllocBytes.vttablet)}</span></td>
                        <td><span>{fixed(data.diff.DiffMetrics.ComponentsMemStatsAllocBytes.vttablet, 2)}</span></td>
                    </tr>
                </tbody>
            </table>
        </div>
    );
});

export default Macrobench;