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

import { fixed, formatByte, secondToMicrosecond } from "../../../utils/Utils";
import { Link } from "react-router-dom";

export default function SearchMacro({ data, gitRef }) {
  return (
    <div className="flex flex-col border border-primary relative rounded-xl bg-foreground bg-opacity-5 shadow-xl">
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
            value={data[1] ? fixed(data[1][0].Result.qps.total, 2) : 0}
          />

          <Row
            title={"QPS Reads"}
            value={data[1] ? fixed(data[1][0].Result.qps.reads, 2) : 0}
          />

          <Row
            title={"QPS Writes"}
            value={data[1] ? fixed(data[1][0].Result.qps.writes, 2) : 0}
          />

          <Row
            title={"QPS Other"}
            value={data[1] ? fixed(data[1][0].Result.qps.other, 2) : 0}
          />

          <Row
            title={"TPS"}
            value={data[1] ? fixed(data[1][0].Result.tps, 2) : 0}
          />

          <Row
            title={"Latency"}
            value={data[1] ? fixed(data[1][0].Result.latency, 2) : 0}
          />

          <Row
            title={"Errors"}
            value={data[1] ? fixed(data[1][0].Result.errors, 2) : 0}
          />

          <Row
            title={"Reconnects"}
            value={data[1] ? fixed(data[1][0].Result.reconnects, 2) : 0}
          />

          <Row
            title={"Time"}
            value={data[1] ? fixed(data[1][0].Result.time, 2) : 0}
          />

          <Row
            title={"Threads"}
            value={data[1] ? fixed(data[1][0].Result.threads, 2) : 0}
          />

          <Row
            title={"Total CPU / query"}
            value={
              data[1] ? secondToMicrosecond(data[1][0].Metrics.TotalComponentsCPUTime) : 0
            }
          />

          <Row
            title={"CPU / query (vtgate)"}
            value={
              data[1]
                ? secondToMicrosecond(data[1][0].Metrics.ComponentsCPUTime.vtgate)
                : 0
            }
          />

          <Row
            title={"CPU / query (vttablet)"}
            value={
              data[1]
                ? secondToMicrosecond(data[1][0].Metrics.ComponentsCPUTime.vttablet)
                : 0
            }
          />

          <Row
            title={"Total Allocated / query"}
            value={
              data[1]
                ? formatByte(
                    data[1][0].Metrics.TotalComponentsMemStatsAllocBytes
                  )
                : 0
            }
          />

          <Row
            title={"Allocated / query (vtgate)"}
            value={
              data[1]
                ? formatByte(
                    data[1][0].Metrics.ComponentsMemStatsAllocBytes.vtgate
                  )
                : 0
            }
          />

          <Row
            title={"Allocated / query (vttablet)"}
            value={
              data[1]
                ? formatByte(
                    data[1][0].Metrics.ComponentsMemStatsAllocBytes.vttablet
                  )
                : 0
            }
          />
        </tbody>
      </table>
    </div>
  );
}

function Row({ title, value }) {
  return (
    <tr className="border-t border-front border-opacity-70 duration-150 hover:bg-foreground hover:bg-opacity-20">
      <td className="flex pt-4 pb-1 px-4 text-lg justify-end border-r border-r-primary font-bold">
        <span>{title}</span>
      </td>
      <td className="px-24 pt-4 pb-1 text-center">
        <span>{value || 0}</span>
      </td>
    </tr>
  );
}

SearchMacro.propTypes = {
  data: PropTypes.arrayOf(
    PropTypes.oneOfType([
      PropTypes.string,
      PropTypes.arrayOf(
        PropTypes.shape({
          Result: PropTypes.shape({
            qps: PropTypes.shape({
              total: PropTypes.number.isRequired,
              reads: PropTypes.number.isRequired,
              writes: PropTypes.number.isRequired,
              other: PropTypes.number.isRequired,
            }).isRequired,
            tps: PropTypes.number.isRequired,
            latency: PropTypes.number.isRequired,
            errors: PropTypes.number.isRequired,
            reconnects: PropTypes.number.isRequired,
            time: PropTypes.number.isRequired,
            threads: PropTypes.number.isRequired,
          }).isRequired,
          Metrics: PropTypes.shape({
            TotalComponentsCPUTime: PropTypes.number.isRequired,
            ComponentsCPUTime: PropTypes.shape({
              vtgate: PropTypes.number.isRequired,
              vttablet: PropTypes.number.isRequired,
            }).isRequired,
            TotalComponentsMemStatsAllocBytes: PropTypes.number.isRequired,
            ComponentsMemStatsAllocBytes: PropTypes.shape({
              vtgate: PropTypes.number.isRequired,
              vttablet: PropTypes.number.isRequired,
            }).isRequired,
          }).isRequired,
        })
      ),
    ])
  ).isRequired,
};
