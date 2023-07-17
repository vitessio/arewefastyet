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
import useApiCall from "../../utils/Hook";

import "../PRcomponents/PRGitInfo.css";

import { closeTables, openTables } from "../../utils/Utils";

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
        <span className="width--5em">{data.ID}</span>
        <span className="width--15em hidden--tablet">{data.Title}</span>
        <span className="width--6em hidden--tablet">{data.Author}</span>
        <span className="width--10em hidden--mobile">{data.CreatedAt}</span>
        <span
          className="linkToCompare width--11em"
          onClick={() => handlePrInfo(data.ID)}
        >
          Click to compare with main
        </span>
        <span className="hidden--desktop">
          <i className="fa-solid fa-circle-info" onClick={handleClick}></i>
        </span>
      </div>
      <div className="prGit__bottom ">
        <table className="hidden--desktop">
          <thead>
            <tr>
              <th>Title</th>
              <th>Author</th>
              <th>Created_at</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>{data.Title}</td>
              <td>{data.Author}</td>
              <td>{data.CreatedAt}</td>
            </tr>
          </tbody>
          </table>
      </div>
    </div>
  );
};

export default PRGitInfo;
