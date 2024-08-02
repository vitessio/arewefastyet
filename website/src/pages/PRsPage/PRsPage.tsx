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

import useApiCall from "@/hooks/useApiCall";

import { Skeleton } from "@/components/ui/skeleton";
import { PrData } from "@/types";
import { columns } from "./components/Columns";
import PRHero from "./components/PRHero";
import PRTable from "./components/PRTable";

export default function PRsPage() {
  const {
    data: dataPRList,
    isLoading: isPRListLoading,
    error: PRListError,
  } = useApiCall<PrData[]>({
    url: `${import.meta.env.VITE_API_URL}pr/list`,
    queryKey: "prList",
  });

  return (
    <>
      <PRHero />

      {isPRListLoading && (
        <div className="lg:mx-auto w-full p-page xl:w-[80vw] my-12 flex flex-col">
          <Skeleton className="h-[712px]"></Skeleton>
        </div>
      )}

      {PRListError ? (
        <div className="text-red-500 text-center my-2">
          {<>{PRListError}</>}
        </div>
      ) : null}

      {!isPRListLoading && dataPRList && (
        <div className="lg:mx-auto w-full p-page xl:w-[80vw] my-12 flex flex-col">
          <PRTable data={dataPRList} columns={columns} />
        </div>
      )}
    </>
  );
}
