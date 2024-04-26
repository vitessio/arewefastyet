/*
Copyright 2023 The Vitess Authors.

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
import React from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

export default function DisplayList(props) {
  const { data } = props;

  return (
    <div className="">
      <div className="widescreen:hidden">
        {data.map((val, key) => (
          <div
            key={key}
            className="cursor-default mb-4 p-4 rounded-lg border border-front border-[1px]"
          >
            {Object.keys(val).map((item, itemKey) => (
              <div
                key={itemKey}
                className="flex justify-between border-b py-2 last:border-b-0 hover:bg-accent "
              >
                <span className="font-semibold">{item}</span>
                <span>{val[item]}</span>
              </div>
            ))}
          </div>
        ))}
      </div>
      <div className="w-full overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow>
              {Object.keys(data[0]).map((item, key) => (
                <TableHead key={key}>{item}</TableHead>
              ))}
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.map((val, key) => (
              <TableRow key={key} className="hover:bg-accent cursor-default">
                {Object.keys(val).map((item, key) => (
                  <TableCell key={key}>{val[item]}</TableCell>
                ))}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
