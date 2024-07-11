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

import DailySummary from "@/common/DailySummary";
import { Skeleton } from "@/components/ui/skeleton";
import useDailySummaryData from "@/hooks/useDailySummaryData";
import { Workloads } from "@/types";

export type DailDailySummaryProps = {
  benchmarkType: string;
  setBenchmarktype: (type: string) => void;
};

export default function DailyDailySummary(props: DailDailySummaryProps) {
  const { benchmarkType, setBenchmarktype } = props;
  const workloads: Workloads[] = [
    "OLTP",
    "OLTP-READONLY",
    "OLTP-SET",
    "TPCC",
    "TPCC_FK",
    "TPCC_UNSHARDED",
    "TPCC_FK_UNMANAGED",
  ];

  const { dataDailySummary, isLoadingDailySummary, dailySummaryError } =
    useDailySummaryData(workloads);
  return (
    <>
      <section className="flex p-page flex-wrap justify-center gap-12 p-4">
        {isLoadingDailySummary && (
          <>
            {workloads.map((_, index) => {
              return (
                <Skeleton
                  key={index}
                  className="w-[310px] h-[124px] md:w-[316px] md:h-[124px] rounded-lg"
                />
              );
            })}
          </>
        )}

        {dailySummaryError && (
          <div className="text-red-500 text-center my-10">
            {dailySummaryError}
          </div>
        )}

        {!dailySummaryError &&
          dataDailySummary &&
          dataDailySummary.length > 0 && (
            <>
              {dataDailySummary.map((dailySummary, index) => {
                return (
                  <DailySummary
                    key={index}
                    data={dailySummary}
                    benchmarkType={benchmarkType}
                    setBenchmarktype={setBenchmarktype}
                  />
                );
              })}
            </>
          )}
      </section>
      <p className="text-center py-8">
        Y-Axis: <span className="text-primary font-medium">Total QPS</span>
        &emsp;&emsp;X-Axis:{" "}
        <span className="text-primary font-medium">Time</span>
      </p>
    </>
  );
}
