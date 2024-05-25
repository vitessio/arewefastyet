import { AxiosRequestConfig } from "axios";
import {
  BenchmarkExecution,
  BenchmarkStatus,
  BenchmarkType,
  DailySummary,
  GitRef,
  MacroBenchmark,
  MacrobenchComparison,
  Macros,
  Micro,
  PR,
  QueryPlanComparison,
} from ".";

type ApiEndpoints = [
  {
    uri: "/recent";
    response: BenchmarkExecution<BenchmarkStatus.Completed>[];
  },
  {
    uri: "/queue";
    response: BenchmarkExecution<BenchmarkStatus.Ongoing>[];
  },
  {
    uri: "/status/stats";
    response: { Total: number; Finished: number; Last30Days: number };
  },
  {
    uri: "/vitess/refs";
    response: GitRef[];
  },
  {
    uri: "/pr/list";
    response: PR[];
  },
  {
    uri: `/pr/info/${string}`;
    response: PR;
  },
  {
    uri: "/search";
    response: { Macros: Macros; Micro: Micro };
    params: { git_ref: string };
  },
  {
    uri: "/macrobench/compare";
    response: MacrobenchComparison[];
    params: { ltag: string; rtag: string };
  },
  {
    uri: "/macrobench/compare/queries";
    response: QueryPlanComparison[];
    params: { ltag: string; rtag: string; type: string };
  },
  {
    uri: "/daily/summary";
    response: DailySummary;
  },
  {
    uri: "/daily";
    response: MacroBenchmark[];
    params: { type: BenchmarkType };
  }
];

export type ApiEndpoint = ApiEndpoints[number]["uri"];

type EndpointByUri<T extends ApiEndpoint> = Extract<
  ApiEndpoints[number],
  { uri: T }
>;

export type ApiResponse<T extends ApiEndpoint> = EndpointByUri<T>["response"];

export type ApiParams<T extends ApiEndpoint> = EndpointByUri<T> extends {
  params: any;
}
  ? EndpointByUri<T>["params"]
  : undefined;

type BaseApiCallConfig = {
  method?: "get" | "post" | "put" | "delete";
};

export type ApiCallConfig<T extends ApiEndpoint> = BaseApiCallConfig &
  (ApiParams<T> extends undefined
    ? AxiosRequestConfig
    : Omit<AxiosRequestConfig, "params"> & { params: ApiParams<T> });
