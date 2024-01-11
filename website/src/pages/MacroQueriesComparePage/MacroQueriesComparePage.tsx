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

import { useEffect, useState } from "react";
import { useNavigate, Link } from "react-router-dom";

import QueryPlan from "./components/QueryPlan";
import Hero from "./components/Hero";
import RingLoader from "react-spinners/RingLoader";
import useApiCall from "../../hooks/useApiCall";

export default function MacroQueriesComparePage() {

  const urlParams = new URLSearchParams(window.location.search);
  const commits = {
    left: urlParams.get("ltag") || "",
    right: urlParams.get("rtag") || "",
  };

  const type = urlParams.get("type") || "";
  const navigate = useNavigate();

  useEffect(() => {
    navigate(`?ltag=${commits.left}&rtag=${commits.right}&type=${type}`);
  }, []);

  const [data, loading, error] = useApiCall("/macrobench/compare/queries", {
    params: { ltag: commits.left, rtag: commits.right, type: type },
  });

  useEffect(() => {
    if (
      !commits.left.length ||
      !commits.right.length ||
      !type ||
      !type.length ||
      commits.left === "null" ||
      commits.right === "null" ||
      type === "null"
    ) {
      throw new Error(
        "Error: Invalid left or right comparison parameters, please check the SHAs provided"
      );
    }
  }, []);

  return (
    <>
      {error && (
        <div className="flex flex-col h-screen fixed top-0 left-0 w-full justify-center items-center bg-background z-10">
          <span className="my-2 text-sm text-red-600">
            An error occured <br /> <br />
            {error}
          </span>
          <Link
            to="/macro"
            className="text-primary my-5 text-lg underline underline-offset-2 hover:no-underline"
          >
            Back to Macro
          </Link>
        </div>
      )}

      {!error && (
        <>
          <Hero commits={commits} />

          <div className="p-page my-8">
            <div className="border border-front" />
          </div>

          <section className="p-page flex flex-col gap-y-8 my-5">
            {loading && (
              <div className="flex items-center justify-center my-10">
                <RingLoader loading={loading} color="#E77002" size={300} />
              </div>
            )}

            {!loading &&
              data &&
              data.length &&
              data.map((queryPlan, index) => {
                return <QueryPlan key={index} data={queryPlan} />;
              })}
          </section>
        </>
      )}
    </>
  );
}
