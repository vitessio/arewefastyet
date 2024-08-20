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

import { Separator } from "@/components/ui/separator";
import useApiCall from "@/hooks/useApiCall";
import { FilterConfigs } from "@/types";
import { columns, MacroQueriesPlan } from "./components/Columns";
import MacroQueriesCompareHero from "./components/MacroQueriesCompareHero";
import { MacroQueriesCompareTable } from "./components/MacroQueriesCompareTable";

export default function MacroQueriesComparePage() {
  const urlParams = new URLSearchParams(window.location.search);
  const commits = {
    oldGitRef: urlParams.get("old") || "",
    newGitRef: urlParams.get("new") || "",
  };
  const workload = urlParams.get("workload") || "";

  const {
    data: data,
    isLoading: isMacrobenchLoading,
    error: macrobenchError,
  } = useApiCall<MacroQueriesPlan[]>({
    url: `${import.meta.env.VITE_API_URL}macrobench/compare/queries?ltag=${
      commits.oldGitRef
    }&rtag=${commits.newGitRef}&workload=${workload}`,
    queryKey: ["compare", commits.oldGitRef, commits.newGitRef, workload],
  });

  let filterConfigs: FilterConfigs[] = [
    {
      column: "query",
      title: "Operators",
      options: ["select", "insert", "update", "delete"].map((value) => {
        return { label: value, value: value };
      }),
    },
  ];

  return (
    <>
      <MacroQueriesCompareHero commits={commits} />
      <Separator className="w-4/5 m-auto" />
      {data && !isMacrobenchLoading && (
        <div className="w-[80vw] xl:w-[60vw] m-auto">
          <MacroQueriesCompareTable
            columns={columns}
            data={data}
            filterConfigs={filterConfigs}
          />
        </div>
      )}
    </>
  );
}
