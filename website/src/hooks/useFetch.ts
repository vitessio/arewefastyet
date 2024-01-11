import axios, { AxiosRequestConfig } from "axios";
import { useEffect, useState } from "react";

export default function useFetch<T, I = undefined>(
  url: string,
  config?: AxiosRequestConfig & {
    initial?: I;
    method?: "get" | "post" | "put" | "delete";
  }
) {
  const [data, setData] = useState<T | I>(config?.initial as I);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<any>(null);

  async function loadData() {
    try {
      const response = await axios[config?.method || "get"]<T>(url, config);

      if (response.statusText !== "OK") {
        throw new Error(`HTTP error status: ${response.status}`);
      }

      setData(response.data);
    } catch (err) {
      setError(err);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadData();
  }, [url, config]);

  return [data, loading, error] as const;
}
