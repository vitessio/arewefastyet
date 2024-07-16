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

import { Cross2Icon } from "@radix-ui/react-icons";
import { Table } from "@tanstack/react-table";

import { Button } from "@/components/ui/button";
import { DataTableViewOptions } from "@/components/ui/data-table-view-options";
import { Input } from "@/components/ui/input";

import { DataTableFacetedFilter } from "@/components/ui/data-table-faceted-filter";

const workloadses = [
  {
    label: "oltp",
    value: "oltp",
  },
  {
    label: "oltp-readonly",
    value: "oltp-readonly",
  },
  {
    label: "oltp-set",
    value: "oltp-set",
  },
  {
    label: "tpcc",
    value: "tpcc",
  },
  {
    label: "tpcc_fk",
    value: "tpcc_fk",
  },
  {
    label: "tpcc_fk_unmanaged",
    value: "tpcc_fk_unmanaged",
  },
  {
    label: "tpcc_unsharded",
    value: "tpcc_unsharded",
  },
];

const sourceses = [
  {
    label: "cron",
    value: "cron",
  },
  {
    label: "cron_pr",
    value: "cron_pr",
  },
  {
    label: "cron_pr_base",
    value: "cron_pr_base",
  },
  {
    label: "cron_tags",
    value: "cron_tags",
  },
];

const filterConfigs = [
  {
    column: "workload",
    title: "Workload",
    options: workloadses,
  },
  {
    column: "source",
    title: "Source",
    options: sourceses,
  },
];

interface DataTableToolbarProps<ExecutionQueueType> {
  table: Table<ExecutionQueueType>;
}

export function DataTableToolbar<ExecutionQueueType>({
  table,
}: DataTableToolbarProps<ExecutionQueueType>) {
  const isFiltered = table.getState().columnFilters.length > 0;

  return (
    <div className="flex items-center justify-between py-4">
      <div className="flex flex-1 items-center space-x-2">
        <Input
          placeholder="Filter executions..."
          value={(table.getColumn("git_ref")?.getFilterValue() as string) ?? ""}
          onChange={(event) =>
            table.getColumn("git_ref")?.setFilterValue(event.target.value)
          }
          className="h-8 w-[150px] lg:w-[250px]"
        />
        {filterConfigs.map((filter) => {
          const column = table.getColumn(filter.column);
          return (
            column && (
              <DataTableFacetedFilter
                key={filter.column}
                column={column}
                title={filter.title}
                options={filter.options}
              />
            )
          );
        })}
        {isFiltered && (
          <Button
            variant="ghost"
            onClick={() => table.resetColumnFilters()}
            className="h-8 px-2 lg:px-3"
          >
            Reset
            <Cross2Icon className="ml-2 h-4 w-4" />
          </Button>
        )}
      </div>
      <DataTableViewOptions table={table} />
    </div>
  );
}
