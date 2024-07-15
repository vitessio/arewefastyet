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
      color: "hsl(var(--primary))",
    },
  };

  const last7days = [234, 283, 210, 371, 234, 283, 210];

  const chartData = last7days.map((executions, index) => ({
    day: index.toString(),
    executions,
  }));

  return (
    <Hero title={heroProps.title} description={heroProps.description}>
      <div className="flex flex-wrap justify-center gap-8 w-[80vw]">
        {info.map(({ content, title }, key) => (
          <Card
            key={key}
            className="w-[300px] h-[150px] md:min-h-[130px] border-border"
          >
            <CardHeader>
              <CardTitle
                className="counter text-4xl md:text-5xl text-primary"
                style={{ ["--num" as string]: dataStatusStats[content] }}
              ></CardTitle>
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
      </div>
    </Hero>
  );
}
