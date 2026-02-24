#!/usr/bin/env bash
# collect_profiles.sh - Collect pprof CPU and allocation profiles from vtgate and vttablets.
#
# Run this WHILE the benchmark is running to capture CPU profiles.
#
# Usage:
#   ./bench/collect_profiles.sh [duration_seconds]
#
# Example:
#   ./bench/collect_profiles.sh 30
#
# After collection:
#   go tool pprof -http=:8888 profiles/<timestamp>/vtgate-cpu.pprof
#   go tool pprof -http=:8889 profiles/<timestamp>/vttablet-1001-cpu.pprof

set -euo pipefail

DURATION="${1:-30}"
TIMESTAMP="$(date +%Y%m%d_%H%M%S)"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUT="${SCRIPT_DIR}/profiles/${TIMESTAMP}"
mkdir -p "${OUT}"

echo "==> Collecting CPU profiles (${DURATION}s) and allocation snapshots from vtgate and vttablet..."
echo "    Output directory: ${OUT}"
echo ""

# Collect concurrently
curl -s "http://localhost:15001/debug/pprof/profile?seconds=${DURATION}" \
  -o "${OUT}/vtgate-cpu.pprof" &
PID_VTGATE=$!

curl -s "http://localhost:16011/debug/pprof/profile?seconds=${DURATION}" \
  -o "${OUT}/vttablet-1001-cpu.pprof" &
PID_TAB1=$!

# Wait for CPU profiles to finish
wait $PID_VTGATE && echo "  vtgate CPU profile saved" || echo "  WARNING: vtgate CPU profile collection failed"
wait $PID_TAB1   && echo "  vttablet-1001 CPU profile saved" || echo "  WARNING: vttablet-1001 CPU profile collection failed"

# Collect allocation snapshots (instant)
curl -sf "http://localhost:15001/debug/pprof/heap" \
  -o "${OUT}/vtgate-heap.pprof" \
  && echo "  vtgate heap profile saved" || echo "  WARNING: vtgate heap profile collection failed"
curl -sf "http://localhost:16011/debug/pprof/heap" \
  -o "${OUT}/vttablet-1001-heap.pprof" \
  && echo "  vttablet-1001 heap profile saved" || echo "  WARNING: vttablet-1001 heap profile collection failed"

echo ""
echo "==> Profiles saved to: ${OUT}/"
echo ""
echo "View CPU profiles:"
echo "  go tool pprof -http=:8888 ${OUT}/vtgate-cpu.pprof"
echo "  go tool pprof -http=:8889 ${OUT}/vttablet-1001-cpu.pprof"
echo ""
echo "View allocation profiles:"
echo "  go tool pprof -http=:8890 ${OUT}/vtgate-heap.pprof"
echo "  go tool pprof -http=:8891 ${OUT}/vttablet-1001-heap.pprof"
