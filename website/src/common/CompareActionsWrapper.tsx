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
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CompareContext } from "@/contexts/CompareContext";
import { VitessRefs } from "@/types";
import useApiCall from "@/utils/Hook";
import { X } from "lucide-react";
import { useContext, useRef } from "react";
import Draggable from "react-draggable";
import { useNavigate } from "react-router-dom";
import CompareActions from "./CompareActions";

export default function CompareActionsWrapper() {
  const {
    oldGitRef,
    newGitRef,
    setOldGitRef,
    setNewGitRef,
    isCompareVisible,
    closeCompare,
  } = useContext(CompareContext);

  const navigate = useNavigate();

  const compareClicked = () => {
    navigate(`/compare?old=${oldGitRef}&new=${newGitRef}`);
    closeCompare();
  };

  const { data: vitessRefs } = useApiCall<VitessRefs>(
    `${import.meta.env.VITE_API_URL}vitess/refs`
  );

  const nodeRef = useRef(null);

  return (
    isCompareVisible && (
      <div className="fixed right-4 bottom-4 max-w-max w-3/5 md:w-max">
        <Draggable nodeRef={nodeRef} cancel=".cancel-drag">
          <div ref={nodeRef}>
            <Card className="border-border shadow-lg">
              <CardHeader className="flex flex-col items-center">
                <CardTitle className="text-primary text-lg font-medium">
                  Compare Versions
                </CardTitle>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={closeCompare}
                  className="absolute top-1 right-1 z-50  cancel-drag"
                >
                  <X className="h-4 w-4" />
                </Button>
              </CardHeader>
              <CardContent>
                <CompareActions
                  className="cancel-drag"
                  oldGitRef={oldGitRef}
                  setOldGitRef={setOldGitRef}
                  newGitRef={newGitRef}
                  setNewGitRef={setNewGitRef}
                  compareClicked={compareClicked}
                  vitessRefs={vitessRefs}
                  keyboardShortcut={{ old: "l", new: "m" }}
                />
              </CardContent>
            </Card>
          </div>
        </Draggable>
      </div>
    )
  );
}
