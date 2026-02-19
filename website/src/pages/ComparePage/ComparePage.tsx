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

import MacroBenchmarkTable, { formatCellValue, getRange } from "@/common/MacroBenchmarkTable";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import useApiCall from "@/hooks/useApiCall";
import { CompareData, MacroBenchmarkTableData, VitessRefs } from "@/types";
import {
  errorApi,
  fixed,
  formatCompareData,
  getGitRefFromRefName,
  getRefName,
} from "@/utils/Utils";
import { PlusCircledIcon } from "@radix-ui/react-icons";
import { useEffect, useRef, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import CompareHero from "./components/CompareHero";
import { Copy, Check } from "lucide-react";

export default function Compare() {
  const navigate = useNavigate();
  const urlParams = new URLSearchParams(window.location.search);

  const [gitRef, setGitRef] = useState({
    old: urlParams.get("old") || "",
    new: urlParams.get("new") || "",
  });

  const shouldFetchCompareData = gitRef.old && gitRef.new;

  const { data: vitessRefs } = useApiCall<VitessRefs>({
    url: `${import.meta.env.VITE_API_URL}vitess/refs`,
    queryKey: ["vitessRefs"],
  });

  const gitOldRef = vitessRefs
    ? getGitRefFromRefName(gitRef.old, vitessRefs)
    : gitRef.old;
  const gitNewRef = vitessRefs
    ? getGitRefFromRefName(gitRef.new, vitessRefs)
    : gitRef.new;

  const {
    data: data,
    isLoading: isMacrobenchLoading,
    error: macrobenchError,
  } = useApiCall<CompareData[]>(
    shouldFetchCompareData
      ? {
          url: `${
            import.meta.env.VITE_API_URL
          }macrobench/compare?new=${gitNewRef}&old=${gitOldRef}`,
          queryKey: ["compare", gitOldRef, gitNewRef],
        }
      : { url: null, queryKey: ["compare", gitOldRef, gitNewRef] }
  );

  useEffect(() => {
    let oldRefName = gitRef.old;
    let newRefName = gitRef.new;
    if (vitessRefs) {
      oldRefName = getRefName(gitRef.old, vitessRefs);
      newRefName = getRefName(gitRef.new, vitessRefs);
    }

    navigate(`?old=${oldRefName}&new=${newRefName}`);
  }, [gitRef.old, gitRef.new, vitessRefs]);

  let formattedData: MacroBenchmarkTableData[] = [];
  if (data !== undefined && data.length > 0) {
    formattedData = formatCompareData(data);
  }

  const [copied, setCopied] = useState(false);
  const copiedTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  // Generates a markdown summary of all workload comparison tables.
  const generateCompareMarkdown = (): string => {
    if (!data || data.length === 0) return "";
    const oldName = vitessRefs
      ? getRefName(gitRef.old, vitessRefs)
      : gitRef.old;
    const newName = vitessRefs
      ? getRefName(gitRef.new, vitessRefs)
      : gitRef.new;
    const lines: string[] = [];

    lines.push("# arewefastyet - Vitess benchmark comparison");
    lines.push("");
    const formatRef = (name: string, hash: string): string => {
      if (name !== hash) {
        return `${name} (\`${hash}\`)`;
      }
      return `\`${hash}\``;
    };
    lines.push(`**Old:** ${formatRef(oldName, gitOldRef)}`);
    lines.push(`**New:** ${formatRef(newName, gitNewRef)}`);
    lines.push("");
    lines.push("---");
    lines.push("");
    // Per-workload tables.
    data.forEach((macro, index) => {
      if (macro.result.missing_results) return;
      const tableData = formattedData[index];
      if (!tableData) return;
      lines.push(`## ${macro.workload}`);
      lines.push("");
      lines.push(
        `| Metric | ${oldName} | ${newName} | P | Delta |`
      );
      lines.push(`|---|---|---|---|---|`);
      const dataKeys = Object.keys(
        tableData
      ) as Array<keyof MacroBenchmarkTableData>;
      dataKeys.forEach((key) => {
        const row = tableData[key];
        const oldVal = `${formatCellValue(key, Number(row.old.center))} (${getRange(row.old.range)})`;
        const newVal = `${formatCellValue(key, Number(row.new.center))} (${getRange(row.new.range)})`;
        const p = fixed(row.p, 3);
        const delta = `${fixed(row.delta, 3)}%`;
        lines.push(
          `| ${row.title} | ${oldVal} | ${newVal} | ${p} | ${delta} |`
        );
      });

      lines.push("");
    });
    return lines.join("\n");
  };

  const handleCopyMarkdown = async () => {
    const markdown = generateCompareMarkdown();

    try {
      await navigator.clipboard.writeText(markdown);
      setCopied(true);

      if (copiedTimerRef.current) {
        clearTimeout(copiedTimerRef.current);
      }

      copiedTimerRef.current = setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error("Failed to copy markdown to clipboard:", err);
    }
  };

  return (
    <>
      <CompareHero
        gitRef={gitRef}
        setGitRef={setGitRef}
        vitessRefs={vitessRefs}
      />
      {macrobenchError && (
        <div className="text-destructive text-center my-2">
          {errorApi}
        </div>
      )}

      <section className="flex flex-col items-center">
        {isMacrobenchLoading && (
          <>
            {[...Array(8)].map((_, index) => {
              return (
                <div key={index} className="w-[80vw] xl:w-[60vw] my-12">
                  <Skeleton className="h-[852px]"></Skeleton>
                </div>
              );
            })}
          </>
        )}
        {!isMacrobenchLoading && data === undefined && (
          <div className="md:text-xl text-primary">
            Chose two commits to compare
          </div>
        )}
        {!isMacrobenchLoading && data !== undefined && data.length > 0 && (
          <>
            <div className="flex justify-end w-[80vw] xl:w-[60vw] mt-8">
              <Button
                variant="outline"
                size="sm"
                className="h-8 border-dashed"
                onClick={handleCopyMarkdown}
              >
                {copied ? (
                  <Check className="mr-2 h-4 w-4 text-primary" />
                ) : (
                  <Copy className="mr-2 h-4 w-4 text-primary" />
                )}
                {copied ? "Copied!" : "Copy as markdown"}
              </Button>
            </div>
            {data.map((macro, index) => {
              return (
                <div className="w-[80vw] xl:w-[60vw] my-12" key={index}>
                  <Card className="border-border">
                    <CardHeader className="flex flex-col gap-4 md:gap-0 md:flex-row justify-between pt-6">
                      <CardTitle className="text-2xl md:text-4xl">
                        {macro.workload}
                      </CardTitle>
                      <Button
                        variant="outline"
                        size="sm"
                        className="h-8 w-fit border-dashed mt-4 md:mt-0"
                        disabled={macro.result.missing_results}
                      >
                        <PlusCircledIcon className="mr-2 h-4 w-4 text-primary" />
                        <Link
                          to={`/macrobench/queries/compare?old=${gitRef.old}&new=${gitRef.new}&workload=${macro.workload}`}
                        >
                          See Query Plan{" "}
                        </Link>
                      </Button>
                    </CardHeader>
                    <CardContent className="w-full p-0">
                      {macro.result.missing_results ? (
                        <div className="text-center md:text-xl text-destructive pb-12">
                          Missing results for this workload
                        </div>
                      ) : (
                        <MacroBenchmarkTable
                          data={formattedData[index]}
                          new={gitRef.new}
                          old={gitRef.old}
                          vitessRefs={vitessRefs}
                        />
                      )}
                    </CardContent>
                  </Card>
                </div>
              );
            })}
          </>
        )}
      </section>
    </>
  );
}
