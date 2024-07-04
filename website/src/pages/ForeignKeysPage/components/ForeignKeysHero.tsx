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
import Dropdown from "@/common/Dropdown";
import Hero, { HeroProps } from "@/common/Hero";
import { twMerge } from "tailwind-merge";

const heroProps: HeroProps = {
  title: "Foreign Keys",
  description: (
    <p className="my-6 leading-loose">
      Support for Foreign Keys have been added to Vitess in v18.0.0. We want
      to be able to compare the performance of Vitess with and without
      Foreign Keys.
      <br />
      For this purpose, we propose four benchmarks:
      <br />
      - TPCC, with a sharded keyspace
      <br />
      - TPCC_UNSHARDED, with an unsharded keyspace
      <br />
      - TPCC_FK, with Foreign Keys enabled and set to vitess managed
      <br />
      - TPCC_FK_UNMANAGED, with Foreign Keys enabled and set to vitess
      unmanaged
      <br />
      Use the dropdown on the right to select which version of Vitess you
      would like to use to compare the performance of our four TPCC
      benchmarks.
    </p>
  ),
};

export default function ForeignKeysHero(props: {
  refs: any;
  gitRef: any;
  setGitRef: any;
}) {
  const { refs, gitRef, setGitRef } = props;

  return (
    <Hero title={heroProps.title} description={heroProps.description}>
      <div className="flex-1 flex flex-col items-center">
        <div className="flex flex-col gap-y-8">
          {refs && refs.length > 0 && (
            <div className="flex gap-x-24">
              <Dropdown.Container
                className="w-[30vw] md:w-[20vw] py-2 border border-primary rounded-md mb-[1px] text-lg shadow-xl"
                defaultIndex={refs
                  .map((r: { Name: any }) => r.Name)
                  .indexOf(gitRef.tag)}
                onChange={(event: { value: any }) => {
                  setGitRef((p: any) => {
                    return { ...p, tag: event.value };
                  });
                }}
              >
                {refs.map((ref: { Name: any }, key: number) => (
                  <Dropdown.Option
                    key={key}
                    className={twMerge(
                      "w-[20vw] relative border-front border border-t-transparent border-opacity-60 bg-background py-2 font-medium hover:bg-accent",
                      key === 0 && "rounded-t border-t-front",
                      key === refs.length - 1 && "rounded-b"
                    )}
                  >
                    {ref.Name}
                  </Dropdown.Option>
                ))}
              </Dropdown.Container>
            </div>
          )}
        </div>
      </div>
    </Hero>
  );
}
