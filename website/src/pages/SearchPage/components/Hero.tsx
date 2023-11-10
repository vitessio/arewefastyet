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

interface HeroProps {
  gitRef: string;
  setGitRef: React.Dispatch<React.SetStateAction<string>>;
}

export default function Hero(props: HeroProps) {
  const { gitRef, setGitRef } = props;

  return (
    <section className="h-[30vh] pt-[5vh] flex justify-center items-center">
      <div className="p-[3px] bg-gradient-to-br from-primary to-accent rounded-full w-1/2 duration-300 focus-within:p-[1px]">
        <input
          type="text"
          onChange={(e) => setGitRef(e.target.value)}
          defaultValue={gitRef}
          className="px-4 py-2 relative text-lg bg-background rounded-inherit w-full text-center focus:outline-none duration-inherit focus:border-none"
          placeholder="Enter commit SHA"
        />
      </div>
    </section>
  );
}
