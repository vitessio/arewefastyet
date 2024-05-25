import Macrobench from "../common/Macrobench";
import { BenchmarkStatus } from "./enums";

export type BenchmarkExecution<
  STATUS extends keyof typeof BenchmarkStatus = "Completed"
> = {
  source: string;
  git_ref: string;
  type_of: string;
  pull_nb: number;
} & (STATUS extends BenchmarkStatus.Completed
  ? {
      uuid: string;
      status: string;
      golang_version: string;
      started_at: string;
      finished_at: string;
    }
  : {});

export interface QPS {
  ID: number;
  RefID: number;
  total: number;
  reads: number;
  writes: number;
  other: number;
}

export interface Result {
  ID: number;
  qps: QPS;
  tps: number;
  latency: number;
  errors: number;
  reconnects: number;
  time: number;
  threads: number;
}

export interface ComponentsCPUTime {
  vtgate: number;
  vttablet: number;
}

export interface ComponentsMemStatsAllocBytes {
  vtgate: number;
  vttablet: number;
}

export interface Metrics {
  TotalComponentsCPUTime: number;
  ComponentsCPUTime: ComponentsCPUTime;
  TotalComponentsMemStatsAllocBytes: number;
  ComponentsMemStatsAllocBytes: ComponentsMemStatsAllocBytes;
}

export interface MacroBenchmark {
  ID: number;
  Source: string;
  CreatedAt: null;
  ExecUUID: string;
  GitRef: string;
  Result: Result;
  Metrics: Metrics;
}

export interface MicroResult {
  Ops: number;
  NSPerOp: number;
  MBPerSec: number;
  BytesPerOp: number;
  AllocsPerOp: number;
}

export interface MicroBenchmark {
  PkgName: string;
  Name: string;
  SubBenchmarkName: string;
  GitRef: string;
  StartedAt: string;
  Result: MicroResult;
}

export interface Macros {
  OLTP: MacroBenchmark[] | null;
  "OLTP-READONLY": MacroBenchmark[] | null;
  "OLTP-READONLY-OLAP": MacroBenchmark[] | null;
  "OLTP-SET": MacroBenchmark[] | null;
  TPCC: MacroBenchmark[] | null;
}

export type Micro = MicroBenchmark[] | null;

export interface MacrobenchDiff {
  Left: MacroBenchmark;
  Right: MacroBenchmark;
  Diff: Result;
  DiffMetrics: Metrics;
}

export interface MacrobenchComparison {
  type: string;
  result: MacrobenchResult;
}
export interface MacrobenchResult {
  total_qps: BenchmarkMetric;
  reads_qps: BenchmarkMetric;
  writes_qps: BenchmarkMetric;
  other_qps: BenchmarkMetric;
  tps: BenchmarkMetric;
  latency: BenchmarkMetric;
  errors: BenchmarkMetric;
  total_components_cpu_time: BenchmarkMetric;
  components_cpu_time: {
    vtgate: BenchmarkMetric;
    vttablet: BenchmarkMetric;
  };
  total_components_mem_stats_alloc_bytes: BenchmarkMetric;
  components_mem_stats_alloc_bytes: {
    vtgate: BenchmarkMetric;
    vttablet: BenchmarkMetric;
  };
}

export interface BenchmarkMetric {
  insignificant: boolean;
  delta: number;
  p: number;
  n1: number;
  n2: number;
  old: MetricDetail;
  new: MetricDetail;
}

export interface MetricDetail {
  center: number;
  confidence: number;
  range: {
    infinite: boolean;
    unknown: boolean;
    value: number;
  };
}

export type DailySummary = [
  {
    Name: "OLTP";
    Data: DailySummaryData;
  },
  {
    Name: "OLTP-READONLY";
    Data: DailySummaryData;
  },
  {
    Name: "OLTP-READONLY-OLAP";
    Data: DailySummaryData;
  },
  {
    Name: "OLTP-SET";
    Data: DailySummaryData;
  },
  {
    Name: "TPCC";
    Data: DailySummaryData;
  }
];

export interface DailySummaryData {
  CreatedAt: string;
  QPSTotal: number;
}

export interface PR {
  ID: number;
  Author: string;
  Title: string;
  CreatedAt: string;
  Base: string;
  Head: string;
}

export interface QueryPlanComparison {
  Left: QueryPlan;
  Right: QueryPlan;
  SamePlan: boolean;
  Key: string;
  ExecCountDiff: number;
  ExecTimeDiff: number;
  RowsReturnedDiff: number;
  ErrorsDiff: number;
}

export interface QueryPlan {
  Key: string;
  Value: {
    QueryType: string;
    Original: string;
    Instructions: string;
    ExecCount: number;
    ExecTime: number;
    ShardQueries: number;
    RowsReturned: number;
    RowsAffected: number;
    Errors: number;
    TablesUsed: string[] | null;
  };
}

export interface GitRef {
  Name: string;
  CommitHash: string;
  Version: { Major: number; Minor: number; Patch: number };
  RCnumber: number;
}
