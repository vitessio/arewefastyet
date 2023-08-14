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

import "../SearchMacroDesktop/searchMacroDesktop.css";

import { fixed, formatByteForGB } from "../../../../utils/Utils";

const SearchMacroDesktop = ({ data, gitRef }) => {
  const renderZeroSpans = <span>0</span>;

  return (
    <div className="flex--column searchMAcro__desktop__data">
      <div className="searchMAcro__desktop__data__header">
        <h3>{data[0]}</h3>
        <a
          target="_blank"
          rel="noopener noreferrer"
          href={`https://github.com/vitessio/vitess/commit/${gitRef}`}
        >
          <span>{gitRef}</span>
        </a>
      </div>
      <table>
        <tbody>
          <tr>
            <td className="sidebar--border">
              <span>QPS Total</span>
            </td>
            <td className="datatd">
              {data[1] ? (
                <span>{fixed(data[1][0].Result.qps.total, 2)}</span>
              ) : (
                renderZeroSpans
              )}
            </td>
          </tr>
          <tr>
            <td className="sidebar--border">
              <span>QPS Reads</span>
            </td>
            <td className="datatd">
              {data[1] ? (
                <span>{data[1][0].Result.qps.reads}</span>
              ) : (
                renderZeroSpans
              )}
            </td>
          </tr>
          <tr>
            <td className="sidebar--border">
              <span>QPS Writes</span>
            </td>
            <td className="datatd">
              {data[1] ? (
                <span>{fixed(data[1][0].Result.qps.writes, 2)}</span>
              ) : (
                renderZeroSpans
              )}
            </td>
          </tr>
          <tr>
            <td className="sidebar--border">
              <span>QPS Other</span>
            </td>
            <td className="datatd">
              {data[1] ? (
                <span>{data[1][0].Result.qps.other.toFixed(2)}</span>
              ) : (
                renderZeroSpans
              )}
            </td>
          </tr>
          <tr>
            <td className="sidebar--border">
              <span>TPS</span>
            </td>
            <td className="datatd">
              {data[1] ? (
                <span>{fixed(data[1][0].Result.tps, 2)}</span>
              ) : (
                renderZeroSpans
              )}
            </td>
          </tr>
          <tr>
            <td className="sidebar--border">
              <span>Latency</span>
            </td>
            <td className="datatd">
              {data[1] ? (
                <span>{data[1][0].Result.latency}</span>
              ) : (
                renderZeroSpans
              )}
            </td>
          </tr>
          <tr>
            <td className="sidebar--border">
              <span>Errors</span>
            </td>
            <td className="datatd">
              {data[1] ? (
                <span>{fixed(data[1][0].Result.errors, 2)}</span>
              ) : (
                renderZeroSpans
              )}
            </td>
          </tr>
          <tr>
            <td className="sidebar--border">
              <span>Reconnects</span>
            </td>
            <td className="datatd">
              {data[1] ? (
                <span>{data[1][0].Result.reconnects}</span>
              ) : (
                renderZeroSpans
              )}
            </td>
          </tr>
          <tr>
            <td className="sidebar--border">
              <span>Time</span>
            </td>
            <td className="datatd">
              {data[1] ? (
                <span>{data[1][0].Result.time}</span>
              ) : (
                renderZeroSpans
              )}
            </td>
          </tr>
          <tr>
            <td className="sidebar--border">
              <span>Threads</span>
            </td>
            <td className="datatd">
              {data[1] ? (
                <span>{data[1][0].Result.threads}</span>
              ) : (
                renderZeroSpans
              )}
            </td>
          </tr>
          <tr>
            <td className="sidebar--border">
              <span>Total CPU time</span>
            </td>
            <td className="datatd">
              {data[1] ? (
                <span>
                  {data[1][0].Metrics.TotalComponentsCPUTime.toFixed(0)}
                </span>
              ) : (
                renderZeroSpans
              )}
            </td>
          </tr>
          <tr>
            <td className="sidebar--border">
              <span>CPU time vtgate</span>
            </td>
            <td className="datatd">
              {data[1] ? (
                <span>
                  {fixed(data[1][0].Metrics.ComponentsCPUTime.vtgate, 2)}
                </span>
              ) : (
                renderZeroSpans
              )}
            </td>
          </tr>
          <tr>
            <td className="sidebar--border">
              <span>CPU time vttablet</span>
            </td>
            <td className="datatd">
              {data[1] ? (
                <span>
                  {fixed(data[1][0].Metrics.ComponentsCPUTime.vttablet, 2)}
                </span>
              ) : (
                renderZeroSpans
              )}
            </td>
          </tr>
          <tr>
            <td className="sidebar--border">
              <span>Total Allocs bytes</span>
            </td>
            <td className="datatd">
              {data[1] ? (
                <span>
                  {formatByteForGB(
                    data[1][0].Metrics.TotalComponentsMemStatsAllocBytes
                  )}
                </span>
              ) : (
                renderZeroSpans
              )}
            </td>
          </tr>
          <tr>
            <td className="sidebar--border">
              <span>Allocs bytes vtgate</span>
            </td>
            <td className="datatd">
              {data[1] ? (
                <span>
                  {formatByteForGB(
                    data[1][0].Metrics.ComponentsMemStatsAllocBytes.vtgate
                  )}
                </span>
              ) : (
                renderZeroSpans
              )}
            </td>
          </tr>
          <tr>
            <td className="sidebar--border">
              <span>Allocs bytes vttablet</span>
            </td>
            <td className="datatd">
              {data[1] ? (
                <span>
                  {formatByteForGB(
                    data[1][0].Metrics.ComponentsMemStatsAllocBytes.vttablet
                  )}
                </span>
              ) : (
                renderZeroSpans
              )}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  );
};

SearchMacroDesktop.propTypes = {
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

export default SearchMacroDesktop;
