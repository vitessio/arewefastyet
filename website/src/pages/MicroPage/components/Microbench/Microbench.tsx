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
import PropTypes from "prop-types";

import "../Microbench/microbench.css";

import { closeTables, openTables, fixed } from "../../../../utils/Utils";

const Microbench = ({ data, className, gitRefLeft, gitRefRight }) => {
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
      className={`microbench  ${className}`}
      style={{ maxHeight: `${maxHeight}px` }}
    >
      <div className="space--between justify--content align--center">
        <span className="width--12em">{data.PkgName}</span>
        <span className="name">{data.SubBenchmarkName}</span>
        <div className="width--18em space--between--flex hiddenMobile">
          <span className="width--100">{fixed(data.Right.Ops, 0)}</span>
          <span className="width--100">{fixed(data.Left.Ops, 0)}</span>
          <span
            className={`width--100 ${
              data.Diff.Ops <= -5
                ? "negatif--Micro"
                : data.Diff.Ops >= 5
                  ? "positif--Micro"
                  : ""
            }`}
          >
            {data.Diff.Ops.toFixed(2)}
          </span>
        </div>
        <div className="width--18em space--between--flex hiddenTablet">
          <span className="width--100">{fixed(data.Right.NSPerOp, 0)}</span>
          <span className="width--100">{fixed(data.Left.NSPerOp, 0)}</span>
          <span
            className={`width--100 ${
              data.Diff.Ops <= -10
                ? "negatif--Micro"
                : data.Diff.Ops >= 10
                  ? "positif--Micro"
                  : ""
            }`}
          >
            {data.Diff.NSPerOp.toFixed(2)}
          </span>
        </div>
        <div className="width--6em">
          <i className="fa-solid fa-circle-info" onClick={handleClick}></i>
        </div>
      </div>
      <div className="microbench__bottom flex--column">
        <div className="space--between">
          <span className="width--18em"></span>
          <span className="width--12em">{gitRefLeft}</span>
          <span className="width--12em">{gitRefRight}</span>
          <span className="width--12em">Diff %</span>
        </div>
        <figure className="microbench__bottom__line"></figure>
        <div className="hiddenDesktop microbenchMore">
          <span className="width--18em">Number of Iterations</span>
          <span className="width--12em">{fixed(data.Right.Ops, 0)}</span>
          <span className="width--12em">{fixed(data.Left.Ops, 0)}</span>
          <span
            className={`width--12em ${
              data.Diff.Ops <= -5
                ? "negatif--Micro"
                : data.Diff.Ops >= 5
                  ? "positif--Micro"
                  : ""
            }`}
          >
            {data.Diff.Ops.toFixed(2)}
          </span>
        </div>
        <div className=" hiddenDesktop microbenchMore">
          <span className="width--18em">Time/op</span>
          <span className="width--12em">{fixed(data.Right.NSPerOp, 0)}</span>
          <span className="width--12em">{fixed(data.Left.NSPerOp, 0)}</span>
          <span
            className={`width--12em ${
              data.Diff.Ops <= -10
                ? "negatif--Micro"
                : data.Diff.Ops >= 10
                  ? "positif--Micro"
                  : ""
            }`}
          >
            {data.Diff.NSPerOp.toFixed(2)}
          </span>
        </div>
        <div className="flex microbenchMore">
          <span className="width--18em">Bytes/op</span>
          <span className="width--12em">{data.Right.BytesPerOp}</span>
          <span className="width--12em">{data.Left.BytesPerOp}</span>
          <span
            className={`width--12em ${
              data.Diff.Ops <= -10
                ? "negatif--Micro"
                : data.Diff.Ops >= 10
                  ? "positif--Micro"
                  : ""
            }`}
          >
            {data.Diff.BytesPerOp.toFixed(2)}
          </span>
        </div>
        <div className="flex microbenchMore">
          <span className="width--18em">Megabytess/s</span>
          <span className="width--12em">{data.Right.MBPerSec}</span>
          <span className="width--12em">{data.Left.MBPerSec}</span>
          <span
            className={`width--12em ${
              data.Diff.Ops <= -10
                ? "negatif--Micro"
                : data.Diff.Ops >= 10
                  ? "positif--Micro"
                  : ""
            }`}
          >
            {data.Diff.MBPerSec.toFixed(2)}
          </span>
        </div>
        <div className="flex microbenchMore">
          <span className="width--18em">Allocations/op</span>
          <span className="width--12em">{data.Right.AllocsPerOp}</span>
          <span className="width--12em">{data.Left.AllocsPerOp}</span>
          <span
            className={`width--12em ${
              data.Diff.Ops <= -10
                ? "negatif--Micro"
                : data.Diff.Ops >= 10
                  ? "positif--Micro"
                  : ""
            }`}
          >
            {data.Diff.AllocsPerOp.toFixed(2)}
          </span>
        </div>
      </div>
    </div>
  );
};

Microbench.propTypes = {
  data: PropTypes.object.isRequired,
  className: PropTypes.string,
  gitRefLeft: PropTypes.string.isRequired,
  gitRefRight: PropTypes.string.isRequired,
};

export default Microbench;
