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
import { twMerge } from "tailwind-merge";
import useTheme from "../../../hooks/useTheme";

export default function Diagram() {
  const theme = useTheme();

  return (
    <section className="p-page flex flex-col items-center my-8">
      <h1 className="bg-primary bg-opacity-20 text-primary text-lg font-semibold px-4 py-2 rounded-full">
        Diagramatic Overview
      </h1>

      <div className="relative duration-1000">
        <img
          className="duration-inherit"
          src="/images/execution-pipeline-dark.png"
          alt="execution pipeline"
        />

        <img
          className={twMerge(
            "absolute-cover duration-inherit",
            theme.current === "dark" && "opacity-0"
          )}
          src="/images/execution-pipeline.png"
          alt="execution pipeline"
        />
      </div>
    </section>
  );
}
