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
import { Button } from "@/components/ui/button";
import { VitessRefs } from "@/types";
import useApiCall from "@/utils/Hook";
import { useEffect, useState } from "react";
import CommandComponent from "./CompareCommand";

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
  const [oldGitRef, setOldGitRef] = useState(gitRef.old);
  const [newGitRef, setNewGitRef] = useState(gitRef.new);
  const isButtonDisabled = !oldGitRef || !newGitRef;

  const { data: vitessRefs } = useApiCall<VitessRefs>(
    `${import.meta.env.VITE_API_URL}vitess/refs`
  );

  const compareClicked = () => {
    setGitRef({ old: oldGitRef, new: newGitRef });
  };

  useEffect(() => {
    setOldGitRef(gitRef.old);
    setNewGitRef(gitRef.new);
  }, [gitRef]);

  return (
    <Hero title={heroProps.title}>
      <div className="flex flex-col md:flex-row gap-4">
        <div className="flex flex-col">
          <label className="text-primary mb-2">Old</label>
          <CommandComponent
            inputLabel={"Search  commit or releases..."}
            gitRef={oldGitRef}
            setGitRef={setOldGitRef}
            vitessRefs={vitessRefs}
            keyboardShortcut="o"
          />
        </div>
        <div className="flex flex-col">
          <label className="text-primary mb-2">New</label>
          <CommandComponent
            inputLabel={"Search  commit or releases..."}
            gitRef={newGitRef}
            setGitRef={setNewGitRef}
            vitessRefs={vitessRefs}
            keyboardShortcut="j"
          />
        </div>
        <div className="flex md:items-end items-center justify-center mt-4 md:mt-0">
          <Button
            onClick={compareClicked}
            disabled={isButtonDisabled}
            className="w-fit md:w-auto"
          >
            Compare
          </Button>
        </div>
      </div>
    </Hero>
  );
}
