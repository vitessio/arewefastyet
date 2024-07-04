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

interface SearchHeroProps {
  setGitRef: (value: string) => void;
}

const searchHeroProps: HeroProps = {
  title: "Search"
};

export default function SearchHero({ setGitRef }: SearchHeroProps) {
  return (
    <Hero title={searchHeroProps.title}>
      <div className="p-[3px] bg-gradient-to-br from-primary to-theme rounded-full duration-300 focus-within:p-[1px]">
        <input
          type="text"
          onChange={(e) => setGitRef(e.target.value)}
          className="px-4 py-2 relative text-lg bg-background rounded-inherit w-full text-center focus:outline-none duration-inherit focus:border-none"
          placeholder="Enter commit SHA"
        />
      </div>
    </Hero>
  );
}
