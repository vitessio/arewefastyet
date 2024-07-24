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
import { FilterConfigs } from "@/types";

interface DataTableToolbarProps<TData> {
  table: Table<TData>;
  filterConfigs: FilterConfigs[];
}

export function DataTableToolbar<TData>({
  table,
  filterConfigs,
}: DataTableToolbarProps<TData>) {
  const isFiltered = table.getState().columnFilters.length > 0;

  return (
    <div className="flex items-center gap-8 md:gap-0 md:justify-between py-4 flex-row">
      <div className="flex flex-1 gap-4 md:flex-none h-full">
        <Input
          placeholder="Filter executions..."
          value={(table.getColumn("query")?.getFilterValue() as string) ?? ""}
          onChange={(event) =>
            table.getColumn("query")?.setFilterValue(event.target.value)
          }
          className="h-full w-full flex-1 md:w-[150px] lg:w-[250px]"
        />
        <div className="hidden w-0 md:flex items-center space-x-2">
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
      </div>
      <div className="md:w-auto justify-end">
        <DataTableViewOptions table={table} />
      </div>
    </div>
  );
}
