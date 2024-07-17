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

import { Badge } from "@/components/ui/badge";
import { ColumnDef } from "@tanstack/react-table";

export type ExecutionQueueType = {
  executions: ExecutionQueueExecution[];
  sources: string[];
  workloads: string[];
  statuses: string[];
};

export type ExecutionQueueExecution = {
  source: string;
  git_ref: string;
  workload: string;
  pull_nb: number;
};

export const columns: ColumnDef<ExecutionQueueExecution>[] = [
  {
    header: "SHA",
    accessorKey: "git_ref",
    cell: ({ row }) => {
      const formatted = row.original.git_ref.slice(0, 8);
      return (
        <div>
          {" "}
          <Badge color="primary">{formatted}</Badge>
        </div>
      );
    },
  },
  {
    header: "Source",
    accessorKey: "source",
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id));
    },
  },
  {
    header: "Workload",
    accessorKey: "workload",
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id));
    },
  },
  {
    header: "PR",
    accessorKey: "pull_nb",
    cell: ({ row }) => {
      const formatted =
        row.original.pull_nb === 0 ? "N/A" : row.original.pull_nb;
      return (
        <div>
          {" "}
          <Badge color="primary">{formatted}</Badge>
        </div>
      );
    },
  },
];
