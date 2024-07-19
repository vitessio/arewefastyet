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

import { Button } from "@/components/ui/button";
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import { cn } from "@/library/utils";
import { VitessRefs, VitessRefsData } from "@/types";
import { ChangeEvent, useCallback, useEffect, useState } from "react";

interface VitessRefsCommandProps {
  inputLabel: string;
  setGitRef: (value: string) => void;
  vitessRefs: VitessRefs | null;
}

export default function VitessRefsCommand({
  inputLabel,
  setGitRef,
  vitessRefs,
  ...props
}: VitessRefsCommandProps) {
  const [open, setOpen] = useState(false);
  const [inputValue, setInputValue] = useState("");
  const [selectedRefName, setSelectedRefName] = useState("");

  const handleSelect = (vitessRef: VitessRefsData) => {
    setInputValue(vitessRef.name);
    setSelectedRefName(vitessRef.name);
    setGitRef(vitessRef.commit_hash);
    setOpen(false);
  };

  const handleInputKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      setGitRef(inputValue);
      setSelectedRefName(inputValue);
      setOpen(false);
    }
  };

  useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if ((e.key === "k" && (e.metaKey || e.ctrlKey)) || e.key === "/") {
        if (
          (e.target instanceof HTMLElement && e.target.isContentEditable) ||
          e.target instanceof HTMLInputElement ||
          e.target instanceof HTMLTextAreaElement ||
          e.target instanceof HTMLSelectElement
        ) {
          return;
        }

        e.preventDefault();
        setOpen((open) => !open);
      }
    };

    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  return (
    <>
      <Button
        variant="outline"
        className={cn(
          "relative h-full w-full justify-start rounded-[0.5rem] bg-muted/50 text-sm font-normal text-muted-foreground shadow-none sm:pr-12 md:w-40 lg:w-64 overflow-hidden"
        )}
        onClick={() => setOpen(true)}
        {...props}
      >
        <span className="hidden lg:inline-flex w-fi">
          {selectedRefName || inputLabel}
        </span>
        <span className="inline-flex lg:hidden w-full">
          {selectedRefName || "Search..."}
        </span>
      </Button>
      <CommandDialog open={open} onOpenChange={setOpen} >
        <CommandInput
          placeholder={inputLabel}
          value={inputValue}
          onInput={(e: ChangeEvent<HTMLInputElement>) =>
            setInputValue(e.target.value)
          }
          onKeyDown={handleInputKeyDown}
        />
        <CommandList>
          <CommandEmpty>No results found.</CommandEmpty>
          <CommandGroup heading="Branches">
            {vitessRefs?.branches?.map((ref) => (
              <CommandItem key={ref.name} onSelect={() => handleSelect(ref)}>
                <span>{ref.name}</span>
                <span className="hidden">{ref.commit_hash}</span>
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
      </CommandDialog>
    </>
  );
}
