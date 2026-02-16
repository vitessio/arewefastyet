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
import { Button } from "@/components/ui/button";
import { ColumnDef } from "@tanstack/react-table";
import { ArrowUpDown } from "lucide-react";

export type MacroQueriesPlanCommitValue = {
  query_type: string;
  original: string;
  instructions: string;
  exec_count: number;
  exec_time: number;
  shard_queries: number;
  rows_returned: number;
  rows_affected: number;
  errors: number;
  tables_used: string;
};

export type MacroQueriesPlanCommit = {
  key: string;
  value: MacroQueriesPlanCommitValue;
};

export type MacroQueriesPlan = {
  key: string;
  exec_time_diff: number;
  exec_count_diff: number;
  errors_diff: number;
  rows_returned_diff: number;
  same_plan: boolean;
  right: MacroQueriesPlanCommit | null;
  left: MacroQueriesPlanCommit | null;
};

export const columns: ColumnDef<MacroQueriesPlan>[] = [
  {
    id: "query",
    enableHiding: true,
    header: () => {
      return <div className="text-left">Query</div>;
    },
    accessorKey: "key",
    cell: ({ row }) => {
      const formatted = row.original.key;
      return <div className="text-left min-w-fit">{formatted}</div>;
    },
    enableColumnFilter: true,
    filterFn: (row, _, value) => {
      const original = row.original.key;
      return original.toString().includes(value);
    },
  },
  {
    header: ({ column }) => {
      return (
        <Button
          className="p-0"
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Execution Time
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
    id: "Execution Time",
    accessorKey: "exec_time_diff",
    cell: ({ row }) => {
      const formatted = row.original.exec_time_diff;
      let variant: "success" | "warning" | "destructive" = "success";
      if (formatted === 0) {
        variant = "warning";
      } else if (formatted < 0) {
        variant = "destructive";
      }

      return (
        <div>
          {" "}
          <Badge variant={variant}>{formatted}%</Badge>
        </div>
      );
    },
  },
];
