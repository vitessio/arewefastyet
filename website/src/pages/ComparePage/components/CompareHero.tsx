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
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import { VitessRefs } from "@/types";
import useApiCall from "@/utils/Hook";
import { formatGitRef } from "@/utils/Utils";
import { ChangeEvent, useState } from "react";

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
  const [oldInputValue, setOldInputValue] = useState("");
  const [newInputValue, setNewInputValue] = useState("");
  const isButtonDisabled = !oldInputValue || !newInputValue;
  const [oldListVisible, setOldListVisible] = useState(false);
  const [newListVisible, setNewListVisible] = useState(false);

  const {
    data: vitessRefs,
    isLoading: vitessRefsLoading,
    error: vitessRefsError,
  } = useApiCall<VitessRefs>(`${import.meta.env.VITE_API_URL}vitess/refs`);

  console.log({vitessRefs});

  const handleOldSelect = (commitHash: string) => {
    setOldInputValue(commitHash);
    setOldListVisible(false);
  };

  const handleNewSelect = (commitHash: string) => {
    setNewInputValue(commitHash);
  };

  const compareClicked = () => {
    setGitRef({ old: oldInputValue, new: newInputValue });
  };

  const handleOldInputFocus = () => setOldListVisible(true);
  const handleNewInputFocus = () => setNewListVisible(true);

  return (
    <Hero title={heroProps.title}>
      <div className="flex flex-row gap-4">
        <Command className="w-[300px] rounded-lg border shadow-md">
          <CommandInput
            placeholder="Search commits or releases..."
            value={formatGitRef(oldInputValue)}
            onInput={(e: ChangeEvent<HTMLInputElement>) =>
              setOldInputValue(e.target.value)
            }
            onFocus={handleOldInputFocus}
          />
          {oldListVisible && (
            <CommandList>
              <CommandEmpty>No results found.</CommandEmpty>
              <CommandGroup heading="Branches">
                {vitessRefs?.branches?.map((ref, index) => (
                  <CommandItem
                    key={index}
                    onSelect={() => handleOldSelect(ref.commit_hash)}
                  >
                    <span>
                      {ref.name}
                      {/* <span className="hidden">{ref.commit_hash}</span> */}
                    </span>
                  </CommandItem>
                ))}
              </CommandGroup>
              <CommandGroup heading="Releases">
                {vitessRefs?.tags?.map((ref, index) => (
                  <CommandItem
                    key={index}
                    onSelect={() => handleOldSelect(ref.commit_hash)}
                  >
                    <span>
                      {ref.name}
                      <span className="hidden">{ref.commit_hash}</span>
                    </span>
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          )}
        </Command>
        {/* <Command className="w-[300px] rounded-lg border shadow-md">
          <CommandInput
            placeholder="Search commits or releases..."
            value={formatGitRef(newInputValue)}
            onInput={(e: ChangeEvent<HTMLInputElement>) =>
              setNewInputValue(e.target.value)
            }
            onFocus={handleNewInputFocus}
          />
          {newListVisible && (
            <CommandList>
              <CommandEmpty>No results found.</CommandEmpty>
              <CommandGroup heading="Releases">
                {vitessRefs?.map((ref, index) => (
                  <CommandItem
                    key={index}
                    onSelect={() => handleNewSelect(ref.CommitHash)}
                  >
                    <span>{ref.Name}</span>
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          )}
        </Command> */}
        <Button onClick={compareClicked} disabled={isButtonDisabled}>
          Compare
        </Button>
      </div>
    </Hero>
  );
}
