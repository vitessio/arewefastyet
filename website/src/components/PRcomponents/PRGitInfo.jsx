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

import React from "react";
import useApiCall from "../../utils/Hook";

import "../PRcomponents/PRGitInfo.css";

const PRGitInfo = ({ pull_nb, setPrNumber, className }) => {
  // const {
  //     data: dataPRGit,
  //     isLoading: isPRGitLoading,
  //     error: PRGitError,
  //   } = useApiCall(`https://api.github.com/repos/vitessio/vitess/pulls/${pull_nb}`, []);
  //    console.log(dataPRGit)

  const dataPRGit = {
    title: "sdgkusdgisudhfvkjhdbfvkjbdfvkjhdfvhb",
    author: "fouifoui",
    created_at: "51651-65814-6145",
  };
  const handlePrInfo = (e) => {
    const number = e.toString();
    setPrNumber(number);
  };

  return (
    <div className={`prGit flex ${className}`}>
      <span className="width--5em">{pull_nb}</span>
      <span className="width--15em">{dataPRGit.title}</span>
      <span className="width--6em">{dataPRGit.author}</span>
      <span className="width--10em">{dataPRGit.created_at}</span>
      <span
        className="linkToCompare width--11em"
        onClick={() => handlePrInfo(pull_nb)}
      >
        Click to compare with main
      </span>
    </div>
  );
};

export default PRGitInfo;
