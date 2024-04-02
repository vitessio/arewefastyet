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

import { formatByte, secondToMicrosecond } from "../../../utils/Utils";
import { Link } from "react-router-dom";
import {getRange} from "@/common/Macrobench";
import PropTypes from "prop-types";

export default function SearchMacro({ data, gitRef }) {
  return (
    <div className="flex flex-col border border-primary relative rounded-xl bg-background bg-opacity-5 shadow-xl">
      <div className="p-5">
        <h3 className="text-xl font-semibold">{data[0]}</h3>
        <Link
          target="_blank"
          className="text-primary"
          to={`https://github.com/vitessio/vitess/commit/${gitRef}`}
        >
          {gitRef}
        </Link>
      </div>
      <table>
        <tbody>
          <Row
            title={"QPS Total"}
            value={data[1].total_qps}
          />

          <Row
            title={"QPS Reads"}
            value={data[1].reads_qps}
          />

          <Row
            title={"QPS Writes"}
            value={data[1].writes_qps}
          />

          <Row
            title={"QPS Other"}
            value={data[1].other_qps}
          />

          <Row
            title={"TPS"}
            value={data[1].tps}
          />

          <Row
            title={"Latency"}
            value={data[1].latency}
          />

          <Row
            title={"Errors"}
            value={data[1].errors}
          />

          <Row
            title={"Total CPU / query"}
            value={data[1].total_components_cpu_time}
            fmt={"time"}
          />

          <Row
            title={"CPU / query (vtgate)"}
            value={data[1].components_cpu_time.vtgate}
            fmt={"time"}
          />

          <Row
            title={"CPU / query (vttablet)"}
            value={data[1].components_cpu_time.vttablet}
            fmt={"time"}
          />

          <Row
            title={"Total Allocated / query"}
            value={data[1].total_components_mem_stats_alloc_bytes}
            fmt={"memory"}
          />

          <Row
            title={"Allocated / query (vtgate)"}
            value={data[1].components_mem_stats_alloc_bytes.vtgate}
            fmt={"memory"}
          />

          <Row
            title={"Allocated / query (vttablet)"}
            value={data[1].components_mem_stats_alloc_bytes.vttablet}
            fmt={"memory"}
          />
        </tbody>
      </table>
    </div>
  );
}

function Row({ title, value, fmt }) {
    var valFmt = value.center
    if (fmt == "time") {
        valFmt = secondToMicrosecond(value.center)
    } else if (fmt == "memory") {
        valFmt = formatByte(value.center)
    }

    return (
    <tr className="border-t border-front border-opacity-70 duration-150 hover:bg-accent">
      <td className="flex pt-4 pb-1 px-4 text-lg justify-end border-r border-r-primary font-bold">
        <span>{title}</span>
      </td>
      <td className="px-24 pt-4 pb-1 text-center">
        <span>{valFmt} ({getRange(value.range)})</span>
      </td>
    </tr>
    );
}

Row.propTypes = {
    title: PropTypes.string.isRequired,
    value: PropTypes.shape({
        center: PropTypes.number.isRequired,
        range: PropTypes.shape({
            infinite: PropTypes.bool.isRequired,
            unknown: PropTypes.bool.isRequired,
            value: PropTypes.number.isRequired,
        }),
    }).isRequired,
    fmt: PropTypes.oneOf(['time', 'memory']),
};