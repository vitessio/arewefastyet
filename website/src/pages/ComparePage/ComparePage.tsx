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
import { CompareData } from "@/types";
import useApiCall from "@/utils/Hook";
import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import RingLoader from "react-spinners/RingLoader";
import CompareHero from "./components/CompareHero";
import { PlusCircledIcon } from "@radix-ui/react-icons";


export default function Compare() {
  const urlParams = new URLSearchParams(window.location.search);

  const [gitRef, setGitRef] = useState({
    old: urlParams.get("old") || "",
    new: urlParams.get("new") || "",
  });

  const {
    data: compareData,
    isLoading: isMacrobenchLoading,
    error: macrobenchError,
  } = useApiCall<CompareData[]>(
    `${import.meta.env.VITE_API_URL}macrobench/compare?new=${gitRef.new}&old=${
      gitRef.old
    }`
  );

  console.log({ compareData });

  const formattedData: MacroBenchmarkTableData[] =
    compareData?.map((data: CompareData) => {
      return {
        qpsTotal: {
          title: "QPS Total",
          old: data.result.total_qps.old.center,
          new: data.result.total_qps.new.center,
          p: data.result.total_qps.p,
          delta: data.result.total_qps.delta,
        },
        qpsReads: {
          title: "Reads",
          old: data.result.reads_qps.old.center,
          new: data.result.reads_qps.new.center,
          p: data.result.reads_qps.p,
          delta: data.result.reads_qps.delta,
        },
        qpsWrites: {
          title: "Writes",
          old: data.result.writes_qps.old.center,
          new: data.result.writes_qps.new.center,
          p: data.result.writes_qps.p,
          delta: data.result.writes_qps.delta,
        },
        qpsOther: {
          title: "Other",
          old: data.result.other_qps.old.center,
          new: data.result.other_qps.new.center,
          p: data.result.other_qps.p,
          delta: data.result.other_qps.delta,
        },
        tps: {
          title: "TPS",
          old: data.result.tps.old.center,
          new: data.result.tps.new.center,
          p: data.result.tps.p,
          delta: data.result.tps.delta,
        },
        latency: {
          title: "Latency",
          old: data.result.latency.old.center,
          new: data.result.latency.new.center,
          p: data.result.latency.p,
          delta: data.result.latency.delta,
        },
        totalComponentsCpuTime: {
          title: "Total CPU / Query",
          old: data.result.total_components_cpu_time.old.center,
          new: data.result.total_components_cpu_time.new.center,
          p: data.result.total_components_cpu_time.p,
          delta: data.result.total_components_cpu_time.delta,
        },
        vtgateCpuTime: {
          title: "vtgate",
          old: data.result.components_cpu_time.vtgate.old.center,
          new: data.result.components_cpu_time.vtgate.new.center,
          p: data.result.components_cpu_time.vtgate.p,
          delta: data.result.components_cpu_time.vtgate.delta,
        },
        vttabletCpuTime: {
          title: "vttablet",
          old: data.result.components_cpu_time.vttablet.old.center,
          new: data.result.components_cpu_time.vttablet.new.center,
          p: data.result.components_cpu_time.vttablet.p,
          delta: data.result.components_cpu_time.vttablet.delta,
        },
        totalComponentsMemStatsAllocBytes: {
          title: "Total Allocated / Query",
          old: data.result.total_components_mem_stats_alloc_bytes.old.center,
          new: data.result.total_components_mem_stats_alloc_bytes.new.center,
          p: data.result.total_components_mem_stats_alloc_bytes.p,
          delta: data.result.total_components_mem_stats_alloc_bytes.delta,
        },
        vtgateMemStatsAllocBytes: {
          title: "vtgate",
          old: data.result.components_mem_stats_alloc_bytes.vtgate.old.center,
          new: data.result.components_mem_stats_alloc_bytes.vtgate.new.center,
          p: data.result.components_mem_stats_alloc_bytes.vtgate.p,
          delta: data.result.components_mem_stats_alloc_bytes.vtgate.delta,
        },
        vttabletMemStatsAllocBytes: {
          title: "vttablet",
          old: data.result.components_mem_stats_alloc_bytes.vttablet.old.center,
          new: data.result.components_mem_stats_alloc_bytes.vttablet.new.center,
          p: data.result.components_mem_stats_alloc_bytes.vttablet.p,
          delta: data.result.components_mem_stats_alloc_bytes.vttablet.delta,
        },
      };
    }) || [];

  const navigate = useNavigate();

  useEffect(() => {
    navigate(`?old=${gitRef.old}&new=${gitRef.new}`);
  }, [gitRef.old, gitRef.new]);

  return (
    <>
      <CompareHero gitRef={gitRef} setGitRef={setGitRef} />
      {macrobenchError && (
        <div className="text-red-500 text-center my-2">{macrobenchError}</div>
      )}

      {isMacrobenchLoading && (
        <div className="flex justify-center items-center">
          <RingLoader
            loading={isMacrobenchLoading}
            color="#E77002"
            size={300}
          />
        </div>
      )}

      {!isMacrobenchLoading && compareData && compareData.length > 0 && (
        <section className="flex flex-col items-center">
          <h3 className="my-6 text-primary text-2xl">Macro Benchmarks</h3>
          {compareData.map((macro, index) => {
            return (
              <div className="w-full p-page my-12" key={index}>
                <Card className="border-border">
                  <CardHeader className="flex flex-row justify-between pt-6">
                    <CardTitle className="text-4xl">{macro.type}</CardTitle>
                    <Button variant="outline" size="sm" className="h-8 border-dashed">
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
        </section>
      )}
    </>
  );
}
