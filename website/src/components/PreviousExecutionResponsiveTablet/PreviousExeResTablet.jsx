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
import moment from 'moment';

import '../PreviousExecutionResponsiveTablet/previousExeResTablet.css'

const PreviousExeResTablet = ({data, className}) => {

    const getStatusClass = (status) => {
        switch (status) {
            case 'finished':
              return 'finished';
            case 'failed':
              return 'failed';
            case 'started':
              return 'started';
            default:
              return 'default';
          }
    }


    // Create a Moment object from the given date string
    const startedDate = moment(data.started_at)
    const finishedDate = moment(data.finished_at)
    
    // Format the date using the format method of Moment.js
    // Here, we use 'DD/MM/YYYY HH:mm:ss' as the desired format
    const formattedStartedDate = startedDate.format('MM/DD/YYYY HH:mm')
    const formattedFinishedDate = finishedDate.format('MM/DD/YYYY HH:mm')


    const [maxHeight, setMaxHeight] = useState(70);

    const handleClick = () => {
        if (maxHeight === 70) {
            setMaxHeight(400);
          } else {
            setMaxHeight(70);
          }
        };
    return (
        <div className={`previousExe__data__tablet ${className}`} style={{ maxHeight: `${maxHeight}px` }}>
            <div className='previousExe__data__tablet__top flex'>
                <span className='width--6em'><a target='_blank' href={`https://github.com/vitessio/vitess/commit/${data.git_ref}`}>{data.git_ref.slice(0,6)}</a></span>
                <span className='width--11em'>{data.source}</span>
                <span className='width--11em'>{formattedStartedDate}</span>
                <span className='width--11em'>{formattedFinishedDate}</span>
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