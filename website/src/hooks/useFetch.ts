import axios, { AxiosRequestConfig } from "axios";
import { useEffect, useState } from "react";

const serverUrl = import.meta.env.VITE_API_URL;

export default function useFetch<T>(
  url: string,
  method: "get" | "post" | "put" | "delete" = "get",
  config?: AxiosRequestConfig
) {
  const [data, setData] = useState<T>();
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<any>(null);

  async function loadData() {
    try {
      const response = await axios[method]<T>(url, {
        baseURL: serverUrl,
        timeout: 32_000,
        headers: {
          "Content-Type": "application/json",
        },
        ...config,
      });

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

type ApiEndpointMap = {
  "/data1": { result: string };
  "/data2": { result: number };
  "/data3": { result: boolean };
};

type ApiEndpoint = keyof ApiEndpointMap;
