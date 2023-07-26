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

import "../Micro/micro.css";

import {
  errorApi,
  closeDropDownValue,
  updateCommitHash,
  openDropDown,
  valueDropDown,
} from "../../utils/Utils";
import Microbench from "../../components/Microbench/Microbench";

const Micro = () => {
  const urlParams = new URLSearchParams(window.location.search);
  const [gitRefLeft, setGitRefLeft] = useState(
    urlParams.get("ltag") == null ? "Left" : urlParams.get("ltag")
  );
  const [gitRefRight, setGitRefRight] = useState(
    urlParams.get("rtag") == null ? "Right" : urlParams.get("rtag")
  );
  const [openDropDownLeft, setOpenDropDownLeft] = useState(closeDropDownValue);
  const [openDropDownRight, setOpenDropDownRight] =
    useState(closeDropDownValue);
  const [dataRefs, setDataRefs] = useState([]);
  const [dataMicrobench, setDataMicrobench] = useState([]);
  const [isFirstCallFinished, setIsFirstCallFinished] = useState(false);
  const [commitHashLeft, setCommitHashLeft] = useState("");
  const [commitHashRight, setCommitHashRight] = useState("");
  const [error, setError] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const responseRefs = await fetch(
          `${import.meta.env.VITE_API_URL}vitess/refs`
        );

        const jsonDataRefs = await responseRefs.json();

        setDataRefs(jsonDataRefs);
        setIsLoading(false);

        updateCommitHash(gitRefLeft, setCommitHashLeft, jsonDataRefs);
        updateCommitHash(gitRefRight, setCommitHashRight, jsonDataRefs);

        setIsFirstCallFinished(true);
      } catch (error) {
        console.log("Error while retrieving data from the API", error);
        setError(errorApi);
        setIsLoading(false);
      }
    };

    fetchData();
  }, []);

  useEffect(() => {
    if (isFirstCallFinished) {
      const fetchData = async () => {
        try {
          const responseMicrobench = await fetch(
            `${
              import.meta.env.VITE_API_URL
            }microbench/compare?rtag=${commitHashRight}&ltag=${commitHashLeft}`
          );
          console.log(commitHashLeft);
          const jsonDataMicrobench = await responseMicrobench.json();
          setDataMicrobench(jsonDataMicrobench);
        } catch (error) {
          console.log("Error while retrieving data from the API", error);
          setError(errorApi);
        }
      };
      fetchData();
    }
  }, [commitHashLeft, commitHashRight]);

  // Changing the URL relative to the reference of a selected benchmark.
  // Storing the carousel position as a URL parameter.
  const navigate = useNavigate();

  useEffect(() => {
    navigate(`?ltag=${gitRefLeft}&rtag=${gitRefRight}`);
  }, [gitRefLeft, gitRefRight]);

  return (
    <div className="micro">
      <div className="micro__top justify--content">
        <div className="micro__top__text">
          <h2>Compare Microbenchmarks</h2>
          <div className="micro__bottom__title justify--content">
            <figure
              className="micro__bottom__DropDownLeft flex--column"
              style={{ maxHeight: `${openDropDownLeft}px` }}
            >
              <span
                className="DropDown__BaseSpan"
                onClick={() =>
                  openDropDown(openDropDownLeft, setOpenDropDownLeft)
                }
              >
                {gitRefLeft} <i className="fa-solid fa-circle-arrow-down"></i>
              </span>
              {dataRefs.map((ref, index) => {
                return (
                  <React.Fragment key={index}>
                    <figure className="dropDown--line"></figure>
                    <span
                      className="dropDown__ref"
                      onClick={() => {
                        valueDropDown(
                          ref,
                          setGitRefLeft,
                          setCommitHashLeft,
                          setOpenDropDownLeft
                        );
                      }}
                    >
                      {ref.Name}
                    </span>
                  </React.Fragment>
                );
              })}
            </figure>

            <figure
              className="micro__bottom__DropDownRight flex--column"
              style={{ maxHeight: `${openDropDownRight}px` }}
            >
              <span
                className="DropDown__BaseSpan"
                onClick={() =>
                  openDropDown(openDropDownRight, setOpenDropDownRight)
                }
              >
                {gitRefRight} <i className="fa-solid fa-circle-arrow-down"></i>
              </span>
              {dataRefs.map((ref, index) => {
                return (
                  <React.Fragment key={index}>
                    <figure className="dropDown--line"></figure>
                    <span
                      className="dropDown__ref"
                      onClick={() =>
                        valueDropDown(
                          ref,
                          setGitRefRight,
                          setCommitHashRight,
                          setOpenDropDownRight
                        )
                      }
                    >
                      {ref.Name}
                    </span>
                  </React.Fragment>
                );
              })}
            </figure>
          </div>
        </div>
        <figure className="microStats"></figure>
      </div>
      <figure className="line"></figure>
      <div className="micro__bottom">
        {error ? (
          <div className="apiError">{error}</div>
        ) : isLoading ? (
          <div className="loadingSpinner">
            <RingLoader loading={isLoading} color="#E77002" size={300} />
          </div>
        ) : (
          <div className="micro__container">
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
                <span className="width--100">{gitRefLeft}</span>
                <span className="width--100">{gitRefRight}</span>
                <span className="width--100">Diff %</span>
              </div>
              <div className="width--18em space--between--flex hiddenTablet">
                <span className="width--100">{gitRefLeft}</span>
                <span className="width--100">{gitRefRight}</span>
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
                    gitRefLeft={gitRefLeft}
                    gitRefRight={gitRefRight}
                  />
                );
              })}
          </div>
        )}
      </div>
    </div>
  );
};

export default Micro;
