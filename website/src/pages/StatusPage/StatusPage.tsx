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

import React, { useState } from "react";
import RingLoader from "react-spinners/RingLoader";
import useApiCall from "../../utils/Hook";

import Hero from "./components/Hero";
import ExecutionQueue from "./components/PreviousExecutions";
import PreviousExecutions from "./components/PreviousExecutions";

export default function StatusPage() {
  const {
    data: dataQueue,
    isLoading: isLoadingQueue,
    error: errorQueue,
  } = useApiCall(`${import.meta.env.VITE_API_URL}queue`);
  const { data: dataPreviousExe, isLoading: isLoadingPreviousExe } = useApiCall(
    `${import.meta.env.VITE_API_URL}recent`
  );

  const [filters, setFilters] = useState({
    type: '',
    status: '',
    source: '',
  });

  const filterData = (data) => {
    return data.filter((item) => {
      const itemType = item.type_of ? item.type_of.toString() : '';
      const itemSource = item.source ? item.source.toString() : '';
      const itemStatus = item.status ? item.status.toString() : '';

      const matchesType =
        filters.type === '' || filters.type === undefined || itemType === filters.type;
      const matchesSource =
        filters.source === '' || filters.source === undefined || itemSource === filters.source;
      const matchesStatus =
        filters.status === '' || filters.status === undefined || itemStatus === filters.status;

      return matchesType && matchesSource && matchesStatus;
    });
  };

  const handleFilterChange = (e) => {
    const { name, value } = e.target;
    setFilters((prevFilters) => ({ ...prevFilters, [name]: value }));
  };

  const filteredData = filterData(dataQueue);

  return (
    <>
      <Hero />

      <div className="border-accent border mt-5" />

      {/* FILTERS OPTIONS*/}
      <div>
        <label>
          Type:
          <select name="type" value={filters.type} onChange={handleFilterChange} className="bg-accent rounded-xl p-1">
            <option value="">All</option>
            {[...new Set(filteredData.map((item) => item.type_of))].map((type) => (
              <option key={type} value={type}>
                {type}
              </option>
            ))}
          </select>
        </label>

        <label>
          Source:
          <select name="source" value={filters.source} onChange={handleFilterChange} className="bg-accent rounded-xl p-1">
            <option value="">All</option>
            {[...new Set(filteredData.map((item) => item.source))].map((source) => (
              <option key={source} value={source}>
                {source}
              </option>
            ))}
          </select>
        </label>

        <label>
          Status:
          <select name="status" value={filters.status} onChange={handleFilterChange} className="bg-accent rounded-xl p-1">
            <option value="">All</option>
            {[...new Set(filteredData.map((item) => item.status))].map((status) => (
              <option key={status} value={status}>
                {status}
              </option>
            ))}
          </select>
        </label>
      </div>

      {/* EXECUTION QUEUE */}
      {!isLoadingQueue && dataQueue && dataQueue.length > 0 && (
        <ExecutionQueue data={filteredData} title={"Execution Queue"} />
      )}

      {/* PREVIOUS EXECUTIONS */}
      {!isLoadingPreviousExe &&
        dataPreviousExe &&
        dataPreviousExe.length > 0 && (
          <PreviousExecutions
            data={filteredData}
            title={"Previous Executions"}
          />
        )}

      {/* SHOW LOADER BENEATH IF EITHER IS LOADING */}
      {(isLoadingPreviousExe || isLoadingQueue) && (
        <div className="flex justify-center w-full my-16">
          <RingLoader
            loading={isLoadingPreviousExe || isLoadingQueue}
            color="#E77002"
            size={300}
          />
        </div>
      )}

      {errorQueue && (
        <div className="my-10 text-center text-red-500">{errorQueue}</div>
      )}
    </>
  );
}
