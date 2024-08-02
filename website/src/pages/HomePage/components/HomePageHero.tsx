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

import DailySummary from "@/common/DailySummary";
import Icon from "@/common/Icon";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import useDailySummaryData from "@/hooks/useDailySummaryData";
import { useState } from "react";
import { Link } from "react-router-dom";

export default function HomePageHero() {
  const [workload, setWorkload] = useState<string>("");

  const { dataDailySummary, isLoadingDailySummary, dailySummaryError } =
    useDailySummaryData(["OLTP", "TPCC"]);

  return (
    <section className="flex flex-col items-center h-fit my-12">
      <h1 className="flex flex-col gap-8 text-3xl md:text-6xl font-semibold text-center">
        <p>Benchmarking</p>
        <p> System for</p>
        <p className="text-primary"> Vitess</p>
      </h1>
      <div className="flex md:flex-row flex-col gap-4 mt-10">
        <div className="flex flex-row gap-4">
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
        </div>
        <div className="flex justify-center">
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
      </div>
      <h2 className="text-lg md:text-2xl font-medium mt-12">
        Historical results on the{" "}
        <Link
          className="text-primary"
          to="https://github.com/vitessio/vitess/tree/main/"
          target="_blank"
        >
          main
        </Link>{" "}
        branch
      </h2>
      <section className="flex md:flex-row flex-col justify-center gap-10 my-20">
        {isLoadingDailySummary && (
          <>
            <Skeleton className="w-[310px] h-[124px] md:w-[316px] md:h-[124px] rounded-lg" />
            <Skeleton className="w-[310px] h-[124px] md:w-[316px] md:h-[124px] rounded-lg" />
          </>
        )}
        {!isLoadingDailySummary &&
          (dailySummaryError ||
            !dataDailySummary ||
            dataDailySummary.length === 0) && (
            <div className="text-destructive text-center my-10">
              {<>{dailySummaryError || "No data available"}</>}
            </div>
          )}
        {!isLoadingDailySummary &&
          dataDailySummary &&
          dataDailySummary.length > 0 && (
            <>
              {dataDailySummary.map((dailySummary, index) => {
                if (
                  dailySummary.name === "OLTP" ||
                  dailySummary.name === "TPCC"
                ) {
                  return (
                    <Link
                      key={index}
                      to={`/daily?workload=${dailySummary.name}`}
                    >
                      <DailySummary
                        data={dailySummary}
                        workload={workload}
                        setWorkload={setWorkload}
                      />
                    </Link>
                  );
                }
                return null;
              })}
            </>
          )}
      </section>
      <Link to="/daily" className="text-primary text-xs md:text-lg">
        See more historical results {">"}
      </Link>
    </section>
  );
}
