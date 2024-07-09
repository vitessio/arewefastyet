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

import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
} from "@/components/ui/chart";
import { DailySummarydata } from "@/types";
import PropTypes from "prop-types";
import { Line, LineChart, XAxis } from "recharts";

export type DailySummaryProps = {
  data: DailySummarydata;
  benchmarkType: string;
  setBenchmarktype: (type: string) => void;
};

export default function DailySummary({ data }: DailySummaryProps) {
  type ChartData = { name: string; totalQps: number };
  const chartData: ChartData[] = [];

  const chartConfig = {
    desktop: {
      label: "Total QPS",
      color: "hsl(var(--primary))",
    },
    mobile: {
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

  return (
    <Card className="w-[250px] h-[150px] md:w-[316px] md:h-[186px] hover:scale-105 duration-300 border-border">
      <CardHeader>
        <CardTitle className="font-light text-sm">{data.name}</CardTitle>
      </CardHeader>
      <CardContent className="max-h-[4vh] md:max-h-[7vh]">
        <ChartContainer
          config={chartConfig}
          className="md:h-[80px] h-[60px] w-full"
        >
          <LineChart data={chartData}>
            <XAxis dataKey="time" tickLine={false} axisLine={false} />
            <ChartTooltip cursor={false} content={<CustomTooltip />} />
            <Line
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
      <CardFooter>
        <i className="h-4 w-4 text-foreground fa-solid fa-arrow-right daily--fa-arrow-right"></i>
      </CardFooter>
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

CustomTooltip.propTypes = {
  active: PropTypes.bool,
  payload: PropTypes.array,
  label: PropTypes.string,
};
