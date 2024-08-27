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
import VitessRefsCommand from "@/common/VitessRefsCommand";
import WorkloadsCommand from "@/common/WorkloadsCommand";
import { Button } from "@/components/ui/button";
import useApiCall from "@/hooks/useApiCall";
import { VitessRefs } from "@/types";
import { useEffect, useState } from "react";

const heroProps: HeroProps = {
  title: "Foreign Keys",
  description: (
    <>
      Foreign Keys Management is a new feature that was added in version 18 of
      Vitess. This page enables the comparison of two different Foreign Keys
      workloads on the same commit, allowing you to see the performance of a
      workflow with Vitess-managed foreign keys versus one where foreign keys
      are either not enabled, or not managed by Vitess.
    </>
  ),
};

export default function ForeignKeysHero(props: {
  gitRef: string;
  setGitRef: React.Dispatch<React.SetStateAction<string>>;
  workload: { old: string; new: string };
  setWorkload: React.Dispatch<
    React.SetStateAction<{ old: string; new: string }>
  >;
  vitessRefs: VitessRefs | undefined;
}) {
  const { gitRef, setGitRef, workload, setWorkload, vitessRefs } = props;
  const [oldWorkload, setOldWorkload] = useState(workload.old);
  const [newWorkload, setNewWorkload] = useState(workload.new);
  const [commit, setCommit] = useState(gitRef);
  let { data: workloads } = useApiCall<string[]>({
    url: `${import.meta.env.VITE_API_URL}workloads`,
    queryKey: ["workloads"],
  });

  workloads =
    workloads !== undefined
      ? (workloads.filter((workload) => workload.includes("TPCC")) ?? [])
      : [];

  const isButtonDisabled = !oldWorkload || !newWorkload || !commit;

  const compareClicked = () => {
    setGitRef(commit);
    setWorkload({ old: oldWorkload, new: newWorkload });
  };

  useEffect(() => {
    setOldWorkload(workload.old);
    setNewWorkload(workload.new);
  }, [workload]);

  useEffect(() => {
    setCommit(gitRef);
  }, [gitRef]);

  return (
    <Hero title={heroProps.title} description={heroProps.description}>
      <div className="flex flex-col md:flex-row gap-4">
        <div className="flex flex-col">
          <label className="text-primary mb-2">Old Workload</label>
          <WorkloadsCommand
            inputLabel="Search workloads..."
            workload={oldWorkload}
            workloads={workloads}
            keyboardShortcut="o"
            setWorkload={setOldWorkload}
          />
        </div>
        <div className="flex flex-col">
          <label className="text-primary mb-2">New Workload</label>
          <WorkloadsCommand
            inputLabel="Search workloads..."
            workload={newWorkload}
            workloads={workloads}
            keyboardShortcut="j"
            setWorkload={setNewWorkload}
          />
        </div>
        <div className="flex flex-col">
          <label className="text-primary mb-2">Commit</label>
          <VitessRefsCommand
            inputLabel={"Search commit or releases..."}
            gitRef={commit}
            setGitRef={setCommit}
            vitessRefs={vitessRefs}
            keyboardShortcut="m"
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
