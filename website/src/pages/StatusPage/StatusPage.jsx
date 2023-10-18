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

  return (
    <>
      <Hero />

      <figure className="p-page w-full">
        <div className="border-front border" />
      </figure>

      {/* EXECUTION QUEUE */}
      {!isLoadingQueue && dataQueue && dataQueue.length > 0 && (
        <ExecutionQueue data={dataQueue} />
      )}

      {/* PREVIOUS EXECUTIONS */}
      {!isLoadingPreviousExe &&
        dataPreviousExe &&
        dataPreviousExe.length > 0 && (
          <PreviousExecutions data={dataPreviousExe} />
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
      
      {errorQueue && <div className="my-10 text-center text-red-500">{errorQueue}</div>}
    </>
  );
}
