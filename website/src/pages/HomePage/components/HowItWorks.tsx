import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { ReactNode } from "react";
import { Link } from "react-router-dom";

type HowItWorksCardItem = {
  title: string;
  content: ReactNode;
};

const items: HowItWorksCardItem[] = [
  {
    title: "Workloads",
    content: (
      <>
        Seven workloads are executed against every commit we decide to benchmark. These workloads define the SQL schema, the queries we execute, the configuration of the Vitess cluster, and the configuration of sysbench. The settings of our workload can be found on the
        <Link to="" className="text-primary"> arewefastyet repository</Link>
        <br />
        <br />
        We have two categories of workloads: OLTP and TPCC, both are a modified version of the official workloads.
      </>
    )
  },
  {
    title: "Frequency",
    content: (
      <>
        There are three cron schedules that enable us to periodically benchmark Vitess, the definition of these schedules is
        <Link to="" className="text-primary"> available in our yaml configuration</Link>
        <br />
        <br />
        Generally, we benchmark the main and release branches of Vitess every night at midnight UTC. We also detect new PRs that need to be benchmarked every five minutes, and new tags/releases every minute.      </>
    )
  },
  {
    title: "Methodology",
    content: (
      <>
        Each commit is benchmarked ten times by each workload. Under the hood, we use sysbench with a custom configuration to perform the benchmark. We then perform statistical analysis on the ten results of a given workload to get a more accurate and reliable result.
        <br />
        <br />
        Every execution is done on the same hardware: 192Gb of RAM and 2x Intel Xeon Silver 4214 Processor 24-Core @ 2.20GHz.      </>
    )
  },
  {
    title: "Results",
    content: (
      <>
We collect the same results as sysbench (QPS, TPS, Error rate, latency, etc), along with several Golang metrics such as the CPU used per query, and the total memory used per query.      </>
    )

  }
];

export default function HowItWorks() {
  return (
    <section className="relative flex flex-col items-center p-page -z-1 pb-14 bg-background text-foreground">
      <h1 className="text-4xl font-semibold my-14 text-primary dark:text-front">How it works</h1>
      <div className="flex flex-col md:flex-row md:flex-wrap justify-center items-center md:justify-between md:items-start gap-y-12 px-5 md:px-10 md:relative z-1">
        {items.map((item, key) => (
          <HowItWorksCard key={key} title={item.title} content={item.content} />
        ))}
      </div>
    </section>
  );
}

function HowItWorksCard({ title, content }: { title: string; content: ReactNode }) {
  return (
    <Card className="w-full h-fit 2xl:h-[334px] md:w-[calc(50%_-_1.25rem)] border-border rounded-xl overflow-hidden">
      <CardHeader>
        <CardTitle className="text-3xl dark:text-primary font-semibold">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-lg mt-2 text-foreground">{content}</p>
      </CardContent>
    </Card>
  );
}
