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

import{ useState } from "react";
import PropTypes from 'prop-types';

import "../PRcomponents/PRGitInfo.css";

import { closeTables, openTables, formatDate } from "../../utils/Utils";

const PRGitInfo = ({ data, setPrNumber, className }) => {

  const handlePrInfo = (e) => {
    const number = e.toString();
    setPrNumber(number);
  };

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
      className={`prGit flex--column ${className}`}
      style={{ maxHeight: `${maxHeight}px` }}
    >
      <div className="prGit__top flex">
        <span className="width--4em">
          <a href={`https://github.com/vitessio/vitess/pull/${data.ID}`}>
            {data.ID}
          </a>
        </span>
        <span className="width--40 hidden--tablet">{data.Title}</span>
        <span className="width--20 hidden--tablet"><a target="_blank" rel="noopener noreferrer" href={`https://github.com/${data.Author}`}>{data.Author}</a></span>
        <span className="width--10em hidden--mobile">
          {formatDate(data.CreatedAt)}
        </span>
        <span className="linkToCompare" >
        <i className="fa-solid fa-arrow-right" onClick={() => handlePrInfo(data.ID)}></i>
        <span id="compareMessage">compare</span>
        </span>
        <span className="hidden--desktop">
          <i className="fa-solid fa-circle-info" onClick={handleClick}></i>
        </span>
      </div>
      <div className="prGit__bottom ">
        <table className="hidden--desktop">
          <thead>
            <tr>
              <th className="thTitle">Title</th>
              <th>Author</th>
              <th>Created_at</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td className="tdTitle">{data.Title}</td>
              <td>{data.Author}</td>
              <td>{formatDate(data.CreatedAt)}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  );
};

PRGitInfo.propTypes = {
  data: PropTypes.shape({
    ID: PropTypes.number.isRequired,
    Title: PropTypes.string.isRequired,
    Author: PropTypes.string.isRequired,
    CreatedAt: PropTypes.string.isRequired,
  }).isRequired,
  setPrNumber: PropTypes.func.isRequired,
  className: PropTypes.string,
};

export default PRGitInfo;
