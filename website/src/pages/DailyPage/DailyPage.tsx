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

import React, { useState, useEffect } from "react";
import RingLoader from "react-spinners/RingLoader";
import { useNavigate } from "react-router-dom";
import useApiCall from "../../utils/Hook";

import ResponsiveChart from "./components/Chart";
import DailySummary from "./components/DailySummary";
import Hero from "./components/Hero";
import { MacroData, MacroDataValue } from "@/types";

import { secondToMicrosecond } from "../../utils/Utils";

interface DailySummarydata {
  name: string;
  data : { total_qps: MacroDataValue }[];
}

type NumberDataPoint = { x: string; y: number};
type StringDataPoint = { x: string; y: string };

type ChartDataItem =
  | { id: string; data: NumberDataPoint[] }
  | { id: string; data: StringDataPoint[] };

export default function DailyPage() {
  const urlParams = new URLSearchParams(window.location.search);
  const [benchmarkType, setBenchmarktype] = useState<string>(
    urlParams.get("type") ?? "OLTP"
  );

  const {
    data: dataDailySummary,
    isLoading: isLoadingDailySummary,
    error: errorDailySummary,
  } = useApiCall<DailySummarydata>(`${import.meta.env.VITE_API_URL}daily/summary`);

  const {
    data: dataDaily,
    error: dailyError,
    textLoading: dailyTextLoading,
  } = useApiCall<MacroData>(`${import.meta.env.VITE_API_URL}daily?type=${benchmarkType}`);

  const navigate = useNavigate();

  useEffect(() => {
    navigate(`?type=${benchmarkType}`);
  }, [benchmarkType]);

  const TPSData: { id: string; data: NumberDataPoint[] }[] = [
    {
      id: "TPS",
      data: [],
    },
  ];

  const QPSData: { id: string; data: NumberDataPoint[] }[] = [
    {
      id: "Reads",
      data: [],
    },
    {
      id: "Total",
      data: [],
    },

    {
      id: "Writes",
      data: [],
    },

    {
      id: "Other",
      data: [],
    },
  ];

  const latencyData: { id: string; data: NumberDataPoint[] }[] = [
    {
      id: "Latency",
      data: [],
    },
  ];

  const CPUTimeData : { id: string; data: StringDataPoint[] }[]= [
    {
      id: "Total",
      data: [],
    },
    {
      id: "vtgate",
      data: [],
    },
    {
      id: "vttablet",
      data: [],
    },
  ];

  const MemBytesData: { id: string; data: NumberDataPoint[] }[] = [
    {
      id: "Total",
      data: [],
    },
    {
      id: "vtgate",
      data: [],
    },
    {
      id: "vttablet",
      data: [],
    },
  ];

  for (const item of dataDaily) {
    const xValue = item.git_ref.slice(0, 8);

    // TPS Data

    TPSData[0].data.push({
      x: xValue,
      y: item.tps.center,
    });

    // QPS Data

    QPSData[0].data.push({
      x: xValue,
      y: item.reads_qps.center,
    });

    QPSData[1].data.push({
      x: xValue,
      y: item.total_qps.center,
    });

    QPSData[2].data.push({
      x: xValue,
      y: item.writes_qps.center,
    });

    QPSData[3].data.push({
      x: xValue,
      y: item.other_qps.center,
    });

    // Latency Data

    latencyData[0].data.push({
      x: xValue,
      y: item.latency.center,
    });

    // CPUTime Data

    CPUTimeData[0].data.push({
      x: xValue,
      y: secondToMicrosecond(item.total_components_cpu_time.center),
    });

    CPUTimeData[1].data.push({
      x: xValue,
      y: secondToMicrosecond(item.components_cpu_time.vtgate.center),
    });

    CPUTimeData[2].data.push({
      x: xValue,
      y: secondToMicrosecond(item.components_cpu_time.vttablet.center),
    });

    //MemStatsAllocBytes Data

    MemBytesData[0].data.push({
      x: xValue,
      y: item.total_components_mem_stats_alloc_bytes.center,
    });

    MemBytesData[1].data.push({
      x: xValue,
      y: item.components_mem_stats_alloc_bytes.vtgate.center,
    });

    MemBytesData[2].data.push({
      x: xValue,
      y: item.components_mem_stats_alloc_bytes.vttablet.center,
    });
  }

  const allChartData: {
    data: ChartDataItem[];
    title: string;
    colors: string[];
  }[] = [
    {
      data: QPSData,
      title: "QPS (Queries per second)",
      colors: ["#fad900", "orange", "brown", "purple"],
    },
    {
      data: TPSData,
      title: "TPS (Transactions per second)",
      colors: ["#fad900"],
    },
    {
      data: latencyData,
      title: "Latency (ms)",
      colors: ["#fad900"],
    },
    {
      data: CPUTimeData,
      title: "CPU / query (μs)",
      colors: ["#fad900", "orange", "brown"],
    },
    {
      data: MemBytesData,
      title: "Allocated / query (bytes)",
      colors: ["#fad900", "orange", "brown"],
    },
  ];

  return (
    <>
      <Hero />

      <figure className="p-page w-full">
        <div className="border-front border" />
      </figure>

      {isLoadingDailySummary && (
        <div className="flex justify-center w-full my-16">
          <RingLoader
            loading={isLoadingDailySummary}
            color="#E77002"
            size={300}
          />
        </div>
      )}

      {!errorDailySummary && dataDailySummary && dataDailySummary.length > 0 &&(
        <>
          <section className="flex p-page justif-center flex-wrap gap-10 py-10">
            {dataDailySummary.map((dailySummary, index) => {
              return (
                <DailySummary
                  key={index}
                  data={dailySummary}
                  setBenchmarktype={setBenchmarktype}
                  benchmarkType={benchmarkType}
                />
              );
            })}
          </section>

          <figure className="p-page w-full">
            <div className="border-front border" />
          </figure>

          {!dailyTextLoading && benchmarkType !== "" && (
            <section className="p-page mt-12 flex flex-col gap-y-8">
              {allChartData.map((chartData, index) => (
                <div key={index} className="relative w-full h-[500px]">
                  <ResponsiveChart
                    data={chartData.data as any}
                    title={chartData.title}
                    colors={chartData.colors}
                    isFirstChart={index === 0}
                  />
                </div>
              ))}
            </section>
          )}
        </>
      )}

      {(errorDailySummary || dailyError) && (
        <div className="text-red-500 text-center my-10">{dailyError}</div>
      )}
    </>
  );
}
