import axios, { AxiosRequestConfig } from "axios";
import { ApiEndpoint, ApiResponse, Prettify } from "../types";

const serverUrl = import.meta.env.VITE_API_URL;

let axiosClient = createApi();

function createApi() {
  const client = axios.create({
    baseURL: serverUrl,
    timeout: 32_000,
    headers: {
      "Content-Type": "application/json",
    },
  });

  return client;
}

const api = {
  client: axiosClient.get,

  async getRecentExecutions() {
    return await client.get("/recent");
  },

  async getExecutionQueue() {
    return await client.get("/queue");
  },
};

const client = {
  async get(url: ApiEndpoint, config?: AxiosRequestConfig) {
    return (await axiosClient.get<ApiResponse<typeof url>>(url, config)).data;
  },
  async post(url: ApiEndpoint, config?: AxiosRequestConfig) {
    return (await axiosClient.post<ApiResponse<typeof url>>(url, config)).data;
  },
  async put(url: ApiEndpoint, config?: AxiosRequestConfig) {
    return (await axiosClient.put<ApiResponse<typeof url>>(url, config)).data;
  },
  async delete(url: ApiEndpoint, config?: AxiosRequestConfig) {
    return (await axiosClient.delete<ApiResponse<typeof url>>(url, config))
      .data;
  },
};

export default api;
