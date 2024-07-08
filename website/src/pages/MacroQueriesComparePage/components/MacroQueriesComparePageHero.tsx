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

import { Link } from "react-router-dom";
import Hero, { HeroProps } from "@/common/Hero";

export default function MacroQueriesComparePageHero(props: {
  commits: { left: any; right: any };
}) {
  const { left, right } = props.commits;
  const heroProps: HeroProps = {
    title: "Compare Query Plans",
    description: (
      <>
        <b className="text-primary text-xl">Old</b> benchmarked commit{" "}
        <Link
          className="text-primary text-xl"
          target="_blank"
          to={`https://github.com/vitessio/vitess/commit/${left}`}
        >
          {left.slice(0, 8)}
        </Link>{" "}
        <br />
        <br />
        <b className="text-primary text-xl">New</b> benchmarked commit{" "}
        <Link
          className="text-primary text-xl"
          target="_blank"
          to={`https://github.com/vitessio/vitess/commit/${right}`}
        >
          {right.slice(0, 8)}
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
