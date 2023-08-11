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

import './exeQueue.css'

import PropTypes from 'prop-types';

const ExeQueue = ({data}) => {

    return (
        <div className='queue__data flex'>
            <span className='sha'><a target='_blank' rel="noopener noreferrer" href={`https://github.com/vitessio/vitess/commit/${data.git_ref}`}>{data.git_ref.slice(0,6)}</a></span>
            <span className='width--11em'>{data.source}</span>
            <span className='width--11em'>{data.type_of}</span>
            <span className='width--5em'>{data.pull_nb}</span>
        </div>
    );
};

ExeQueue.propTypes = {
    data: PropTypes.shape({
        git_ref: PropTypes.string.isRequired,
        source: PropTypes.string.isRequired,
        type_of: PropTypes.string.isRequired,
        pull_nb: PropTypes.number.isRequired,
    }).isRequired,
};

export default ExeQueue;