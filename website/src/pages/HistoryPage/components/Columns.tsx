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

import { Badge, Variant } from "@/components/ui/badge";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { formatGitRef } from "@/utils/Utils";
import { ColumnDef } from "@tanstack/react-table";
import { format, formatDistanceToNow } from "date-fns";

export type HistoryType = {
  sha: string;
  source: string;
  workloads_benchmarked: number;
  started_at: Date;
};

export const columns: ColumnDef<HistoryType>[] = [
  {
    header: "SHA",
    accessorKey: "sha",
    cell: ({ row }) => {
      const formatted = formatGitRef(row.original.sha);
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
    header: "Workloads Benchmarked",
    accessorKey: "workloads_benchmarked",
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id));
    },
  },
  {
    header: "Started",
    accessorKey: "started_at",
    cell: ({ row }) => {
      const date = new Date(row.original.started_at);
      const formatted = formatDistanceToNow(date, {
        addSuffix: true,
      });
      return (
        <TooltipProvider delayDuration={200}>
          <Tooltip>
            <TooltipTrigger asChild>
              <div className="underline">{formatted}</div>
            </TooltipTrigger>
            <TooltipContent>
              <p>{format(date, "MMM d, yyyy, h:mm a 'GMT'XXX")}</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      );
    },
  },
];
