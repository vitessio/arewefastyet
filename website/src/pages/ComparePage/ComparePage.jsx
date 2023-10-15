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
import { SwiperSlide } from "swiper/react";
import useApiCall from "../../utils/Hook";

import Macrobench from "../../components/MacroComponents/Macrobench/Macrobench";
import MacrobenchMobile from "../../components/MacroComponents/MacrobenchMobile/MacrobenchMobile";
import Microbench from "../../components/Microbench/Microbench";
import Hero from "./components/Hero";

export default function Compare() {
  const urlParams = new URLSearchParams(window.location.search);

  const [gitRef, setGitRef] = useState({
    left: urlParams.get("ltag") || "Left",
    right: urlParams.get("rtag") || "Right",
  });

  const [currentSlideIndexMobile, setCurrentSlideIndexMobile] = useState(
    urlParams.get("ptagM") || "0"
  );

  const {
    data: dataMacrobench,
    isLoading: isMacrobenchLoading,
    error: macrobenchError,
    textLoading: macroTextLoading,
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

  // Changing the URL relative to the reference of a selected benchmark.
  // Storing the carousel position as a URL parameter.
  const navigate = useNavigate();

  useEffect(() => {
    navigate(
      `?ltag=${gitRef.left}&rtag=${gitRef.right}&ptagM=${currentSlideIndexMobile}`
    );
  }, [gitRef.left, gitRef.right, currentSlideIndexMobile]);

  const handleSlideChange = (swiper) => {
    setCurrentSlideIndexMobile(swiper.realIndex);
  };

  return (
    <>
      <Hero setGitRef={setGitRef} />
      {macrobenchError && <div className="apiError">{macrobenchError}</div>}

      {isMacrobenchLoading && (
        <div className="flex justify-center items-center">
          <RingLoader
            loading={isMacrobenchLoading}
            color="#E77002"
            size={300}
          />
        </div>
      )}

      {!isMacrobenchLoading && dataMacrobench && (
        <section className="flex flex-col items-center">
          <h3 className="my-6 text-primary text-2xl">Macro Benchmarks</h3>
          <div className="compare__macrobench__Container flex">
            <div className="compare__carousel__container">
              {dataMacrobench.map((macro, index) => {
                return (
                  <div key={index}>
                    <Macrobench
                      data={macro}
                      textLoading={macroTextLoading}
                      gitRefLeft={gitRef.left.slice(0, 8)}
                      gitRefRight={gitRef.right.slice(0, 8)}
                      swiperSlide={SwiperSlide}
                      commitHashLeft={gitRef.left}
                      commitHashRight={gitRef.right}
                    />
                    <MacrobenchMobile
                      data={macro}
                      gitRefLeft={gitRef.left.slice(0, 8)}
                      gitRefRight={gitRef.right.slice(0, 8)}
                      swiperSlide={SwiperSlide}
                      textLoading={macroTextLoading}
                      handleSlideChange={handleSlideChange}
                      setCurrentSlideIndexMobile={setCurrentSlideIndexMobile}
                      currentSlideIndexMobile={currentSlideIndexMobile}
                      commitHashLeft={gitRef.left}
                      commitHashRight={gitRef.right}
                    />
                  </div>
                );
              })}
            </div>
          </div>
        </section>
      )}

      {!isMicrobenchLoading && dataMicrobench && (
        <section className="flex flex-col items-center">
          <h3 className="my-6 text-primary text-2xl">Micro benchmarks</h3>
          <div className="micro__thead space--between">
            <span className="width--12em">Package</span>
            <span className="width--14em">Benchmark Name</span>
            <span className="width--18em hiddenMobile">
              Number of Iterations
            </span>
            <span className="width--18em hiddenTablet">Time/op</span>
            <span className="width--6em">More</span>
          </div>
          <figure className="micro__thead__line"></figure>
          <div className="space--between--flex data__top hiddenMobile">
            <div className="width--12em"></div>
            <div className="width--14em"></div>
            <div className="width--18em space--between--flex">
              <span className="width--100">{gitRef.left.slice(0, 8)}</span>
              <span className="width--100">{gitRef.right.slice(0, 8)}</span>
              <span className="width--100">Diff %</span>
            </div>
            <div className="width--18em space--between--flex hiddenTablet">
              <span className="width--100">{gitRef.left.slice(0, 8)}</span>
              <span className="width--100">{gitRef.right.slice(0, 8)}</span>
              <span className="width--100">Diff %</span>
            </div>
            <div className="width--6em"></div>
          </div>
          {dataMicrobench.length > 0 &&
            dataMicrobench[0].PkgName !== "" &&
            dataMicrobench[0].Name !== "" &&
            dataMicrobench[0].SubBenchmarkName !== "" &&
            dataMicrobench.map((micro, index) => {
              const isEvenIndex = index % 2 === 0;
              const backgroundGrey = isEvenIndex ? "grey--background" : "";
              return (
                <Microbench
                  data={micro}
                  key={index}
                  className={backgroundGrey}
                  gitRefLeft={gitRef.left.slice(0, 8)}
                  gitRefRight={gitRef.right.slice(0, 8)}
                />
              );
            })}
        </section>
      )}
    </>
  );
}
