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
import Hero, { HeroProps } from "@/common/Hero";

const heroProps: HeroProps = {
  title: "Pull Request",
  description: (
    <p>
      If a given Pull Request on vitessio/vitess is labelled with the
      <span className="bg-red-800 text-white px-2 py-1 rounded-2xl ml-2">Benchmark Me</span> label
      the Pull Request will be handled and benchmarked by arewefastyet. For each commit on the Pull Request there will be two benchmarks: one on the Pull Request's HEAD and another on the base of the Pull Request.
      <br />
      <br />
      On this page you can find all benchmarked Pull Requests.
    </p>
  ),
};

export default function PRHero() {
  return (
    <Hero title={heroProps.title} description={heroProps.description} />
  );
}
