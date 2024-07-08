import useApiCall from "@/utils/Hook";
import { DailySummarydata, Workloads } from "@/types";

export default function useDailySummaryData(workloads: Workloads[]) {

  const workloadsQuery = workloads.join('&workloads=');

  const {
    data: dataDailySummary,
    isLoading: isLoadingDailySummary,
    error: errorDailySummary,
  } = useApiCall<DailySummarydata>(`${import.meta.env.VITE_API_URL}daily/summary?workloads=${workloadsQuery}`);

  return {
    dataDailySummary,
    isLoadingDailySummary,
    errorDailySummary,
  };
};