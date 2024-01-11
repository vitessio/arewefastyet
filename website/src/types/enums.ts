export enum BenchmarkStatus {
  Ongoing = "Ongoing",
  Completed = "Completed",
}

export enum BenchmarkType {
  "OLTP",
  "OLTP-READONLY",
  "OLTP-READONLY-OLAP",
  "OLTP-SET",
  "TPCC",
}
