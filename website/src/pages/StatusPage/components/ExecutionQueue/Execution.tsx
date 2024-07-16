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

import { ExecutionQueueType, columns } from "./Columns";
import { data } from "./data.json";
import { ExecutionQueueTable } from "./ExecutionTable";

export default function ExecutionQueue() {
  const executionQueueData: ExecutionQueueType[] = data.map(
    (value): ExecutionQueueType => {
      return {
        source: value.source,
        git_ref: value.git_ref,
        workload: value.workload,
        pull_nb: value.pull_nb,
      };
    }
  );

  return (
    <>
      <div className="xl:px-96 p-page my-12 flex flex-col">
        <h3 className="text-4xl md:text-5xl font-semibold text-primary mb-4 self-center">
          Execution Queue
        </h3>
        <ExecutionQueueTable data={executionQueueData} columns={columns} />
      </div>
    </>
  );
}
