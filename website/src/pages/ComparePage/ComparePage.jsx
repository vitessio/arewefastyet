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
import Microbench from "../MicroPage/components/Microbench/Microbench";

export default function Compare() {
  const urlParams = new URLSearchParams(window.location.search);

  const [gitRef, setGitRef] = useState({
    left: urlParams.get("ltag") || "Left",
    right: urlParams.get("rtag") || "Right",
  });

  const {
    data: dataMacrobench,
    isLoading: isMacrobenchLoading,
    textLoading: macrobenchTextLoading,
    error: macrobenchError,
  } = useApiCall(
    `${import.meta.env.VITE_API_URL}macrobench/compare?rtag=${
      gitRef.right
    }&ltag=${gitRef.left}`
  );

  const { data: dataMicrobench, isLoading: isMicrobenchLoading } = useApiCall(
    `${import.meta.env.VITE_API_URL}microbench/compare?rtag=${
      gitRef.right
    }&ltag=${gitRef.left}`
  );

  const navigate = useNavigate();

  useEffect(() => {
    navigate(`?ltag=${gitRef.left}&rtag=${gitRef.right}`);
  }, [gitRef.left, gitRef.right]);

  return (
    <>
      <Hero gitRef={gitRef} setGitRef={setGitRef} />
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

      {!isMacrobenchLoading && !macrobenchTextLoading && dataMacrobench && (
        <section className="flex flex-col items-center">
          <h3 className="my-6 text-primary text-2xl">Macro Benchmarks</h3>
          <div className="flex flex-col gap-y-20">
            {dataMacrobench.map((macro, index) => {
              return (
                <div key={index}>
                  <Macrobench
                    data={macro}
                    gitRef={{
                      left: gitRef.left.slice(0, 8),
                      right: gitRef.right.slice(0, 8),
                    }}
                    commits={{ left: gitRef.left, right: gitRef.right }}
                  />
                </div>
              );
            })}
          </div>
        </section>
      )}

      {/*{!isMicrobenchLoading && dataMicrobench && (*/}
      {/*  <section className="flex flex-col items-center">*/}
      {/*    <h3 className="my-6 text-primary text-2xl">Micro benchmarks</h3>*/}
      {/*    <div className="micro__thead space--between">*/}
      {/*      <span className="width--12em">Package</span>*/}
      {/*      <span className="width--14em">Benchmark Name</span>*/}
      {/*      <span className="width--18em hiddenMobile">*/}
      {/*        Number of Iterations*/}
      {/*      </span>*/}
      {/*      <span className="width--18em hiddenTablet">Time/op</span>*/}
      {/*      <span className="width--6em">More</span>*/}
      {/*    </div>*/}
      {/*    <figure className="micro__thead__line"></figure>*/}
      {/*    <div className="space--between--flex data__top hiddenMobile">*/}
      {/*      <div className="width--12em"></div>*/}
      {/*      <div className="width--14em"></div>*/}
      {/*      <div className="width--18em space--between--flex">*/}
      {/*        <span className="width--100">{gitRef.left.slice(0, 8)}</span>*/}
      {/*        <span className="width--100">{gitRef.right.slice(0, 8)}</span>*/}
      {/*        <span className="width--100">Diff %</span>*/}
      {/*      </div>*/}
      {/*      <div className="width--18em space--between--flex hiddenTablet">*/}
      {/*        <span className="width--100">{gitRef.left.slice(0, 8)}</span>*/}
      {/*        <span className="width--100">{gitRef.right.slice(0, 8)}</span>*/}
      {/*        <span className="width--100">Diff %</span>*/}
      {/*      </div>*/}
      {/*      <div className="width--6em"></div>*/}
      {/*    </div>*/}
      {/*    {dataMicrobench.length > 0 &&*/}
      {/*      dataMicrobench[0].PkgName !== "" &&*/}
      {/*      dataMicrobench[0].Name !== "" &&*/}
      {/*      dataMicrobench[0].SubBenchmarkName !== "" &&*/}
      {/*      dataMicrobench.map((micro, index) => {*/}
      {/*        const isEvenIndex = index % 2 === 0;*/}
      {/*        const backgroundGrey = isEvenIndex ? "grey--background" : "";*/}
      {/*        return (*/}
      {/*          <Microbench*/}
      {/*            data={micro}*/}
      {/*            key={index}*/}
      {/*            className={backgroundGrey}*/}
      {/*            gitRefLeft={gitRef.left.slice(0, 8)}*/}
      {/*            gitRefRight={gitRef.right.slice(0, 8)}*/}
      {/*          />*/}
      {/*        );*/}
      {/*      })}*/}
      {/*  </section>*/}
      {/*)}*/}
    </>
  );
}
