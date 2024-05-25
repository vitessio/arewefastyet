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

import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import RingLoader from "react-spinners/RingLoader";
import useApiCall from "../../hooks/useApiCall";

import Hero from "./components/Hero";
import SearchMacro from "./components/SearchMacro";

export default function SearchPage() {
  const urlParams = new URLSearchParams(window.location.search);
  const [gitRef, setGitRef] = useState(urlParams.get("git_ref") || "");

  const [data, loading, error] = useApiCall("/search", {
    params: { git_ref: gitRef },
  });

  const navigate = useNavigate();
  useEffect(() => {
    navigate(`?git_ref=${gitRef}`);
  }, [gitRef]);

  return (
    <>
      <Hero gitRef={gitRef} setGitRef={setGitRef} />

      {error && <div className="text-red-500 text-center my-2">{error}</div>}

      {loading && (
        <div className="flex my-10 justify-center items-center">
          <RingLoader loading={loading} color="#E77002" size={300} />
        </div>
      )}

      {!loading && data && (
        <section className="flex flex-col items-center p-page">
          <div className="w-1/2 flex flex-col gap-y-16">
            {Object.entries(data.Macros).map(
              ([macroName, macroData], index) => (
                <SearchMacro
                  key={index}
                  data={macroData}
                  macroName={macroName}
                  gitRef={gitRef}
                />
              )
            )}
          </div>
        </section>
      )}
    </>
  );
}
