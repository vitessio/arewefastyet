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

import { DailySummarydata } from "@/types";
import useApiCall from "@/utils/Hook";

export default function useDailySummaryData(workloads: string[]) {
  const workloadsQuery = workloads.join("&workload=");

  const {
    data: dataDailySummary,
    isLoading: isLoadingDailySummary,
    error: dailySummaryError,
  } = useApiCall<DailySummarydata>(
    `${import.meta.env.VITE_API_URL}daily/summary?workload=${workloadsQuery}`
  );

  if (dailySummaryError || dataDailySummary.length === 0) {
    return {
      dataDailySummary: null,
      isLoadingDailySummary,
      dailySummaryError:
        dailySummaryError || "An error occurred while fetching data.",
    };
  }

  return {
    dataDailySummary,
    isLoadingDailySummary,
    dailySummaryError: null,
  };
}
