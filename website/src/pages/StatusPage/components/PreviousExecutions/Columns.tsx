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

import { ColumnDef } from "@tanstack/react-table";
import { Badge, Variant } from "@/components/ui/badge";
import { formatDistanceToNow, format } from "date-fns";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

export type PreviousExecution = {
  uuid: string;
  source: string;
  git_ref: string;
  status: string;
  type_of: string;
  pull_nb: string;
  golang_version: string;
  started_at: string;
  finished_at: string;
};

export const columns: ColumnDef<PreviousExecution>[] = [
  {
    header: "UUID",
    accessorKey: "uuid",
    cell: ({ row }) => {
      const formatted = row.original.uuid.slice(0, 8);
      return <div>{formatted}</div>;
    },
  },
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
    header: "Started At",
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
  {
    header: "Finished At",
    accessorKey: "finished_at",
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
  {
    header: "PR",
    accessorKey: "pull_nb",
    cell: ({ row }) => {
      // if pull_nb is 0, display N/A
      const formatted = row.original.pull_nb == "0" ? "N/A" : row.original.pull_nb;
      console.log({ formatted });
      return (
        <div>
          {" "}
          <Badge color="primary">{formatted}</Badge>
        </div>
      );
    },
  },
  {
    header: "Go Version",
    accessorKey: "golang_version",
  },
  {
    header: "Status",
    accessorKey: "status",
    cell: ({ row }) => {
      const status = row.original.status;
      let variant: Variant = "success";
      if (status === "failed") {
        variant = "destructive";
      } else if (status === "canceled") {
        variant = "warning";
      } else if (status === "started") {
        variant = "progress";
      }
      return (
        <div>
          <Badge variant={variant}>{status}</Badge>
        </div>
      );
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id));
    },
  },
];
