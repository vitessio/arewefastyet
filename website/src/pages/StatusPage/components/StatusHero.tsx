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

import Hero, { HeroProps } from "@/common/Hero";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import useApiCall from "@/utils/Hook";
import { Bar, BarChart, CartesianGrid, XAxis, YAxis } from "recharts";

const info = [
  { title: "Benchmark in total", content: "Total" },
  { title: "Benchmark this month", content: "Last30Days" },
  { title: `Commits Benchmarked`, content: "Commits" },
];

type DataStatusType = {
  Total: number;
  Last30Days: number;
  Last7Days: number[];
  Commits: number;
};

const heroProps: HeroProps = {
  title: "Status",
  description: (
    <>
      Arewefastyet has a single execution queue with tasks that are executed
      sequentially on our benchmarking server. Each execution has a status that
      can be either: <b>started</b>started, <b>failed</b> or <b>finished</b>.
      When a benchmark is marked as finished it means that it ran successfully.
    </>
  ),
};

export default function StatusHero() {
  const { data: dataStatusStats } = useApiCall<DataStatusType>(
    `${import.meta.env.VITE_API_URL}status/stats`
  );

  const chartConfig = {
    day: {
      label: "Day",
      color: "hsl(var(--chart-1))",
    },
  };

  const last7days = [234, 283, 210, 371, 234, 283, 210];

  const chartData = last7days.map((commits, index) => ({
    day: index.toString(),
    commits,
  }));

  return (
    <Hero title={heroProps.title} description={heroProps.description}>
      <div className="flex flex-col md:flex-row justify-between items-center w-full md:w-fit md:gap-x-24">
        {info.map(({ content, title }, key) => (
          <Card
            key={key}
            className="w-[300px] h-[150px] md:min-h-[130px] border-border"
          >
            <CardHeader>
              <CardTitle
                className="counter text-3xl md:text-5xl text-primary"
                // TODO: Fix the type error caused by useApiCall that returns an array everytime even if the data is not an array
                style={{ ["--num" as string]: dataStatusStats[content] }}
              ></CardTitle>
            </CardHeader>
            <CardFooter className="text-xs md:text-xl font-light md:font-medium">
              {title}
            </CardFooter>
          </Card>
        ))}
        <Card className=" w-[300px] max-h-[150px] border-border">
          <CardContent className="w-[300px] max-h-[150px]">
            <ChartContainer config={chartConfig}>
              <BarChart
                margin={{ top: 12, right: 0, left: -32, bottom: 0 }}
                barSize={10}
                accessibilityLayer
                data={chartData}
              >
                <YAxis tickMargin={0} />
                <XAxis tickMargin={0} tick={false} />
                <CartesianGrid vertical={false} />
                <ChartTooltip
                  cursor={true}
                  content={<ChartTooltipContent hideLabel />}
                />
                <Bar dataKey="commits" fill="var(--color-day)" />
              </BarChart>
            </ChartContainer>
          </CardContent>
          <CardFooter className="text-xs text-primary flex justify-center md:text-sm font-light md:font-medium mt-[-28px]">
            Benchmarks over the last 7 days
          </CardFooter>
        </Card>
      </div>
    </Hero>
  );
}
