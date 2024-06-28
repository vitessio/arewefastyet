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
import { statusDataTypes } from "@/types";

interface DisplayList {
  [key: string]: string | React.ReactNode;
}

interface PreviousExecutionsProps {
  data: statusDataTypes[];
  title: string;
}

/**
 * This component renders the execution queue with data provided.
 *
 * @param {object} props - The props for the ExecutionQueue component.
 * @param {statusDataTypes[]} props.data - An array of status data types representing the execution queue.
 * @param {string} props.title - The title of the execution queue section.
 * @returns {JSX.Element} The rendered ExecutionQueue component.
 */

export default function ExecutionQueue({
  data,
  title,
}: PreviousExecutionsProps): JSX.Element {
  
  const [executionQueue, setExecutionQueue] = useState<DisplayList[]>([]);

  useEffect(() => {
    const transformedData = data.map((entry) => {
      const newData: DisplayList = {};

      newData["SHA"] = (
        <Link
          target="__blank"
          rel="noopener noreferrer"
          className="text-primary"
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
      return newData;
    });

    setExecutionQueue(transformedData);
  }, [data]);

  return (
    <section className="p-page mt-20 flex flex-col">
      <h1 className="text-primary text-3xl my-5 text-center">Execution Queue</h1>
      {executionQueue.length > 0 && <DisplayList data={executionQueue} />}
    </section>
  );
}
