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

import React from 'react';
import '../SearchMacro/searchMacro.css';

const SearchMacro = ({ data }) => {
  const result = data[1]?.[0]?.Result;
  const defaultValue = '0';

  const renderRow = (label, value) => (
    <tr key={label}>
      <td>
        <span>{value || defaultValue}</span>
      </td>
    </tr>
  );

  return (
    <div className='search__macro__data'>
      <h3>{data[0]}</h3>
      {result ? (
        <table>
          <tbody>
            {Object.entries(result.qps).map(([key, value]) => renderRow(key, value))}
            {renderRow('tps', result.tps)}
            {renderRow('latency', result.latency)}
            {renderRow('errors', result.errors)}
            {renderRow('reconnects', result.reconnects)}
            {renderRow('time', result.time)}
            {renderRow('threads', result.threads)}
            {renderRow('TotalComponentsCPUTime', data[1][0].Metrics.TotalComponentsCPUTime?.toFixed(0))}
            {renderRow('ComponentsCPUTime.vtgate', data[1][0].Metrics.ComponentsCPUTime?.vtgate)}
            {renderRow('ComponentsCPUTime.vttablet', data[1][0].Metrics.ComponentsCPUTime?.vttablet)}
            {renderRow('TotalComponentsMemStatsAllocBytes', data[1][0].Metrics.TotalComponentsMemStatsAllocBytes)}
            {renderRow('ComponentsMemStatsAllocBytes.vtgate', data[1][0].Metrics.ComponentsMemStatsAllocBytes?.vtgate)}
            {renderRow('ComponentsMemStatsAllocBytes.vttablet', data[1][0].Metrics.ComponentsMemStatsAllocBytes?.vttablet)}
          </tbody>
        </table>
      ) : (
        <table>
          <tbody>
            {Array.from({ length: 16 }, (_, index) => renderRow(`default_${index}`, null))}
          </tbody>
        </table>
      )}
    </div>
  );
};

export default SearchMacro;
