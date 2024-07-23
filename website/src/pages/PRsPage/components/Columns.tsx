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

import Icon from "@/common/Icon";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { PrData } from "@/types";
import { ColumnDef } from "@tanstack/react-table";
import { format, formatDistanceToNow } from "date-fns";
import { Link } from "react-router-dom";

export const columns: ColumnDef<PrData>[] = [
  {
    header: "#",
    accessorKey: "ID",
    cell: ({ row }) => {
      const formatted = row.original.ID;
      return (
        <Link
          target="__blank"
          to={`https://github.com/vitessio/vitess/pull/${formatted}`}
        >
          <p className="text-primary text-left">{formatted}</p>
        </Link>
      );
    },
  },
  {
    header: "Title",
    accessorKey: "Title",
    cell: ({ row }) => {
      const formatted = row.original.Title;
      return (
        <p className="text-left min-w-fit whitespace-nowrap">{formatted}</p>
      );
    },
  },
  {
    header: "Author",
    accessorKey: "Author",
    cell: ({ row }) => {
      const formatted = row.original.Author;
      return (
        <Link target="__blank" to={`https://github.com/${formatted}`}>
          <p className="text-primary text-left">{formatted}</p>
        </Link>
      );
    },
  },
  {
    header: "Opened At",
    accessorKey: "CreatedAt",
    cell: ({ row }) => {
      const date = new Date(row.original.CreatedAt);
      const formatted = formatDistanceToNow(date, {
        addSuffix: true,
      });
      return (
        <TooltipProvider delayDuration={200}>
          <Tooltip>
            <TooltipTrigger asChild>
              <div className="underline text-left">{formatted}</div>
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
    header: "Details",
    accessorKey: "Base",
    cell: ({ row }) => {
      const formatted = row.original.ID;
      return (
        <Link
          to={`/pr/${formatted}`}
          className="flex justify-center text-lg text-primary duration-300 hover:scale-105"
        >
          <Icon icon="open_in_new" />
        </Link>
      );
    },
  },
];
