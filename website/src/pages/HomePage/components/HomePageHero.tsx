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

import { Link, To } from "react-router-dom";
import Icon from "@/common/Icon";
import { Button } from "@/components/ui/button";
import DailySummary from "@/common/DailySummary";
import useDailySummaryData from "@/hooks/useDailySummaryData";
import { useState } from "react";

export default function HomePageHero() {
  const [benchmarkType, setBenchmarktype] = useState<string>("");

  function getBenchmarkType(type: string) {
    setBenchmarktype(type);
  }

  const {
    dataDailySummary,
  } = useDailySummaryData();

  return (
    <section className="flex flex-col items-center h-screen p-page">
      <h1 className="text-6xl font-semibold text-center mt-10 leading-normal">
        Benchmarking <br />
        System for <br />
        <span className="text-orange-500"> Vitess</span>
      </h1>
      <div className="flex gap-x-4 mt-10">
        <Button
          asChild
          size={"lg"}
          variant={"default"}
          className="bg-background hover:bg-muted/90 dark:bg-front dark:hover:bg-front/90 dark:text-background text-foreground rounded-lg border"
        >
          <Link
            className="rounded-2xl p-5 flex items-center gap-x-2"
            to="https://vitess.io/blog/2021-07-08-announcing-vitess-arewefastyet"
            target="__blank"
          >
            Blog Post
            <Icon className="text-2xl" icon="bookmark" />
          </Link>
        </Button>
        <Button
          asChild
          size={"lg"}
          variant={"default"}
          className="bg-front text-background hover:bg-front/90 rounded-lg"
        >
          <Link
            className="rounded-2xl p-5 flex items-center gap-x-2"
            to="https://github.com/vitessio/arewefastyet"
            target="__blank"
          >
            GitHub
            <Icon className="text-2xl" icon="github" />
          </Link>
        </Button>
        <Button
          asChild
          size={"lg"}
          variant={"default"}
          className="bg-background hover:bg-muted/90 dark:bg-front dark:hover:bg-front/90 dark:text-background text-foreground rounded-lg border"
        >
          <Link
            className="rounded-2xl p-5 flex items-center gap-x-2"
            to="https://www.vitess.io"
            target="__blank"
          >
            Vitess
            <Icon className="text-2xl" icon="vitess" />
          </Link>
        </Button>
      </div>
      <h2 className="text-2xl font-medium mt-20">
        Historical results on the <Link className="text-primary" to="https://github.com/vitessio/arewefastyet/tree/main" target="_blank">main</Link>{" "}
        branch
      </h2>
      <section className="flex p-page justif-center flex-wrap gap-10 py-10">
        {dataDailySummary.map((dailySummary, index) => {
          if (dailySummary.name === "OLTP" || dailySummary.name === "TPCC") {
            return (
              <Link to={`/daily?type=${dailySummary.name}`}>
                <DailySummary
                  key={index}
                  data={dailySummary}
                  benchmarkType={benchmarkType}
                  setBenchmarktype={setBenchmarktype}
                />
              </Link>

            );
          }
        })}
      </section>
      <Link
        to="/daily"
        className="text-primary text-lg mt-10"
      >
        See more historical results {">"}
      </Link>
    </section>
  );
}
