#!/usr/bin/env bash
set -euo pipefail

SYSBENCH_WORKLOAD="${SYSBENCH_WORKLOAD:-oltp_read_write}"
SYSBENCH_WORKING_DIR="${SYSBENCH_WORKING_DIR:-}"
SYSBENCH_EXTRA="${SYSBENCH_EXTRA:-}"
THREADS="${THREADS:-42}"
TABLES="${TABLES:-10}"
TABLE_SIZE="${TABLE_SIZE:-10000}"
SCALE="${SCALE:-}"
WARMUP_TIME="${WARMUP_TIME:-20}"
RUN_TIME="${RUN_TIME:-60}"

PHASE="${1:-all}"

if [ -n "${SYSBENCH_WORKING_DIR}" ]; then
  cd "${SYSBENCH_WORKING_DIR}"
fi

SYSBENCH_COMMON=(
  "${SYSBENCH_WORKLOAD}"
  --mysql-host=vtgate
  --mysql-port=13306
  --mysql-db=main
  --mysql-user=root
  --mysql-password=
  --threads="${THREADS}"
  --tables="${TABLES}"
  --db-ps-mode=disable
  --rand-type=uniform
  --rand-seed=1
)

# Add table_size or scale depending on workload
if [ -n "${SCALE}" ]; then
  SYSBENCH_COMMON+=(--scale="${SCALE}")
else
  SYSBENCH_COMMON+=(--table-size="${TABLE_SIZE}")
fi

# Add extra workload-specific flags
if [ -n "${SYSBENCH_EXTRA}" ]; then
  read -ra EXTRA_ARGS <<< "${SYSBENCH_EXTRA}"
  SYSBENCH_COMMON+=("${EXTRA_ARGS[@]}")
fi

case "${PHASE}" in
  prepare-and-warmup)
    echo "==> sysbench prepare (workload=${SYSBENCH_WORKLOAD}, threads=${THREADS}, tables=${TABLES})"
    sysbench "${SYSBENCH_COMMON[@]}" prepare

    echo "==> sysbench warmup (time=${WARMUP_TIME}s)"
    sysbench "${SYSBENCH_COMMON[@]}" \
      --time="${WARMUP_TIME}" \
      run > /dev/null
    ;;
  run)
    echo "==> sysbench run (time=${RUN_TIME}s)"
    sysbench "${SYSBENCH_COMMON[@]}" \
      --time="${RUN_TIME}" \
      --report-interval=5 \
      --verbosity=0 \
      run
    ;;
  all)
    echo "==> sysbench prepare (workload=${SYSBENCH_WORKLOAD}, threads=${THREADS}, tables=${TABLES})"
    sysbench "${SYSBENCH_COMMON[@]}" prepare

    echo "==> sysbench warmup (time=${WARMUP_TIME}s)"
    sysbench "${SYSBENCH_COMMON[@]}" \
      --time="${WARMUP_TIME}" \
      run > /dev/null

    echo "==> sysbench run (time=${RUN_TIME}s)"
    sysbench "${SYSBENCH_COMMON[@]}" \
      --time="${RUN_TIME}" \
      --report-interval=5 \
      --verbosity=0 \
      run
    ;;
  *)
    echo "Unknown phase: ${PHASE}" >&2
    echo "Usage: entrypoint.sh [all|prepare-and-warmup|run]" >&2
    exit 1
    ;;
esac
