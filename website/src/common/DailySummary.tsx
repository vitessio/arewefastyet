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
  ChartTooltipContent,
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
  const chartData: any[] | undefined = [];

  const chartConfig = {
    desktop: {
      label: "Desktop",
      color: "hsl(var(--primary))",
    },
    mobile: {
      label: "Mobile",
      color: "hsl(var(--primary))",
    },
  } satisfies ChartConfig;

  if (data.data !== null) {
    data.data.map((item) => ({
      totalQps: chartData.push({
        totalQps: item.total_qps.center / 2,
      }),
    }));
  }

  return (
    <Card className="w-[250px] h-[150px] md:w-[316px] md:h-[186px] hover:scale-105 duration-300 border-border">
      <CardHeader>
        <CardTitle className="font-light text-sm">{data.name}</CardTitle>
      </CardHeader>
      <CardContent className="max-h-[4vh] md:max-h-[7vh]">
        <ChartContainer config={chartConfig}>
          <LineChart data={chartData}>
            <XAxis dataKey="time" tickLine={false} axisLine={false} />
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent hideLabel />}
            />
            <Line
              dataKey="totalQps"
              type="natural"
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
      <CardFooter className="bottom-0">
        <i className="fa-solid fa-arrow-right daily--fa-arrow-right"></i>
      </CardFooter>
    </Card>
  );
}

DailySummary.propTypes = {
  data: PropTypes.shape({
    name: PropTypes.string.isRequired,
    data: PropTypes.array,
  }),
  setBenchmarktype: PropTypes.func.isRequired,
  isSelected: PropTypes.bool,
};
