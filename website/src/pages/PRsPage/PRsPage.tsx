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

import Hero from "./components/Hero";

import RingLoader from "react-spinners/RingLoader";

import PRTable from "./components/PRTable";
import useApiCall from "../../hooks/useApiCall";

export default function PRsPage() {
  const [data, loading, error] = useApiCall("/pr/list");

  return (
    <>
      <Hero />

      {loading && (
        <div className="flex justify-center w-full my-16">
          <RingLoader loading={loading} color="#E77002" size={300} />
        </div>
      )}

      {error ? (
        <div className="text-red-500 text-center my-2">{error}</div>
      ) : null}

      {!loading && data && <PRTable data={data} />}
    </>
  );
}
