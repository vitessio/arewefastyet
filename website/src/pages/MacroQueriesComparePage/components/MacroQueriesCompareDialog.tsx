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

import {
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import ReactJson from "react-json-pretty";
import "react-json-pretty/themes/monikai.css";
import { MacroQueriesPlan } from "./Columns";

export type MacroQueriesCompareDialogProps = {
  data: MacroQueriesPlan;
};

export default function MacroQueriesCompareDialog<MacroQueriesPlan, TValue>(
  props: MacroQueriesCompareDialogProps
) {
  const { data } = props;

  return (
    <DialogContent className="max-w-[80vw] md:max-w-[60vw] max-h-[80vh] flex flex-col overflow-scroll">
      <div className="flex flex-col items-center justify-center gap-4 md:gap-8">
        <DialogHeader>
          <DialogTitle className="text-primary text-center text-lg md:text-4xl">
            Statistics
          </DialogTitle>
        </DialogHeader>
        <div className="w-full md:w-3/5">
          <Table className="border-border border rounded-lg">
            <TableHeader className="text-center">
              <TableRow>
                <TableHead></TableHead>
                <TableHead className="text-center text-primary">Old</TableHead>
                <TableHead className="text-center text-primary">New</TableHead>
                <TableHead className="text-center">Delta (%)</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow>
                <TableCell className="text-center border-r border-border">
                  Execution Time
                </TableCell>
                <TableCell className="text-center">
                  {data.left.value.exec_time}
                </TableCell>
                <TableCell className="text-center">
                  {data.right.value.exec_time}
                </TableCell>
                <TableCell className="text-center">
                  {data.exec_time_diff}
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell className="text-center border-r border-border">
                  Rows Returned
                </TableCell>
                <TableCell className="text-center">
                  {data.left.value.rows_returned}
                </TableCell>
                <TableCell className="text-center">
                  {data.right.value.rows_returned}
                </TableCell>
                <TableCell className="text-center">
                  {data.rows_returned_diff}
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell className="text-center border-r border-border">
                  Errors
                </TableCell>
                <TableCell className="text-center">
                  {data.left.value.errors}
                </TableCell>
                <TableCell className="text-center">
                  {data.right.value.errors}
                </TableCell>
                <TableCell className="text-center">
                  {data.errors_diff}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
        <DialogTitle className="text-primary text-center text-lg md:text-4xl">
          Old query plan
        </DialogTitle>
        <div className="w-full px-4">
          <ScrollArea className="max-w-full whitespace-nowrap rounded-md border">
            <ReactJson
              data={data.left.value.instructions}
              className="text-sm md:text-base"
            />
            <ScrollBar orientation="horizontal" />
          </ScrollArea>
        </div>
        <DialogTitle className="text-primary text-center text-lg md:text-4xl">
          New query plan
        </DialogTitle>
        <div className="w-full px-4">
          <ScrollArea className="max-w-full whitespace-nowrap rounded-md border">
            <ReactJson
              data={data.right.value.instructions}
              className="text-sm md:text-base"
            />
            <ScrollBar orientation="horizontal" />
          </ScrollArea>
        </div>
      </div>
    </DialogContent>
  );
}
