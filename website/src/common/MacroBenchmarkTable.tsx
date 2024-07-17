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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { MacroDataValue } from "@/types";
import {
  fixed,
  formatByte,
  formatGitRef,
  secondToMicrosecond,
} from "@/utils/Utils";

type MacroBenchmarkTableDataRow = {
  title: string;
  old: MacroDataValue;
  new: MacroDataValue;
  p: number;
  delta: number;
  insignificant: boolean;
};

export type MacroBenchmarkTableData = {
  qpsTotal: MacroBenchmarkTableDataRow;
  qpsReads: MacroBenchmarkTableDataRow;
  qpsWrites: MacroBenchmarkTableDataRow;
  qpsOther: MacroBenchmarkTableDataRow;
  tps: MacroBenchmarkTableDataRow;
  latency: MacroBenchmarkTableDataRow;
  totalComponentsCpuTime: MacroBenchmarkTableDataRow;
  vtgateCpuTime: MacroBenchmarkTableDataRow;
  vttabletCpuTime: MacroBenchmarkTableDataRow;
  totalComponentsMemStatsAllocBytes: MacroBenchmarkTableDataRow;
  vtgateMemStatsAllocBytes: MacroBenchmarkTableDataRow;
  vttabletMemStatsAllocBytes: MacroBenchmarkTableDataRow;
};

export type MacroBenchmarkTableProps = {
  data: MacroBenchmarkTableData;
  oldGitRef: string;
  newGitRef: string;
};

const getBadgeVariant = (delta: number) => {
  if (delta < 0) {
    return "success";
  } else if (delta === 0) {
    return "warning";
  } else {
    return "destructive";
  }
};

const formatCellValue = (key: string, value: any) => {
  if (key.includes("CpuTime")) {
    return secondToMicrosecond(value);
  } else if (key.includes("MemStatsAllocBytes")) {
    return formatByte(value);
  }
  return value;
};

export default function MacroBenchmarkTable({
  data,
  newGitRef,
  oldGitRef,
}: MacroBenchmarkTableProps) {
  if (!data) {
    return null;
  }
  const dataKeys = Object.keys(data);
  const classNameMap: { [key: string]: string } = {
    qpsTotal: "bg-background border-b light:border-foreground",
    qpsReads: "bg-muted hover:bg-muted/120",
    qpsWrites: "bg-muted hover:bg-muted/120",
    qpsOther: "bg-muted hover:bg-muted/120 border-b light:border-foreground",
    tps: "bg-background",
    latency: "bg-background",
    totalComponentsCpuTime: "bg-background border-b light:border-foreground",
    vtgateCpuTime: "bg-muted hover:bg-muted/120 ",
    vttabletCpuTime:
      "bg-muted hover:bg-muted/120 border-b light:border-foreground",
    totalComponentsMemStatsAllocBytes:
      "bg-background border-b light:border-foreground",
    vtgateMemStatsAllocBytes: "bg-muted hover:bg-muted/120",
    vttabletMemStatsAllocBytes:
      "bg-muted hover:bg-muted/120 border-b light:border-foreground",
  };

  return (
    <Table>
      <TableHeader>
        <TableRow className="hover:bg-background border-b">
          <TableHead className="w-[200px]"></TableHead>
          <TableHead className="text-center text-primary font-semibold">
            {formatGitRef(oldGitRef) || "N/A"}
          </TableHead>
          <TableHead className="text-center text-primary font-semibold">
            {formatGitRef(newGitRef) || "N/A"}
          </TableHead>
          <TableHead className="text-center font-semibold">P</TableHead>
          <TableHead className="text-center font-semibold">Delta</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {dataKeys.map((key, index) => (
          <TableRow key={index} className={classNameMap[key]}>
            <TableCell className="w-[200px] font-medium text-right border-r border-border">
              {data[key as keyof MacroBenchmarkTableData].title}
            </TableCell>
            <TableCell className="text-center">
              {formatCellValue(
                key,
                data[key as keyof MacroBenchmarkTableData].old.center
              )}
            </TableCell>
            <TableCell className="text-center border-r border-border">
              {formatCellValue(
                key,
                data[key as keyof MacroBenchmarkTableData].new.center
              )}
            </TableCell>
            <TableCell className="text-center">
              {fixed(data[key as keyof MacroBenchmarkTableData].p, 3)}
            </TableCell>
            <TableCell className="text-center">
              <Badge
                variant={getBadgeVariant(
                  data[key as keyof MacroBenchmarkTableData].delta
                )}
              >
                {data[key as keyof MacroBenchmarkTableData].delta}
              </Badge>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
