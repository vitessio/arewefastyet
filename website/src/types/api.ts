import { AxiosRequestConfig } from "axios";
import { BenchmarkExecution, BenchmarkStatus, Macros, Micro } from ".";

type SimpleApiEndpoints = {
  "/recent": {
    response: BenchmarkExecution<BenchmarkStatus.Completed>[];
  };

  "/queue": {
    response: BenchmarkExecution<BenchmarkStatus.Ongoing>[];
  };

  "/status/stats": {
    response: {
      Total: number;
      Finished: number;
      Last30Days: number;
    };
  };
};

export type ApiEndpointsWithParams = {
  "/search": {
    response: { Macros: Macros; Micro: Micro };
    params: { git_ref: string };
  };
};

type ApiEndpoints = SimpleApiEndpoints & ApiEndpointsWithParams;

export type ApiEndpoint = keyof ApiEndpoints;

export type ApiResponse<T extends ApiEndpoint> = ApiEndpoints[T]["response"];

export type ApiParams<T extends ApiEndpoint> =
  T extends keyof ApiEndpointsWithParams
    ? ApiEndpointsWithParams[T]["params"]
    : undefined;

type BaseApiCallConfig = {
  method?: "get" | "post" | "put" | "delete";
};

export type ApiCallConfig<T extends ApiEndpoint> = BaseApiCallConfig &
  (ApiParams<T> extends undefined
    ? AxiosRequestConfig
    : Omit<AxiosRequestConfig, "params"> & { params: ApiParams<T> });
