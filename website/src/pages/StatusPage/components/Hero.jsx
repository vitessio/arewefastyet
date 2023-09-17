import React from "react";
import useApiCall from "../../../utils/Hook";

const info = [
  { title: "Total Benchmarks", content: "Total" },
  { title: "Successful Benchmarks", content: "Finished" },
  { title: "Total Benchmarks (last 30 days)", content: "Last30Days" },
];

export default function Hero() {
  const { data: dataStatusStats } = useApiCall(
    `${import.meta.env.VITE_API_URL}status/stats`
  );

  return (
    <section className="flex h-[90vh] items-center p-page">
      <div className="flex basis-1/2 flex-col">
        <h2 className="text-8xl text-primary">Status</h2>
        <p className="my-6 leading-loose">
          Arewefastyet has a single execution queue, each element in the queue
          is executed one after the other. In the status page, we can see the
          content of the execution queue, along with the 50 last executions.
          Each execution has a status that can be either: <b>started</b>,{" "}
          <b>failed</b> or <b>finished</b>. When a benchmark is marked as{" "}
          <b>finished</b> it means that it successfully finished.
        </p>
      </div>

      <div className="flex-1 flex flex-col items-center">
        <div className="flex flex-col gap-y-8">
          {info.map((item, key) => (
            <div key={key} className="">
              <h4
                className="counter text-7xl text-primary"
                style={{ "--num": dataStatusStats[item.content] }}
              ></h4>
              <p className="text-lg font-medium">{item.title}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
