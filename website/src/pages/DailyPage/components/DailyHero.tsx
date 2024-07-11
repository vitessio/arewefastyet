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
  title: "Daily",
  description: (
    <>
      We run all macro benchmark workloads against the <i>main</i> branch every
      day. This is done to ensure the consistency of the results over time. On
      this page, you can find graphs that show you the results of all workloads
      over the last 30 days. Click on a workload to see the historical results
      of that workload.
    </>
  ),
};

export default function DailyHero() {
  return <Hero title={heroProps.title} description={heroProps.description} />;
}
