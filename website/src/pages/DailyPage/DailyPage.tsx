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

import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

import { Separator } from "@/components/ui/separator";
import DailyCharts from "./components/DailyCharts";
import DailDailySummary from "./components/DailyDailySummary";
import DailyHero from "./components/DailyHero";

export default function DailyPage() {
  const urlParams = new URLSearchParams(window.location.search);
  const [workload, setWorkload] = useState<string>(
    urlParams.get("workload") ?? "OLTP"
  );

  const navigate = useNavigate();

  useEffect(() => {
    navigate(`?workload=${workload}`);
  }, [workload]);

  return (
    <>
      <DailyHero />
      <DailDailySummary workload={workload} setWorkload={setWorkload} />
      <Separator className="mx-auto w-[80%] foreground" />
      <DailyCharts workload={workload} />
    </>
  );
}
