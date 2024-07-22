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
import { ChangeEvent, useEffect, useState } from "react";

type VitessRefsCommandProps = {
  inputLabel: string;
  gitRef: string;
  setGitRef: (value: string) => void;
  vitessRefs: VitessRefs | null;
  keyboardShortcut?: string;
};

export default function VitessRefsCommand({
  inputLabel,
  setGitRef,
  vitessRefs,
  gitRef,
  keyboardShortcut = "k",
}: VitessRefsCommandProps) {
  const [open, setOpen] = useState(false);
  const [inputValue, setInputValue] = useState(gitRef);
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
      if (
        (e.key === keyboardShortcut && (e.metaKey || e.ctrlKey)) ||
        e.key === "/"
      ) {
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

  useEffect(() => {
    setInputValue(gitRef);
  }, [gitRef]);

  return (
    <>
      <Button
        variant="outline"
        className={cn(
          "relative justify-between rounded-[0.5rem] bg-muted/50 text-sm font-normal text-muted-foreground shadow-none sm:pr-4 w-full md:w-40 lg:w-64"
        )}
        onClick={() => setOpen(true)}
      >
        <span className="hidden lg:inline-flex overflow-hidden">
          {selectedRefName || gitRef || inputLabel}
        </span>
        <span className="inline-flex lg:hidden overflow-hidden">
          {selectedRefName || gitRef || "Search..."}
        </span>
        <kbd className="pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100">
          <span className="text-xs">⌘</span>
          {keyboardShortcut}
        </kbd>
      </Button>
      <CommandDialog open={open} onOpenChange={setOpen}>
        <CommandInput
          placeholder={inputLabel}
          value={inputValue}
          onInput={(e: ChangeEvent<HTMLInputElement>) =>
            setInputValue(e.target.value)
          }
          onKeyDown={handleInputKeyDown}
          className="w-full max-w-md mx-auto sm:max-w-lg lg:max-w-xl"
        />
        <CommandList>
          <CommandEmpty>No results found.</CommandEmpty>
          {(vitessRefs?.branches ?? []).length > 0 && (
            <CommandGroup heading="Branches">
              {vitessRefs?.branches?.map((ref) => (
                <CommandItem key={ref.name} onSelect={() => handleSelect(ref)}>
                  <span>{ref.name}</span>
                  <span className="hidden">{ref.commit_hash}</span>
                </CommandItem>
              ))}
            </CommandGroup>
          )}
          {(vitessRefs?.tags ?? []).length > 0 && (
            <CommandGroup heading="Releases">
              {vitessRefs?.tags?.map((ref, index) => (
                <CommandItem key={index} onSelect={() => handleSelect(ref)}>
                  <span>{ref.name}</span>
                  <span hidden>{ref.commit_hash}</span>
                </CommandItem>
              ))}
            </CommandGroup>
          )}
        </CommandList>
      </CommandDialog>
    </>
  );
}
