import useApiCall from "@/utils/Hook";
import { DailySummarydata } from "@/types";

export default function useDailySummaryData() {
  const {
    data: dataDailySummary,
    isLoading: isLoadingDailySummary,
    error: errorDailySummary,
  } = useApiCall<DailySummarydata>(`${import.meta.env.VITE_API_URL}daily/summary`);

  return {
    dataDailySummary,
    isLoadingDailySummary,
    errorDailySummary,
  };
};