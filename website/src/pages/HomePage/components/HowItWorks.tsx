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

import React, { useEffect, useRef } from "react";

interface WorkItem {
  title: string;
  content: string;
}

const items: WorkItem[] = [
  {
    title: "The Execution Engine",
    content:
      "At the heart of arewefastyet is the Execution engine. It orchestrates the entire benchmarking process, ensuring accuracy and reproducibility on a large scale. Each benchmark run is initiated by new releases, new PRs, and new commits on main.",
  },
  {
    title: "Dedicated Benchmarking Servers",
    content:
      "arewefastyet relies on dedicated hardware provided by CNCF and Equinix Metal. Our benchmarking infrastructure uses large bare-metal servers, boosting benchmark reliability and accuracy.",
  },
  {
    title: "Customized Benchmark Settings",
    content:
      "Different benchmarks demand distinct configurations. For instance, a macro-benchmark necessitates the setup of a Vitess cluster, while a micro-benchmark does not. The default setup for macro-benchmarks examines Vitess performance in a sharded keyspace with six VTGates and two VTTablets.",
  },
  {
    title: "Starting Benchmark Runs",
    content:
      "Once the server is ready, the final step is initiating the benchmark run. Ansible triggers arewefastyet's CLI to set the benchmark in motion. This comprehensive process, from YAML-based pipeline configuration to dynamic server setup, ensures that every benchmark run is accurate, reproducible, and adaptable to the unique demands of each benchmark type. arewefastyet streamlines the complexities of executing benchmarks against Vitess, offering a robust and precise benchmarking solution at scale.",
  },
];

export default function HowItWorks() {
  return (
    <section className="relative flex flex-col items-center p-page -z-1 pb-14 bg-black text-white">
      <div className="absolute-cover bg-white bg-opacity-5 z-0" />
      <h1 className="text-4xl font-semibold my-14">How it works</h1>
      <div className="flex flex-col md:flex-row md:flex-wrap justify-center items-center gap-5 md:justify-between md:items-start md:gap-y-12 px-5 md:px-10 md:relative z-1">
        {items.map((item, key) => (
          <HowItWorksCard key={key} title={item.title} content={item.content} />
        ))}
      </div>
    </section>
  );
}

function HowItWorksCard({ title, content }: WorkItem) {
  const cardRef = useRef<HTMLDivElement>(null);

  const glowRef = useRef<HTMLDivElement>(null);

  function glowWithMouse(event: { x: number; y: number }) {
    if (!cardRef.current || !glowRef.current) return;

    const cardRect = cardRef.current.getBoundingClientRect();
    glowRef.current.style.setProperty("--x", `${event.x - cardRect.x}px`);
    glowRef.current.style.setProperty("--y", `${event.y - cardRect.y}px`);
  }

  useEffect(() => {
    const isSmallScreen = window.innerWidth < 768;

    if (!isSmallScreen) {
      window.addEventListener("mousemove", glowWithMouse);
    }

    return () => {
      if (!isSmallScreen) {
        window.removeEventListener("mousemove", glowWithMouse);
      }
    };
  }, []);

  return (
    <div
      ref={cardRef}
      className="w-full min-h-[200px] md:relative flex flex-col items-center p-4 md:w-[calc(50%_-_1.25rem)] bg-black 
      rounded-lg border border-white border-opacity-20 gap-y-4 overflow-hidden"
    >
      <div
        ref={glowRef}
        className="absolute bg-white bg-opacity-10 z-1 md:top-[var(--y)] md:left-[var(--x)] blur-3xl
      md:w-[20vw] w-full h-[20vw] rounded-full -translate-x-1/2 -translate-y-1/2"
      />
      <h3 className="text-2xl font-semibold">{title}</h3>
      <p className="text-sm mt-2 text-white text-opacity-70">{content}</p>
    </div>
  );
}