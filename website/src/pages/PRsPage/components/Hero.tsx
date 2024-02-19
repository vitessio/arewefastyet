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

export default function Hero() {
  return (
    <section className="flex h-[60vh] pt-[15vh] items-center p-page">
      <div className="flex flex-col justify-start items-start gap-5 max-w-screen-md">
        <h2 className="text-8xl text-primary">Pull Request</h2>
        <p className="md:my-6 leading-loose">
          If a given Pull Request on vitessio/vitess is labelled with the
          Benchmark me label the Pull Request will be handled and benchmark by
          arewefastyet. For each commit on the Pull Request there will be two
          benchmarks: one on the Pull Request's HEAD and another on the base of
          the Pull Request.
          <br />
          <br />
          On this page you can find all benchmarked Pull Requests.
        </p>
      </div>
    </section>
  );
}
