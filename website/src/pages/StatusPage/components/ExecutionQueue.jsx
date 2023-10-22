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
    <section className="p-page mt-20 flex flex-col">
      <h1 className="text-primary text-3xl my-5 text-center">
        Execution Queue
      </h1>
      {executionQueue.length > 0 && <DisplayList data={executionQueue} />}
    </section>
  );
}
