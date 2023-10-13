import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import DisplayList from "../../../common/DisplayList";

export default function ExecutionQueue(props) {
  const { data } = props;

  const [executionQueue, setExecutionQueue] = useState([]);

  useEffect(() => {
    for (const entry of data) {
      const newData = {};

      newData["SHA"] = (
        <Link
          target="__blank"
          rel="noopener noreferrer"
          className="text-primary text"
          to={`https://github.com/vitessio/vitess/commit/${entry.git_ref}`}
        >
          {entry.git_ref.slice(0, 6)}
        </Link>
      );

      newData["Source"] = entry.source;

      if (entry.type_of) newData["Type"] = entry.type_of;

      if (entry.pull_nb) {
        newData["PR"] = (
          <Link
            target="__blank"
            rel="noopener noreferrer"
            className="text-primary"
            to={`https://github.com/vitessio/vitess/pull/${entry.pull_nb}`}
          >
            {entry.pull_nb}
          </Link>
        );
      } else {
        newData["PR"] = <span></span>;
      }

      setExecutionQueue((p) => [...p, newData]);
    }
  }, []);

  return (
    <section className="p-page mt-20">
      <h1 className="text-primary text-3xl my-5 text-center">
        Execution Queue
      </h1>
      {executionQueue.length > 0 && <DisplayList data={executionQueue} />}
    </section>
  );
}
