/*
Copyright 2024 The Vitess Authors.

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
import { formatByte, fixed, secondToMicrosecond } from "../../../utils/Utils";
import { getRange } from "@/common/Macrobench";
import PropTypes from "prop-types";

export default function FK({ data }) {
  return (
    <div className="w-full border border-primary rounded-xl relative shadow-lg">
      <div className="p-5 flex flex-col gap-y-3"></div>
      <table>
        <thead>
          <tr>
            <th />
            {Object.entries(data).map(([key, value]) => {
              return <th>{key}</th>;
            })}
          </tr>
        </thead>

        <tbody>
          <Row
            title={"QPS Total"}
            data={data}
            extract={function ({ value }) {
              return value.total_qps;
            }}
          />

          <Row
            title="QPS Reads"
            data={data}
            extract={function ({ value }) {
              return value.reads_qps;
            }}
          />

          <Row
            title="QPS Writes"
            data={data}
            extract={function ({ value }) {
              return value.writes_qps;
            }}
          />

          <Row
            title="QPS Other"
            data={data}
            extract={function ({ value }) {
              return value.other_qps;
            }}
          />

          <Row
            title="TPS"
            data={data}
            extract={function ({ value }) {
              return value.tps;
            }}
          />

          <Row
            title="Latency"
            data={data}
            extract={function ({ value }) {
              return value.latency;
            }}
          />

          <Row
            title="Errors"
            data={data}
            extract={function ({ value }) {
              return value.errors;
            }}
          />

          <Row
            title="Total CPU / query"
            fmt={"time"}
            data={data}
            extract={function ({ value }) {
              return value.total_components_cpu_time;
            }}
          />

          <Row
            title="CPU / query (vtgate)"
            fmt={"time"}
            data={data}
            extract={function ({ value }) {
              return value.components_cpu_time.vtgate;
            }}
          />

          <Row
            title="CPU / query (vttablet)"
            fmt={"time"}
            data={data}
            extract={function ({ value }) {
              return value.components_cpu_time.vttablet;
            }}
          />

          <Row
            title="Total Allocated / query"
            fmt={"memory"}
            data={data}
            extract={function ({ value }) {
              return value.total_components_mem_stats_alloc_bytes;
            }}
          />

          <Row
            title="Allocated / query (vtgate)"
            fmt={"memory"}
            data={data}
            extract={function ({ value }) {
              return value.components_mem_stats_alloc_bytes.vtgate;
            }}
          />

          <Row
            title="Allocated / query (vttablet)"
            fmt={"memory"}
            data={data}
            extract={function ({ value }) {
              return value.components_mem_stats_alloc_bytes.vttablet;
            }}
          />
        </tbody>
      </table>
    </div>
  );
}

function fmtString(value, fmt) {
  var valFmt = value.center;
  if (fmt == "time") {
    valFmt = secondToMicrosecond(value.center);
  } else if (fmt == "memory") {
    valFmt = formatByte(value.center);
  }
  return valFmt;
}

function Row({ title, data, extract, fmt }) {
  return (
    <tr className="border-t border-front border-opacity-70 duration-150 hover:bg-accent">
      <td className="flex pt-4 pb-2 px-4 justify-end border-r border-r-primary font-semibold text-end">
        <span>{title}</span>
      </td>
      {Object.entries(data).map(([key, value]) => {
        var val = extract({ value: value });
        return (
          <td className="px-24 pt-4 pb-2 text-center">
            <span>
              {fmtString(val, fmt)} ({getRange(val.range)})
            </span>
          </td>
        );
      })}
    </tr>
  );
}

Row.propTypes = {
  title: PropTypes.string.isRequired,
  extract: PropTypes.func.isRequired,
  data: PropTypes.object.isRequired,
  fmt: PropTypes.oneOf(["time", "memory"]),
};
