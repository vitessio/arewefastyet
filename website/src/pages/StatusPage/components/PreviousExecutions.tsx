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
import { formatDate } from "../../../utils/Utils";
import { twMerge } from "tailwind-merge";
import DisplayList from "../../../common/DisplayList";
import { statusDataTypes } from "@/types";

interface TransformedData {
  [key: string]: string | React.ReactNode | null;
}

interface PreviousExecutionsProps {
  data: statusDataTypes[];
  title: string;
}

/**
 * The PreviousExecutions component displays a list of previous executions.
 * @param {PreviousExecutionsProps} props - The props for the PreviousExecutions component.
 * @returns {JSX.Element} - The rendered JSX element.
 */

export default function PreviousExecutions({
  data,
  title,
}: PreviousExecutionsProps) {
  const [previousExecutions, setPreviousExecutions] = useState<
    TransformedData[]
  >([]);

  useEffect(() => {
    const transformedData = data.map((entry) => {
      const newData: TransformedData = {
        UUID: entry.uuid.slice(0, 8),
        SHA: (
          <Link
            target="__blank"
            rel="noopener noreferrer"
            className="text-primary text"
            to={`https://github.com/vitessio/vitess/commit/${entry.git_ref}`}
          >
            {entry.git_ref.slice(0, 6)}
          </Link>
        ),
        Source: entry.source,
        Started: formatDate(entry.started_at) || "N/A",
        Finished: formatDate(entry.finished_at) || "N/A",
        Type: entry.type_of,
        PR: entry.pull_nb ? (
          <Link
            target="__blank"
            rel="noopener noreferrer"
            className="text-primary"
            to={`https://github.com/vitessio/vitess/pull/${entry.pull_nb}`}
          >
            {entry.pull_nb}
          </Link>
        ) : (
          <span></span>
        ),
        "Go version": entry.golang_version,
        Status: (
          <span
            className={twMerge(
              "text-lg text-white px-4 rounded-full",
              entry.status === "failed" && "bg-[#dd1a2a]",
              entry.status === "finished" && "bg-[#00aa00]",
              entry.status === "started" && "bg-[#3a3aed]"
            )}
          >
            {entry.status}
          </span>
        ),
      };

      return newData;
    });

    setPreviousExecutions(transformedData);
  }, [data]);

  return (
    <section className="p-page mt-20 flex flex-col overflow-y-scroll">
      <h1 className="text-primary text-3xl my-5 text-center">{title}</h1>
      {previousExecutions.length > 0 && (
        <DisplayList data={previousExecutions} />
      )}
    </section>
  );
}
