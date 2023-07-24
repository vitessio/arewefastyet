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

import React from "react";

import "../SearchMacroDesktop/searchMacroDesktop.css";

import { formatByteForGB, fixed } from "../../../../utils/Utils";

const SearchMacroDesktop = ({ data }) => {
  const renderZeroSpans = Array(16)
    .fill(0)
    .map((_, index) => <span key={index}>0</span>);

  return (
    <div className="flex--column searchMAcro__desktop__data">
      {data[1] ? (
        <>
          <span>{fixed(data[1][0].Result.qps.total, 2)}</span>
          <span>{data[1][0].Result.qps.reads}</span>
          <span>{fixed(data[1][0].Result.qps.writes, 2)}</span>
          <span>{data[1][0].Result.qps.other.toFixed(2)}</span>
          <span>{fixed(data[1][0].Result.tps, 2)}</span>
          <span>{data[1][0].Result.latency}</span>
          <span>{fixed(data[1][0].Result.errors, 2)}</span>
          <span>{data[1][0].Result.reconnects}</span>
          <span>{data[1][0].Result.time}</span>
          <span>{data[1][0].Result.threads}</span>
          <span>{data[1][0].Metrics.TotalComponentsCPUTime.toFixed(0)}</span>
          <span>{fixed(data[1][0].Metrics.ComponentsCPUTime.vtgate, 2)}</span>
          <span>{fixed(data[1][0].Metrics.ComponentsCPUTime.vttablet, 2)}</span>
          <span>{data[1][0].Metrics.TotalComponentsMemStatsAllocBytes}</span>
          <span>
            {fixed(data[1][0].Metrics.ComponentsMemStatsAllocBytes.vtgate, 2)}
          </span>
          <span>
            {fixed(data[1][0].Metrics.ComponentsMemStatsAllocBytes.vttablet, 2)}
          </span>
        </>
      ) : (
        renderZeroSpans
      )}
    </div>
  );
};

export default SearchMacroDesktop;
