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
import ReactJson from "react-json-pretty";

import "react-json-pretty/themes/monikai.css";
import "../QueryPlan/queryPlan.css";

const QueryPlan = ({ data, isOpen, togglePlan }) => {
    
  // The queryPlanKeyStyle and queryTitleStyle are used to define the background of elements based on the open state and window width.
  // If isOpen is true and the window width is greater than 1225 pixels, the background will be orange for queryPlanKeyStyle.
  // If isOpen is true and the window width is less than 1225 pixels, the background will be orange for queryTitleStyle.
  // These styles dynamically adjust the appearance of elements based on these conditions.
  const queryPlanKeyStyle = {
    background: isOpen && window.innerWidth > 1225 ? "orange" : "initial",
  };

  const queryTitleStyle = {
    background: isOpen && window.innerWidth < 1225 ? "orange" : "initial",
  };

  return (
    <div className="queryPlan">
      <div
        className="queryPlan__key"
        style={queryPlanKeyStyle}
        onClick={togglePlan}
      >
        <span className="queryPlan__key__span" style={queryTitleStyle}>
          {data.Key}
        </span>
        <div className="badge__container">
          {data.Left &&
            data.Right &&
            data.Left.Value.Errors > 0 &&
            data.Right.Value.Errors > 0 && (
              <span className="badge badge--danger">Both Have Errors</span>
            )}
          {data.Left && data.Right && data.Left.Value.Errors > 0 && (
            <span className="badge badge--danger">A Has Errors</span>
          )}
          {data.Left && data.Right && data.Right.Value.Errors > 0 && (
            <span className="badge badge--danger">B Has Errors</span>
          )}
          {!(data.Right.Value.Errors > 0) && (
            <span className="badge badge--danger">B Has Errors</span>
          )}
          {!(<span className="badge badge--info">Only B</span>)}
          {!data.Right && data.Left && data.Left.Value.Errors > 0 && (
            <span className="badge badge--danger">A Has Errors</span>
          )}
          {!data.Right && data.Left && (
            <span className="badge badge--info">Only A</span>
          )}
          {data.ExecTimeDiff > 5 && (
            <span className="badge badge--warning">A fastest</span>
          )}
          {data.ExecTimeDiff < -5 && (
            <span className="badge badge--succes">B fastest</span>
          )}
        </div>
        {isOpen && (
          <div className="plan">
            <div className="statistics">
              <h3>Statistics</h3>
              <div>
                <div className="statistics__top">
                  <span className="statistics__key__title">
                    <b>Query</b>
                  </span>
                  <span className="statistics__key">{data.Key}</span>
                </div>
                <div className="statistics__table">
                  <table>
                    <thead>
                      <tr>
                        <th colSpan="1"></th>
                        <th colSpan="1">A</th>
                        <th colSpan="1">B</th>
                        <th colSpan="1">Diff</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr>
                        <td>
                          <b>ExecutionsCount</b>
                        </td>
                        {<td>{data.Left && data.Left.Value.ExecCount}</td>}
                        {<td>{data.Right && data.Right.Value.ExecCount}</td>}
                        <td>{data.ExecCountDiff}</td>
                      </tr>
                      <tr>
                        <td>
                          <b>Execution Time</b>
                        </td>
                        {<td>{data.Left && data.Left.Value.ExecTime}</td>}
                        {<td>{data.Right && data.Right.Value.ExecTime}</td>}
                        <td>{data.ExecTimeDiff}</td>
                      </tr>
                      <tr>
                        <td>
                          <b>Rows Returned</b>
                        </td>
                        {<td>{data.Left && data.Left.Value.RowsReturned}</td>}
                        {<td>{data.Right && data.Right.Value.RowsReturned}</td>}
                        <td>{data.RowsReturnedDiff}</td>
                      </tr>
                      <tr>
                        <td>
                          <b>Errors</b>
                        </td>
                        {<td>{data.Left && data.Left.Value.Errors}</td>}
                        {<td>{data.Right && data.Right.Value.Errors}</td>}
                        <td>{data.ErrorsDiff}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
            {data.Left &&
            data.right &&
            data.Left.Value.Instructions === data.Right.Value.Instructions ? (
              <div className="planEquivalent">Plans are equivalent.</div>
            ) : (
              <div className="planDifferent">Plans are different</div>
            )}
            {data.Left &&
            data.Right &&
            data.Left.Value.Instructions === data.Right.Value.Instructions ? (
              <div className="plan__queryPlan">
                <h3>Query Plan</h3>
                <ReactJson
                  data={data.Left.Value.Instructions}
                  className="json"
                />
              </div>
            ) : (
              <>
                {data.Left && (
                  <div className="plan__queryPlan">
                    <h3>Query Plan A</h3>
                    <ReactJson
                      data={data.Left.Value.Instructions}
                      className="json"
                    />
                  </div>
                )}
                {data.Right && (
                  <div className="plan__queryPlan">
                    <h3>Query Plan B</h3>
                    <ReactJson
                      data={data.Right.Value.Instructions}
                      className="json"
                    />
                  </div>
                )}
              </>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default QueryPlan;
