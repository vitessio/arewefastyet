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

import { Skeleton } from "@/components/ui/skeleton";
import { FilterConfigs } from "@/types";
import useApiCall from "@/utils/Hook";
import { useSearchParams } from "react-router-dom";
import { columns, HistoryType } from "./components/Columns";
import HistoryHero from "./components/HistoryHero";
import { HistoryTable } from "./components/HistoryTable";

export default function HistoryPage() {
  let {
    data: dataHistory,
    isLoading,
    error,
  } = useApiCall<HistoryType[]>(`${import.meta.env.VITE_API_URL}history`);
  const [searchParams] = useSearchParams();
  const gitRef = searchParams.get("gitRef") ?? "";

  const sources = dataHistory?.map((source) => source.source);
  const uniqueSources = Array.from(new Set(sources));
  const filterConfigs: FilterConfigs[] = [
    {
      column: "source",
      title: "Source",
      options:
        uniqueSources.map((source) => {
          return { label: source, value: source };
        }) || [],
    },
  ];

  return (
    <>
      <HistoryHero />
      <section className="mx-auto p-page lg:w-[60vw] my-12 flex flex-col">
        {isLoading && <Skeleton className="h-[732px]"></Skeleton>}
        {error && (
          <div className="text-destructive text-center my-2">{error}</div>
        )}
        {dataHistory && (
          <HistoryTable
            columns={columns}
            data={dataHistory}
            filterConfigs={filterConfigs}
            initialGitRef={gitRef}
          />
        )}
      </section>
    </>
  );
}
