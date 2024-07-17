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

import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import { VitessRefs, VitessRefsData } from "@/types";
import { ChangeEvent, useState } from "react";

interface VitessRefsCommandProps {
  inputLabel: string;
  setGitRef: (value: string) => void;
  listVisible: boolean;
  setListVisible: (visible: boolean) => void;
  vitessRefs: VitessRefs | null;
}

export default function VitessRefsCommand({
  inputLabel,
  setGitRef,
  listVisible,
  setListVisible,
  vitessRefs,
}: VitessRefsCommandProps) {
  const handleSelect = (vitessRef: VitessRefsData) => {
    setInputValue(vitessRef.name);
    setGitRef(vitessRef.commit_hash);
    setTimeout(() => setListVisible(false), 100);
  };

  const [inputValue, setInputValue] = useState("");

  const handleInputFocus = () => setListVisible(true);
  const handleInputBlur = () => setTimeout(() => setListVisible(false), 200);

  return (
    <div className="relative w-[300px]">
      <Command className="rounded-lg border shadow-md">
        <CommandInput
          placeholder={inputLabel}
          value={inputValue}
          onInput={(e: ChangeEvent<HTMLInputElement>) =>
            setInputValue(e.target.value)
          }
          onFocus={handleInputFocus}
          onBlur={handleInputBlur}
        />
        {listVisible && (
          <CommandList className="absolute z-10 w-full bg-white border border-t-0 rounded-b-lg shadow-md mt-12">
            <CommandEmpty>No results found.</CommandEmpty>
            <CommandGroup heading="Branches">
              {vitessRefs?.branches?.map((ref, index) => (
                <CommandItem key={ref.name} onSelect={() => handleSelect(ref)}>
                  <span>{ref.name}</span>
                  <span className="hidden" >{ref.commit_hash}</span>
                </CommandItem>
              ))}
            </CommandGroup>
            <CommandGroup heading="Releases">
              {vitessRefs?.tags?.map((ref, index) => (
                <CommandItem key={index} onSelect={() => handleSelect(ref)}>
                  <span>{ref.name}</span>
                  <span hidden>{ref.commit_hash}</span>
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        )}
      </Command>
    </div>
  );
}
