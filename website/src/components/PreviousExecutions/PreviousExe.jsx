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

const PreviousExe = ({data, handleClick}) => {
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

        <tr className='previousExetd'>
            <td className='hiddenResponsiveMobile'>{data.uuid.slice(0, 8)}</td>
            <td className='hiddenResponsiveMobile'><a target='_blank' href={`https://github.com/vitessio/vitess/commit/${data.git_ref}`}>{data.git_ref.slice(0,6)}</a></td>
            <td className='tdSource'>{data.source}</td>
            <td className='hiddenResponsiveMobile'>{formattedStartedDate}</td>
            <td className='hiddenResponsiveMobile'>{formattedFinishedDate}</td>
            <td className='hiddenResponsiveMobile'>{data.type_of}</td>
            <td className='tdPR hiddenResponsiveMobile'>{data.pull_nb}</td>
            <td className='hiddenResponsiveMobile'>{data.golang_version}</td>
            <td ><span className={`data ${getStatusClass(data.status)} tdStatus`}>{data.status}</span></td>
            <td  className='tdInfos hiddenResponsiveDesktop '><i id={data.uuid} onClick={handleClick} className='fa-sharp fa-solid fa-circle-info'></i></td>
        </tr>
        
    );
};

export default PreviousExe;