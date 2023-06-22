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

import '../PreviousExecutionResponsiveTablet/previousExeResTablet.css'
import { getStatusClass, formatDate, closeTables, openTables} from '../../utils/utils';

const PreviousExeResTablet = ({data, className}) => {

    const [maxHeight, setMaxHeight] = useState(closeTables);

    const handleClick = () => {
        if (maxHeight === closeTables) {
            setMaxHeight(openTables);
          } else {
            setMaxHeight(closeTables);
        }
    };

    return (
        <div className={`previousExe__data__tablet ${className}`} style={{ maxHeight: `${maxHeight}px` }}>
            <div className='previousExe__data__tablet__top flex'>
                <span className='width--6em'><a target='_blank' href={`https://github.com/vitessio/vitess/commit/${data.git_ref}`}>{data.git_ref.slice(0,6)}</a></span>
                <span className='width--11em'>{data.source}</span>
                <span className='width--11em'>{formatDate(data.started_at)}</span>
                <span className='width--11em'>{formatDate(data.finished_at)}</span>
                <span  className={`data ${getStatusClass(data.status)} spanStatus width--6em`}>{data.status}</span>
                <span className='width--3em'><i className="fa-solid fa-circle-info" onClick={handleClick}></i></span>

            </div>
            <div className='previousExe__data__tablet__bottom flex'>
                <div className='previousExe__data__tablet__bottom__more flex'>
                        <span>UUID</span>
                        <i className="fa-solid fa-arrow-right"></i>
                        <span >{data.uuid.slice(0, 8)}</span>
                    </div>
                <div className='previousExe__data__tablet__bottom__more flex'>
                    <span>Type</span>
                    <i className="fa-solid fa-arrow-right"></i>
                    <span>{data.type_of}</span>
                </div>
                <div className='previousExe__data__tablet__bottom__more flex'>
                    <span>PR</span>
                    <i className="fa-solid fa-arrow-right"></i>
                    <span>{data.pull_nb}</span>
                </div>
                <div className='previousExe__data__tablet__bottom__more flex'>
                    <span>Go version</span>
                    <i className="fa-solid fa-arrow-right"></i>
                    <span>{data.golang_version}</span>
                </div>
            </div>
        </div>
    );
};

export default PreviousExeResTablet;