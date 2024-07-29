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
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { CompareContext } from "@/contexts/CompareContext";
import { MoreHorizontal } from "lucide-react";
import { useContext } from "react";
import { useNavigate } from "react-router-dom";

export type CompareRowActionsProps = {
  gitRef: string;
};

export default function CompareRowActions(props: CompareRowActionsProps) {
  const { gitRef } = props;
  const { setOldGitRef, setNewGitRef, showCompare } =
    useContext(CompareContext);
  const navigate = useNavigate();

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" className="h-10 w-10 p-0">
          <span className="sr-only">Open menu</span>
          <MoreHorizontal className="h-6 w-6" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuLabel>Actions</DropdownMenuLabel>
        <DropdownMenuSub>
          <DropdownMenuSubTrigger>Add to compare</DropdownMenuSubTrigger>
          <DropdownMenuSubContent>
            <DropdownMenuGroup>
              <DropdownMenuItem
                onClick={() => {
                  setOldGitRef(gitRef);
                  showCompare();
                }}
              >
                Add as old
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => {
                  setNewGitRef(gitRef);
                  showCompare();
                }}
              >
                Add as new
              </DropdownMenuItem>
            </DropdownMenuGroup>
          </DropdownMenuSubContent>
        </DropdownMenuSub>
        <DropdownMenuItem
          onClick={() => {
            navigate(`/history?gitRef=${gitRef}`);
          }}
        >
          Benchmarks History
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
