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

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ChartContainer, ChartTooltip } from "@/components/ui/chart";
import { Skeleton } from "@/components/ui/skeleton";
import useApiCall from "@/hooks/useApiCall";
import { MacroData } from "@/types";
import { errorApi, formatGitRef } from "@/utils/Utils";
import { useEffect, useState } from "react";
import { CartesianGrid, Legend, Line, LineChart, XAxis, YAxis } from "recharts";

export type DailyChartsProps = {
  workload: string;
};

type DailyDataType = {
  gitRef: string;
  qpsReads: number;
  qpsWrites: number;
  qpsOther: number;
  qpsTotal: number;
  tps: number;
  latency: number;
  cpuTimeTotal: number;
  cpuTimeVtgate: number;
  cpuTimeVttablet: number;
  memBytesTotal: number;
  memBytesVtgate: number;
  memBytesVttablet: number;
};

const chartConfig: { [key: string]: { label: string; color: string } } = {
  qpsReads: {
    label: "Reads",
    color: "hsl(var(--chart-qpsReads))",
  },
  qpsWrites: {
    label: "Writes",
    color: "hsl(var(--chart-qpsWrites))",
  },
  qpsOther: {
    label: "Other",
    color: "hsl(var(--chart-qpsOther))",
  },
  qpsTotal: {
    label: "Total",
    color: "hsl(var(--chart-qpsTotal))",
  },
  tps: {
    label: "TPS",
    color: "hsl(var(--chart-tps))",
  },
  latency: {
    label: "Latency",
    color: "hsl(var(--chart-latency))",
  },
  cpuTimeTotal: {
    label: "Total",
    color: "hsl(var(--chart-cpuTimeTotal))",
  },
  cpuTimeVtgate: {
    label: "Vtgate",
    color: "hsl(var(--chart-cpuTimeVtgate))",
  },
  cpuTimeVttablet: {
    label: "Vttablet",
    color: "hsl(var(--chart-cpuTimeVttablet))",
  },
  memBytesVtgate: {
    label: "Vtgate",
    color: "hsl(var(--chart-memBytesVtgate))",
  },
  memBytesVttablet: {
    label: "Vttablet",
    color: "hsl(var(--chart-memBytesVttablet))",
  },
  memBytesTotal: {
    label: "Total",
    color: "hsl(var(--chart-memBytesTotal))",
  },
};

const chartMetadas = [
  {
    title: "QPS (Queries per second)",
    metrics: [
      { dataKey: "qpsReads", legend: "Reads" },
      { dataKey: "qpsWrites", legend: "Writes" },
      { dataKey: "qpsOther", legend: "Other" },
      { dataKey: "qpsTotal", legend: "Total" },
    ],
  },
  {
    title: "TPS (Transactions per second)",
    metrics: [{ dataKey: "tps", legend: "TPS" }],
  },
  {
    title: "Latency (ms)",
    metrics: [{ dataKey: "latency", legend: "Latency" }],
  },
  {
    title: "CPU / query (μs)",
    metrics: [
      { dataKey: "cpuTimeVtgate", legend: "Vtgate" },
      { dataKey: "cpuTimeVttablet", legend: "Vttablet" },
      { dataKey: "cpuTimeTotal", legend: "Total" },
    ],
  },
  {
    title: "Allocated / query (bytes)",
    metrics: [
      { dataKey: "memBytesVtgate", legend: "Vtgate" },
      { dataKey: "memBytesVttablet", legend: "Vttablet" },
      { dataKey: "memBytesTotal", legend: "Total" },
    ],
  },
];

