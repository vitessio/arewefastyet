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

import '../PreviousExecutions/previousexe.css'

const PreviousExe = ({data, className}) => {
        /*BACKGROUND STATUS*/ 

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
          

          
          
    return (

        <div className={`previousExe__data flex ${className}`}>
          <span className='width--6em '>{data.uuid.slice(0, 8)}</span>
          <span className='width--6em'><a target='_blank' href={`https://github.com/vitessio/vitess/commit/${data.git_ref}`}>{data.git_ref.slice(0,6)}</a></span>
          <span className='tdSource width--11em'>{data.source}</span>
          <span className='tdSource width--11em'>{formattedStartedDate}</span>
          <span className='tdSource width--11em'>{formattedFinishedDate}</span>
          <span className='width--11em'>{data.type_of}</span>
          <span className='width--5em'>{data.pull_nb}</span>
          <span className='width--6em'>{data.golang_version}</span>
          <span  className={`data ${getStatusClass(data.status)} spanStatus width--6em`}>{data.status}</span>
        </div>
        
    );
};

export default PreviousExe;