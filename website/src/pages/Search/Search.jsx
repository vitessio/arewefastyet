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

import "../Search/search.css";
import "swiper/css";
import "swiper/css/pagination";

import SearchMicro from "../../components/SearchComponents/SearchMicro/SearchMicro";
import SearchMacroDesktop from "../../components/SearchComponents/SearchMacro/SearchMacroDesktop/SearchMacroDesktop";

const Search = () => {
  const urlParams = new URLSearchParams(window.location.search);
  const [gitRef, setGitRef] = useState(urlParams.get("git_ref") || "");
  const [isFormSubmitted, setIsFormSubmitted] = useState(false);
  const [currentSlideIndex, setCurrentSlideIndex] = useState(
    urlParams.get("ptag") == null ? "0" : urlParams.get("ptag")
  );

  const {
    data: dataSearch,
    isLoading: isSearchLoading,
    error: searchError,
  } = useApiCall(`${import.meta.env.VITE_API_URL}search?git_ref=${gitRef}`, [
    isFormSubmitted,
  ]);

  const navigate = useNavigate();

  useEffect(() => {
    navigate(`?git_ref=${gitRef}&ptag=${currentSlideIndex}`);
  }, [gitRef, currentSlideIndex]);

  const handleInputChange = (e) => {
    setGitRef(e.target.value);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    setIsFormSubmitted((prevState) => !prevState);
  };

  const handleSlideChange = (swiper) => {
    setCurrentSlideIndex(swiper.realIndex);
  };

  const tableHeaders = dataSearch.Macros ? Object.keys(dataSearch.Macros) : [];

  return (
    <div className="search">
      <div className="research">
        <form className="justify--content" onSubmit={handleSubmit}>
          <input
            type="text"
            value={gitRef}
            onChange={handleInputChange}
            placeholder="Search using commit SHA"
            className="research__input"
          />
          <button type="submit">Search</button>
        </form>
      </div>
      {searchError ? (
        <div className="apiError">{searchError}</div>
      ) : (
        <>
          {isSearchLoading ? (
            <div className="loadingSpinner">
              <RingLoader
                loading={isSearchLoading}
                color="#E77002"
                size={300}
              />
            </div>
          ) : (
            <>
              <div className="search__macro justify--content ">
                <div className="searchMacro__desktop">
                  <div className="flex searchMacro__desktop__tbody flex--column">
                    {dataSearch.Macros &&
                      typeof dataSearch.Macros === "object" &&
                      Object.entries(dataSearch.Macros).map(function (
                        searchMacro,
                        index
                      ) {
                        return (
                          <SearchMacroDesktop key={index} data={searchMacro} gitRef={gitRef} />
                        );
                      })}
                  </div>
                </div>
              </div>
              <div className="search__micro">
                {dataSearch.Micro && typeof dataSearch.Micro === "object" && (
                  <>
                    <div className="micro__thead space--between">
                      <span className="width--12em">Package</span>
                      <span className="width--14em">Benchmark Name</span>
                      <span className="width--18em hiddenMobile">
                        Number of Iterations
                      </span>
                      <span className="width--18em hiddenTablet">Time/op</span>
                      <span className="width--6em">More</span>
                    </div>
                    <figure className="micro__thead__line"></figure>
                  </>
                )}
                {dataSearch.Micro &&
                  typeof dataSearch.Micro === "object" &&
                  Object.entries(dataSearch.Micro).map(function (
                    searchMicro,
                    index
                  ) {
                    const isEvenIndex = index % 2 === 0;
                    const backgroundGrey = isEvenIndex
                      ? "grey--background"
                      : "";
                    return (
                      <SearchMicro
                        key={index}
                        data={searchMicro}
                        className={backgroundGrey}
                      />
                    );
                  })}
              </div>
            </>
          )}
        </>
      )}
    </div>
  );
};

export default Search;
