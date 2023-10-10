import React from "react";
import useApiCall from "../../../utils/Hook";

export default function Hero() {
  return (
    <section className="flex h-[74vh] items-center p-page">
      <div className="flex basis-1/2 flex-col">
        <h2 className="text-8xl text-primary">Daily</h2>
        <p className="my-6 leading-loose">
          We run all macro benchmark workloads against the <i>main</i> branch
          every day. This is done to ensure the consistency of the results over
          time on <i>main</i>. On this page, you can find graphs that show you the
          results of all five macro benchmark workload over the last 30 days.
          Click on a macro benchmark workload to see all the results for that
          workload.
        </p>
      </div>
    </section>
  );
}
