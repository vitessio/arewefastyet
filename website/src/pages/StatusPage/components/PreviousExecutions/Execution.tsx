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

import useApiCall from "@/hooks/useApiCall";
import { FilterConfigs } from "@/types";
import { columns, type PreviousExecution } from "./Columns";
import { PreviousExecutionQueueTable } from "./ExecutionTable";

export default function PreviousExecution() {
  const { data: dataPreviousExecution, isLoading } =
    useApiCall<PreviousExecution>({
      url: `${import.meta.env.VITE_API_URL}recent`,
      queryKey: "recent",
    });

  const filterConfigs: FilterConfigs[] = [
    {
      column: "source",
      title: "Source",
      options:
        dataPreviousExecution?.sources?.map((source) => {
          return { label: source, value: source };
        }) || [],
    },
    {
      column: "status",
      title: "Status",
      options:
        dataPreviousExecution?.statuses.map((status) => {
          return { label: status, value: status };
        }) || [],
    },
    {
      column: "workload",
      title: "Workload",
      options:
        dataPreviousExecution?.workloads.map((workload) => {
          return { label: workload, value: workload };
        }) || [],
    },
  ];

  return (
    <>
      <div className="p-page my-12 flex flex-col">
        <h3 className="text-4xl md:text-5xl font-semibold text-primary mb-4 self-center">
          Previous Executions
        </h3>
        {dataPreviousExecution !== undefined && (
          <PreviousExecutionQueueTable
            columns={columns}
            data={dataPreviousExecution.executions}
            filterConfigs={filterConfigs}
          />
        )}
      </div>
    </>
  );
}
