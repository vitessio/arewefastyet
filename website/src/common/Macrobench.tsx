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

import { Link } from "react-router-dom";
import { formatByteForGB, fixed } from "../utils";
import { MacrobenchComparison } from "../types";

interface MacrobenchProps {
  data: MacrobenchComparison;
  gitRef: { left: string; right: string };
  commits: { left: string; right: string };
}

export default function Macrobench(props: MacrobenchProps) {
  const { data, gitRef, commits } = props;
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

      <table className="w-full">
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
            left={fixed(data.result.total_qps.old.center, 2)}
            right={fixed(data.result.total_qps.new.center, 2)}
            diffMetric={fixed(data.result.total_qps.delta, 2)}
          />
          <Row
            title="QPS Reads"
            left={fixed(data.result.reads_qps.old.center, 2)}
            right={fixed(data.result.reads_qps.new.center, 2)}
            diffMetric={fixed(data.result.reads_qps.delta, 2)}
          />

          <Row
            title="QPS Writes"
            left={fixed(data.result.writes_qps.old.center, 2)}
            right={fixed(data.result.writes_qps.new.center, 2)}
            diffMetric={fixed(data.result.writes_qps.delta, 2)}
          />

          <Row
            title="QPS Other"
            left={fixed(data.result.other_qps.old.center, 2)}
            right={fixed(data.result.other_qps.new.center, 2)}
            diffMetric={fixed(data.result.other_qps.delta, 2)}
          />

          <Row
            title="TPS"
            left={fixed(data.result.tps.old.center, 2)}
            right={fixed(data.result.tps.new.center, 2)}
            diffMetric={fixed(data.result.tps.delta, 2)}
          />

          <Row
            title="Latency"
            left={fixed(data.result.latency.old.center, 2)}
            right={fixed(data.result.latency.new.center, 2)}
            diffMetric={fixed(data.result.latency.delta, 2)}
          />

          <Row
            title="Errors"
            left={fixed(data.result.errors.old.center, 2)}
            right={fixed(data.result.errors.new.center, 2)}
            diffMetric={fixed(data.result.errors.delta, 2)}
          />

          {/* <Row
            title="Reconnects"
            left={fixed(data.result., 2)}
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
          /> */}

          <Row
            title="Total CPU time"
            left={fixed(data.result.total_components_cpu_time.old.center, 2)}
            right={fixed(data.result.total_components_cpu_time.new.center, 2)}
            diffMetric={fixed(data.result.total_components_cpu_time.delta, 2)}
          />

          <Row
            title="CPU time vtgate"
            left={fixed(data.result.components_cpu_time.vtgate.old.center, 2)}
            right={fixed(data.result.components_cpu_time.vtgate.new.center, 2)}
            diffMetric={fixed(data.result.components_cpu_time.vtgate.delta, 2)}
          />

          <Row
            title="CPU time vttablet"
            left={fixed(data.result.components_cpu_time.vttablet.old.center, 2)}
            right={fixed(
              data.result.components_cpu_time.vttablet.new.center,
              2
            )}
            diffMetric={fixed(
              data.result.components_cpu_time.vttablet.delta,
              2
            )}
          />

          <Row
            title="Total Allocs bytes"
            left={formatByteForGB(
              data.result.total_components_mem_stats_alloc_bytes.old.center
            )}
            right={formatByteForGB(
              data.result.total_components_mem_stats_alloc_bytes.new.center
            )}
            diffMetric={fixed(
              data.result.total_components_mem_stats_alloc_bytes.old.center,
              2
            )}
          />

          <Row
            title="Allocs bytes vtgate"
            left={formatByteForGB(
              data.result.components_mem_stats_alloc_bytes.vtgate.old.center
            )}
            right={formatByteForGB(
              data.result.components_mem_stats_alloc_bytes.vtgate.new.center
            )}
            diffMetric={fixed(
              data.result.components_mem_stats_alloc_bytes.vtgate.delta,
              2
            )}
          />

          <Row
            title="Allocs bytes vttablet"
            left={formatByteForGB(
              data.result.components_mem_stats_alloc_bytes.vttablet.old.center
            )}
            right={formatByteForGB(
              data.result.components_mem_stats_alloc_bytes.vttablet.new.center
            )}
            diffMetric={fixed(
              data.result.components_mem_stats_alloc_bytes.vtgate.delta,
              2
            )}
          />
        </tbody>
      </table>
    </div>
  );
}

interface RowProps {
  title: string;
  left?: string | number;
  right?: string | number;
  diffMetric?: string | number;
}

function Row(props: RowProps) {
  return (
    <tr className="border-t border-front border-opacity-70 duration-150 hover:bg-foreground hover:bg-opacity-20">
      <td className="flex pt-4 pb-2 px-4 justify-end border-r border-r-primary font-semibold text-end">
        <span>{props.title}</span>
      </td>
      <td className="px-24 pt-4 pb-2 text-center">
        <span>{props.left || 0}</span>
      </td>
      <td className="px-24 pt-4 pb-2 text-center">
        <span>{props.right || 0}</span>
      </td>
      <td className="px-24 pt-4 pb-2 text-center">
        <span>{props.diffMetric || 0}</span>
      </td>
    </tr>
  );
}
