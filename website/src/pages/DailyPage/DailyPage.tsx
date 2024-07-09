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

import { useState, useEffect } from "react";
import RingLoader from "react-spinners/RingLoader";
import { useNavigate } from "react-router-dom";
import useApiCall from "@/utils/Hook";

import DailySummary from "@/common/DailySummary";
import { MacroData, MacroDataValue, Workloads } from "@/types";

import { secondToMicrosecond } from "@/utils/Utils";
import DailyHero from "./components/DailyHero";
import useDailySummaryData from "@/hooks/useDailySummaryData";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { CartesianGrid, XAxis, Line, LineChart } from "recharts";
import {
  ChartTooltipContent,
  ChartTooltip,
  ChartContainer,
} from "@/components/ui/chart";
import ResponsiveChart from "./components/Chart";

interface DailySummarydata {
  name: string;
  data: { total_qps: MacroDataValue }[];
}

type NumberDataPoint = { x: string; y: number };
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
    data: dataDaily,
    error: dailyError,
    textLoading: dailyTextLoading,
  } = useApiCall<MacroData>(
    `${import.meta.env.VITE_API_URL}daily?workloads=${benchmarkType}`
  );

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

  const CPUTimeData: { id: string; data: StringDataPoint[] }[] = [
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

  const workloads: Workloads[] = [
    "OLTP",
    "OLTP-READONLY",
    "OLTP-SET",
    "TPCC",
    "TPCC_FK",
    "TPCC_UNSHARDED",
    "TPCC_FK_UNMANAGED",
  ];

  const { dataDailySummary, isLoadingDailySummary, errorDailySummary } =
    useDailySummaryData(workloads);

  const [expandedStates, setExpandedStates] = useState(
    Array(allChartData.length).fill(true)
  );

  const toggleExpand = (index: number) => {
    setExpandedStates((prevStates) => {
      const newStates = [...prevStates];
      newStates[index] = !newStates[index];
      return newStates;
    });
  };

  return (
    <>
      <DailyHero />
      <section className="flex p-page flex-wrap justify-center gap-12 p-4">
        {isLoadingDailySummary && (
          <>
            {workloads.map((_, index) => {
              return (
                <Skeleton
                  key={index}
                  className="w-[250px] h-[150px] md:w-[316px] md:h-[186px] rounded-lg"
                />
              );
            })}
          </>
        )}

        {!errorDailySummary &&
          dataDailySummary &&
          dataDailySummary.length > 0 && (
            <>
              {dataDailySummary.map((dailySummary, index) => {
                return (
                  <DailySummary
                    key={index}
                    data={dailySummary}
                    benchmarkType={benchmarkType}
                    setBenchmarktype={setBenchmarktype}
                  />
                );
              })}

              <Separator className="mx-auto w-[80%] foreground" />
            </>
          )}
      </section>

      {dailyTextLoading && (
        <div className="flex justify-center w-full my-16">
          <RingLoader loading={dailyTextLoading} color="#E77002" size={300} />
        </div>
      )}

      {!dailyTextLoading && benchmarkType !== "" && (
        <section className="p-page mt-12 flex flex-col gap-y-8">
          {allChartData.map((chartData, index) => (
            <Card key={index} className="w-full">
              <CardHeader
                className="cursor-pointer hover:bg-muted duration-300 py-0"
                onClick={() => toggleExpand(index)}
              >
                <div className="flex items-center justify-between">
                  <CardTitle className="my-10 text-xl font-medium text-primary">
                    {chartData.title}
                  </CardTitle>
                  <i className={`h-4 w-4 text-foreground fa-solid ${expandedStates[index] ? 'fa-chevron-up' : 'fa-chevron-down'} daily--fa-chevron-right`}></i>
                </div>
              </CardHeader>
              {expandedStates[index] && (
                <CardContent>
                  <div className="relative w-full h-[500px]">
                    <ResponsiveChart
                      data={chartData.data as any}
                      title={chartData.title}
                      colors={chartData.colors}
                      isFirstChart={index === 0}
                    />
                  </div>
                </CardContent>
              )}
            </Card>
          ))}
        </section>
      )}

      {(errorDailySummary || dailyError) && (
        <div className="text-red-500 text-center my-10">{dailyError}</div>
      )}
    </>
  );
}
