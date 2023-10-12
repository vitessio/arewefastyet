import React, { useEffect, useState } from "react";
import DisplayTable from "./DisplayTable";
import { Link } from "react-router-dom";
import { formatDate } from "../../../utils/Utils";
import { twMerge } from "tailwind-merge";

export default function PreviousExecutions(props) {
  const { data } = props;

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

      newData["Finished"] = formatDate(entry.finished_at) || "In Progress";

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
    <section className="p-page mt-20">
      <h1 className="text-primary text-3xl my-5 text-center">
        Previous Executions
      </h1>
      {previousExecutions.length > 0 && (
        <DisplayTable data={previousExecutions} />
      )}
    </section>
  );
}
