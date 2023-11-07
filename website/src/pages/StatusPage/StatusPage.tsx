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

import RingLoader from "react-spinners/RingLoader";
import useApiCall from "../../hooks/useApiCall";

import Hero from "./components/Hero";
import ExecutionQueue from "./components/ExecutionQueue";
import PreviousExecutions from "./components/PreviousExecutions";

export default function StatusPage() {
  const [queue, queueLoading, queueError] = useApiCall("/queue");
  const [executions, executionsLoading, executionsError] =
    useApiCall("/recent");

  return (
    <>
      <Hero />

      <figure className="p-page w-full">
        <div className="border-front border" />
      </figure>

      {!queueLoading && queue?.length && <ExecutionQueue data={queue} />}

      {!executionsLoading && executions?.length && (
        <PreviousExecutions data={executions} />
      )}

      {/* SHOW LOADER BENEATH IF EITHER IS LOADING */}
      {(executionsLoading || queueLoading) && (
        <div className="flex justify-center w-full my-16">
          <RingLoader
            loading={executionsLoading || queueLoading}
            color="#E77002"
            size={300}
          />
        </div>
      )}

      {(queueError || executionsError) && (
        <div className="my-10 text-center text-red-500">
          {JSON.stringify([queueError, executionsError], null, 2)}
        </div>
      )}
    </>
  );
}
