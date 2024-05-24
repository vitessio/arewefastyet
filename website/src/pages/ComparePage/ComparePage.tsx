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

import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import RingLoader from "react-spinners/RingLoader";
import useApiCall from "../../utils/Hook";
import Hero from "./components/Hero";
import Macrobench from "../../common/Macrobench";
import { CompareData } from '@/types'
 
export default function Compare() {
  const urlParams = new URLSearchParams(window.location.search);

  const [gitRef, setGitRef] = useState({
    old: urlParams.get("old") || "",
    new: urlParams.get("new") || "",
  });

  const {
    data: compareData,
    isLoading: isMacrobenchLoading,
    textLoading: macrobenchTextLoading,
    error: macrobenchError,
  } = useApiCall<CompareData>(
    `${import.meta.env.VITE_API_URL}macrobench/compare?new=${gitRef.new}&old=${gitRef.old}`
  );

  const navigate = useNavigate();

  useEffect(() => {
    navigate(`?old=${gitRef.old}&new=${gitRef.new}`);
  }, [gitRef.old, gitRef.new]);

  return (
    <>
      <Hero gitRef={gitRef} setGitRef={setGitRef} />
      {macrobenchError && (
        <div className="text-red-500 text-center m-2">{macrobenchError}</div>
      )}

      {(isMacrobenchLoading || macrobenchTextLoading) && (
        <div className="flex justify-center items-center">
          <RingLoader
            loading={isMacrobenchLoading || macrobenchTextLoading}
            color="#E77002"
            size={300}
          />
        </div>
      )}

      {!isMacrobenchLoading && !macrobenchTextLoading && compareData && compareData.length > 0 && (
        <section className="flex flex-col items-center">
          <h3 className="my-6 text-primary text-2xl">Macro Benchmarks</h3>
          <div className="flex flex-col gap-y-20">
            {compareData.map((macro, index) => {
              return (
                <div key={index}>
                  <Macrobench
                    data={macro}
                    gitRef={{
                      old: gitRef.old.slice(0, 8),
                      new: gitRef.new.slice(0, 8),
                    }}
                    commits={{ old: gitRef.old, new: gitRef.new }}
                  />
                </div>
              );
            })}
          </div>
        </section>
      )}
    </>
  );
}
