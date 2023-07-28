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
import PulseLoader from "react-spinners/PulseLoader";
import { formatByteForGB, fixed } from "../../../utils/Utils";

import "./macrobench.css";

const Macrobench = React.memo(
  ({
    data,
    gitRefLeft,
    gitRefRight,
    commitHashLeft,
    commitHashRight,
    textLoading,
  }) => {
    const renderDataOrLoader = (data, loading) => {
      if (loading) {
        return (
          <span>
            <PulseLoader loading={true} size={5} color="#ffffff" />
          </span>
        );
      } else {
        return <span>{data}</span>;
      }
    };

    return (
      <div className="macrobench__component">
        <div className="macrobench__component__header">
          <h3>{data.type}</h3>
          <span className="linkQueryT">
            Click{" "}
            <Link
              to="/macrobench/queries/compare"
              state={{
                data: data,
                commitHashLeft: commitHashLeft,
                commitHashRight: commitHashRight,
              }}
            >
              here
            </Link>{" "}
            to see the query plans comparison for this benchmark.
          </span>
        </div>

        <table>
          <thead>
            <tr>
              <th></th>
              <th>
                <h4>{gitRefLeft ? gitRefLeft : "Left"}</h4>
              </th>
              <th>
                <h4>{gitRefRight ? gitRefRight : "Right"}</h4>
              </th>
              <th>
                <h4>Impoved by %</h4>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr className="border--bottom">
              <td className="sidebar--border">
                <span className="sidebar">QPS Total</span>
              </td>
              <td>
                {renderDataOrLoader(
                  data.diff.Left.Result.qps.total,
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  data.diff.Right.Result.qps.total,
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  fixed(data.diff.Diff.qps.total, 2),
                  textLoading
                )}
              </td>
            </tr>
            <tr className="border--bottom">
              <td className="sidebar--border">
                <span className="sidebar">QPS Reads</span>
              </td>
              <td>
                {renderDataOrLoader(
                  data.diff.Left.Result.qps.reads,
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  data.diff.Right.Result.qps.reads,
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  fixed(data.diff.Diff.qps.reads, 2),
                  textLoading
                )}
              </td>
            </tr>
            <tr className="border--bottom">
              <td className="sidebar--border">
                <span className="sidebar">QPS Writes</span>
              </td>
              <td>
                {renderDataOrLoader(
                  data.diff.Left.Result.qps.writes,
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  data.diff.Right.Result.qps.writes,
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  fixed(data.diff.Diff.qps.writes, 2),
                  textLoading
                )}
              </td>
            </tr>
            <tr className="border--bottom">
              <td className="sidebar--border">
                <span className="sidebar">QPS Other</span>
              </td>
              <td>
                {renderDataOrLoader(
                  data.diff.Left.Result.qps.other,
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  data.diff.Right.Result.qps.other,
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  fixed(data.diff.Diff.qps.other, 2),
                  textLoading
                )}
              </td>
            </tr>
            <tr className="border--bottom">
              <td className="sidebar--border">
                <span className="sidebar">TPS</span>
              </td>
              <td>
                {renderDataOrLoader(data.diff.Left.Result.tps, textLoading)}
              </td>
              <td>
                {renderDataOrLoader(data.diff.Right.Result.tps, textLoading)}
              </td>
              <td>
                {renderDataOrLoader(fixed(data.diff.Diff.tps, 2), textLoading)}
              </td>
            </tr>
            <tr className="border--bottom">
              <td className="sidebar--border">
                <span className="sidebar">Latency</span>
              </td>
              <td>
                {renderDataOrLoader(data.diff.Left.Result.latency, textLoading)}
              </td>
              <td>
                {renderDataOrLoader(
                  data.diff.Right.Result.latency,
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  fixed(data.diff.Diff.latency, 2),
                  textLoading
                )}
              </td>
            </tr>
            <tr className="border--bottom">
              <td className="sidebar--border">
                <span className="sidebar">Errors</span>
              </td>
              <td>
                {renderDataOrLoader(data.diff.Left.Result.errors, textLoading)}
              </td>
              <td>
                {renderDataOrLoader(data.diff.Right.Result.errors, textLoading)}
              </td>
              <td>
                {renderDataOrLoader(
                  fixed(data.diff.Diff.errors, 2),
                  textLoading
                )}
              </td>
            </tr>
            <tr className="border--bottom">
              <td className="sidebar--border">
                <span className="sidebar">Reconnects</span>
              </td>
              <td>
                {renderDataOrLoader(
                  data.diff.Left.Result.reconnects,
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  data.diff.Right.Result.reconnects,
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  fixed(data.diff.Diff.reconnects, 2),
                  textLoading
                )}
              </td>
            </tr>
            <tr className="border--bottom">
              <td className="sidebar--border">
                <span className="sidebar">Time</span>
              </td>
              <td>
                {renderDataOrLoader(data.diff.Left.Result.time, textLoading)}
              </td>
              <td>
                {renderDataOrLoader(data.diff.Right.Result.time, textLoading)}
              </td>
              <td>
                {renderDataOrLoader(fixed(data.diff.Diff.time, 2), textLoading)}
              </td>
            </tr>
            <tr className="border--bottom">
              <td className="sidebar--border">
                <span className="sidebar">Threads</span>
              </td>
              <td>
                {renderDataOrLoader(data.diff.Left.Result.threads, textLoading)}
              </td>
              <td>
                {renderDataOrLoader(
                  data.diff.Right.Result.threads,
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  fixed(data.diff.Diff.threads, 2),
                  textLoading
                )}
              </td>
            </tr>
            <tr className="border--bottom">
              <td className="sidebar--border">
                <span className="sidebar">Total CPU time</span>
              </td>
              <td>
                {renderDataOrLoader(
                  fixed(data.diff.Left.Metrics.TotalComponentsCPUTime, 0),
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  fixed(data.diff.Right.Metrics.TotalComponentsCPUTime, 0),
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  fixed(data.diff.DiffMetrics.TotalComponentsCPUTime, 2),
                  textLoading
                )}
              </td>
            </tr>
            <tr className="border--bottom">
              <td className="sidebar--border">
                <span className="sidebar">CPU time vtgate</span>
              </td>
              <td>
                {renderDataOrLoader(
                  data.diff.Left.Metrics.ComponentsCPUTime.vtgate,
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  data.diff.Right.Metrics.ComponentsCPUTime.vtgate,
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  fixed(data.diff.DiffMetrics.ComponentsCPUTime.vtgate, 2),
                  textLoading
                )}
              </td>
            </tr>
            <tr className="border--bottom">
              <td className="sidebar--border">
                <span className="sidebar">CPU time vttablet</span>
              </td>
              <td>
                {renderDataOrLoader(
                  data.diff.Left.Metrics.ComponentsCPUTime.vttablet,
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  data.diff.Right.Metrics.ComponentsCPUTime.vttablet,
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  fixed(data.diff.DiffMetrics.ComponentsCPUTime.vttablet, 2),
                  textLoading
                )}
              </td>
            </tr>
            <tr className="border--bottom">
              <td className="sidebar--border">
                <span className="sidebar">Total Allocs bytes</span>
              </td>
              <td>
                {renderDataOrLoader(
                  formatByteForGB(
                    data.diff.Left.Metrics.TotalComponentsMemStatsAllocBytes
                  ),
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  formatByteForGB(
                    data.diff.Right.Metrics.TotalComponentsMemStatsAllocBytes
                  ),
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  fixed(
                    data.diff.DiffMetrics.TotalComponentsMemStatsAllocBytes,
                    2
                  ),
                  textLoading
                )}
              </td>
            </tr>
            <tr className="border--bottom">
              <td className="sidebar--border">
                <span className="sidebar">Allocs bytes vtgate</span>
              </td>
              <td>
                {renderDataOrLoader(
                  formatByteForGB(
                    data.diff.Left.Metrics.ComponentsMemStatsAllocBytes.vtgate
                  ),
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  formatByteForGB(
                    data.diff.Right.Metrics.ComponentsMemStatsAllocBytes.vtgate
                  ),
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  fixed(
                    data.diff.DiffMetrics.ComponentsMemStatsAllocBytes.vtgate,
                    2
                  ),
                  textLoading
                )}
              </td>
            </tr>
            <tr>
              <td className="sidebar--border">
                <span className="sidebar">Allocs bytes vttablet</span>
              </td>
              <td>
                {renderDataOrLoader(
                  formatByteForGB(
                    data.diff.Left.Metrics.ComponentsMemStatsAllocBytes.vttablet
                  ),
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  formatByteForGB(
                    data.diff.Right.Metrics.ComponentsMemStatsAllocBytes
                      .vttablet
                  ),
                  textLoading
                )}
              </td>
              <td>
                {renderDataOrLoader(
                  fixed(
                    data.diff.DiffMetrics.ComponentsMemStatsAllocBytes.vttablet,
                    2
                  ),
                  textLoading
                )}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    );
  }
);

export default Macrobench;
