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

import useApiCall from "@/utils/Hook";
import { DailySummarydata, Workloads } from "@/types";

export default function useDailySummaryData(workloads: Workloads[]) {

  const workloadsQuery = workloads.join('&workloads=');

  const {
    data: dataDailySummary,
    isLoading: isLoadingDailySummary,
    error: dailySummaryError,
  } = useApiCall<DailySummarydata>(`${import.meta.env.VITE_API_URL}daily/summary?workloads=${workloadsQuery}`);

  return {
    dataDailySummary,
    isLoadingDailySummary,
    dailySummaryError,
  };
};