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
import { MacroBenchmarkTableData, Range, VitessRefs } from "@/types";
import {
  fixed,
  formatByte,
  getRefName,
  secondToMicrosecond,
} from "@/utils/Utils";
import { Link } from "react-router-dom";

export type MacroBenchmarkTableProps = {
  data: MacroBenchmarkTableData;
  old: string;
  new: string;
  isGitRef?: boolean;
  vitessRefs: VitessRefs | undefined;
};

const getSignificanceBadge = (p: number) => {
  let backgroundColor = "";
  let textColor = "";
  let val = fixed(p, 3)

  if (p <= 0.01) {
    backgroundColor = "#2E7D32"
    textColor = "#FFFFFF"
  } else if (p <= 0.05) {
    backgroundColor = "#388E3C"
    textColor = "#FFFFFF"
  } else if (p <= 0.10) {
    backgroundColor = "#6A9A1F"
    textColor = "#FFFFFF"
  } else {
    backgroundColor = "#9E9E9E"
    textColor = "#000000"
  }

  return (
      <Badge style={{backgroundColor: backgroundColor, color: textColor}}>
        {val}
      </Badge>
  );
};

const getSignificanceText = (p: number) => {
  if (p <= 0.01) {
    return "Statistically Significant";
  } else if (p <= 0.05) {
    return "Moderate Significance";
  } else if (p <= 0.10) {
    return "Marginal Significance";
  }
  return "Not Statistically Significant";
};

const getDeltaBadgeVariant = (key: string, delta: number, p: number) => {
  if (delta === 0) {
    return "warning";
  }
  if (
    key.includes("CpuTime") ||
    key.includes("Mem") ||
    key.includes("latency")
  ) {
    if (delta < 0) {
      return "success";
    }
  }
  if (key.includes("qps") || key.includes("tps")) {
    if (delta > 0) {
      return "success";
    }
  }
  return "destructive";
};

const formatCellValue = (key: string, value: number) => {
  if (key.includes("CpuTime")) {
    return secondToMicrosecond(value);
  } else if (key.includes("MemStatsAllocBytes")) {
    return formatByte(Number(fixed(value, 2)));
  } else if (key.includes("latency")) {
    return fixed(value, 2) + "ms";
  }
  return fixed(value, 2);
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
  new: newGitRef,
  old,
  isGitRef = true,
  vitessRefs,
}: MacroBenchmarkTableProps) {
  if (!data) {
    return null;
  }
  const dataKeys = Object.keys(data) as Array<keyof MacroBenchmarkTableData>;

  // Use keyof MacroBenchmarkTableData to enforce key types
  const classNameMap: { [key in keyof MacroBenchmarkTableData]: string } = {
    qpsTotal: "border-b border-foreground dark:border-none",
    qpsReads: "bg-muted/80",
    qpsWrites: "bg-muted/80",
    qpsOther: "bg-muted/80 border-b border-foreground dark:border-none",
    tps: "bg-background",
    latency: "bg-background",
    errors: "bg-background",
    totalComponentsCpuTime: "border-b border-foreground dark:border-none",
    vtgateCpuTime: "bg-muted/80 ",
    vttabletCpuTime: "bg-muted/80 border-b border-foreground dark:border-none",
    totalComponentsMemStatsAllocBytes:
      "border-b border-foreground dark:border-none",
    vtgateMemStatsAllocBytes: "bg-muted/80",
    vttabletMemStatsAllocBytes:
      "bg-muted/80 border-b border-foreground dark:border-none",
  };

  return (
    <Table>
      <TableHeader>
        <TableRow className="hover:bg-background border-b">
          <TableHead className="w-[200px]"></TableHead>
          <TableHead className="text-center text-primary font-semibold min-w-[150px]">
            {isGitRef && vitessRefs ? (
              <Link
                to={`https://github.com/vitessio/vitess/commit/${old}`}
                target="__blank"
              >
                {getRefName(old, vitessRefs) || "N/A"}
              </Link>
            ) : (
              <>{old}</>
            )}
          </TableHead>
          <TableHead className="text-center text-primary font-semibold min-w-[150px]">
            {isGitRef && vitessRefs ? (
              <Link
                to={`https://github.com/vitessio/vitess/commit/${newGitRef}`}
                target="__blank"
              >
                {getRefName(newGitRef, vitessRefs) || "N/A"}
              </Link>
            ) : (
              <>{newGitRef}</>
            )}
          </TableHead>
          <TableHead className="lg:w-[150px] text-center font-semibold">
            P
          </TableHead>
          <TableHead className="lg:w-[150px] text-center font-semibold">
            Delta
          </TableHead>
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
              <TableCell className="text-center text-front">
                {formatCellValue(key, Number(row.old.center))} (
                {getRange(row.old.range)})
              </TableCell>
              <TableCell className="text-center text-front border-r border-border">
                {formatCellValue(key, Number(row.new.center))} (
                {getRange(row.new.range)})
              </TableCell>
              <TableCell className="lg:w-[150px] text-center text-front">
                <TooltipProvider delayDuration={200}>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <div>
                        {getSignificanceBadge(row.p)}
                      </div>
                    </TooltipTrigger>
                    <TooltipContent>
                      <p>
                        {getSignificanceText(row.p)}
                      </p>
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              </TableCell>
              <TableCell className="lg:w-[150px] text-center text-front">
                {row.p > 0.1 && <>{fixed(row.delta, 3)}%</>}
                {row.p <= 0.1 && (
                  <Badge variant={getDeltaBadgeVariant(key, row.delta, row.p)}>
                    {fixed(row.delta, 3)}%
                  </Badge>
                )}
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}
