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
import { Link } from "react-router-dom";
import { formatByteForGB, fixed } from "../../../utils/Utils";

export default function FK({ data }) {
  return (
    <div className="w-full border border-primary rounded-xl relative shadow-lg">
      <div className="p-5 flex flex-col gap-y-3">
        <h3 className="text-start text-3xl font-medium">FKs benchmark Comparison</h3>
      </div>

      <table>
        <thead>
          <tr>
            <th />
            {Object.entries(data).map(([key, value]) => {
              return (
                  <th>
                    {key}
                  </th>
              )
            })}
          </tr>
        </thead>

        <tbody>
          <Row
            title={"QPS Total"}
            data={data}
            extract={function ({value}) {
              return value.Result.qps.total
            }}
          />

          <Row
            title="QPS Reads"
            data={data}
            extract={function ({value}) {
              return value.Result.qps.reads
            }}
          />

          <Row
            title="QPS Writes"
            data={data}
            extract={function ({value}) {
              return value.Result.qps.writes
            }}
          />

          <Row
            title="QPS Other"
            data={data}
            extract={function ({value}) {
              return value.Result.qps.other
            }}
          />

          <Row
            title="TPS"
            data={data}
            extract={function ({value}) {
              return value.Result.tps
            }}
          />

          <Row
            title="Latency"
            data={data}
            extract={function ({value}) {
              return value.Result.latency
            }}
          />

          <Row
            title="Errors"
            data={data}
            extract={function ({value}) {
              return value.Result.errors
            }}
          />

          <Row
            title="Reconnects"
            data={data}
            extract={function ({value}) {
              return value.Result.reconnects
            }}
          />

          <Row
            title="Time"
            data={data}
            extract={function ({value}) {
              return value.Result.time
            }}
          />

          <Row
            title="Threads"
            data={data}
            extract={function ({value}) {
              return value.Result.threads
            }}
          />

          <Row
            title="Total CPU time"
            data={data}
            extract={function ({value}) {
              return value.Metrics.TotalComponentsCPUTime
            }}
          />

          <Row
            title="CPU time vtgate"
            data={data}
            extract={function ({value}) {
              return value.Metrics.ComponentsCPUTime.vtgate
            }}
          />

          <Row
            title="CPU time vttablet"
            data={data}
            extract={function ({value}) {
              return value.Metrics.ComponentsCPUTime.vttablet
            }}
          />

          <Row
            title="Total Allocs bytes"
            data={data}
            extract={function ({value}) {
              return formatByteForGB(value.Metrics.TotalComponentsMemStatsAllocBytes)
            }}
          />

          <Row
            title="Allocs bytes vtgate"
            data={data}
            extract={function ({value}) {
              return formatByteForGB(value.Metrics.ComponentsMemStatsAllocBytes.vtgate)
            }}
          />

          <Row
            title="Allocs bytes vttablet"
            data={data}
            extract={function ({value}) {
              return formatByteForGB(value.Metrics.ComponentsMemStatsAllocBytes.vttablet)
            }}
          />
        </tbody>
      </table>
    </div>
  );
}

function Row({ title, data, extract }) {
  return (
    <tr className="border-t border-front border-opacity-70 duration-150 hover:bg-foreground hover:bg-opacity-20">
      <td className="flex pt-4 pb-2 px-4 justify-end border-r border-r-primary font-semibold text-end">
        <span>{title}</span>
      </td>
      {Object.entries(data).map(([key, value]) => {
        return (
            <td className="px-24 pt-4 pb-2 text-center">
              <span>{extract({value: value}) || 0}</span>
            </td>
        )
      })}
    </tr>
  );
}
