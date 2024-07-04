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
import useApiCall from "@/utils/Hook";
import { statusDataTypes } from "@/types";
import Hero, { HeroProps } from "@/common/Hero";

const info = [
  { title: "Total Benchmarks", content: "Total" },
  { title: "Successful", content: "Finished" },
  { title: `Total (last 30 days)`, content: "Last30Days" },
];

const heroProps: HeroProps = {
  title: "Status",
  description: (
    <p>
      Arewefastyet has a single execution queue, each element in the queue
      is executed one after the other. In the status page, we can see the
      content of the execution queue, along with the 50 last executions.
      Each execution has a status that can be either: <b>started</b>,{" "}
      <b>failed</b> or <b>finished</b>. When a benchmark is marked as{" "}
      <b>finished</b> it means that it successfully finished.
    </p>
  ),
};

export default function StatusHero() {
  const { data: dataStatusStats } = useApiCall<statusDataTypes>(
    `${import.meta.env.VITE_API_URL}status/stats`
  );

  return (
    <Hero title={heroProps.title} description={heroProps.description}>
      <div className="flex flex-col md:flex-row justify-between items-center w-full md:w-fit md:gap-x-24">
        {info.map((item, key) => (
          <div
            key={key}
            className="w-full h-[6em] md:min-h-[5em] flex flex-col justify-start md:items-start items-center gap-5"
          >
            <h4
              className="counter text-3xl md:text-7xl text-primary"
              style={{ "--num": dataStatusStats[item.content] }}
            ></h4>
            <p className="text-xs md:text-lg font-light md:font-medium">
              {item.title}
            </p>
          </div>
        ))}
      </div>
    </Hero>
  );
}
