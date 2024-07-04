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
import Macrobench from "../../common/Macrobench";
import { CompareData } from '@/types'
import Hero, { HeroProps } from "@/common/Hero";
import { twMerge } from "tailwind-merge";

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

  function ComparisonInput(props: {
    className: any;
    gitRef: any;
    setGitRef: any;
    name: any;
  }) {
    const { className, gitRef, setGitRef, name } = props;

    return (
      <input
        type="text"
        name={name}
        className={twMerge(
          className,
          "relative text-xl px-6 py-2 bg-background focus:border-none focus:outline-none border border-primary"
        )}
        defaultValue={gitRef[name]}
        placeholder={`${name} SHA`}
        onChange={(event) =>
          setGitRef((p: any) => {
            return { ...p, [name]: event.target.value };
          })
        }
      />
    );
  }

  const heroProps: HeroProps = {
    title: "Compare versions",
    children: (
      <div>
        <h1 className="mb-3 text-front text-opacity-70">
          Enter SHAs to compare commits
        </h1>
        <div className="flex overflow-hidden bg-gradient-to-br from-primary to-theme p-[2px] rounded-full">
          <ComparisonInput
            name="old"
            className="rounded-l-full"
            setGitRef={setGitRef}
            gitRef={gitRef}
          />
          <ComparisonInput
            name="new"
            className="rounded-r-full "
            setGitRef={setGitRef}
            gitRef={gitRef}
          />
        </div>
      </div>
    )
  };

  return (
    <>
      <Hero title={heroProps.title} description={heroProps.description}>
        {heroProps.children}
      </Hero>
      {macrobenchError && (
        <div className="text-red-500 text-center my-2">{macrobenchError}</div>
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
