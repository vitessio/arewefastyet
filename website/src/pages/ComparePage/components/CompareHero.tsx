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

import Hero, { HeroProps } from "@/common/Hero";
import CommandComponent from "@/common/VitessRefsCommand";
import { Button } from "@/components/ui/button";
import { VitessRefs } from "@/types";
import useApiCall from "@/utils/Hook";
import { useState } from "react";

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
};

export default function CompareHero(props: CompareHeroProps) {
  const { gitRef, setGitRef } = props;
  const [oldGitRef, setOldGitRef] = useState("");
  const [newGitRef, setNewGitRef] = useState("");
  const isButtonDisabled = !oldGitRef || !newGitRef;
  const [oldListVisible, setOldListVisible] = useState(false);
  const [newListVisible, setNewListVisible] = useState(false);

  const {
    data: vitessRefs,
  } = useApiCall<VitessRefs>(`${import.meta.env.VITE_API_URL}vitess/refs`);

  const compareClicked = () => {
    setGitRef({ old: oldGitRef, new: newGitRef });
  };

  return (
    <Hero title={heroProps.title}>
      <div className="flex flex-row gap-4">
        <CommandComponent
          inputLabel="Search commits or releases..."
          setGitRef={setOldGitRef}
          vitessRefs={vitessRefs}
        />
        <CommandComponent
          inputLabel="Search commits or releases..."
          setGitRef={setNewGitRef}
          vitessRefs={vitessRefs}
        />
        <Button onClick={compareClicked} disabled={isButtonDisabled}>
          Compare
        </Button>
      </div>
    </Hero>
  );
}
