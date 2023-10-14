import React from "react";
import { Link } from "react-router-dom";

export default function Hero(props) {
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
