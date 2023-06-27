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

import React, {useState}from 'react';
import ReactJson from 'react-json-pretty';

import 'react-json-pretty/themes/monikai.css';
import '../QueryPlan/queryPlan.css'

const QueryPlan = ({data, isOpen, togglePlan}) => {

    const queryPlanKeyStyle = {
        background: isOpen ? 'orange' : 'initial',
    };
    
    return (
        <div className='queryPlan'>
            <div className='queryPlan__key' style={queryPlanKeyStyle} onClick={togglePlan}>
                <span>{data.Key}</span>
                {isOpen && (
                <div className='plan'>
                    <div className='statistics'>
                        <h3>Statistics</h3>
                        <div>
                            <div className='statistics__top'>
                                <span className='statistics__key__title'><b>Query</b></span>
                                <span className='statistics__key'>{data.Key}</span>
                            </div>
                            <div className='statistics__table'>
                                <table>
                                    <thead>
                                        <tr>
                                            <th colSpan="1"></th>
                                            <th colSpan="1">A</th>
                                            <th colSpan="1">B</th>
                                            <th colSpan="1">Diff</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        <tr>
                                            <td><b>ExecutionsCount</b></td>
                                            <td>{data.Left.Value.ExecCount}</td>
                                            <td>{data.Right.Value.ExecCount}</td>
                                            <td>{data.ExecCountDiff}</td>
                                        </tr>
                                        <tr>
                                            <td><b>Execution Time</b></td>
                                            <td>{data.Left.Value.ExecTime}</td>
                                            <td>{data.Right.Value.ExecTime}</td>
                                            <td>{data.ExecTimeDiff}</td>
                                        </tr>
                                        <tr>
                                            <td><b>Rows Returned</b></td>
                                            <td>{data.Left.Value.RowsReturned}</td>
                                            <td>{data.Right.Value.RowsReturned}</td>
                                            <td>{data.RowsReturnedDiff}</td>
                                        </tr>
                                        <tr>
                                            <td><b>Errors</b></td>
                                            <td>{data.Left.Value.Errors}</td>
                                            <td>{data.Right.Value.Errors}</td>
                                            <td>{data.ErrorsDiff}</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                    {data.Left.Value.Plan === data.Right.Value.Plan ? (
                        <div className='planEquivalent'>Plans are equivalent.</div>
                    ) : (
                        <div className='planDifferent'>Plans are différent</div>
                    )}
                    <div className='plan__queryPlan'>
                        <h3>Query Plan</h3>
                        <ReactJson data={data.Left.Value.Instructions} className='json'/>
                    </div>
                    {data.Left.Value.Plan === data.Right.Value.Plan ? null : (
                        <div className='plan__queryPlan'>
                            <h3>Query Plan</h3>
                            <ReactJson data={data.Right.Value.Instructions} className='json'/>
                        </div>
                    )}
                </div>
                )}
            </div>
        </div>
    );
};

export default QueryPlan;