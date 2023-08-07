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

import React, { useEffect, useState } from "react";
import { useNavigate, Link } from "react-router-dom";

import "../MacroQueriesCompare/macroQueriesCompare.css";

import { errorApi } from "../../utils/Utils";
import QueryPlan from "../../components/QueryPlan/QueryPlan";

const MacroQueriesCompare = () => {
  const [error, setError] = useState(null);

  const urlParams = new URLSearchParams(window.location.search);
  const [commitLeft, setcommitLeft] = useState(
    urlParams.get("ltag") == null ? "" : urlParams.get("ltag")
  );
  const [commitRight, setcommitRight] = useState(
    urlParams.get("rtag") == null ? "" : urlParams.get("rtag")
  );
  const [type, settype] = useState(
    urlParams.get("type") == null ? "" : urlParams.get("type")
  );
  const navigate = useNavigate();

  useEffect(() => {
    navigate(`?ltag=${commitLeft}&rtag=${commitRight}&type=${type}`);
  }, []);

  const [dataQueryPlan, setDataQueryPlan] = useState([]);

  useEffect(() => {
    if (!commitLeft || !commitRight || !type) {
      setError(
        "Error: Some URL parameters are missing. Please provide both 'ltag', 'rtag' and type parameters."
      );
    } else {
      const fetchData = async () => {
        try {
          const responseQueryPlan = await fetch(
            `${
              import.meta.env.VITE_API_URL
            }macrobench/compare/queries?ltag=${commitLeft}&rtag=${commitRight}&type=${type}`
          );

          const jsonDataQueryPlan = await responseQueryPlan.json();

          setDataQueryPlan(jsonDataQueryPlan);
        } catch (error) {
          console.log("Error while retrieving data from the API", error);
          setError(errorApi);
        }
      };
      fetchData();
    }
  }, []);

  const [openPlanIndex, setOpenPlanIndex] = useState(0); // Index of the currently open plan

  const togglePlan = (index) => {
    setOpenPlanIndex((prevIndex) => (prevIndex === index ? -1 : index)); // If the plan is already open, close it; otherwise, open it
  };

  return (
    <div className="query">
      {error ? (
        <div className="query__error flex--column">
          <span>{error}</span>
          <Link to="/macro" className="macro-link">
            Back to Macro
          </Link>
        </div>
      ) : (
        <>
          <div className="query__top flex--column">
            <h2>Compare Query Plans</h2>
            <span>
              Comparing the query plans of two OLTP benchmarks: A and B.
            </span>
            <span>
              <b>A</b> benchmarked commit{" "}
              <a
                href={`https://github.com/vitessio/vitess/commit/${commitLeft}`}
              >
                {commitLeft.slice(0, 8)}
              </a>{" "}
              using the Gen4 query planner.
            </span>
            <span>
              <b>B</b> benchmarked commit{" "}
              <a
                href={`https://github.com/vitessio/vitess/commit/${commitRight}`}
              >
                {commitRight.slice(0, 8)}
              </a>{" "}
              using the Gen4 query planner.
            </span>
            <span>
              Queries are ordered from the worst regression in execution time to
              the best. All executed queries are shown below.
            </span>
          </div>
          <figure className="line"></figure>
          <div className="queryPlan__container">
            {dataQueryPlan.map((queryPlan, index) => {
              return (
                <QueryPlan
                  key={index}
                  data={queryPlan}
                  isOpen={index === openPlanIndex}
                  togglePlan={() => togglePlan(index)}
                />
              );
            })}
          </div>
        </>
      )}
    </div>
  );
};

export default MacroQueriesCompare;
