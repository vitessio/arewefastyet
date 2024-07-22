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

import MacroBenchmarkTable, {
  MacroBenchmarkTableData,
} from "@/common/MacroBenchmarkTable";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { CompareData } from "@/types";
import useApiCall from "@/utils/Hook";
import { PlusCircledIcon } from "@radix-ui/react-icons";
import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import CompareHero from "./components/CompareHero";

export default function Compare() {
  const urlParams = new URLSearchParams(window.location.search);

  const [gitRef, setGitRef] = useState({
    old: urlParams.get("old") || "",
    new: urlParams.get("new") || "",
  });

  const navigate = useNavigate();

  useEffect(() => {
    navigate(`?old=${gitRef.old}&new=${gitRef.new}`);
  }, [gitRef.old, gitRef.new]);

  let {
    data: data,
    isLoading: isMacrobenchLoading,
    error: macrobenchError,
  } = useApiCall<CompareData[]>(
    `${import.meta.env.VITE_API_URL}macrobench/compare?new=${gitRef.new}&old=${
      gitRef.old
    }`
  );

  let formattedData: MacroBenchmarkTableData[] = [];

  if (data !== null && data.length > 0) {
    formattedData =
      data.map((data: CompareData) => {
        return {
          qpsTotal: {
            title: "QPS Total",
            old: data.result.total_qps.old,
            new: data.result.total_qps.new,
            p: data.result.total_qps.p,
            delta: data.result.total_qps.delta,
            insignificant: data.result.total_qps.insignificant,
          },
          qpsReads: {
            title: "Reads",
            old: data.result.reads_qps.old,
            new: data.result.reads_qps.new,
            p: data.result.reads_qps.p,
            delta: data.result.reads_qps.delta,
            insignificant: data.result.reads_qps.insignificant,
          },
          qpsWrites: {
            title: "Writes",
            old: data.result.writes_qps.old,
            new: data.result.writes_qps.new,
            p: data.result.writes_qps.p,
            delta: data.result.writes_qps.delta,
            insignificant: data.result.writes_qps.insignificant,
          },
          qpsOther: {
            title: "Other",
            old: data.result.other_qps.old,
            new: data.result.other_qps.new,
            p: data.result.other_qps.p,
            delta: data.result.other_qps.delta,
            insignificant: data.result.other_qps.insignificant,
          },
          tps: {
            title: "TPS",
            old: data.result.tps.old,
            new: data.result.tps.new,
            p: data.result.tps.p,
            delta: data.result.tps.delta,
            insignificant: data.result.tps.insignificant,
          },
          latency: {
            title: "P95 Latency",
            old: data.result.latency.old,
            new: data.result.latency.new,
            p: data.result.latency.p,
            delta: data.result.latency.delta,
            insignificant: data.result.latency.insignificant,
          },
          errors: {
            title: "Errors / Second",
            old: data.result.errors.old,
            new: data.result.errors.new,
            p: data.result.errors.p,
            delta: data.result.errors.delta,
            insignificant: data.result.errors.insignificant,
          },
          totalComponentsCpuTime: {
            title: "Total CPU / Query",
            old: data.result.total_components_cpu_time.old,
            new: data.result.total_components_cpu_time.new,
            p: data.result.total_components_cpu_time.p,
            delta: data.result.total_components_cpu_time.delta,
            insignificant: data.result.total_components_cpu_time.insignificant,
          },
          vtgateCpuTime: {
            title: "vtgate",
            old: data.result.components_cpu_time.vtgate.old,
            new: data.result.components_cpu_time.vtgate.new,
            p: data.result.components_cpu_time.vtgate.p,
            delta: data.result.components_cpu_time.vtgate.delta,
            insignificant: data.result.components_cpu_time.vtgate.insignificant,
          },
          vttabletCpuTime: {
            title: "vttablet",
            old: data.result.components_cpu_time.vttablet.old,
            new: data.result.components_cpu_time.vttablet.new,
            p: data.result.components_cpu_time.vttablet.p,
            delta: data.result.components_cpu_time.vttablet.delta,
            insignificant:
              data.result.components_cpu_time.vttablet.insignificant,
          },
          totalComponentsMemStatsAllocBytes: {
            title: "Total Allocated / Query",
            old: data.result.total_components_mem_stats_alloc_bytes.old,
            new: data.result.total_components_mem_stats_alloc_bytes.new,
            p: data.result.total_components_mem_stats_alloc_bytes.p,
            delta: data.result.total_components_mem_stats_alloc_bytes.delta,
            insignificant:
              data.result.total_components_mem_stats_alloc_bytes.insignificant,
          },
          vtgateMemStatsAllocBytes: {
            title: "vtgate",
            old: data.result.components_mem_stats_alloc_bytes.vtgate.old,
            new: data.result.components_mem_stats_alloc_bytes.vtgate.new,
            p: data.result.components_mem_stats_alloc_bytes.vtgate.p,
            delta: data.result.components_mem_stats_alloc_bytes.vtgate.delta,
            insignificant:
              data.result.components_mem_stats_alloc_bytes.vtgate.insignificant,
          },
          vttabletMemStatsAllocBytes: {
            title: "vttablet",
            old: data.result.components_mem_stats_alloc_bytes.vttablet.old,
            new: data.result.components_mem_stats_alloc_bytes.vttablet.new,
            p: data.result.components_mem_stats_alloc_bytes.vttablet.p,
            delta: data.result.components_mem_stats_alloc_bytes.vttablet.delta,
            insignificant:
              data.result.components_mem_stats_alloc_bytes.vttablet
                .insignificant,
          },
        };
      }) || [];
  }

  return (
    <>
      <CompareHero gitRef={gitRef} setGitRef={setGitRef} />
      {macrobenchError && (
        <div className="text-red-500 text-center my-2">{macrobenchError}</div>
      )}

      <section className="flex flex-col items-center">
      {isMacrobenchLoading && (
        <>
         {[...Array(8)].map((_, index) => {
            return (
              <div key={index} className="w-full p-page lg:w-[60vw] my-12">
                <Skeleton className="h-[852px]"></Skeleton>
              </div>
            );
          })}
        </>
      )}
        {!isMacrobenchLoading && data !== null && data.length > 0 && (
          <>
            {data.map((macro, index) => {
              return (
                <div className="w-full p-page lg:w-[60vw] my-12" key={index}>
                  <Card className="border-border">
                    <CardHeader className="flex flex-col gap-4 md:gap-0 md:flex-row justify-between pt-6">
                      <CardTitle className="text-2xl md:text-4xl">
                        {macro.type}
                      </CardTitle>
                      <Button
                        variant="outline"
                        size="sm"
                        className="h-8 w-fit border-dashed mt-4 md:mt-0"
                      >
                        <PlusCircledIcon className="mr-2 h-4 w-4 text-primary" />
                        <Link
                          to={`/macrobench/queries/compare?ltag=${gitRef.old}&rtag=${gitRef.new}&type=${macro.type}`}
                        >
                          See Query Plan{" "}
                        </Link>
                      </Button>
                    </CardHeader>
                    <CardContent className="w-full p-0">
                      <MacroBenchmarkTable
                        data={formattedData[index]}
                        newGitRef={gitRef.new}
                        oldGitRef={gitRef.old}
                      />
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
