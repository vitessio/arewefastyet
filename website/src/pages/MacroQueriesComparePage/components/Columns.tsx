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
  QueryType: string;
  Original: string;
  Instructions: string;
  ExecCount: number;
  ExecTime: number;
  ShardQueries: number;
  RowsReturned: number;
  RowsAffected: number;
  Errors: number;
  TablesUsed: string;
};

export type MacroQueriesPlanCommit = {
  Key: string;
  Value: MacroQueriesPlanCommitValue;
};

export type MacroQueriesPlan = {
  Key: string;
  ExecTimeDiff: number;
  ExecCountDiff: number;
  ErrorsDiff: number;
  RowsReturnedDiff: number;
  SamePlan: boolean;
  Right: MacroQueriesPlanCommit;
  Left: MacroQueriesPlanCommit;
};

export const columns: ColumnDef<MacroQueriesPlan>[] = [
  {
    header: ({ column }) => {
      return <div className="text-left">Query</div>;
    },
    accessorKey: "Key",
    cell: ({ row }) => {
      const formatted = row.original.Key;
      return <div className="text-left">{formatted}</div>;
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
    accessorKey: "ExecTimeDiff",
    cell: ({ row }) => {
      const formatted = row.original.ExecTimeDiff;
      let variant: "success" | "warning" | "destructive" = "success";
      if (formatted === 0) {
        variant = "warning";
      } else if (formatted < 0) {
        variant = "destructive";
      }

      return (
        <div>
          {" "}
          <Badge variant={variant}>{formatted}</Badge>
        </div>
      );
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
          Execution Count
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
    accessorKey: "ExecCountDiff",
    cell: ({ row }) => {
      const formatted = row.original.ExecCountDiff;
      return <>{formatted}%</>;
    },
  },
];
