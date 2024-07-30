/*
Copyright 2023 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import {
  CompareData,
  CompareResult,
  MacroBenchmarkTableData,
  VitessRefs,
} from "@/types";
import bytes from "bytes";

//FORMATTING BYTES TO Bytes
export const formatByte = (byte: number) => {
  const byteValue = bytes(byte);
  if (byteValue === null) {
    return "0";
  }
  return byteValue;
};

export function fixed(value: number, fractionDigits: number): string {
  if (
    value === null ||
    typeof value === "undefined" ||
    isNaN(value) ||
    value === 0
  ) {
    return "0";
  }
  return value.toFixed(fractionDigits);
}

export function secondToMicrosecond(value: number): string {
  return fixed(value * 1000000, 2) + "μs";
}

//ERROR API MESSAGE ERROR

export const errorApi: string =
  "An error occurred while retrieving data from the API. Please try again.";

//NUMBER OF PIXELS TO OPEN AND CLOSE THE DROP-DOWN
export const openDropDownValue = 1000;
export const closeDropDownValue = 58;

export function formatGitRef(gitRef: string): string {
  return gitRef.slice(0, 8);
}

export function formatCompareData(
  data: CompareData[],
): MacroBenchmarkTableData[] {
  return (
    data.map((data: CompareData) => {
      return {
        qpsTotal: {
          title: "QPS Total",
          old: data.result.total_qps.old,
          new: data.result.total_qps.new,
          p: data.result.total_qps.p,
          delta: data.result.total_qps.delta,
          insignificant: data.result.total_qps.insignificant,
        },
        qpsReads: {
          title: "Reads",
          old: data.result.reads_qps.old,
          new: data.result.reads_qps.new,
          p: data.result.reads_qps.p,
          delta: data.result.reads_qps.delta,
          insignificant: data.result.reads_qps.insignificant,
        },
        qpsWrites: {
          title: "Writes",
          old: data.result.writes_qps.old,
          new: data.result.writes_qps.new,
          p: data.result.writes_qps.p,
          delta: data.result.writes_qps.delta,
          insignificant: data.result.writes_qps.insignificant,
        },
        qpsOther: {
          title: "Other",
          old: data.result.other_qps.old,
          new: data.result.other_qps.new,
          p: data.result.other_qps.p,
          delta: data.result.other_qps.delta,
          insignificant: data.result.other_qps.insignificant,
        },
        tps: {
          title: "TPS",
          old: data.result.tps.old,
          new: data.result.tps.new,
          p: data.result.tps.p,
          delta: data.result.tps.delta,
          insignificant: data.result.tps.insignificant,
        },
        latency: {
          title: "P95 Latency",
          old: data.result.latency.old,
          new: data.result.latency.new,
          p: data.result.latency.p,
          delta: data.result.latency.delta,
          insignificant: data.result.latency.insignificant,
        },
        errors: {
          title: "Errors / Second",
          old: data.result.errors.old,
          new: data.result.errors.new,
          p: data.result.errors.p,
          delta: data.result.errors.delta,
          insignificant: data.result.errors.insignificant,
        },
        totalComponentsCpuTime: {
          title: "Total CPU / Query",
          old: data.result.total_components_cpu_time.old,
          new: data.result.total_components_cpu_time.new,
          p: data.result.total_components_cpu_time.p,
          delta: data.result.total_components_cpu_time.delta,
          insignificant: data.result.total_components_cpu_time.insignificant,
        },
        vtgateCpuTime: {
          title: "vtgate",
          old: data.result.components_cpu_time.vtgate.old,
          new: data.result.components_cpu_time.vtgate.new,
          p: data.result.components_cpu_time.vtgate.p,
          delta: data.result.components_cpu_time.vtgate.delta,
          insignificant: data.result.components_cpu_time.vtgate.insignificant,
        },
        vttabletCpuTime: {
          title: "vttablet",
          old: data.result.components_cpu_time.vttablet.old,
          new: data.result.components_cpu_time.vttablet.new,
          p: data.result.components_cpu_time.vttablet.p,
          delta: data.result.components_cpu_time.vttablet.delta,
          insignificant: data.result.components_cpu_time.vttablet.insignificant,
        },
        totalComponentsMemStatsAllocBytes: {
          title: "Total Allocated / Query",
          old: data.result.total_components_mem_stats_alloc_bytes.old,
          new: data.result.total_components_mem_stats_alloc_bytes.new,
          p: data.result.total_components_mem_stats_alloc_bytes.p,
          delta: data.result.total_components_mem_stats_alloc_bytes.delta,
          insignificant:
            data.result.total_components_mem_stats_alloc_bytes.insignificant,
        },
        vtgateMemStatsAllocBytes: {
          title: "vtgate",
          old: data.result.components_mem_stats_alloc_bytes.vtgate.old,
          new: data.result.components_mem_stats_alloc_bytes.vtgate.new,
          p: data.result.components_mem_stats_alloc_bytes.vtgate.p,
          delta: data.result.components_mem_stats_alloc_bytes.vtgate.delta,
          insignificant:
            data.result.components_mem_stats_alloc_bytes.vtgate.insignificant,
        },
        vttabletMemStatsAllocBytes: {
          title: "vttablet",
          old: data.result.components_mem_stats_alloc_bytes.vttablet.old,
          new: data.result.components_mem_stats_alloc_bytes.vttablet.new,
          p: data.result.components_mem_stats_alloc_bytes.vttablet.p,
          delta: data.result.components_mem_stats_alloc_bytes.vttablet.delta,
          insignificant:
            data.result.components_mem_stats_alloc_bytes.vttablet.insignificant,
        },
      };
    }) || []
  );
}

export function formatCompareResult(
  data: CompareResult,
): MacroBenchmarkTableData {
  return {
    qpsTotal: {
      title: "QPS Total",
      old: data.total_qps.old,
      new: data.total_qps.new,
      p: data.total_qps.p,
      delta: data.total_qps.delta,
      insignificant: data.total_qps.insignificant,
    },
    qpsReads: {
      title: "Reads",
      old: data.reads_qps.old,
      new: data.reads_qps.new,
      p: data.reads_qps.p,
      delta: data.reads_qps.delta,
      insignificant: data.reads_qps.insignificant,
    },
    qpsWrites: {
      title: "Writes",
      old: data.writes_qps.old,
      new: data.writes_qps.new,
      p: data.writes_qps.p,
      delta: data.writes_qps.delta,
      insignificant: data.writes_qps.insignificant,
    },
    qpsOther: {
      title: "Other",
      old: data.other_qps.old,
      new: data.other_qps.new,
      p: data.other_qps.p,
      delta: data.other_qps.delta,
      insignificant: data.other_qps.insignificant,
    },
    tps: {
      title: "TPS",
      old: data.tps.old,
      new: data.tps.new,
      p: data.tps.p,
      delta: data.tps.delta,
      insignificant: data.tps.insignificant,
    },
    latency: {
      title: "P95 Latency",
      old: data.latency.old,
      new: data.latency.new,
      p: data.latency.p,
      delta: data.latency.delta,
      insignificant: data.latency.insignificant,
    },
    errors: {
      title: "Errors / Second",
      old: data.errors.old,
      new: data.errors.new,
      p: data.errors.p,
      delta: data.errors.delta,
      insignificant: data.errors.insignificant,
    },
    totalComponentsCpuTime: {
      title: "Total CPU / Query",
      old: data.total_components_cpu_time.old,
      new: data.total_components_cpu_time.new,
      p: data.total_components_cpu_time.p,
      delta: data.total_components_cpu_time.delta,
      insignificant: data.total_components_cpu_time.insignificant,
    },
    vtgateCpuTime: {
      title: "vtgate",
      old: data.components_cpu_time.vtgate.old,
      new: data.components_cpu_time.vtgate.new,
      p: data.components_cpu_time.vtgate.p,
      delta: data.components_cpu_time.vtgate.delta,
      insignificant: data.components_cpu_time.vtgate.insignificant,
    },
    vttabletCpuTime: {
      title: "vttablet",
      old: data.components_cpu_time.vttablet.old,
      new: data.components_cpu_time.vttablet.new,
      p: data.components_cpu_time.vttablet.p,
      delta: data.components_cpu_time.vttablet.delta,
      insignificant: data.components_cpu_time.vttablet.insignificant,
    },
    totalComponentsMemStatsAllocBytes: {
      title: "Total Allocated / Query",
      old: data.total_components_mem_stats_alloc_bytes.old,
      new: data.total_components_mem_stats_alloc_bytes.new,
      p: data.total_components_mem_stats_alloc_bytes.p,
      delta: data.total_components_mem_stats_alloc_bytes.delta,
      insignificant: data.total_components_mem_stats_alloc_bytes.insignificant,
    },
    vtgateMemStatsAllocBytes: {
      title: "vtgate",
      old: data.components_mem_stats_alloc_bytes.vtgate.old,
      new: data.components_mem_stats_alloc_bytes.vtgate.new,
      p: data.components_mem_stats_alloc_bytes.vtgate.p,
      delta: data.components_mem_stats_alloc_bytes.vtgate.delta,
      insignificant: data.components_mem_stats_alloc_bytes.vtgate.insignificant,
    },
    vttabletMemStatsAllocBytes: {
      title: "vttablet",
      old: data.components_mem_stats_alloc_bytes.vttablet.old,
      new: data.components_mem_stats_alloc_bytes.vttablet.new,
      p: data.components_mem_stats_alloc_bytes.vttablet.p,
      delta: data.components_mem_stats_alloc_bytes.vttablet.delta,
      insignificant:
        data.components_mem_stats_alloc_bytes.vttablet.insignificant,
    },
  };
}

export const getRefName = (gitRef: string, vitessRefs: VitessRefs): string => {
  if (gitRef === "") {
    return "";
  }
  let title = gitRef;
  vitessRefs.branches.forEach((branch) => {
    if (branch.commit_hash.match(gitRef)) {
      title = branch.name;
    }
  });
  vitessRefs.tags.forEach((branch) => {
    if (branch.commit_hash.match(gitRef)) {
      title = branch.name;
    }
  });
  return title;
};

export const getGitRefFromRefName = (refName: string, vitessRefs: VitessRefs): string => {
  if (refName === "") {
    return "";
  }
  let gitRef = refName;
  vitessRefs.branches.forEach((branch) => {
    if (branch.name === refName) {
      gitRef = branch.commit_hash;
    }
  });
  vitessRefs.tags.forEach((tag) => {
    if (tag.name === refName) {
      gitRef = tag.commit_hash;
    }
  });
  return gitRef;
};
