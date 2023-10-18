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
import { formatByteForGB, fixed } from "../utils/Utils";

export default function Macrobench({ data, gitRef, commits }) {
  return (
    <div className="w-full border border-primary rounded-xl relative shadow-lg">
      <div className="p-5 flex flex-col gap-y-3">
        <h3 className="text-start text-3xl font-medium">{data.type}</h3>
        <span className="flex gap-x-1">
          Click
          <Link
            className="text-primary"
            to={`/macrobench/queries/compare?ltag=${commits.left}&rtag=${commits.right}&type=${data.type}`}
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
                to={`https://github.com/vitessio/vitess/commit/${commits.left}`}
                className="text-primary"
              >
                {gitRef.left || "Left"}
              </Link>
            </th>
            <th>
              <Link
                target="_blank"
                to={`https://github.com/vitessio/vitess/commit/${commits.right}`}
                className="text-primary"
              >
                {gitRef.right || "Right"}
              </Link>
            </th>
            <th>
              <h4>Impoved by %</h4>
            </th>
          </tr>
        </thead>

        <tbody>
          <Row
            title={"QPS Total"}
            left={fixed(data.diff.Left.Result.qps.total, 2)}
            right={fixed(data.diff.Right.Result.qps.total, 2)}
            diffMetric={fixed(data.diff.Diff.qps.total, 2)}
          />
          <Row
            title="QPS Reads"
            left={fixed(data.diff.Left.Result.qps.reads, 2)}
            right={fixed(data.diff.Right.Result.qps.reads, 2)}
            diffMetric={fixed(data.diff.Diff.qps.reads, 2)}
          />

          <Row
            title="QPS Writes"
            left={fixed(data.diff.Left.Result.qps.writes, 2)}
            right={fixed(data.diff.Right.Result.qps.writes, 2)}
            diffMetric={fixed(data.diff.Diff.qps.writes, 2)}
          />

          <Row
            title="QPS Other"
            left={fixed(data.diff.Left.Result.qps.other, 2)}
            right={fixed(data.diff.Right.Result.qps.other, 2)}
            diffMetric={fixed(data.diff.Diff.qps.other, 2)}
          />

          <Row
            title="TPS"
            left={fixed(data.diff.Left.Result.tps, 2)}
            right={fixed(data.diff.Right.Result.tps, 2)}
            diffMetric={fixed(data.diff.Diff.tps, 2)}
          />

          <Row
            title="Latency"
            left={fixed(data.diff.Left.Result.latency, 2)}
            right={fixed(data.diff.Right.Result.latency, 2)}
            diffMetric={fixed(data.diff.Diff.latency, 2)}
          />

          <Row
            title="Errors"
            left={fixed(data.diff.Left.Result.errors, 2)}
            right={fixed(data.diff.Right.Result.errors, 2)}
            diffMetric={fixed(data.diff.Diff.errors, 2)}
          />

          <Row
            title="Reconnects"
            left={fixed(data.diff.Left.Result.reconnects, 2)}
            right={fixed(data.diff.Right.Result.reconnects, 2)}
            diffMetric={fixed(data.diff.Diff.reconnects, 2)}
          />

          <Row
            title="Time"
            left={fixed(data.diff.Left.Result.time, 2)}
            right={fixed(data.diff.Right.Result.time, 2)}
            diffMetric={fixed(data.diff.Diff.time, 2)}
          />

          <Row
            title="Threads"
            left={fixed(data.diff.Left.Result.threads, 2)}
            right={fixed(data.diff.Right.Result.threads, 2)}
            diffMetric={fixed(data.diff.Diff.threads, 2)}
          />

          <Row
            title="Total CPU time"
            left={fixed(data.diff.Left.Metrics.TotalComponentsCPUTime, 2)}
            right={fixed(data.diff.Right.Metrics.TotalComponentsCPUTime, 2)}
            diffMetric={fixed(data.diff.DiffMetrics.TotalComponentsCPUTime, 2)}
          />

          <Row
            title="CPU time vtgate"
            left={fixed(data.diff.Left.Metrics.ComponentsCPUTime.vtgate, 2)}
            right={fixed(data.diff.Right.Metrics.ComponentsCPUTime.vtgate, 2)}
            diffMetric={fixed(
              data.diff.DiffMetrics.ComponentsCPUTime.vtgate,
              2
            )}
          />

          <Row
            title="CPU time vttablet"
            left={fixed(data.diff.Left.Metrics.ComponentsCPUTime.vttablet, 2)}
            right={fixed(data.diff.Right.Metrics.ComponentsCPUTime.vttablet, 2)}
            diffMetric={fixed(
              data.diff.DiffMetrics.ComponentsCPUTime.vttablet,
              2
            )}
          />

          <Row
            title="Total Allocs bytes"
            left={formatByteForGB(
              data.diff.Left.Metrics.TotalComponentsMemStatsAllocBytes
            )}
            right={formatByteForGB(
              data.diff.Right.Metrics.TotalComponentsMemStatsAllocBytes
            )}
            diffMetric={fixed(
              data.diff.DiffMetrics.TotalComponentsMemStatsAllocBytes,
              2
            )}
          />

          <Row
            title="Allocs bytes vtgate"
            left={formatByteForGB(
              data.diff.Left.Metrics.ComponentsMemStatsAllocBytes.vtgate
            )}
            right={formatByteForGB(
              data.diff.Right.Metrics.ComponentsMemStatsAllocBytes.vtgate
            )}
            diffMetric={fixed(
              data.diff.DiffMetrics.ComponentsMemStatsAllocBytes.vtgate,
              2
            )}
          />

          <Row
            title="Allocs bytes vttablet"
            left={formatByteForGB(
              data.diff.Left.Metrics.ComponentsMemStatsAllocBytes.vttablet
            )}
            right={formatByteForGB(
              data.diff.Right.Metrics.ComponentsMemStatsAllocBytes.vttablet
            )}
            diffMetric={fixed(
              data.diff.DiffMetrics.ComponentsMemStatsAllocBytes.vttablet,
              2
            )}
          />
        </tbody>
      </table>
    </div>
  );
}

function Row({ title, left, right, diffMetric }) {
  return (
    <tr className="border-t border-front border-opacity-70 duration-150 hover:bg-foreground hover:bg-opacity-20">
      <td className="flex pt-4 pb-2 px-4 justify-end border-r border-r-primary font-semibold">
        <span>{title}</span>
      </td>
      <td className="px-24 pt-4 pb-2 text-center">
        <span>{left || 0}</span>
      </td>
      <td className="px-24 pt-4 pb-2 text-center">
        <span>{right || 0}</span>
      </td>
      <td className="px-24 pt-4 pb-2 text-center">
        <span>{diffMetric || 0}</span>
      </td>
    </tr>
  );
}
