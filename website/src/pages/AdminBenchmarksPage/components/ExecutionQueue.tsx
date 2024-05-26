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
import { BenchmarkExecution, BenchmarkStatus } from "../../../types";
import admin from "../../../utils/admin";
import { twMerge } from "tailwind-merge";

interface ExecutionQueueProps {
  data: BenchmarkExecution<BenchmarkStatus.Ongoing>[];
}

export default function ExecutionQueue(props: ExecutionQueueProps) {
  const [executionQueue, setExecutionQueue] = useState<DataType[]>([]);
  const [deleteLoading, setDeleteLoading] = useState(false);

  function deleteHandler(uuid: string, sha: string) {
    setDeleteLoading(false);

    admin.deleteRun(uuid, sha);

    setDeleteLoading(true);
  }

  useEffect(() => {
    for (const entry of props.data) {
      const newData: Partial<DataType> = {};

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

      if (entry.type_of) newData["Workload"] = entry.type_of;

      newData[""] = (
        <button
          disabled={deleteLoading}
          onClick={() => {
            deleteHandler(entry.git_ref, entry.git_ref);
          }}
          className={twMerge(
            "px-5 py-1 rounded bg-red-600 text-white duration-300 hover:scale-105 hover:shadow hover:-translate-y-[2px]",
            deleteLoading && "opacity-50 cursor-not-allowed"
          )}
        >
          Cancel
        </button>
      );

      setExecutionQueue((p) => [...p, newData as DataType]);
    }
  }, []);

  return (
    <section className="p-page my-16">
      <div className="flex flex-col">
        <h2 className="text-primary text-3xl">Current Execution Queue</h2>

        <article className="flex flex-col my-8 rounded-2xl overflow-hidden bg-foreground bg-opacity-10">
          {executionQueue.length > 0 && <DisplayList data={executionQueue} />}
        </article>
      </div>
    </section>
  );
}

interface DataType {
  SHA: React.ReactNode;
  Source: string;
  Workload: string;
  PR: React.ReactNode;
  "": React.ReactNode;
}
