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

import ResponsiveChart from "./components/Chart";
import DailySummary from "./components/DailySummary";
import Hero from "./components/Hero";

import { formatByteForGB } from "../../utils";
import useApiCall from "../../hooks/useApiCall";

export default function DailyPage() {
  const urlParams = new URLSearchParams(window.location.search);
  const [benchmarkType, setBenchmarkType] = useState(
    urlParams.get("type") == null ? "" : urlParams.get("type")
  );

  const [summary, summaryLoading, summaryError] = useApiCall("/daily/summary");

  const [data, error, loading] = useApiCall(`/daily`);

  const navigate = useNavigate();

  useEffect(() => {
    navigate(`?type=${benchmarkType}`);
  }, [benchmarkType]);

  interface DataPoint {
    x: string;
    y: number | string;
  }

  const TPSData: { id: "TPS"; data: DataPoint[] }[] = [
    {
      id: "TPS",
      data: [],
    },
  ];

  const QPSData: {
    id: "Reads" | "Total" | "Writes" | "Other";
    data: DataPoint[];
  }[] = [
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

  const latencyData: { id: "Latency"; data: DataPoint[] }[] = [
    {
      id: "Latency",
      data: [],
    },
  ];

  const CPUTimeData: {
    id: "Total" | "vtgate" | "vttablet";
    data: DataPoint[];
  }[] = [
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

  const MemBytesData: {
    id: "Total" | "vtgate" | "vttablet";
    data: DataPoint[];
  }[] = [
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

  if (data)
    for (const item of data) {
      const xValue = item.GitRef.slice(0, 8);

      // TPS Data

      TPSData[0].data.push({
        x: xValue,
        y: item.Result.tps,
      });

      // QPS Data

      QPSData[0].data.push({
        x: xValue,
        y: item.Result.qps.reads,
      });

      QPSData[1].data.push({
        x: xValue,
        y: item.Result.qps.total,
      });

      QPSData[2].data.push({
        x: xValue,
        y: item.Result.qps.writes,
      });

      QPSData[3].data.push({
        x: xValue,
        y: item.Result.qps.other,
      });

      // Latency Data

      latencyData[0].data.push({
        x: xValue,
        y: item.Result.latency,
      });

      // CPUTime Data

      CPUTimeData[0].data.push({
        x: xValue,
        y: item.Metrics.TotalComponentsCPUTime,
      });

      CPUTimeData[1].data.push({
        x: xValue,
        y: item.Metrics.ComponentsCPUTime.vtgate,
      });

      CPUTimeData[2].data.push({
        x: xValue,
        y: item.Metrics.ComponentsCPUTime.vttablet,
      });

      //MemStatsAllocBytes Data

      MemBytesData[0].data.push({
        x: xValue,
        y: formatByteForGB(item.Metrics.TotalComponentsMemStatsAllocBytes),
      });

      MemBytesData[1].data.push({
        x: xValue,
        y: formatByteForGB(item.Metrics.ComponentsMemStatsAllocBytes.vtgate),
      });

      MemBytesData[2].data.push({
        x: xValue,
        y: formatByteForGB(item.Metrics.ComponentsMemStatsAllocBytes.vttablet),
      });
    }

  const allChartData = [
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
      title: "Latency (Milliseconds)",
      colors: ["#fad900"],
    },
    {
      data: CPUTimeData,
      title: "CPU Time",
      colors: ["#fad900", "orange", "brown"],
    },
    {
      data: MemBytesData,
      title: "Allocated Bytes",
      colors: ["#fad900", "orange", "brown"],
    },
  ];

  return (
    <>
      <Hero />

      <figure className="p-page w-full">
        <div className="border-front border" />
      </figure>

      {summaryLoading && (
        <div className="flex justify-center w-full my-16">
          <RingLoader loading={summaryLoading} color="#E77002" size={300} />
        </div>
      )}

      {!summaryError && summary && (
        <>
          <section className="flex p-page justify-center flex-wrap gap-10 py-10">
            {summary.map((dailySummary, index) => {
              return (
                <DailySummary
                  key={index}
                  data={summary}
                  setBenchmarkType={setBenchmarkType}
                  benchmarkType={benchmarkType}
                />
              );
            })}
          </section>

          <figure className="p-page w-full">
            <div className="border-front border" />
          </figure>

          {!loading && benchmarkType !== "" && (
            <section className="p-page mt-12 flex flex-col gap-y-8">
              {allChartData.map((chartData, index) => (
                <div key={index} className="relative w-full h-[500px]">
                  <ResponsiveChart
                    /* @ts-ignore */
                    data={chartData.data}
                    title={chartData.title}
                    colors={chartData.colors}
                  />
                </div>
              ))}
            </section>
          )}
        </>
      )}

      {(error || summaryError) && (
        <div className="text-red-500 text-center my-10">
          {error || summaryError}
        </div>
      )}
    </>
  );
}
