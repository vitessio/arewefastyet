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

import React, {useState} from 'react';


import './previousExeResMobile.css'
import { getStatusClass, formatDate, closeTables, openTables} from '../../../utils/Utils';

const PreviousExeRes = ({data, className}) => {

    const [maxHeight, setMaxHeight] = useState(closeTables);

    const handleClick = () => {
        if (maxHeight === closeTables) {
            setMaxHeight(openTables);
          } else {
            setMaxHeight(closeTables);
          }
        };

    return (
        <div className={`previousExe__data__mobile ${className}`} style={{ maxHeight: `${maxHeight}px` }}>
            <div className='previousExe__data__mobile__top flex'>
                <span className='width--11em'>{data.source}</span>
                <span  className={`data ${getStatusClass(data.status)} spanStatus width--6em`}>{data.status}</span>
                <span className='width--3em'><i className="fa-solid fa-circle-info" onClick={handleClick}></i></span>
            </div>
            <div className='previousExe__data__mobile__bottom flex'>
                <table className='justify--content'>
                    <thead >
                        <tr className='flex--column'>
                            <th>UUID</th>
                            <th>SHA</th>
                            <th>Type</th>
                            <th>Started</th>
                            <th>Finished</th>
                            <th>PR</th>
                            <th>Go version</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr className='flex--column'>
                            <td><span >{data.uuid.slice(0, 8)}</span></td>
                            <td><span ><a target='_blank' href={`https://github.com/vitessio/vitess/commit/${data.git_ref}`}>{data.git_ref.slice(0,6)}</a></span></td>
                            <td> <span>{data.type_of}</span></td>
                            <td><span>{formatDate(data.started_at)}</span></td>
                            <td><span>{formatDate(data.finished_at)}</span></td>
                            <td> <span>{data.pull_nb}</span></td>
                            <td><span>{data.golang_version}</span></td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default PreviousExeRes;