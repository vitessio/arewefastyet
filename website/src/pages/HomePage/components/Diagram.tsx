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

import { Card, CardContent } from "@/components/ui/card";

export default function Diagram() {
  return (
    <section className="p-page flex flex-col items-center my-8">
      <h2 className="text-4xl font-semibold my-14 text-primary dark:text-front">
        Architecture
      </h2>
      <Card className="w-full max-w-screen-xl my-8 border-border">
        <CardContent className="flex justify-center">
          <div className="relative duration-1000">
            <img
              className="duration-inherit"
              src="/images/execution-pipeline-dark.png"
              alt="execution pipeline"
            />

            <img
              className={"absolute-cover duration-inherit dark:opacity-0"}
              src="/images/execution-pipeline.png"
              alt="execution pipeline"
            />
          </div>
        </CardContent>
      </Card>
    </section>
  );
}
