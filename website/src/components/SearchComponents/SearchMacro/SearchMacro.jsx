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

import '../SearchMacro/searchMacro.css'

const SearchMacro = ({data}) => {
    const result = data[1] && Array.isArray(data[1]) && data[1].length > 0
    ? data[1][0].Result
    : null;
    return (
        <div className='search__macro__data'>
           <h3>{data[0]}</h3>
           {result ? (
                <table>
                    <tbody>
                        <tr ><td id='marginTr'><span className='marginTr'>{data[1][0].Result.qps.total}</span></td></tr>
                        <tr><td><span>{data[1][0].Result.qps.reads}</span></td></tr>
                        <tr><td><span>{data[1][0].Result.qps.writes}</span></td></tr>
                        <tr><td><span>{data[1][0].Result.qps.other}</span></td></tr>
                        <tr><td><span>{data[1][0].Result.tps}</span></td></tr>
                        <tr><td><span>{data[1][0].Result.latency}</span></td></tr>
                        <tr><td><span>{data[1][0].Result.errors}</span></td></tr>
                        <tr><td><span>{data[1][0].Result.reconnects}</span></td></tr>
                        <tr><td><span>{data[1][0].Result.time}</span></td></tr>
                        <tr><td><span>{data[1][0].Result.threads}</span></td></tr>
                        <tr><td><span>{data[1][0].Metrics.TotalComponentsCPUTime.toFixed(0)}</span></td></tr>
                        <tr><td><span>{data[1][0].Metrics.ComponentsCPUTime.vtgate}</span></td></tr>
                        <tr><td><span>{data[1][0].Metrics.ComponentsCPUTime.vttablet}</span></td></tr>
                        <tr><td><span>{data[1][0].Metrics.TotalComponentsMemStatsAllocBytes}</span></td></tr>
                        <tr><td><span>{data[1][0].Metrics.ComponentsMemStatsAllocBytes.vtgate}</span></td></tr>
                        <tr><td><span>{data[1][0].Metrics.ComponentsMemStatsAllocBytes.vttablet}</span></td></tr>
                    </tbody>
                </table>
           ) : (
            <table>
            <tbody>
                <tr ><td><span>0</span></td></tr>
                <tr><td><span>0</span></td></tr>
                <tr><td><span>0</span></td></tr>
                <tr><td><span>0</span></td></tr>
                <tr><td><span>0</span></td></tr>
                <tr><td><span>0</span></td></tr>
                <tr><td><span>0</span></td></tr>
                <tr><td><span>0</span></td></tr>
                <tr><td><span>0</span></td></tr>
                <tr><td><span>0</span></td></tr>
                <tr><td><span>0</span></td></tr>
                <tr><td><span>0</span></td></tr>
                <tr><td><span>0</span></td></tr>
                <tr><td><span>0</span></td></tr>
                <tr><td><span>0</span></td></tr>
                <tr><td><span>0</span></td></tr>
            </tbody>
        </table>
           )}
           
        </div>
    );
    
};

export default SearchMacro;