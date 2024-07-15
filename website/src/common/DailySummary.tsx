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
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
} from "@/components/ui/chart";
import { DailySummarydata } from "@/types";
import { Line, LineChart, XAxis } from "recharts";
import { twMerge } from "tailwind-merge";

export type DailySummaryProps = {
  data: DailySummarydata;
  workload: string;
  setWorkload: (type: string) => void;
};

export default function DailySummary(props: DailySummaryProps) {
  const { data, workload, setWorkload } = props;
  type ChartData = { name: string; totalQps: number };
  const chartData: ChartData[] = [];

  const chartConfig = {
    desktop: {
      label: "Total QPS",
      color: "hsl(var(--primary))",
    },
  } satisfies ChartConfig;

  if (data.data !== null) {
    data.data.map((item) => ({
      totalQps: chartData.push({
        name: "Total QPS",
        totalQps: item.total_qps.center,
      }),
    }));
  }

  const getBenchmarkType = () => {
    setWorkload(data.name);
  };

  return (
    <Card
      className={twMerge(
        "w-[310px] h-[124px] md:w-[316px] md:h-[124px] hover:scale-105 duration-300 hover:bg-muted border-border",
        workload === data.name && "border-2 border-front"
      )}
      onClick={() => getBenchmarkType()}
    >
      <CardHeader className="flex flex-row justify-between">
        <CardTitle className="font-light text-sm">{data.name}</CardTitle>
        <i className="h-4 w-4 text-foreground fa-solid fa-arrow-right daily--fa-arrow-right"></i>
      </CardHeader>
      <CardContent className="max-h-[4vh] md:max-h-[7vh] pt-0 pb-0">
        <ChartContainer
          config={chartConfig}
          className="md:h-[80px] h-[60px] w-full"
        >
          <LineChart data={chartData}>
            <XAxis dataKey="time" tickLine={false} axisLine={false} />
            <ChartTooltip cursor={false} content={<CustomTooltip />} />
            <Line
              className="pb-0 pt-0"
              dataKey="totalQps"
              type="monotone"
              label="Total QPS"
              stroke="var(--color-desktop)"
              strokeWidth={1}
              dot={{
                fill: "var(--color-desktop)",
              }}
              activeDot={{
                r: 6,
              }}
            />
          </LineChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}

const CustomTooltip = ({
  active,
  payload,
}: {
  active?: boolean;
  payload?: { value: number }[];
}) => {
  if (active && payload && payload.length) {
    return (
      <div className="custom-tooltip p-2 border border-border shadow-lg rounded">
        <p className="label flex items-center">
          <span
            className="inline-block w-2 h-2 mr-2"
            style={{ backgroundColor: "hsl(var(--primary))" }}
          ></span>
          {`Total QPS: ${payload[0].value.toFixed(0)}`}
        </p>
      </div>
    );
  }

  return null;
};
