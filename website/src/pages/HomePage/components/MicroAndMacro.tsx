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

import React from "react";

interface Point {
  title: string;
  content: string;
}

interface Item {
  title: string;
  points: Point[];
}

const items: Item[] = [
  {
    title: "Gaining Functional Insights",
    points: [
      {
        title: "Focused Evaluation",
        content:
          "Micro benchmarks dissect specific functional units within Vitess, allowing precise assessment of individual components.",
      },
      {
        title: "Golang Advantage",
        content:
          "Vitess leverages the Go standard library's testing framework and micro-benchmarking tools, ensuring accurate measurements and consistent results. Relying on native go features produces better accuracy.",
      },
      {
        title: "Execution Efficiency",
        content:
          "Micro-benchmarks are effortlessly executed using the default go test runner and arewefastyet's microbench command, facilitating streamlined testing processes.",
      },
      {
        title: "Critical Metrics",
        content:
          "Key performance indicators, such as iteration time measured in nanoseconds and memory allocation in bytes, are derived from micro-benchmark results. Multiple metrics allow better understanding for various audiences.",
      },
      {
        title: "Structured Analysis",
        content:
          "Extracted metrics are meticulously analyzed and stored in a MySQL database, forming a valuable repository for future reference and comparison. Basically analogous to unit tests",
      },
      {
        title: "Granular Performance Assessment",
        content:
          "Micro benchmarks provide an unparalleled level of granularity, enabling a meticulous examination of individual code units, ensuring that even the smallest performance nuances are captured.",
      },
    ],
  },
  {
    title: "Real-World Performance Insights",
    points: [
      {
        title: "Comprehensive Overview",
        content:
          "Macro benchmarks provide a comprehensive view of Vitess' performance, simulating real-world production conditions for accurate evaluations.",
      },
      {
        title: "Cluster Configuration",
        content:
          "Benchmark Vitess clusters are thoughtfully assembled, encompassing vtgates, vttablets, etcd clusters, and vtctld servers, creating an environment closely resembling real deployments.",
      },
      {
        title: "Multi-Step Process",
        content:
          "Macro benchmarks comprise three sequential stages: preparation, warm-up, and the actual run, systematically capturing the performance trajectory.",
      },
      {
        title: "Custom Benchmarking",
        content:
          "Sysbench, tailored to benchmark various data stores, is utilized for the main benchmarking process, capturing critical metrics such as latency, transactions per second (TPS), and queries per second (QPS).",
      },
      {
        title: "Incorporating Insights",
        content:
          "Internal cluster metrics and operating system metrics are integrated and processed through a Prometheus backend, enriching the assessment with deeper performance context.",
      },
      {
        title: "Informed Optimization",
        content:
          "The combined insights from macro and micro benchmarks empower users to optimize Vitess effectively for their specific production environments, making well-informed decisions based on both granular functional details and broader performance trends.",
      },
    ],
  },
];

export default function MicroAndMacro() {
  return (
    <section className="flex flex-col items-center p-page py-14">
      <h1 className="bg-primary bg-opacity-20 text-primary text-lg font-semibold px-4 py-2 rounded-full">
        Micro and Macro Benchmarks
      </h1>

      <div className="flex justify-between mt-12">
        {items.map((item, key) => (
          <div key={key} className="w-[48%] border border-primary rounded-lg">
            <h3 className="text-center my-4 text-primary font-semibold text-lg">
              {item.title}
            </h3>
            <ul className="list-none p-4">
              {item.points.map((point, i) => (
                <li key={i} className="m-4">
                  <h5 className="font-medium my-1">{point.title}</h5>
                  <p className="text-xs text-front text-opacity-80">
                    {point.content}
                  </p>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </div>
    </section>
  );
}
