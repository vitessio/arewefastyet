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

import React, { useState } from "react";
import { Link } from "react-router-dom";
import { formatByteForGB, fixed } from "../../../utils/Utils";
import { Swiper, SwiperSlide } from "swiper/react";
import PulseLoader from "react-spinners/PulseLoader";

import "./macrobenchmobile.css";
import "swiper/css";
import "swiper/css/effect-cards";

import { EffectCards } from "swiper";

const MacrobenchMobile = ({
  data,
  gitRefLeft,
  gitRefRight,
  setCurrentSlideIndexMobile,
  currentSlideIndexMobile,
  commitHashLeft,
  commitHashRight,
  textLoading,
}) => {
  const handleSlideChange = (swiper) => {
    setCurrentSlideIndexMobile(swiper.realIndex);
  };

  const renderDataOrLoader = (data, loading) => {
    if (loading) {
      return (
        <span id="textloader">
          <PulseLoader loading={true} size={5} color="var(--font-color)" />
        </span>
      );
    } else {
      return <span>{data}</span>;
    }
  };

  return (
    <div className="macrobench__mobile">
      <div className="macrobench__mobile__header">
        <h3>{data.type}</h3>
        <div className="linkQuery">
          Click{" "}
          <Link
            to={`/macrobench/queries/compare?ltag=${commitHashLeft}&rtag=${commitHashRight}&type=${data.type}`}
          >
            here
          </Link>{" "}
          to see the query plans
        </div>
      </div>
      <div className="macrobenchMobile__container flex">
        <div className="macrobench__Sidebar flex--column">
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

        <div className="macrobench__component__container flex">
          <Swiper
            effect={"cards"}
            grabCursor={true}
            modules={[EffectCards]}
            onSlideChange={handleSlideChange}
            initialSlide={currentSlideIndexMobile}
            className="mySwiper"
          >
            <SwiperSlide>
              <div className="macrobench__data flex--column">
                <h4>Impoved by %</h4>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Diff.qps.total, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Diff.qps.reads, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Diff.qps.writes, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Diff.qps.other, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Diff.tps, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Diff.latency, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Diff.errors, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Diff.reconnects, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Diff.time, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Diff.threads, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.DiffMetrics.TotalComponentsCPUTime, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.DiffMetrics.ComponentsCPUTime.vtgate, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.DiffMetrics.ComponentsCPUTime.vttablet, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    formatByteForGB(
                      data.diff.DiffMetrics.TotalComponentsMemStatsAllocBytes
                    ),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(
                      data.diff.DiffMetrics.ComponentsMemStatsAllocBytes.vtgate,
                      2
                    ),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span" id="borderNone">
                  {renderDataOrLoader(
                    fixed(
                      data.diff.DiffMetrics.ComponentsMemStatsAllocBytes
                        .vttablet,
                      0
                    ),
                    textLoading
                  )}
                </span>
              </div>
            </SwiperSlide>
            <SwiperSlide>
              <div className="macrobench__data flex--column">
                {gitRefLeft == "" || gitRefLeft == "Left" ? (
                  <h4>{gitRefLeft ? gitRefLeft : "Left"}</h4>
                ) : (
                  <a
                    target="blank"
                    href={`https://github.com/vitessio/vitess/commit/${commitHashLeft}`}
                  >
                    <h4>{gitRefLeft ? gitRefLeft : "Left"}</h4>
                  </a>
                )}
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Left.Result.qps.total, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Left.Result.qps.reads, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Left.Result.qps.writes, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Left.Result.qps.other, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Left.Result.tps, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Left.Result.latency, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Left.Result.errors, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Left.Result.reconnects, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Left.Result.time, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Left.Result.threads, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Left.Metrics.TotalComponentsCPUTime, 0),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Left.Metrics.ComponentsCPUTime.vtgate, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Left.Metrics.ComponentsCPUTime.vttablet, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    formatByteForGB(
                      data.diff.Left.Metrics.TotalComponentsMemStatsAllocBytes
                    ),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    formatByteForGB(
                      data.diff.Left.Metrics.ComponentsMemStatsAllocBytes.vtgate
                    ),
                    textLoading
                  )}
                </span>
                <span id="borderNone" className="macrobench__data__span">
                  {renderDataOrLoader(
                    formatByteForGB(
                      data.diff.Left.Metrics.ComponentsMemStatsAllocBytes
                        .vttablet
                    ),
                    textLoading
                  )}
                </span>
              </div>
            </SwiperSlide>
            <SwiperSlide>
              <div className="macrobench__data flex--column">
                {gitRefRight == "" || gitRefRight == "Right" ? (
                  <h4>{gitRefRight ? gitRefRight : "Right"}</h4>
                ) : (
                  <a
                    target="blank"
                    href={`https://github.com/vitessio/vitess/commit/${commitHashRight}`}
                  >
                    <h4>{gitRefRight ? gitRefRight : "Right"}</h4>
                  </a>
                )}
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Right.Result.qps.total, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Right.Result.qps.reads, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Right.Result.qps.writes, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Right.Result.qps.other, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Right.Result.tps, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Right.Result.latency, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Right.Result.errors, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Right.Result.reconnects, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Right.Result.time, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Right.Result.threads, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Right.Metrics.TotalComponentsCPUTime, 0),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(data.diff.Right.Metrics.ComponentsCPUTime.vtgate, 2),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    fixed(
                      data.diff.Right.Metrics.ComponentsCPUTime.vttablet,
                      2
                    ),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    formatByteForGB(
                      data.diff.Right.Metrics.TotalComponentsMemStatsAllocBytes
                    ),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span">
                  {renderDataOrLoader(
                    formatByteForGB(
                      data.diff.Right.Metrics.ComponentsMemStatsAllocBytes
                        .vtgate
                    ),
                    textLoading
                  )}
                </span>
                <span className="macrobench__data__span" id="borderNone">
                  {renderDataOrLoader(
                    formatByteForGB(
                      data.diff.Right.Metrics.ComponentsMemStatsAllocBytes
                        .vttablet
                    ),
                    textLoading
                  )}
                </span>
              </div>
            </SwiperSlide>
          </Swiper>
        </div>
      </div>
    </div>
  );
};

export default MacrobenchMobile;
