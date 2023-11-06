import axios, { AxiosRequestConfig } from "axios";
import { useEffect, useState } from "react";
import { ApiEndpoint, ApiResponse } from "../types";

const serverUrl = import.meta.env.VITE_API_URL;

export default function useApiCall<T extends ApiEndpoint>(
  uri: T,
  config?: AxiosRequestConfig & {
    method?: "get" | "post" | "put" | "delete";
  }
) {
  const [data, setData] = useState<ApiResponse<T>>();
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<any>(null);

  async function loadData() {
    try {
      const response = await axios[config?.method || "get"]<ApiResponse<T>>(
        uri,
        {
          baseURL: serverUrl,
          timeout: 32_000,
          headers: {
            "Content-Type": "application/json",
          },
          ...config,
        }
      );

      setData(response.data);
    } catch (err) {
      setError(err);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadData();
  }, [uri, config]);

  return [data, loading, error] as const;
}
