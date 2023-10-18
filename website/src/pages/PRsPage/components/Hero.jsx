import React from "react";

export default function Hero() {
  return (
    <section className="flex pt-[15vh] items-center p-page">
      <div className="flex basis-1/2 flex-col">
        <h2 className="text-8xl text-primary">Pull Request</h2>
        <p className="my-6 leading-loose">
          If a given Pull Request on vitessio/vitess is labelled with the
          Benchmark me label the Pull Request will be handled and benchmark by
          arewefastyet. For each commit on the Pull Request there will be two
          benchmarks: one on the Pull Request's HEAD and another on the base of
          the Pull Request.
          <br />
          <br />
          On this page you can find all benchmarked Pull Requests.
        </p>
      </div>
    </section>
  );
}
