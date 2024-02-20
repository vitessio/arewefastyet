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
import { Link } from "react-router-dom";

export default function Hero(props: { commits: { left: any; right: any } }) {
  const { left, right } = props.commits;

  return (
    <section className="flex flex-col h-[50vh] pt-[10vh] items-center justify-evenly cursor-default">
      <h2 className="text-5xl font-bold">Compare Query Plans</h2>
      <span className="text-front text-opacity-70">
        Comparing the query plans of two OLTP benchmarks: A and B.
      </span>
      <span>
        <b className="text-primary text-xl">A</b> benchmarked commit{" "}
        <Link
          className="text-primary text-xl"
          target="_blank"
          to={`https://github.com/vitessio/vitess/commit/${left}`}
        >
          {left.slice(0, 8)}
        </Link>{" "}
        using the Gen4 query planner.
      </span>
      <span>
        <b className="text-primary text-xl">B</b> benchmarked commit{" "}
        <Link
          className="text-primary text-xl"
          target="_blank"
          to={`https://github.com/vitessio/vitess/commit/${right}`}
        >
          {right.slice(0, 8)}
        </Link>{" "}
        using the Gen4 query planner.
      </span>
      <span className="text-sm text-front text-opacity-70">
        Queries are ordered from the worst regression in execution time to the
        best. All executed queries are shown below.
      </span>
    </section>
  );
}
