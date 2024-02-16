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
import Dropdown from "../../../common/Dropdown";
import { twMerge } from "tailwind-merge";

export default function Hero(props: {
  refs: any;
  gitRef: any;
  setGitRef: any;
}) {
  const { refs, gitRef, setGitRef } = props;

  return (
    <section className="flex flex-col gap-y-[10vh] pt-[5vh] justify-center items-center h-[30vh]">
      {refs && refs.length > 0 && (
        <div className="flex gap-x-24">
          <Dropdown.Container
            className="w-[20vw] py-2 border border-primary rounded-md mb-[1px] text-lg shadow-xl"
            defaultIndex={refs
              .map((r: { Name: any }) => r.Name)
              .indexOf(gitRef.left)}
            onChange={(event: { value: any }) => {
              setGitRef((p: any) => {
                return { ...p, left: event.value };
              });
            }}
          >
            {refs.map((ref: { Name: any }, key: number) => (
              <Dropdown.Option
                key={key}
                className={twMerge(
                  "w-[20vw] relative border-front border border-t-transparent border-opacity-60 bg-background py-2 after:duration-150 after:absolute-cover after:bg-foreground after:bg-opacity-0 hover:after:bg-opacity-10 font-medium",
                  key === 0 && "rounded-t border-t-front",
                  key === refs.length - 1 && "rounded-b"
                )}
              >
                {ref.Name}
              </Dropdown.Option>
            ))}
          </Dropdown.Container>

          <Dropdown.Container
            className="w-[20vw] py-2 border border-primary rounded-md mb-[1px] text-lg shadow-xl"
            defaultIndex={refs
              .map((r: { Name: any }) => r.Name)
              .indexOf(gitRef.right)}
            onChange={(event: { value: any }) => {
              setGitRef((p: any) => {
                return { ...p, right: event.value };
              });
            }}
          >
            {refs.map((ref: { Name: any }, key: number) => (
              <Dropdown.Option
                key={key}
                className={twMerge(
                  "w-[20vw] relative border-front border border-t-transparent border-opacity-60 bg-background py-2 after:duration-150 after:absolute-cover after:bg-foreground after:bg-opacity-0 hover:after:bg-opacity-10 font-medium",
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
    </section>
  );
}
