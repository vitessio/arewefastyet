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
import { Swiper, SwiperSlide } from "swiper/react";
import useApiCall from "../../utils/Hook";

import "../Compare/compare.css";
import "swiper/css";
import "swiper/css/pagination";

import { Mousewheel, Pagination, Keyboard } from "swiper";
import Macrobench from "../../components/MacroComponents/Macrobench/Macrobench";
import MacrobenchMobile from "../../components/MacroComponents/MacrobenchMobile/MacrobenchMobile";
import Microbench from "../../components/Microbench/Microbench";

const Compare = () => {
  const urlParams = new URLSearchParams(window.location.search);
  // The following code sets up state variables `gitRefLeft` and `gitRefRight` using the `useState` hook.
  // The values of these variables are based on the query parameters extracted from the URL.

  // If the 'ltag' query parameter is null or undefined, set the initial value of `gitRefLeft` to 'Left',
  // otherwise, use the value of the 'ltag' query parameter.
  const [gitRefLeft, setGitRefLeft] = useState(
    urlParams.get("ltag") == null ? "Left" : urlParams.get("ltag")
  );
  const [gitRefRight, setGitRefRight] = useState(
    urlParams.get("rtag") == null ? "Right" : urlParams.get("rtag")
  );
  const [currentSlideIndex, setCurrentSlideIndex] = useState(
    urlParams.get("ptag") == null ? "0" : urlParams.get("ptag")
  );
  const [currentSlideIndexMobile, setCurrentSlideIndexMobile] = useState(
    urlParams.get("ptagM") == null ? "0" : urlParams.get("ptagM")
  );

  const {
    data: dataMacrobench,
    isLoading: isMacrobenchLoading,
    error: macrobenchError,
  } = useApiCall(
    `${
      import.meta.env.VITE_API_URL
    }macrobench/compare?rtag=${gitRefRight}&ltag=${gitRefLeft}`,
    [gitRefLeft, gitRefRight]
  );

  const {
    data: dataMicrobench,
    isLoading: isMicrobenchLoading,
    error: microbenchError,
  } = useApiCall(
    `${
      import.meta.env.VITE_API_URL
    }microbench/compare?rtag=${gitRefRight}&ltag=${gitRefLeft}`,
    [gitRefLeft, gitRefRight]
  );

  // Changing the URL relative to the reference of a selected benchmark.
  // Storing the carousel position as a URL parameter.
  const navigate = useNavigate();

  useEffect(() => {
    navigate(
      `?ltag=${gitRefLeft}&rtag=${gitRefRight}&ptag=${currentSlideIndex}&ptagM=${currentSlideIndexMobile}`
    );
  }, [gitRefLeft, gitRefRight, currentSlideIndex, currentSlideIndexMobile]);

  const handleInputChangeLeft = (e) => {
    setGitRefLeft(e.target.value);
  };

  const handleInputChangeRight = (e) => {
    setGitRefRight(e.target.value);
  };

  const handleSlideChange = (swiper) => {
    setCurrentSlideIndex(swiper.realIndex);
  };

  const slicedRef = gitRefLeft.slice(0, 8);
  return (
    <div className="compare">
      <div className="compare__top">
        <div className="justify--content form__container">
          <h3>
            Comparing{" "}
            <a href={`https://github.com/vitessio/vitess/commit/${gitRefLeft}`}>
              {gitRefLeft ? gitRefLeft.slice(0, 8) : "Left"}
            </a>{" "}
            with{" "}
            <a
              href={`https://github.com/vitessio/vitess/commit/${gitRefRight}`}
            >
              {gitRefRight ? gitRefRight.slice(0, 8) : "Right"}
            </a>
          </h3>
          <form className="justify--content">
            <input
              type="text"
              value={gitRefLeft === "Left" ? "" : gitRefLeft}
              onChange={handleInputChangeLeft}
              placeholder="Left commit SHA"
              className="form__inputLeft"
            ></input>
            <input
              type="text"
              value={gitRefRight === "Right" ? "" : gitRefRight}
              onChange={handleInputChangeRight}
              placeholder="Right commit SHA"
              className="form__inputRight"
            ></input>
          </form>
        </div>

        {macrobenchError ? (
          <div className="macrobench__apiError">{macrobenchError}</div>
        ) : isMacrobenchLoading ? (
          <div className="loadingSpinner">
            <RingLoader
              loading={isMacrobenchLoading}
              color="#E77002"
              size={300}
            />
          </div>
        ) : (
          <>
            <h3 className="compare__macrobench__title">Macro Benchmarks</h3>
            <div className="compare__macrobench__Container flex">
              <div className="compare__macrobench__Sidebar flex--column">
                <span>QPS Total</span>
                <figure className="macrobench__Sidebar__line"></figure>
                <span>QPS Reads</span>
                <figure className="macrobench__Sidebar__line"></figure>
                <span>QPS Writes</span>
                <figure className="macrobench__Sidebar__line"></figure>
                <span>QPS Other</span>
                <figure className="macrobench__Sidebar__line"></figure>
                <span>TPS</span>
                <figure className="macrobench__Sidebar__line"></figure>
                <span>Latency</span>
                <figure className="macrobench__Sidebar__line"></figure>
                <span>Errors</span>
                <figure className="macrobench__Sidebar__line"></figure>
                <span>Reconnects</span>
                <figure className="macrobench__Sidebar__line"></figure>
                <span>Time</span>
                <figure className="macrobench__Sidebar__line"></figure>
                <span>Threads</span>
                <figure className="macrobench__Sidebar__line"></figure>
                <span>Total CPU time</span>
                <figure className="macrobench__Sidebar__line"></figure>
                <span>CPU time vtgate</span>
                <figure className="macrobench__Sidebar__line"></figure>
                <span>CPU time vttablet</span>
                <figure className="macrobench__Sidebar__line"></figure>
                <span>Total Allocs bytes</span>
                <figure className="macrobench__Sidebar__line"></figure>
                <span>Allocs bytes vtgate</span>
                <figure className="macrobench__Sidebar__line"></figure>
                <span>Allocs bytes vttablet</span>
              </div>
              <div className="compare__carousel__container">
                <Swiper
                  direction={"vertical"}
                  slidesPerView={1}
                  spaceBetween={30}
                  mousewheel={true}
                  keyboard={{
                    enabled: true,
                  }}
                  pagination={{
                    clickable: true,
                  }}
                  modules={[Mousewheel, Pagination, Keyboard]}
                  onSlideChange={handleSlideChange}
                  initialSlide={currentSlideIndex}
                  className="mySwiper"
                >
                  {dataMacrobench.map((macro, index) => {
                    return (
                      <SwiperSlide key={index}>
                        <Macrobench
                          data={macro}
                          gitRefLeft={gitRefLeft.slice(0, 8)}
                          gitRefRight={gitRefRight.slice(0, 8)}
                          swiperSlide={SwiperSlide}
                        />
                        <MacrobenchMobile
                          data={macro}
                          gitRefLeft={gitRefLeft}
                          gitRefRight={gitRefRight}
                          swiperSlide={SwiperSlide}
                          handleSlideChange={handleSlideChange}
                          setCurrentSlideIndexMobile={
                            setCurrentSlideIndexMobile
                          }
                          currentSlideIndexMobile={currentSlideIndexMobile}
                        />
                      </SwiperSlide>
                    );
                  })}
                </Swiper>
              </div>
            </div>

            <div className="compare__micro__container">
              <h3>Micro benchmarks</h3>
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
                  <span className="width--100">{gitRefLeft.slice(0, 8)}</span>
                  <span className="width--100">{gitRefRight.slice(0, 8)}</span>
                  <span className="width--100">Diff %</span>
                </div>
                <div className="width--18em space--between--flex hiddenTablet">
                  <span className="width--100">{gitRefLeft.slice(0, 8)}</span>
                  <span className="width--100">{gitRefRight.slice(0, 8)}</span>
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
                      gitRefLeft={gitRefLeft.slice(0, 8)}
                      gitRefRight={gitRefRight.slice(0, 8)}
                    />
                  );
                })}
            </div>
          </>
        )}
      </div>
    </div>
  );
};

export default Compare;
