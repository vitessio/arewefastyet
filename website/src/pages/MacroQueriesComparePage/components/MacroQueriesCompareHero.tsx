/*
Copyright 2024 The Vitess Authors.

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

import Hero, { HeroProps } from "@/common/Hero";
import { formatGitRef } from "@/utils/Utils";
import { Link } from "react-router-dom";

export type MacroQueriesCompareHeroProps = {
  commits: { oldGitRef: string; newGitRef: string };
};

export default function MacroQueriesCompareHero(
  props: MacroQueriesCompareHeroProps,
) {
  const { oldGitRef, newGitRef } = props.commits;
  const heroProps: HeroProps = {
    title: "Compare Query Plans",
    description: (
      <>
        <b className="text-primary text-xl">Old</b> benchmarked commit{" "}
        <Link
          className="text-primary text-xl"
          target="_blank"
          to={`https://github.com/vitessio/vitess/commit/${oldGitRef}`}
        >
          {formatGitRef(oldGitRef)}
        </Link>{" "}
        <br />
        <br />
        <b className="text-primary text-xl">New</b> benchmarked commit{" "}
        <Link
          className="text-primary text-xl"
          target="_blank"
          to={`https://github.com/vitessio/vitess/commit/${newGitRef}`}
        >
          {formatGitRef(newGitRef)}
        </Link>{" "}
        <br />
        <br />
        Queries are ordered from the worst regression in execution time to the
        best. All executed queries are shown below.
      </>
    ),
  };

  return <Hero title={heroProps.title} description={heroProps.description} />;
}
