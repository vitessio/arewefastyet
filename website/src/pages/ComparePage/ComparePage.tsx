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

import MacroBenchmarkTable from "@/common/MacroBenchmarkTable";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { CompareData, MacroBenchmarkTableData } from "@/types";
import useApiCall from "@/utils/Hook";
import { formatCompareData } from "@/utils/Utils";
import { PlusCircledIcon } from "@radix-ui/react-icons";
import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import CompareHero from "./components/CompareHero";

export default function Compare() {
  const urlParams = new URLSearchParams(window.location.search);

  const [gitRef, setGitRef] = useState({
    old: urlParams.get("old") || "",
    new: urlParams.get("new") || "",
  });

  const navigate = useNavigate();

  useEffect(() => {
    navigate(`?old=${gitRef.old}&new=${gitRef.new}`);
  }, [gitRef.old, gitRef.new]);

  const {
    data: data,
    isLoading: isMacrobenchLoading,
    error: macrobenchError,
  } = useApiCall<CompareData[]>(
    `${import.meta.env.VITE_API_URL}macrobench/compare?new=${gitRef.new}&old=${
      gitRef.old
    }`
  );

  let formattedData: MacroBenchmarkTableData[] = [];

  if (data !== null && data.length > 0) {
    formattedData = formatCompareData(data);
  }

  return (
    <>
      <CompareHero gitRef={gitRef} setGitRef={setGitRef} />
      {macrobenchError && (
        <div className="text-red-500 text-center my-2">{macrobenchError}</div>
      )}

      <section className="flex flex-col items-center">
        {isMacrobenchLoading && (
          <>
            {[...Array(8)].map((_, index) => {
              return (
                <div key={index} className="w-[80vw] xl:w-[60vw] my-12">
                  <Skeleton className="h-[852px]"></Skeleton>
                </div>
              );
            })}
          </>
        )}
        {!isMacrobenchLoading && data !== null && data.length > 0 && (
          <>
            {data.map((macro, index) => {
              return (
                <div className="w-[80vw] xl:w-[60vw] my-12" key={index}>
                  <Card className="border-border">
                    <CardHeader className="flex flex-col gap-4 md:gap-0 md:flex-row justify-between pt-6">
                      <CardTitle className="text-2xl md:text-4xl">
                        {macro.workload}
                      </CardTitle>
                      <Button
                        variant="outline"
                        size="sm"
                        className="h-8 w-fit border-dashed mt-4 md:mt-0"
                      >
                        <PlusCircledIcon className="mr-2 h-4 w-4 text-primary" />
                        <Link
                          to={`/macrobench/queries/compare?old=${gitRef.old}&new=${gitRef.new}&workload=${macro.workload}`}
                        >
                          See Query Plan{" "}
                        </Link>
                      </Button>
                    </CardHeader>
                    <CardContent className="w-full p-0">
                      <MacroBenchmarkTable
                        data={formattedData[index]}
                        new={gitRef.new}
                        old={gitRef.old}
                      />
                    </CardContent>
                  </Card>
                </div>
              );
            })}
          </>
        )}
      </section>
    </>
  );
}
