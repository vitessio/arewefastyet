import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ChartContainer, ChartTooltip } from "@/components/ui/chart";
import { Skeleton } from "@/components/ui/skeleton";
import { MacroData, Workloads } from "@/types";
import useApiCall from "@/utils/Hook";
import { useEffect, useState } from "react";
import { CartesianGrid, Legend, Line, LineChart, XAxis, YAxis } from "recharts";

export type DailyChartsProps = {
  benchmarkType: string;
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

export default function DailyCharts(props: DailyChartsProps) {
  const { benchmarkType } = props;
  const workloads: Workloads[] = [
    "OLTP",
    "OLTP-READONLY",
    "OLTP-SET",
    "TPCC",
    "TPCC_FK",
    "TPCC_UNSHARDED",
    "TPCC_FK_UNMANAGED",
  ];

  const workloadsQuery = workloads.join("&workloads=");

  const {
    data: dataDaily,
    error: dailyError,
    isLoading: dailyLoading,
  } = useApiCall<MacroData>(
    `${import.meta.env.VITE_API_URL}daily?workloads=${workloadsQuery}`
  );

  let chartData: DailyDataType[] = [];

  if (dataDaily) {
    chartData = dataDaily.map((item) => ({
      gitRef: item.git_ref.slice(0, 8),
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

  const chartMetadas = [
    {
      title: "QPS (Queries per second)",
      dataKeys: ["qpsReads", "qpsTotal", "qpsWrites", "qpsOther"],
    },
    {
      title: "TPS (Transactions per second)",
      dataKeys: ["tps"],
    },
    {
      title: "Latency (ms)",
      dataKeys: ["latency"],
    },
    {
      title: "CPU / query (μs)",
      dataKeys: ["cpuTimeTotal", "cpuTimeVtgate", "cpuTimeVttablet"],
    },
    {
      title: "Allocated / query (bytes)",
      dataKeys: ["memBytesTotal", "memBytesVtgate", "memBytesVttablet"],
    },
  ];

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
            {benchmarkType}
          </h2>
        </div>
        {dailyLoading ? (
          chartMetadas.map((_, index) => (
            <Skeleton key={index} className="w-full border-border h-[400px]" />
          ))
        ) : dailyError || !chartData || chartData.length === 0 ? (
          <div className="text-red-500 text-center my-10">
            {dailyError || "No data available"}
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
                        {chartMetadata.dataKeys.map((dataKey, dataKeyIndex) => (
                          <Line
                            key={dataKeyIndex}
                            className="pb-0 pt-0"
                            dataKey={dataKey}
                            type="natural"
                            label={chartConfig[dataKey].label}
                            stroke={chartConfig[dataKey].color}
                            strokeWidth={2}
                            dot={{
                              fill: chartConfig[dataKey].color,
                            }}
                            activeDot={{
                              r: 6,
                            }}
                          />
                        ))}
                        <Legend />
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
