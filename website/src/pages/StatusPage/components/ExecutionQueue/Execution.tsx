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

import { FilterConfigs } from "@/types";
import useApiCall from "@/utils/Hook";
import {
  ExecutionQueueExecution,
  ExecutionQueueType,
  columns,
} from "./Columns";
import { ExecutionQueueTable } from "./ExecutionTable";

export default function ExecutionQueue() {
  const { data: dataExecutionQueue, isLoading } =
    useApiCall<ExecutionQueueType>(`${import.meta.env.VITE_API_URL}queue`);

  if (
    dataExecutionQueue === null ||
    dataExecutionQueue.executions.length === 0
  ) {
    return <></ >;
  }

  let filterConfigs: FilterConfigs[] = [
    {
      column: "source",
      title: "Source",
      options:
        dataExecutionQueue?.sources?.map((source) => {
          return { label: source, value: source };
        }) || [],
    },
    {
      column: "workload",
      title: "Workload",
      options:
        dataExecutionQueue?.workloads.map((workload) => {
          return { label: workload, value: workload };
        }) || [],
    },
  ];

  const executionQueueData: ExecutionQueueExecution[] | undefined =
    dataExecutionQueue?.executions?.map((value) => {
      return {
        source: value.source,
        git_ref: value.git_ref,
        workload: value.workload,
        pull_nb: value.pull_nb,
      };
    });

  return (
    <>
      <div className="mx-auto w-[50vw] my-12 flex flex-col">
        <h3 className="text-4xl md:text-5xl font-semibold text-primary mb-4 self-center">
          Execution Queue
        </h3>
        {isLoading && <div>Loading...</div>}
        {executionQueueData === null && <div>No data found</div>}
        {executionQueueData && (
          <ExecutionQueueTable
            columns={columns}
            data={executionQueueData}
            filterConfigs={filterConfigs}
          />
        )}{" "}
      </div>
    </>
  );
}
