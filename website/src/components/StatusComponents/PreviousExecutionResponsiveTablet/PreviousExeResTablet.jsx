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

import React, { useState } from "react";
import PropTypes from 'prop-types';


import "./previousExeResTablet.css";
import {
  getStatusClass,
  formatDate,
  closeTables,
  openTables,
} from "../../../utils/Utils";

const PreviousExeResTablet = ({ data, className }) => {
  const [maxHeight, setMaxHeight] = useState(closeTables);

  const handleClick = () => {
    if (maxHeight === closeTables) {
      setMaxHeight(openTables);
    } else {
      setMaxHeight(closeTables);
    }
  };

  return (
    <div
      className={`previousExe__data__tablet ${className}`}
      style={{ maxHeight: `${maxHeight}px` }}
    >
      <div className="previousExe__data__tablet__top flex">
        <span className="width--6em">
          <a
            target="_blank"
            rel="noopener noreferrer"
            href={`https://github.com/vitessio/vitess/commit/${data.git_ref}`}
          >
            {data.git_ref.slice(0, 6)}
          </a>
        </span>
        <span className="width--11em">{data.source}</span>
        <span className="width--11em">{formatDate(data.started_at)}</span>
        <span className="width--11em">{formatDate(data.finished_at)}</span>
        <span
          className={`data ${getStatusClass(
            data.status
          )} spanStatus width--6em`}
        >
          {data.status}
        </span>
        <span className="width--3em">
          <i className="fa-solid fa-circle-info" onClick={handleClick}></i>
        </span>
      </div>
      <div className="previousExe__data__tablet__bottom flex">
        <table>
          <thead>
            <tr>
              <th><span>UUID</span></th>
              <th><span>Type</span></th>
              <th><span>PR</span></th>
              <th><span>Go version</span></th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><span>{data.uuid.slice(0, 8)}</span></td>
              <td><span>{data.type_of}</span></td>
              <td><span>{data.pull_nb}</span></td>
              <td><span>{data.golang_version}</span></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  );
};

PreviousExeResTablet.propTypes = {
  data: PropTypes.shape({
    git_ref: PropTypes.string.isRequired,
    source: PropTypes.string.isRequired,
    started_at: PropTypes.string.isRequired,
    finished_at: PropTypes.string.isRequired,
    status: PropTypes.string.isRequired,
    uuid: PropTypes.string.isRequired,
    type_of: PropTypes.string.isRequired,
    pull_nb: PropTypes.number.isRequired,
    golang_version: PropTypes.string.isRequired,
  }).isRequired,
  className: PropTypes.string,
};

export default PreviousExeResTablet;
