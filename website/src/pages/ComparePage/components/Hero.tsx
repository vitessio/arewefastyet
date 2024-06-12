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
import { twMerge } from "tailwind-merge";

interface GitRef {
  old: string;
  new: string;
}

interface CompareDataType {
  gitRef: GitRef;
  setGitRef: React.Dispatch<React.SetStateAction<GitRef>>;
}

interface ComparisonInputDataType extends CompareDataType {
  className: string;
  name: "old" | "new";
}

/**
 * Renders a section for comparing Git SHAs.
 * @param {CompareDataType} props - The props containing the git reference and reference setter function.
 * @returns {JSX.Element} Render hero component to compare git sha.
 */

export default function Hero({
  gitRef,
  setGitRef,
}: CompareDataType): JSX.Element {
  return (
    <section className="flex flex-col h-[32vh] pt-[8vh] justify-center items-center">
      <h1 className="mb-3 text-front text-opacity-70">
        Enter SHAs to compare commits
      </h1>
      <div className="flex overflow-hidden bg-gradient-to-br from-primary to-theme p-[2px] rounded-full">
        <ComparisonInput
          name="old"
          className="rounded-l-full"
          setGitRef={setGitRef}
          gitRef={gitRef}
        />
        <ComparisonInput
          name="new"
          className="rounded-r-full "
          setGitRef={setGitRef}
          gitRef={gitRef}
        />
      </div>
    </section>
  );
}

/**
 * Comparison Input component for inputting Git SHA values.
 * @param {ComparisonInputDataType} props - The props containing the git reference and reference setter function, className and Input name..
 * @returns {JSX.Element} The input element component.
 */

function ComparisonInput({
  className,
  gitRef,
  setGitRef,
  name,
}: ComparisonInputDataType): JSX.Element {
  return (
    <input
      type="text"
      name={name}
      className={twMerge(
        className,
        "relative text-xl px-6 py-2 bg-background focus:border-none focus:outline-none border border-primary"
      )}
      defaultValue={gitRef[name]}
      placeholder={`${name} SHA`}
      onChange={(event) =>
        setGitRef((p: any) => {
          return { ...p, [name]: event.target.value };
        })
      }
    />
  );
}
