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
import { Skeleton } from "@/components/ui/skeleton";
import useApiCall from "@/hooks/useApiCall";
import { errorApi, fixed } from "@/utils/Utils";
import { Bar, BarChart, CartesianGrid, XAxis, YAxis } from "recharts";

const info = [
  { title: "Benchmark in total", content: "Total" },
  { title: "Benchmark this month", content: "Last30Days" },
  { title: "Commits benchmarked", content: "Commits" },
  { title: "Execution duration (avg)", content: "AvgDuration" },
];

type DataStatusType = {
  Total: number;
  Last30Days: number;
  Last7Days: number[];
  Commits: number;
  AvgDuration: number;
};

const heroProps: HeroProps = {
  title: "Status",
  description: (
    <>
      Arewefastyet has a single execution queue with tasks that are executed
      sequentially on our benchmarking server. Each execution has a status that
      can be either: <b>started</b>started, <b>failed</b> or <b>finished</b>.
      When a benchmark is marked as finished it means that it ran successfully.
      Only the last 1000 executions are shown in the previous executions.
    </>
  ),
};

export default function StatusHero() {
  const {
    data: dataStatusMetrics,
    error: dataStatusError,
    isLoading: isDataStatusLoading,
  } = useApiCall<DataStatusType>({
    url: `${import.meta.env.VITE_API_URL}status/stats`,
    queryKey: ["statusStats"],
  });

  const chartConfig = {
    day: {
      label: "Day",
      color: "hsl(var(--primary))",
    },
  };

  let last7days: number[] = [];

  if (dataStatusMetrics !== undefined) {
    last7days = dataStatusMetrics.Last7Days;
  }

  const chartData = last7days.map((executions, index) => ({
    day: index.toString(),
    executions,
  }));

  return (
    <Hero title={heroProps.title} description={heroProps.description}>

      <div className="flex flex-wrap justify-center gap-8 w-[80vw]">
      {isDataStatusLoading && (
        <>
            {Array.from([, , , , ,]).map((_, index) => {
              return (
                <Skeleton key={index} className="w-[300px] h-[150px] md:min-h-[130px]"></Skeleton>

              );
            })}
        </>
      )}

      {!isDataStatusLoading &&
        (dataStatusError || !dataStatusMetrics) && (
          <div className="text-destructive text-center my-10">{errorApi}</div>
        )}
        {!isDataStatusLoading && dataStatusMetrics && (
          <>
            {info.map(({ content, title }, key) => (
              <Card
                key={key}
                className="w-[300px] h-[150px] md:min-h-[130px] border-border"
              >
                <CardHeader>
                  {content == "AvgDuration" ? (
                    <CardTitle className="text-4xl md:text-5xl text-primary">
                      {fixed(
                        dataStatusMetrics?.[
                          content as keyof DataStatusType
                        ] as number,
                        2
                      )}{" "}
                      min
                    </CardTitle>
                  ) : (
                    <CardTitle
                      className="counter text-4xl md:text-5xl text-primary"
                      style={{
                        ["--num" as string]:
                          dataStatusMetrics?.[content as keyof DataStatusType],
                      }}
                    ></CardTitle>
                  )}
                </CardHeader>
                <CardFooter className="text-lg md:text-xl font-light md:font-medium">
                  {title}
                </CardFooter>
              </Card>
            ))}
            <Card className="w-[300px] h-[150px] border-border">
              <CardContent className="w-full h-full">
                <ChartContainer config={chartConfig}>
                  <BarChart
                    margin={{ top: 12, right: 0, left: -22, bottom: 0 }}
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
                    <Bar dataKey="executions" fill="var(--color-day)" />
                  </BarChart>
                </ChartContainer>
              </CardContent>
              <CardFooter className="text-xs text-primary flex justify-center md:text-sm font-light md:font-medium mt-[-28px]">
                Benchmarks over the last 7 days
              </CardFooter>
            </Card>
          </>
        )}
      </div>
    </Hero>
  );
}
