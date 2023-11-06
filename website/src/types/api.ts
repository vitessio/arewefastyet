import { BenchmarkExecution, BenchmarkStatus } from ".";

type ApiEndpointResponses = {
  "/recent": BenchmarkExecution<BenchmarkStatus.Completed>[];
  "/queue": BenchmarkExecution<BenchmarkStatus.Ongoing>[];
};

export type ApiEndpoint = keyof ApiEndpointResponses;

export type ApiResponse<T extends ApiEndpoint> = ApiEndpointResponses[T];
