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
import useApiCall from "../../utils/Hook";
import RingLoader from "react-spinners/RingLoader";

import { errorApi } from "../../utils/Utils";
import PRGitInfo from "../../components/PRcomponents/PRGitInfo";

const PR = () => {
  const [prNumber, setPrNumber] = useState("");
  const [dataPRInfo, setDataPRInfo] = useState([]);

  const {
    data: dataPRList,
    isLoading: isPRListLoading,
    error: PRListError,
  } = useApiCall(`${import.meta.env.VITE_API_URL}pr/list`, []);

  useEffect(() => {
    const fetchData = async () => {
      if (prNumber === "") {
        return;
      }
      try {
        const responsePRInfo = await fetch(
          `${import.meta.env.VITE_API_URL}pr/info/${prNumber}`
        );

        const jsonDataPRInfo = await responsePRInfo.json();
        console.log(jsonDataPRInfo.Main);
        setDataPRInfo(jsonDataPRInfo);

        if (jsonDataPRInfo.Main) {
          window.location.href = `/compare?ltag=${jsonDataPRInfo.Main}&rtag=${jsonDataPRInfo.PR}`;
        }
      } catch (error) {
        console.log("Error while retrieving data from the API", error);
        setError(errorApi);
      }
    };

    fetchData();
  }, [prNumber]);

  return (
    <div className="pr">
      <div className="pr__top justify--content">
        <h2 className="header--title">Pull Request</h2>
        <span className="header--text">
          If a given Pull Request on <b>vitessio/vitess</b> is labelled with the <b>Benchmark me</b> label
          the Pull Request will be handled and benchmark by arewefastyet. For each commit on the Pull Request
          there will be two benchmarks: one on the Pull Request's HEAD and another on the base of the Pull Request.
        </span>
        <br></br>
        <span className="header--text">
          On this page you can find all benchmarked Pull Requests.
        </span>
      </div>
      <figure className='line'></figure>
      {PRListError ? (
        <div className="apiError">{PRListError}</div>
      ) : (
        <>
          {isPRListLoading ? (
            <div className="PrLoading justify--content">
              <RingLoader
                loading={isPRListLoading}
                color="#E77002"
                size={300}
              />
            </div>
          ) : (
            <>
              <div className="pr__sidebar">
                <span className="pullnbTitle width--4em">#</span>
                <span className="width--40 hidden--tablet">Title</span>
                <span className="width--20 hidden--tablet">Author</span>
                <span className="width--10em hidden--mobile">Opened At</span>
                <span className="linkSidebar"></span>
                <span className="hidden--desktop">More</span>
              </div>

              <div className="pr__container justify--content">
                {dataPRList.map((PRList, index) => {
                  return (
                    <PRGitInfo
                      key={index}
                      data={PRList}
                    />
                  );
                })}
              </div>
            </>
          )}
        </>
      )}
    </div>
  );
};

export default PR;
