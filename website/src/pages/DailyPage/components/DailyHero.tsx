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
import Hero, { HeroProps } from "@/common/Hero";

const heroProps: HeroProps = {
  title: "Daily",
  description: (
    <p>
      We run all macro benchmark workloads against the <i>main</i> branch
      every day. This is done to ensure the consistency of the results over
      time on <i>main</i>. On this page, you can find graphs that show you the
      results of all five macro benchmark workload over the last 30 days. Click
      on a macro benchmark workload to see all the results for that workload.
    </p>
  ),
};


export default function DailyHero() {
  return (
    <Hero title={heroProps.title} description={heroProps.description} />
  );
}
