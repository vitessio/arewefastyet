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

import { Table } from "@tanstack/react-table";

import { DataTableViewOptions } from "@/components/ui/data-table-view-options";
import { Input } from "@/components/ui/input";

interface DataTableToolbarProps<TData> {
  table: Table<TData>;
}

export function DataTableToolbar<TData>({
  table,
}: DataTableToolbarProps<TData>) {
  return (
    <div className="flex items-center gap-8 md:gap-0 md:justify-between py-4 flex-row">
      <div className="flex flex-1 gap-4 md:flex-none h-full">
        <Input
          placeholder="Filter Pull Requests..."
          value={(table.getColumn("Title")?.getFilterValue() as string) ?? ""}
          onChange={(event) =>
            table.getColumn("Title")?.setFilterValue(event.target.value)
          }
          className="h-full w-full flex-1 md:w-[150px] lg:w-[250px]"
        />
      </div>
      <div className="md:w-auto justify-end">
        <DataTableViewOptions table={table} />
      </div>
    </div>
  );
}
