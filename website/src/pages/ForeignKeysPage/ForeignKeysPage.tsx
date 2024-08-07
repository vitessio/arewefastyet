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

import { CompareResult, VitessRefs } from "@/types";
import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";

import MacroBenchmarkTable from "@/common/MacroBenchmarkTable";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import useApiCall from "@/hooks/useApiCall";
import { formatCompareResult, getRefName } from "@/utils/Utils";
import { PlusCircledIcon } from "@radix-ui/react-icons";
import ForeignKeysHero from "./components/ForeignKeysHero";

export default function ForeignKeys() {
  const urlParams = new URLSearchParams(window.location.search);

  const [gitRef, setGitRef] = useState<string>(urlParams.get("sha") || "");
  const [workload, setWorkload] = useState<{ old: string; new: string }>({
    old: urlParams.get("oldWorkload") || "",
    new: urlParams.get("newWorkload") || "",
  });
  const navigate = useNavigate();

  const { data: vitessRefs } = useApiCall<VitessRefs>({
    url: `${import.meta.env.VITE_API_URL}vitess/refs`,
    queryKey: ["vitessRefs"],
  });

  const shouldFetchCompareData = workload.old && workload.new && gitRef;

  const {
    data: data,
    isLoading: isMacrobenchLoading,
    error: macrobenchError,
  } = useApiCall<CompareResult>(
    shouldFetchCompareData
      ? {
          url: `${
            import.meta.env.VITE_API_URL
          }fk/compare?sha=${gitRef}&newWorkload=${workload.new}&oldWorkload=${
            workload.old
          }`,
          queryKey: ["compareResult", gitRef, workload.old, workload.new],
        }
      : {
          url: null,
          queryKey: ["compareResult", gitRef, workload.old, workload.new],
        }
  );

  let formattedData = data !== undefined ? formatCompareResult(data) : null;

  useEffect(() => {
    navigate(
      `?sha=${gitRef}&oldWorkload=${workload.old}&newWorkload=${workload.new}`
    );
  }, [gitRef, workload.old, workload.new]);

  return (
    <>
      {vitessRefs && (
        <ForeignKeysHero
          gitRef={gitRef}
          setGitRef={setGitRef}
          workload={workload}
          setWorkload={setWorkload}
          vitessRefs={vitessRefs}
        />
      )}

      <section className="flex flex-col items-center">
        {macrobenchError && (
          <div className="text-destructive">{<>{macrobenchError}</>}</div>
        )}

        {!isMacrobenchLoading && data === undefined && (
          <div className="md:text-xl text-primary">
            Chose two commits to compare
          </div>
        )}

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
        {!isMacrobenchLoading &&
          formattedData &&
          data !== undefined &&
          vitessRefs && (
            <>
              <div className="w-[80vw] xl:w-[60vw] my-12">
                <Card className="border-border">
                  <CardHeader className="flex flex-col gap-4 md:gap-0 md:flex-row justify-between pt-6">
                    <CardTitle className="text-2xl md:text-4xl text-primary">
                      {gitRef == "" ? (
                        "N/A"
                      ) : (
                        <Link
                          to={`https://github.com/vitessio/vitess/commit/${gitRef}`}
                          target="__blank"
                        >
                          {getRefName(gitRef, vitessRefs)}
                        </Link>
                      )}
                    </CardTitle>
                    <Button
                      variant="outline"
                      size="sm"
                      className="h-8 w-fit border-dashed mt-4 md:mt-0"
                      disabled
                    >
                      <PlusCircledIcon className="mr-2 h-4 w-4 text-primary" />
                      {/* <Link
                          to={`/macrobench/queries/compare?ltag=${gitRef.old}&rtag=${gitRef.new}&workload=${macro.workload}`}
                        > */}
                      See Query Plan {/* </Link> */}
                    </Button>
                  </CardHeader>
                  <CardContent className="w-full p-0">
                    {data.missing_results ? (
                      <div className="text-center md:text-xl text-destructive pb-12">
                        Missing results for this workloads
                      </div>
                    ) : (
                      <MacroBenchmarkTable
                        data={formattedData}
                        new={workload.new}
                        old={workload.old}
                        isGitRef={false}
                        vitessRefs={undefined}
                      />
                    )}
                  </CardContent>
                </Card>
              </div>
            </>
          )}
      </section>
    </>
  );
}
