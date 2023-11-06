import { BenchmarkExecution, BenchmarkStatus } from ".";

type ApiEndpointResponses = {
  "/recent": BenchmarkExecution<BenchmarkStatus.Completed>[];

  "/queue": BenchmarkExecution<BenchmarkStatus.Ongoing>[];

  "/status/stats": {
    Total: number;
    Finished: number;
    Last30Days: number;
  };
};

export type ApiEndpoint = keyof ApiEndpointResponses;

export type ApiResponse<T extends ApiEndpoint> = ApiEndpointResponses[T];
