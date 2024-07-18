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
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { MacroDataValue, Range } from "@/types";
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
  errors: MacroBenchmarkTableDataRow;
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

const getPBadgeVariant = (p: number) => {
  if (p < 0.05) {
    return "success";
  }
  return "destructive";
};

const formatCellValue = (key: string, value: any) => {
  if (key.includes("CpuTime")) {
    return secondToMicrosecond(value);
  } else if (key.includes("MemStatsAllocBytes")) {
    return formatByte(value);
  }
  return value;
};

export function getRange(range: Range) {
  if (range.infinite == true) {
    return "∞";
  }
  if (range.unknown == true) {
    return "?";
  }
  return "±" + fixed(range.value, 1) + "%";
}

// Type-safe function to access properties
const getValue = <T, K extends keyof T>(obj: T, key: K): T[K] => obj[key];

export default function MacroBenchmarkTable({
  data,
  newGitRef,
  oldGitRef,
}: MacroBenchmarkTableProps) {
  if (!data) {
    return null;
  }
  const dataKeys = Object.keys(data) as Array<keyof MacroBenchmarkTableData>;
  const classNameMap: { [key: string]: string } = {
    qpsTotal: "bg-background border-b light:border-foreground",
    qpsReads: "bg-muted hover:bg-muted/120",
    qpsWrites: "bg-muted hover:bg-muted/120",
    qpsOther: "bg-muted hover:bg-muted/120 border-b light:border-foreground",
    tps: "bg-background",
    latency: "bg-background",
    errors: "bg-background",
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
        {dataKeys.map((key, index) => {
          const row = getValue(data, key);
          return (
            <TableRow key={index} className={classNameMap[key]}>
              <TableCell className="w-[200px] font-medium text-right border-r border-border">
                {row.title}
              </TableCell>
              <TableCell className="text-center">
                {formatCellValue(key, row.old.center)} (
                {getRange(row.old.range)})
              </TableCell>
              <TableCell className="text-center border-r border-border">
                {formatCellValue(key, row.new.center)} (
                {getRange(row.new.range)})
              </TableCell>
              <TableCell className="text-center">
                <TooltipProvider delayDuration={200}>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <div>
                        <Badge variant={getPBadgeVariant(row.p)}>
                          {fixed(row.p, 3)}
                        </Badge>
                      </div>
                    </TooltipTrigger>
                    <TooltipContent>
                      <p>
                        {row.insignificant ? "Significant" : "Insignificant"}
                      </p>
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              </TableCell>
              <TableCell className="text-center">
                {row.p <= 0.05 && (
                  <Badge variant="success">{fixed(row.delta, 3)}</Badge>
                )}
                {row.p > 0.05 && <>{fixed(row.delta, 3)}</>}
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}
