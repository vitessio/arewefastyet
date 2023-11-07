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

import React from "react";
import { Link, useParams } from "react-router-dom";
import useApiCall from "../../utils/Hook";
import RingLoader from "react-spinners/RingLoader";

import { formatDate, errorApi } from "../../utils";

export default function PRPage() {
  const { pull_nb } = useParams();

  const {
    data: dataSinglePr,
    isLoading: singlePrLoading,
    error: singlePrError,
  } = useApiCall(`${import.meta.env.VITE_API_URL}pr/info/${pull_nb}`, []);

  const isComparisonAvailable =
    dataSinglePr.Base !== "" && dataSinglePr.Head !== "";

  if (
    dataSinglePr.error ==
    "GET https://api.github.com/repos/vitessio/vitess/pulls/13675: 404 Not Found []"
  ) {
    return (
      <div className="errorPullRequest">Pull request {pull_nb} not found</div>
    );
  }

  return (
    <>
      <section className="flex h-screen flex-col items-center py-[20vh] p-page">
        {singlePrLoading && (
          <div className="loadingSpinner">
            <RingLoader loading={singlePrLoading} color="#E77002" size={300} />
          </div>
        )}

        {singlePrError && <div className="text-red-500 text-center my-2">{errorApi}</div>}

        {!singlePrLoading && dataSinglePr && (
          <div className="flex flex-col border border-front rounded-3xl w-11/12 bg-foreground bg-opacity-5">
            <div className="flex justify-between p-5 border-b border-front">
              <div className="flex flex-col justify-evenly">
                <h2 className="text-2xl font-semibold">
                  <Link
                    className="text-primary"
                    target="_blank"
                    to={`https://github.com/vitessio/vitess/pull/${pull_nb}`}
                  >
                    [#{pull_nb}]
                  </Link>{" "}
                  {dataSinglePr.Title}
                </h2>
                <span>
                  By {dataSinglePr.Author} at{" "}
                  {formatDate(dataSinglePr.CreatedAt)}{" "}
                </span>
              </div>

              {isComparisonAvailable && (
                <div className="flex justify-center items-center">
                  <Link
                    className="text-primary p-6 border border-primary rounded-xl duration-300 hover:bg-primary hover:bg-opacity-20 hover:scale-105 whitespace-nowrap"
                    to={`/compare?ltag=${dataSinglePr.Base}&rtag=${dataSinglePr.Head}`}
                  >
                    Compare with base commit
                  </Link>
                </div>
              )}
            </div>
            <div className="flex flex-col justify-between p-5 text-lg leading-loose">
              {isComparisonAvailable ? (
                <>
                  <span>Base: {dataSinglePr.Base}</span>
                  <span>Head: {dataSinglePr.Head}</span>
                </>
              ) : (
                <div>
                  The Base and Head commit information is not available for this
                  pull request.
                </div>
              )}
              <span>{dataSinglePr.error}</span>
            </div>
          </div>
        )}
      </section>
    </>
  );
}