export default function DailyCharts(props: DailyChartsProps) {
  const { workload } = props;

  const {
    data: dataDaily,
    error: dailyError,
    isLoading: dailyLoading,
  } = useApiCall<MacroData[]>({
    url: `${import.meta.env.VITE_API_URL}daily?workload=${workload}`,
    queryKey: ["dailyWorkload", workload],
  });

  let chartData: DailyDataType[] = [];

  if (dataDaily !== undefined && dataDaily.length > 0) {
    chartData = dataDaily.map((item) => ({
      gitRef: formatGitRef(item.git_ref),
      qpsReads: item.reads_qps.center,
      qpsWrites: item.writes_qps.center,
      qpsOther: item.other_qps.center,
      qpsTotal: item.total_qps.center,
      tps: item.tps.center,
      latency: item.latency.center,
      cpuTimeTotal: Number(
        (item.total_components_cpu_time.center * 1000000).toFixed(2)
      ),
      cpuTimeVtgate: Number(
        (item.components_cpu_time.vtgate.center * 1000000).toFixed(2)
      ),
      cpuTimeVttablet: Number(
        (item.components_cpu_time.vttablet.center * 1000000).toFixed(2)
      ),
      memBytesTotal: item.total_components_mem_stats_alloc_bytes.center,
      memBytesVtgate: item.components_mem_stats_alloc_bytes.vtgate.center,
      memBytesVttablet: item.components_mem_stats_alloc_bytes.vttablet.center,
    }));
  }

  const [expandedStates, setExpandedStates] = useState<boolean[]>(
    Array(chartMetadas.length).fill(true)
  );

  const toggleExpand = (index: number) => {
    setExpandedStates((prevStates) => {
      const newStates = [...prevStates];
      newStates[index] = !newStates[index];
      return newStates;
    });
  };

  useEffect(() => {
    setExpandedStates(Array(chartMetadas.length).fill(true));
  }, [chartMetadas.length]);

  return (
    <>
      <section className="p-page my-12 flex flex-col gap-y-8">
        <div className="flex flex-col items-center gap-4">
          <h2 className="text-4xl md:text-6xl font-semibold text-primary mb-4">
            {workload}
          </h2>
        </div>
        {dailyLoading ? (
          chartMetadas.map((_, index) => (
            <Skeleton key={index} className="w-full border-border h-[400px]" />
          ))
        ) : dailyError || !chartData || chartData.length === 0 ? (
          <div className="text-destructive text-center my-10">
            {errorApi}
          </div>
        ) : (
          chartMetadas.map((chartMetadata, chartMetadataIndex) => (
            <Card key={chartMetadataIndex} className="w-full border-border">
              <CardHeader
                className="cursor-pointer hover:bg-muted duration-300 py-0"
                onClick={() => toggleExpand(chartMetadataIndex)}
              >
                <div className="flex items-center justify-between">
                  <CardTitle className="my-10 text-xl font-medium text-primary">
                    {chartMetadata.title}
                  </CardTitle>
                  <i
                    className={`h-4 w-4 text-foreground fa-solid ${
                      expandedStates[chartMetadataIndex]
                        ? "fa-chevron-up"
                        : "fa-chevron-down"
                    } daily--fa-chevron-right`}
                  ></i>
                </div>
              </CardHeader>
              {expandedStates[chartMetadataIndex] && (
                <CardContent>
                  <div className="relative w-full h-[400px]">
                    <ChartContainer
                      config={chartConfig}
                      className="w-full h-[400px] mx-auto mt-10"
                    >
                      <LineChart data={chartData}>
                        <XAxis
                          dataKey="gitRef"
                          tickLine={true}
                          axisLine={true}
                        />
                        <YAxis />
                        <CartesianGrid vertical={true} />
                        <ChartTooltip
                          cursor={true}
                          content={<CustomTooltip />}
                        />
                        {chartMetadata.metrics.map((metric, dataKeyIndex) => (
                          <Line
                            key={dataKeyIndex}
                            className="pb-0 pt-0"
                            dataKey={metric.dataKey}
                            type="natural"
                            label={chartConfig[metric.dataKey].label}
                            stroke={chartConfig[metric.dataKey].color}
                            strokeWidth={2}
                            dot={{
                              fill: chartConfig[metric.dataKey].color,
                            }}
                            activeDot={{
                              r: 6,
                            }}
                          />
                        ))}
                        <Legend
                          payload={chartMetadata.metrics.map((metric) => ({
                            id: metric.dataKey,
                            type: "circle",
                            value: metric.legend,
                            color: chartConfig[metric.dataKey].color,
                          }))}
                        />
                      </LineChart>
                    </ChartContainer>
                  </div>
                </CardContent>
              )}
            </Card>
          ))
        )}
      </section>
    </>
  );
}

const CustomTooltip = ({
  active,
  payload,
  label,
}: {
  active?: boolean;
  payload?: { color: string; name: string; value: number }[];
  label?: string;
}) => {
  if (payload) {
    // Sort payload to have the same order as the chartConfig
    payload.sort((a, b) => {
      return (
        Object.keys(chartConfig).indexOf(a.name) -
        Object.keys(chartConfig).indexOf(b.name)
      );
    });
  }
  if (active && payload && payload.length) {
    return (
      <div className="custom-tooltip p-2 border border-border shadow-lg rounded bg-background">
        <p className="label font-bold mb-2">Commit: {label}</p>
        {payload.map((entry, index) => (
          <p
            className="w-full label flex items-center mb-1"
            key={`item-${index}`}
          >
            <span
              className="inline-block w-2 h-2 mr-2 rounded-full"
              style={{ backgroundColor: entry.color }}
            ></span>
            <span className="w-full flex items-center justify-between">
              <span>{`${chartConfig[entry.name].label}: `}</span>{" "}
              <span> {`${entry.value.toFixed(0)}`}</span>
            </span>
          </p>
        ))}
      </div>
    );
  }

  return null;
};
