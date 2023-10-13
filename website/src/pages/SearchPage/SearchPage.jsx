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
import { useNavigate } from "react-router-dom";
import RingLoader from "react-spinners/RingLoader";
import useApiCall from "../../utils/Hook";

import Hero from "./components/Hero";
import SearchMacroDesktop from "./components/SearchMacroDesktop";

export default function SearchPage() {
  const urlParams = new URLSearchParams(window.location.search);
  const [gitRef, setGitRef] = useState(urlParams.get("git_ref") || "");

  const {
    data: dataSearch,
    isLoading: isSearchLoading,
    error: searchError,
  } = useApiCall(`${import.meta.env.VITE_API_URL}search?git_ref=${gitRef}`);

  const navigate = useNavigate();

  useEffect(() => {
    navigate(`?git_ref=${gitRef}`);
  }, [gitRef]);

  return (
    <>
      <Hero setGitRef={setGitRef} />

      {searchError && <div className="apiError">{searchError}</div>}

      {isSearchLoading && (
        <div className="flex my-10 justify-center items-center">
          <RingLoader loading={isSearchLoading} color="#E77002" size={300} />
        </div>
      )}

      {!isSearchLoading && dataSearch && (
        <section className="flex flex-col items-center p-page">
          <div className="w-1/2 flex flex-col gap-y-16">
            {dataSearch.Macros &&
              typeof dataSearch.Macros === "object" &&
              Object.entries(dataSearch.Macros).map(function (
                searchMacro,
                index
              ) {
                return (
                  <SearchMacroDesktop
                    key={index}
                    data={searchMacro}
                    gitRef={gitRef}
                  />
                );
              })}
          </div>
        </section>
      )}
    </>
  );
}
