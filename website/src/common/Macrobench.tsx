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
import PropTypes from "prop-types";
import { Link } from "react-router-dom";
import { formatByte, fixed, secondToMicrosecond } from "../utils/Utils";
import { twMerge } from "tailwind-merge";

export default function Macrobench({ data, gitRef, commits }) {
  return (
    <div className="w-full border border-primary rounded-xl relative shadow-lg">
      <div className="p-5 flex flex-col gap-y-3">
        <h3 className="text-start text-3xl font-medium">{data.type}</h3>
        <span className="flex gap-x-1">
          Click
          <Link
            className="text-primary"
            to={`/macrobench/queries/compare?ltag=${commits.old}&rtag=${commits.new}&type=${data.type}`}
          >
            here
          </Link>
          to see the query plans comparison for this benchmark.
        </span>
      </div>

      <table>
        <thead>
          <tr>
            <th />
            <th>
              <Link
                target="_blank"
                to={`https://github.com/vitessio/vitess/commit/${commits.old}`}
                className="text-primary"
              >
                {gitRef.old || "Old"}
              </Link>
            </th>
            <th>
              <Link
                target="_blank"
                to={`https://github.com/vitessio/vitess/commit/${commits.new}`}
                className="text-primary"
              >
                {gitRef.new || "New"}
              </Link>
            </th>
            <th>
              <h4>P</h4>
            </th>
            <th>
              <h4>Delta</h4>
            </th>
            <th>
              <h4>Significant</h4>
            </th>
          </tr>
        </thead>

        <tbody>
          <Row
            title={"QPS Total"}
            oldVal={data.result.total_qps.old}
            newVal={data.result.total_qps.new}
            delta={data.result.total_qps.delta}
            insignificant={data.result.total_qps.insignificant}
            p={fixed(data.result.total_qps.p, 3)}
            fmt={"none"}
          />
          <Row
            title="QPS Reads"
            oldVal={data.result.reads_qps.old}
            newVal={data.result.reads_qps.new}
            delta={data.result.reads_qps.delta}
            insignificant={data.result.reads_qps.insignificant}
            p={fixed(data.result.reads_qps.p, 3)}
            fmt={"none"}
          />

          <Row
            title="QPS Writes"
            oldVal={data.result.writes_qps.old}
            newVal={data.result.writes_qps.new}
            delta={data.result.writes_qps.delta}
            insignificant={data.result.writes_qps.insignificant}
            p={fixed(data.result.writes_qps.p, 3)}
            fmt={"none"}
          />

          <Row
            title="QPS Other"
            oldVal={data.result.other_qps.old}
            newVal={data.result.other_qps.new}
            delta={data.result.other_qps.delta}
            insignificant={data.result.other_qps.insignificant}
            p={fixed(data.result.other_qps.p, 3)}
            fmt={"none"}
          />

          <Row
            title="TPS"
            oldVal={data.result.tps.old}
            newVal={data.result.tps.new}
            delta={data.result.tps.delta}
            insignificant={data.result.tps.insignificant}
            p={fixed(data.result.tps.p, 3)}
            fmt={"none"}
          />

          <Row
            title="Latency"
            oldVal={data.result.latency.old}
            newVal={data.result.latency.new}
            delta={data.result.latency.delta}
            insignificant={data.result.latency.insignificant}
            p={fixed(data.result.latency.p, 3)}
            fmt={"none"}
          />

          <Row
            title="Errors"
            oldVal={data.result.errors.old}
            newVal={data.result.errors.new}
            delta={data.result.errors.delta}
            insignificant={data.result.errors.insignificant}
            p={fixed(data.result.errors.p, 3)}
            fmt={"none"}
          />

          <Row
            title="Total CPU / query"
            oldVal={data.result.total_components_cpu_time.old}
            newVal={data.result.total_components_cpu_time.new}
            delta={data.result.total_components_cpu_time.delta}
            insignificant={data.result.total_components_cpu_time.insignificant}
            p={fixed(data.result.total_components_cpu_time.p, 3)}
            fmt={"time"}
          />

          <Row
            title="CPU / query (vtgate)"
            oldVal={data.result.components_cpu_time.vtgate.old}
            newVal={data.result.components_cpu_time.vtgate.new}
            delta={data.result.components_cpu_time.vtgate.delta}
            insignificant={data.result.components_cpu_time.vtgate.insignificant}
            p={fixed(data.result.components_cpu_time.vtgate.p, 3)}
            fmt={"time"}
          />

          <Row
            title="CPU / query (vttablet)"
            oldVal={data.result.components_cpu_time.vttablet.old}
            newVal={data.result.components_cpu_time.vttablet.new}
            delta={data.result.components_cpu_time.vttablet.delta}
            insignificant={
              data.result.components_cpu_time.vttablet.insignificant
            }
            p={fixed(data.result.components_cpu_time.vttablet.p, 3)}
            fmt={"time"}
          />

          <Row
            title="Total Allocated / query"
            oldVal={data.result.total_components_mem_stats_alloc_bytes.old}
            newVal={data.result.total_components_mem_stats_alloc_bytes.new}
            delta={data.result.total_components_mem_stats_alloc_bytes.delta}
            insignificant={
              data.result.total_components_mem_stats_alloc_bytes.insignificant
            }
            p={fixed(data.result.total_components_mem_stats_alloc_bytes.p, 3)}
            fmt={"memory"}
          />

          <Row
            title="Allocated / query (vtgate)"
            oldVal={data.result.components_mem_stats_alloc_bytes.vtgate.old}
            newVal={data.result.components_mem_stats_alloc_bytes.vtgate.new}
            delta={data.result.components_mem_stats_alloc_bytes.vtgate.delta}
            insignificant={
              data.result.components_mem_stats_alloc_bytes.vtgate.insignificant
            }
            p={fixed(data.result.components_mem_stats_alloc_bytes.vtgate.p, 3)}
            fmt={"memory"}
          />

          <Row
            title="Allocated / query (vttablet)"
            oldVal={data.result.components_mem_stats_alloc_bytes.vttablet.old}
            newVal={data.result.components_mem_stats_alloc_bytes.vttablet.new}
            delta={data.result.components_mem_stats_alloc_bytes.vttablet.delta}
            insignificant={
              data.result.components_mem_stats_alloc_bytes.vttablet
                .insignificant
            }
            p={fixed(
              data.result.components_mem_stats_alloc_bytes.vttablet.p,
              3,
            )}
            fmt={"memory"}
          />
        </tbody>
      </table>
    </div>
  );
}

