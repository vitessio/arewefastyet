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


import './previousexe.css'
import { getStatusClass, formatDate} from '../../../utils/Utils';

const PreviousExe = ({data, className}) => {
          
          
    return (

        <div className={`previousExe__data  ${className}`}>
          <span className='width--6em '>{data.uuid.slice(0, 8)}</span>
          <span className='width--6em'><a target='_blank' href={`https://github.com/vitessio/vitess/commit/${data.git_ref}`}>{data.git_ref.slice(0,6)}</a></span>
          <span className='tdSource width--11em'>{data.source}</span>
          <span className='tdSource width--11em'>{formatDate(data.started_at)}</span>
          <span className='tdSource width--11em'>{formatDate(data.finished_at)}</span>
          <span className='width--11em'>{data.type_of}</span>
          <span className='width--5em'>{data.pull_nb}</span>
          <span className='width--6em'>{data.golang_version}</span>
          <span  className={`data ${getStatusClass(data.status)} spanStatus width--6em`}>{data.status}</span>
        </div>
        
    );
};

export default PreviousExe;