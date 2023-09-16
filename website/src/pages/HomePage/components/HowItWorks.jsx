import React, { useEffect, useRef } from "react";

const items = [
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
      <div className="flex flex-wrap justify-between gap-y-12 px-10 relative z-1">
        {items.map((item, key) => (
          <HowItWorksCard key={key} title={item.title} content={item.content} />
        ))}
      </div>
    </section>
  );
}

function HowItWorksCard(props) {
  const cardRef = useRef();
  const glowRef = useRef();

  useEffect(() => {
    window.addEventListener("mousemove", (event) => {
      const cardRect = cardRef.current.getBoundingClientRect();
      glowRef.current.style.setProperty("--x", `${event.x - cardRect.x}px`);
      glowRef.current.style.setProperty("--y", `${event.y - cardRect.y}px`);
    });
  }, []);

  return (
    <div
      ref={cardRef}
      className="relative flex flex-col items-center p-4 w-[calc(50%_-_1.25rem)] bg-black 
      rounded-lg border border-white border-opacity-20 gap-y-4 overflow-hidden"
    >
      <div
        ref={glowRef}
        className="absolute bg-white bg-opacity-10 z-1 top-[var(--y)] left-[var(--x)] blur-3xl
      w-[20vw] h-[20vw] rounded-full -translate-x-1/2 -translate-y-1/2"
      />
      <h3 className="text-2xl font-semibold">{props.title}</h3>
      <p className="text-sm mt-2 text-white text-opacity-70">{props.content}</p>
    </div>
  );
}