export function getRange(range) {
  if (range.infinite == true) {
    return "∞";
  }
  if (range.unknown == true) {
    return "?";
  }
  return "±" + fixed(range.value, 1) + "%";
}

function Row({ title, oldVal, newVal, delta, insignificant, p, fmt }) {
  let status = (
    <span
      className={twMerge(
        "text-lg text-white px-4 rounded-full",
        insignificant == true && "bg-[#dd1a2a]",
        insignificant == false && "bg-[#00aa00]",
      )}
    >
      {insignificant == true ? "No" : "Yes"}
    </span>
  );

  var oldValFmt = oldVal.center;
  var newValFmt = newVal.center;
  if (fmt == "time") {
    oldValFmt = secondToMicrosecond(oldVal.center);
    newValFmt = secondToMicrosecond(newVal.center);
  } else if (fmt == "memory") {
    oldValFmt = formatByte(oldVal.center);
    newValFmt = formatByte(newVal.center);
  }

  return (
    <tr className="border-t border-front border-opacity-70 duration-150 hover:bg-accent">
      <td className="flex pt-4 pb-2 px-4 justify-end border-r border-r-primary font-semibold text-end">
        <span>{title}</span>
      </td>
      <td className="px-24 pt-4 pb-2 text-center">
        <span>
          {oldValFmt} ({getRange(oldVal.range)})
        </span>
      </td>
      <td className="px-24 pt-4 pb-2 text-center">
        <span>
          {newValFmt} ({getRange(newVal.range)})
        </span>
      </td>
      <td className="px-24 pt-4 pb-2 text-center">
        <span>{p || "?"}</span>
      </td>
      <td className="px-24 pt-4 pb-2 text-center">
        <span>{fixed(delta, 3) || 0}%</span>
      </td>
      <td className="px-24 pt-4 pb-2 text-center">{status}</td>
    </tr>
  );
}

Row.propTypes = {
  title: PropTypes.string.isRequired,
  oldVal: PropTypes.shape({
    center: PropTypes.number.isRequired,
    range: PropTypes.shape({
      infinite: PropTypes.bool.isRequired,
      unknown: PropTypes.bool.isRequired,
      value: PropTypes.number.isRequired,
    }),
  }).isRequired,
  newVal: PropTypes.shape({
    center: PropTypes.number.isRequired,
    range: PropTypes.shape({
      infinite: PropTypes.bool.isRequired,
      unknown: PropTypes.bool.isRequired,
      value: PropTypes.number.isRequired,
    }),
  }).isRequired,
  delta: PropTypes.number.isRequired,
  insignificant: PropTypes.bool.isRequired,
  p: PropTypes.string.isRequired,
  fmt: PropTypes.oneOf(["none", "time", "memory"]).isRequired,
};
