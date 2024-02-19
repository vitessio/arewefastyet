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
              return fixed(value.Result.qps.total, 2);
            }}
          />

          <Row
            title="QPS Reads"
            data={data}
            extract={function ({ value }) {
              return fixed(value.Result.qps.reads, 2);
            }}
          />

          <Row
            title="QPS Writes"
            data={data}
            extract={function ({ value }) {
              return fixed(value.Result.qps.writes, 2);
            }}
          />

          <Row
            title="QPS Other"
            data={data}
            extract={function ({ value }) {
              return fixed(value.Result.qps.other, 2);
            }}
          />

          <Row
            title="TPS"
            data={data}
            extract={function ({ value }) {
              return fixed(value.Result.tps, 2);
            }}
          />

          <Row
            title="Latency"
            data={data}
            extract={function ({ value }) {
              return fixed(value.Result.latency, 2);
            }}
          />

          <Row
            title="Errors"
            data={data}
            extract={function ({ value }) {
              return fixed(value.Result.errors, 2);
            }}
          />

          <Row
            title="Reconnects"
            data={data}
            extract={function ({ value }) {
              return fixed(value.Result.reconnects, 2);
            }}
          />

          <Row
            title="Time"
            data={data}
            extract={function ({ value }) {
              return fixed(value.Result.time, 2);
            }}
          />

          <Row
            title="Threads"
            data={data}
            extract={function ({ value }) {
              return fixed(value.Result.threads, 2);
            }}
          />

          <Row
            title="Total CPU / query"
            data={data}
            extract={function ({ value }) {
              return secondToMicrosecond(value.Metrics.TotalComponentsCPUTime);
            }}
          />

          <Row
            title="CPU / query (vtgate)"
            data={data}
            extract={function ({ value }) {
              return secondToMicrosecond(
                value.Metrics.ComponentsCPUTime.vtgate
              );
            }}
          />

          <Row
            title="CPU / query (vttablet)"
            data={data}
            extract={function ({ value }) {
              return secondToMicrosecond(
                value.Metrics.ComponentsCPUTime.vttablet
              );
            }}
          />

          <Row
            title="Total Allocated / query"
            data={data}
            extract={function ({ value }) {
              return formatByte(
                value.Metrics.TotalComponentsMemStatsAllocBytes
              );
            }}
          />

          <Row
            title="Allocated / query (vtgate)"
            data={data}
            extract={function ({ value }) {
              return formatByte(
                value.Metrics.ComponentsMemStatsAllocBytes.vtgate
              );
            }}
          />

          <Row
            title="Allocated / query (vttablet)"
            data={data}
            extract={function ({ value }) {
              return formatByte(
                value.Metrics.ComponentsMemStatsAllocBytes.vttablet
              );
            }}
          />
        </tbody>
      </table>
    </div>
  );
}

function Row({ title, data, extract }) {
  return (
    <tr className="border-t border-front border-opacity-70 duration-150 hover:bg-accent">
      <td className="flex pt-4 pb-2 px-4 justify-end border-r border-r-primary font-semibold text-end">
        <span>{title}</span>
      </td>
      {Object.entries(data).map(([key, value]) => {
        return (
          <td className="px-24 pt-4 pb-2 text-center">
            <span>{extract({ value: value }) || 0}</span>
          </td>
        );
      })}
    </tr>
  );
}
