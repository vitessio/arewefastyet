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

export default function Hero(props: { gitRef: any; setGitRef: any }) {
  const { gitRef, setGitRef } = props;

  return (
    <section className="flex flex-col h-[32vh] pt-[8vh] justify-center items-center">
      <h1 className="mb-3 text-front text-opacity-70">
        Enter SHAs to compare commits
      </h1>
      <div className="flex overflow-hidden bg-gradient-to-br from-primary to-theme p-[2px] rounded-full m-4">
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

function ComparisonInput(props: {
  className: any;
  gitRef: any;
  setGitRef: any;
  name: any;
}) {
  const { className, gitRef, setGitRef, name } = props;

  return (
    <input
      type="text"
      name={name}
      className={twMerge(
        className,
        "relative text-xl px-6 py-2 bg-background focus:border-none focus:outline-none border border-primary w-full sm:w-48 md:w-64 lg:w-80"
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