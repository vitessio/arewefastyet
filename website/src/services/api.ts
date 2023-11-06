// import axios, { AxiosRequestConfig, AxiosResponse } from "axios";
// import { ApiEndpoint, ApiResponse } from "../types";

// const serverUrl = import.meta.env.VITE_API_URL;

// let axiosClient = createApi();

// function createApi() {
//   const client = axios.create({
//     baseURL: serverUrl,
//     timeout: 32_000,
//     headers: {
//       "Content-Type": "application/json",
//     },
//   });

//   return client;
// }

// const api = {
//   client: axiosClient.get,

//   async getRecentExecutions() {
//     return await client.get("/recent");
//   },

//   async getExecutionQueue() {
//     return await client.get("/queue");
//   },

//   async getStats() {
//     const response = await client.get("/status/stats");

//     return response.data;
//   },
// };

// const client = {
//   get<T extends ApiEndpoint>(
//     url: T,
//     config?: AxiosRequestConfig
//   ): Promise<AxiosResponse<ApiResponse<T>>> {
//     return axiosClient.get<ApiResponse<T>>(url, config);
//   },
//   post<T extends ApiEndpoint>(
//     url: T,
//     config?: AxiosRequestConfig
//   ): Promise<AxiosResponse<ApiResponse<T>>> {
//     return axiosClient.post<ApiResponse<T>>(url, config);
//   },
//   put<T extends ApiEndpoint>(
//     url: T,
//     config?: AxiosRequestConfig
//   ): Promise<AxiosResponse<ApiResponse<T>>> {
//     return axiosClient.put<ApiResponse<T>>(url, config);
//   },
//   delete<T extends ApiEndpoint>(
//     url: T,
//     config?: AxiosRequestConfig
//   ): Promise<AxiosResponse<ApiResponse<T>>> {
//     return axiosClient.delete<ApiResponse<T>>(url, config);
//   },
// };

// export default api;
