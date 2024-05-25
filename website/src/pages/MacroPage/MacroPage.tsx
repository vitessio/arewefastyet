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

import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import RingLoader from "react-spinners/RingLoader";

import Macrobench from "../../common/Macrobench";

import Hero from "./components/Hero";
import useApiCall from "../../hooks/useApiCall";

const Macro = () => {
  const urlParams = new URLSearchParams(window.location.search);

  const [gitRef, setGitRef] = useState({
    left: urlParams.get("ltag") || "Left",
    right: urlParams.get("rtag") || "Right",
  });
  const [commits, setCommits] = useState({ left: "", right: "" });

  const [refs, refsLoading, refsError] = useApiCall("/vitess/refs");
  const [macrobench, loading, error] = useApiCall("/macrobench/compare", {
    params: { ltag: commits.left, rtag: commits.right },
  });

  const navigate = useNavigate();
  useEffect(() => {
    navigate(`?ltag=${gitRef.left}&rtag=${gitRef.right}`);

    if (refs && refs.length > 0) {
      const leftRef = refs.find((ref) => ref.Name === gitRef.left);
      const rightRef = refs.find((ref) => ref.Name === gitRef.right);
      if (leftRef && rightRef) {
        setCommits({ left: leftRef.CommitHash, right: rightRef.CommitHash });
      }
    }
  }, [gitRef.left, gitRef.right, refs, navigate]);

  return (
    <>
      <Hero refs={refs} gitRef={gitRef} setGitRef={setGitRef} />

      <div className="p-page">
        <div className="border border-front" />
      </div>

      <div className="flex flex-col items-center py-20">
        {refsError && (
          <div className="text-sm w-1/2 text-red-500 text-center">
            {refsError}
          </div>
        )}
        {refsLoading && (
          <RingLoader loading={refsLoading} color="#E77002" size={300} />
        )}

        {error && (
          <div className="text-sm w-1/2 text-red-500 text-center">{error}</div>
        )}
        {loading && <RingLoader loading={loading} color="#E77002" size={300} />}

        {!loading && macrobench && (
          <div className="flex flex-col gap-y-20 ">
            {gitRef.left &&
              gitRef.right &&
              macrobench.map((macro, index) => (
                <div key={index}>
                  <Macrobench data={macro} gitRef={gitRef} commits={commits} />
                </div>
              ))}
          </div>
        )}
      </div>
    </>
  );
};

export default Macro;
