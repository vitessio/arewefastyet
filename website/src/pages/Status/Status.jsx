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

import React, { useRef } from "react";
import RingLoader from "react-spinners/RingLoader";
import { v4 as uuidv4 } from "uuid";
import useApiCall from "../../utils/Hook";
import CountUp from "react-countup";

import PreviousExe from "../../components/StatusComponents/PreviousExecutions/PreviousExe";
import ExeQueue from "../../components/StatusComponents/ExecutionQueue/ExeQueue";
import PreviousExeResMobile from "../../components/StatusComponents/PreviousExeResponsiveMobile/PreviousExeResMobile";
import PreviousExeResTablet from "../../components/StatusComponents/PreviousExecutionResponsiveTablet/PreviousExeResTablet";
import "./status.css";

const Status = () => {
  const {
    data: dataQueue,
    isLoading: isLoadingQueue,
    error: errorQueue,
  } = useApiCall(`${import.meta.env.VITE_API_URL}queue`);
  const {
    data: dataPreviousExe,
    isLoading: isLoadingPreviousExe,
    error: errorPreviousExe,
  } = useApiCall(`${import.meta.env.VITE_API_URL}recent`);
  const {
    data: dataStatusStats,
    isLoading: isLoadingStatusStats,
    error: errorStatusStats,
  } = useApiCall(`${import.meta.env.VITE_API_URL}status/stats`);

  return (
    <div className="status">
      <article className="status__top justify--content">
        <div className="status__top__text">
          <h2>Status</h2>
          <span>
            Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer a
            augue mi. Etiam sed imperdiet ligula, vel elementum velit. Phasellus
            sodales felis eu condimentum convallis. Suspendisse sodales
            malesuada iaculis. Mauris molestie placerat ex non malesuada.
            Curabitur eget sagittis eros. Aliquam aliquam sem non tincidunt
            volutpat.
          </span>
        </div>

        <div className="statusStats flex--column">
          <CountUp
            start={0}
            end={dataStatusStats.Total}
            style={{ fontSize: "5rem", color: "#E77002" }}
            duration={3}
          />
          <span className="countUp__title">Total Benchmark</span>
          <CountUp
            start={0}
            end={dataStatusStats.Finished}
            style={{ fontSize: "5rem", color: "#E77002" }}
            duration={3}
          />
          <span className="countUp__title">Successful benchmarks</span>
          <CountUp
            start={0}
            end={dataStatusStats.Last30Days}
            style={{ fontSize: "5rem", color: "#E77002" }}
            duration={3}
          />
          <span className="countUp__title">
            Total Benchmark (last 30 days)
          </span>
        </div>
      </article>
      <figure className="line"></figure>

      {/* EXECUTION QUEUE  */}
      {isLoadingQueue ? (
        <div className="loadingSpinner">
          <RingLoader loading={isLoadingQueue} color="#E77002" size={300} />
        </div>
      ) : dataQueue.length > 0 ? (
        <>
          <article className="queue">
            <h3>Executions Queue</h3>
            <div className="queue__top flex">
              <span className="width--6em">SHA</span>
              <span className="width--11em"> Source</span>
              <span className="width--11em">Type</span>
              <span className="width--5em">Pull Request</span>
            </div>
            <figure className="queue__top__line"></figure>
            {dataQueue.map((queue, index) => {
              return <ExeQueue data={queue} key={index} />;
            })}
          </article>
          <figure className="line"></figure>
        </>
      ) : null}

      {/* PREVIOUS EXECUTIONS */}

      {isLoadingPreviousExe ? (
        <div className="loadingSpinner">
          <RingLoader
            loading={isLoadingPreviousExe}
            color="#E77002"
            size={300}
          />
        </div>
      ) : dataPreviousExe.length > 0 ? (
        <article className="previousExe">
          <h3> Previous Executions</h3>
          <div className="previousExe__top flex">
            <span className="width--6em hiddenMobile hiddenTablet">UUID</span>
            <span className="width--6em hiddenMobile">SHA</span>
            <span className="width--11em">Source</span>
            <span className="width--11em hiddenMobile">Started</span>
            <span className="width--11em hiddenMobile">Finished</span>
            <span className="width--11em hiddenMobile hiddenTablet">Type</span>
            <span className="width--5em hiddenMobile hiddenTablet">PR</span>
            <span className="width--6em hiddenMobile hiddenTablet">
              Go version
            </span>
            <span className="width--6em">Status</span>
            <span className="hiddenDesktop width--3em">More</span>
          </div>
          <figure className="previousExe__top__line"></figure>

          {dataPreviousExe.map((previousExe, index) => {
            return (
              <React.Fragment key={uuidv4()}>
                <PreviousExe data={previousExe} key={index} />
                <PreviousExeResMobile data={previousExe} key={uuidv4()} />
                <PreviousExeResTablet data={previousExe} key={uuidv4()} />
              </React.Fragment>
            );
          })}
        </article>
      ) : null}

      {errorQueue ? <div className="apiError">{errorQueue}</div> : null}
    </div>
  );
};

export default Status;
