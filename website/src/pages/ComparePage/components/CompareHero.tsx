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

import CompareActions from "@/common/CompareActions";
import Hero, { HeroProps } from "@/common/Hero";
import { VitessRefs } from "@/types";
import useApiCall from "@/utils/Hook";
import { useEffect, useState } from "react";

const heroProps: HeroProps = {
  title: "Compare Versions",
};

export type CompareHeroProps = {
  gitRef: { old: string; new: string };
  setGitRef: React.Dispatch<
    React.SetStateAction<{
      old: string;
      new: string;
    }>
  >;
  vitessRefs: VitessRefs | null;
};

export default function CompareHero(props: CompareHeroProps) {
  const { gitRef, setGitRef, vitessRefs } = props;
  const [oldGitRef, setOldGitRef] = useState(gitRef.old);
  const [newGitRef, setNewGitRef] = useState(gitRef.new);

  const compareClicked = () => {
    setGitRef({ old: oldGitRef, new: newGitRef });
  };

  useEffect(() => {
    setOldGitRef(gitRef.old);
    setNewGitRef(gitRef.new);
  }, [gitRef]);

  return (
    <Hero title={heroProps.title}>
      <CompareActions
        compareClicked={compareClicked}
        newGitRef={newGitRef}
        oldGitRef={oldGitRef}
        setNewGitRef={setNewGitRef}
        vitessRefs={vitessRefs}
        setOldGitRef={setOldGitRef}
      />
    </Hero>
  );
}
