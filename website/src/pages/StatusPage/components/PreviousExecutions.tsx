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

export default function PreviousExecutions(props) {
  const { data, title } = props;

  const [previousExecutions, setPreviousExecutions] = useState([]);

  useEffect(() => {
    for (const entry of data) {
      const newData = {};

      newData["UUID"] = entry.uuid.slice(0, 8);

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

      newData["Started"] = formatDate(entry.started_at) || "N/A";

      newData["Finished"] = formatDate(entry.finished_at) || "N/A";

      newData["Type"] = entry.type_of;

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

      newData["Go version"] = entry.golang_version;

      newData["Status"] = (
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
      );

      setPreviousExecutions((p) => [...p, newData]);
    }
  }, []);

  return (
    <section className="p-page mt-20 flex flex-col">
      <h1 className="text-primary text-3xl my-5 text-center">
        {title}
      </h1>
      {previousExecutions.length > 0 && (
        <DisplayList data={previousExecutions} />
      )}
    </section>
  );
}
