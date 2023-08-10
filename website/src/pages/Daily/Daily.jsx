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
import RingLoader from "react-spinners/RingLoader";
import { useNavigate } from "react-router-dom";
import useApiCall from "../../utils/Hook";

import "../Daily/daily.css";

import { formatByteForGB } from "../../utils/Utils";
import ResponsiveChart from "../../components/DailyComponents/Chart/Chart";
import DailySummary from "../../components/DailyComponents/DailySummary/DailySummary";

const Daily = () => {
  const urlParams = new URLSearchParams(window.location.search);
  const [benchmarkType, setBenchmarktype] = useState(
    urlParams.get("type") == null ? "" : urlParams.get("type")
  );

  const {
    data: dataDailySummary,
    isLoading: isLoadingDailySummary,
    error: errorDailySummary,
  } = useApiCall(`${import.meta.env.VITE_API_URL}daily/summary`);

  const {
    data: dataDaily,
    error: dailyError,
    textLoading: dailyTextLoading,
  } = useApiCall(`${import.meta.env.VITE_API_URL}daily?type=${benchmarkType}`, [
    benchmarkType,
  ]);


  // Changing the URL relative to the reference of a selected benchmark.
  // Storing the carousel position as a URL parameter.
  const navigate = useNavigate();

  useEffect(() => {
    navigate(`?type=${benchmarkType}`);
  }, [benchmarkType]);

  const TPSData = [
    {
      id: "TPS",
      data: [],
    },
  ];

  const QPSData = [
    {
      id: "Reads",
      data: [],
    },
    {
      id: "Total",
      data: [],
    },

    {
      id: "Writes",
      data: [],
    },

    {
      id: "Other",
      data: [],
    },
  ];

  const latencyData = [
    {
      id: "Latency",
      data: [],
    },
  ];

  const CPUTimeData = [
    {
      id: "Total",
      data: [],
    },
    {
      id: "vtgate",
      data: [],
    },
    {
      id: "vttablet",
      data: [],
    },
  ];

  const MemBytesData = [
    {
      id: "Total",
      data: [],
    },
    {
      id: "vtgate",
      data: [],
    },
    {
      id: "vttablet",
      data: [],
    },
  ];

  for (const item of dataDaily) {
    const xValue = item.GitRef.slice(0, 8);

    // TPS Data

    TPSData[0].data.push({
      x: xValue,
      y: item.Result.tps,
    });

    // QPS Data

    QPSData[0].data.push({
      x: xValue,
      y: item.Result.qps.reads,
    });

    QPSData[1].data.push({
      x: xValue,
      y: item.Result.qps.total,
    });

    QPSData[2].data.push({
      x: xValue,
      y: item.Result.qps.writes,
    });

    QPSData[3].data.push({
      x: xValue,
      y: item.Result.qps.other,
    });

    // Latency Data

    latencyData[0].data.push({
      x: xValue,
      y: item.Result.latency,
    });

    // CPUTime Data

    CPUTimeData[0].data.push({
      x: xValue,
      y: item.Metrics.TotalComponentsCPUTime,
    });

    CPUTimeData[1].data.push({
      x: xValue,
      y: item.Metrics.ComponentsCPUTime.vtgate,
    });

    CPUTimeData[2].data.push({
      x: xValue,
      y: item.Metrics.ComponentsCPUTime.vttablet,
    });

    //MemStatsAllocBytes Data

    MemBytesData[0].data.push({
      x: xValue,
      y: formatByteForGB(item.Metrics.TotalComponentsMemStatsAllocBytes),
    });

    MemBytesData[1].data.push({
      x: xValue,
      y: formatByteForGB(item.Metrics.ComponentsMemStatsAllocBytes.vtgate),
    });

    MemBytesData[2].data.push({
      x: xValue,
      y: formatByteForGB(item.Metrics.ComponentsMemStatsAllocBytes.vttablet),
    });
  }

  const allChartData = [
    {
      data: QPSData,
      title: "QPS (Queries per second)",
      colors: ["#fad900", "orange", "brown", "purple"],
    },
    {
      data: TPSData,
      title: "TPS (Transactions per second)",
      colors: ["#fad900"],
    },
    {
      data: latencyData,
      title: "Latency (Milliseconds)",
      colors: ["#fad900"],
    },
    {
      data: CPUTimeData,
      title: "CPU Time",
      colors: ["#fad900", "orange", "brown"],
    },
    {
      data: MemBytesData,
      title: "Allocated Bytes",
      colors: ["#fad900", "orange", "brown"],
    },
  ];

  return (
    <div className="daily">
      <div className="daily__top">
        <h2 className="header--title">Daily</h2>
        <span className="header--text">
          Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer a
          augue mi. Etiam sed imperdiet ligula, vel elementum velit. Phasellus
          sodales felis eu condimentum convallis. Suspendisse sodales malesuada
          iaculis. Mauris molestie placerat ex non malesuada. Curabitur eget
          sagittis eros. Aliquam aliquam sem non tincidunt volutpat.
        </span>
      </div>
      <figure className="line"></figure>
      {errorDailySummary || dailyError ? (
        <div className="apiError">{error}</div>
      ) : isLoadingDailySummary ? (
        <div className="loadingSpinner">
          <RingLoader
            loading={isLoadingDailySummary}
            color="#E77002"
            size={300}
          />
        </div>
      ) : (
        <>
          <div className="daily__summary__container justify--content">
            {dataDailySummary.map((dailySummary, index) => {
              return (
                <DailySummary
                  key={index}
                  data={dailySummary}
                  setBenchmarktype={setBenchmarktype}
                  benchmarkType={benchmarkType}
                />
              );
            })}
          </div>
          <figure className="line"></figure>
          {dailyTextLoading ? (
            <div className="loadingSpinner">
              <RingLoader
                loading={dailyTextLoading}
                color="#E77002"
                size={300}
              />
            </div>
          ) : benchmarkType !== "" ? (
            <div className="daily__container">
              {allChartData.map((chartData, index) => (
                <ResponsiveChart
                  key={index}
                  data={chartData.data}
                  title={chartData.title}
                  colors={chartData.colors}
                  isFirstChart={index === 0}
                />
              ))}
            </div>
          ) : null}
        </>
      )}
    </div>
  );
};

export default Daily;
