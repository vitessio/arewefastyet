import axios from "axios";
import { useEffect, useState } from "react";
import { ApiCallConfig, ApiEndpoint, ApiResponse } from "../types";

const serverUrl = import.meta.env.VITE_API_URL;

export default function useApiCall<T extends ApiEndpoint>(
  uri: T,
  config?: ApiCallConfig<T>
) {
  const [data, setData] = useState<ApiResponse<T>>();
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<any>(null);

  async function loadData() {
    setLoading(true);
    try {
      const response = await axios[config?.method || "get"](uri, {
        baseURL: serverUrl,
        timeout: 32_000,
        headers: {
          "Content-Type": "application/json",
        },
        ...config,
      });

      setData(response.data);
    } catch (err: any) {
      setError(
        err.message || "There was a proplem while communicating with the server"
      );
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadData();
  }, [uri, JSON.stringify(config)]);

  return [data, loading, error] as const;
}
