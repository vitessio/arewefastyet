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

import React, { useEffect } from "react";
import Dropdown from "../../../common/Dropdown";
import { twMerge } from "tailwind-merge";

interface Ref {
  Name: string;
  CommitHash: string;
  Version: {
    Major: number;
    Minor: number;
    Patch: number;
  };
  RCnumber: number;
}

interface GitRef {
  old: string;
  new: string;
}

interface HeroProps {
  refs: Ref[];
  gitRef: GitRef;
  setGitRef: React.Dispatch<React.SetStateAction<GitRef>>;
}

/**
 * The Hero component displays git references and allows selecting old and new references using dropdowns.
 * @param {HeroProps} props - The props for the Hero component.
 * @returns {JSX.Element} - The rendered JSX element.
 */

export default function Hero({ refs, gitRef, setGitRef }: HeroProps): JSX.Element {
  useEffect(() => {
    if (refs.length > 0) {
      if (!gitRef.old || !refs.some((ref) => ref.Name === gitRef.old)) {
        setGitRef((prev) => ({ ...prev, old: refs[0].Name }));
      }
      if (!gitRef.new || !refs.some((ref) => ref.Name === gitRef.new)) {
        setGitRef((prev) => ({ ...prev, new: refs[0].Name }));
      }
    }
  }, [refs, gitRef, setGitRef]);

  return (
    <section className="flex flex-col gap-y-[10vh] pt-[5vh] justify-center items-center h-[30vh]">
      {refs && refs.length > 0 && (
        <div className="flex gap-x-24">
          <Dropdown.Container
            className="w-[20vw] py-2 border border-primary rounded-md mb-[1px] text-lg shadow-xl"
            defaultIndex={refs.findIndex((r) => r.Name === gitRef.old)}
            onChange={(event) => {
              setGitRef((p) => {
                return { ...p, old: event.value };
              });
            }}
          >
            {refs.map((ref, key) => (
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

          <Dropdown.Container
            className="w-[20vw] py-2 border border-primary rounded-md mb-[1px] text-lg shadow-xl"
            defaultIndex={refs.findIndex((r) => r.Name === gitRef.new)}
            onChange={(event) => {
              setGitRef((p) => {
                return { ...p, new: event.value };
              });
            }}
          >
            {refs.map((ref, key) => (
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
    </section>
  );
}
