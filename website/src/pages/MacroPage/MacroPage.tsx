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

  const [refs] = useApiCall("/vitess/refs");
  const [macrobench, loading, error] = useApiCall("/macrobench/compare", {
    params: { ltag: commits.left, rtag: commits.right },
  });

  const navigate = useNavigate();
  useEffect(() => {
    navigate(`?ltag=${gitRef.left}&rtag=${gitRef.right}`);

    if (refs) {
      const commits = {
        left: refs.filter((r) => r.Name === gitRef.left)[0].CommitHash,
        right: refs.filter((r) => r.Name === gitRef.right)[0].CommitHash,
      };
      setCommits(commits);
    }
  }, [gitRef.left, gitRef.right, refs]);

  return (
    <>
      <Hero refs={refs} gitRef={gitRef} setGitRef={setGitRef} />

      <div className="p-page">
        <div className="border border-front" />
      </div>

      <div className="flex flex-col items-center py-20">
        {error && (
          <div className="text-sm w-1/2 text-red-500 text-center">{error}</div>
        )}

        {loading && <RingLoader loading={loading} color="#E77002" size={300} />}

        {!loading && macrobench && (
          <div className="flex flex-col gap-y-20 ">
            {gitRef.left &&
              gitRef.right &&
              macrobench.map((macro, index) => {
                return (
                  <div key={index}>
                    <Macrobench
                      data={macro}
                      gitRef={gitRef}
                      commits={commits}
                    />
                  </div>
                );
              })}
          </div>
        )}
      </div>
    </>
  );
};

export default Macro;
