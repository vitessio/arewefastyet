#!/usr/bin/env bash
set -euo pipefail

VTGATE_WORKLOAD_MODE="${VTGATE_WORKLOAD_MODE:-}"
VTGATE_MAX_MEMORY_ROWS="${VTGATE_MAX_MEMORY_ROWS:-}"

VTGATE_ARGS=(
  --cell=local
  --cells-to-watch=local
  --mysql-server-port=13306
  --grpc-port=15306
  --port=15001
  --mysql-auth-server-impl=none
  --planner-version=gen4
  --topo-implementation=etcd2
  --topo-global-server-address=etcd:2379
  --topo-global-root=/vitess/main
  --enable-buffer
  --buffer-size=300
  --vschema-ddl-authorized-users="%"
  --tablet-types-to-wait=PRIMARY
  --pprof-http
)

if [ -n "${VTGATE_WORKLOAD_MODE}" ]; then
  VTGATE_ARGS+=(--mysql-default-workload="${VTGATE_WORKLOAD_MODE}")
fi

if [ -n "${VTGATE_MAX_MEMORY_ROWS}" ]; then
  VTGATE_ARGS+=(--max-memory-rows="${VTGATE_MAX_MEMORY_ROWS}")
fi

echo "Starting vtgate..."
exec gosu vitess vtgate "${VTGATE_ARGS[@]}"
