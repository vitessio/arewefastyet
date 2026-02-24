#!/usr/bin/env bash
# generate.sh - Generate self-contained benchmark directories from benchmarks.yml.
#
# Usage:
#   ./bench/generate.sh [options] [benchmark-name]
#
# Options:
#   --ref <git-ref>   Vitess git ref to benchmark (default: main)
#   --all             Generate all benchmarks
#   -h, --help        Show this help message
#
# Examples:
#   ./bench/generate.sh --all                     # Generate all benchmarks for main
#   ./bench/generate.sh oltp --ref v21.0          # Generate oltp for a specific ref
#   ./bench/generate.sh --all --ref my-branch     # Generate all for a branch

set -euo pipefail

VITESS_REF="main"
BENCHMARK=""
GENERATE_ALL=false

usage() {
  head -16 "$0" | grep '^#' | sed 's/^# \{0,1\}//'
  exit 0
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --ref)     VITESS_REF="$2"; shift 2 ;;
    --all)     GENERATE_ALL=true; shift ;;
    -h|--help) usage ;;
    *)
      if [ -z "${BENCHMARK}" ]; then
        BENCHMARK="$1"; shift
      else
        echo "Unknown option: $1"; usage
      fi
      ;;
  esac
done

if ! $GENERATE_ALL && [ -z "${BENCHMARK}" ]; then
  echo "Error: specify a benchmark name or --all"
  usage
fi

# Resolve paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG="${SCRIPT_DIR}/benchmarks.yml"
OUTPUT_DIR="${SCRIPT_DIR}/benchmarks"

if ! command -v yq &>/dev/null; then
  echo "Error: yq is required. Install from https://github.com/mikefarah/yq"
  exit 1
fi

# Get list of benchmarks to generate
if $GENERATE_ALL; then
  BENCHMARKS=$(yq '.benchmarks | keys | .[]' "${CONFIG}")
else
  # Validate the benchmark exists
  if ! yq -e ".benchmarks.\"${BENCHMARK}\"" "${CONFIG}" &>/dev/null; then
    echo "Error: unknown benchmark '${BENCHMARK}'"
    echo "Available benchmarks:"
    yq '.benchmarks | keys | .[]' "${CONFIG}" | sed 's/^/  /'
    exit 1
  fi
  BENCHMARKS="${BENCHMARK}"
fi

# --- Helper: read a field from benchmarks.yml ---
cfg() {
  local bench="$1" path="$2"
  yq ".benchmarks.\"${bench}\".${path}" "${CONFIG}"
}

# --- Generate one benchmark ---
generate_benchmark() {
  local name="$1"
  local dir="${OUTPUT_DIR}/${name}"

  echo "Generating ${name}..."

  # Read config
  local sharded
  sharded=$(cfg "${name}" "vitess.sharded")
  local vschema
  vschema=$(cfg "${name}" "vitess.vschema")
  local vtgate_workload_mode
  vtgate_workload_mode=$(cfg "${name}" "vitess.vtgate_workload_mode")
  # yq returns "null" for empty strings in some cases
  [ "${vtgate_workload_mode}" = "null" ] && vtgate_workload_mode=""
  local vtgate_max_memory_rows
  vtgate_max_memory_rows=$(cfg "${name}" "vitess.vtgate_max_memory_rows")
  [ "${vtgate_max_memory_rows}" = "null" ] && vtgate_max_memory_rows=""
  local vttablet_extra_flags
  vttablet_extra_flags=$(cfg "${name}" "vitess.vttablet_extra_flags")
  [ "${vttablet_extra_flags}" = "null" ] && vttablet_extra_flags=""

  local sb_workload sb_tables sb_table_size sb_scale sb_threads sb_extra sb_working_dir sb_warmup sb_runtime
  sb_workload=$(cfg "${name}" "sysbench.workload")
  sb_tables=$(cfg "${name}" "sysbench.tables")
  sb_table_size=$(cfg "${name}" "sysbench.table_size")
  sb_scale=$(cfg "${name}" "sysbench.scale")
  sb_threads=$(cfg "${name}" "sysbench.threads")
  sb_extra=$(cfg "${name}" "sysbench.extra")
  sb_working_dir=$(cfg "${name}" "sysbench.working_dir")
  sb_warmup=$(cfg "${name}" "sysbench.warmup_time")
  sb_runtime=$(cfg "${name}" "sysbench.run_time")

  # Normalize nulls
  [ "${sb_table_size}" = "null" ] && sb_table_size=""
  [ "${sb_scale}" = "null" ] && sb_scale=""
  [ "${sb_extra}" = "null" ] && sb_extra=""
  [ "${sb_working_dir}" = "null" ] && sb_working_dir=""
  [ "${vtgate_workload_mode}" = "\"\"" ] && vtgate_workload_mode=""
  [ "${sb_extra}" = "\"\"" ] && sb_extra=""
  [ "${sb_working_dir}" = "\"\"" ] && sb_working_dir=""

  # vttablet extra flags env (computed early since shard2_services references it)
  local vttablet_env=""
  if [ -n "${vttablet_extra_flags}" ]; then
    vttablet_env="
      VTTABLET_EXTRA_FLAGS: \"${vttablet_extra_flags}\""
  fi

  # Resolve sharding
  local shard_1 shard_2 shards vschema_file sharding_label
  if [ "${sharded}" = "true" ]; then
    shard_1="-80"
    shard_2="80-"
    shards="-80:1001 80-:2001"
    vschema_file="/vschema/${vschema}"
    sharding_label="sharded"
  else
    shard_1="0"
    shard_2=""
    shards="0:1001"
    vschema_file="/vschema/${vschema}"
    sharding_label="unsharded"
  fi

  # Create output directory
  mkdir -p "${dir}/prometheus"

  # --- Generate Prometheus agent configs ---
  cat > "${dir}/prometheus/vtgate.yml" <<PROMEOF
global:
  scrape_interval: 15s
  external_labels:
    workload: "${name}"
    sharding: "${sharding_label}"
    vitess_ref: "${VITESS_REF}"
scrape_configs:
  - job_name: vtgate
    static_configs:
      - targets: ["vtgate:15001"]
remote_write:
  - url: http://prometheus:9090/api/v1/write
PROMEOF

  cat > "${dir}/prometheus/vttablet-1001.yml" <<PROMEOF
global:
  scrape_interval: 15s
  external_labels:
    workload: "${name}"
    sharding: "${sharding_label}"
    vitess_ref: "${VITESS_REF}"
scrape_configs:
  - job_name: vttablet-1001
    static_configs:
      - targets: ["mysql-1001:16011", "mysql-1001:9104"]
remote_write:
  - url: http://prometheus:9090/api/v1/write
PROMEOF

  if [ "${sharded}" = "true" ]; then
    cat > "${dir}/prometheus/vttablet-2001.yml" <<PROMEOF
global:
  scrape_interval: 15s
  external_labels:
    workload: "${name}"
    sharding: "${sharding_label}"
    vitess_ref: "${VITESS_REF}"
scrape_configs:
  - job_name: vttablet-2001
    static_configs:
      - targets: ["mysql-2001:16021", "mysql-2001:9104"]
remote_write:
  - url: http://prometheus:9090/api/v1/write
PROMEOF
  fi

  # --- Generate docker-compose.yml ---
  # Build the shard-2 services block conditionally
  local shard2_services=""
  local shard2_volumes=""
  local shard2_alloy=""
  local shard2_prom=""
  local vtctld_init_tablet_2001_dep=""

  if [ "${sharded}" = "true" ]; then
    shard2_services="
  mysql-2001:
    <<: *vitess-image
    hostname: mysql-2001
    cpus: 4
    depends_on:
      vtctld:
        condition: service_healthy
    environment:
      TABLET_ID: \"2001\"
    volumes:
      - ./mysql/entrypoint.sh:/entrypoint.sh:ro
      - tablet-2001-data:/vt
      - tablet-2001-socket:/vt/socket
    entrypoint: [\"/bin/bash\", \"/entrypoint.sh\"]
    ports:
      - \"16021:16021\"
    networks:
      - bench
    security_opt:
      - apparmor=unconfined
    healthcheck:
      test: [\"CMD-SHELL\", \"test -S /vt/socket/mysql2001.sock\"]
      interval: 5s
      timeout: 5s
      retries: 60

  mysqld-exporter-2001:
    <<: *vitess-image
    network_mode: \"service:mysql-2001\"
    depends_on:
      mysql-2001:
        condition: service_healthy
    environment:
      TABLET_ID: \"2001\"
      MYSQLD_EXPORTER_PORT: \"9104\"
      DATA_SOURCE_NAME: \"root@unix(/vt/socket/mysql2001.sock)/\"
    volumes:
      - ./mysqld-exporter/entrypoint.sh:/entrypoint.sh:ro
      - tablet-2001-socket:/vt/socket
    entrypoint: [\"/bin/bash\", \"/entrypoint.sh\"]

  tablet-2001:
    <<: *vitess-image
    cpus: 4
    network_mode: \"service:mysql-2001\"
    depends_on:
      mysql-2001:
        condition: service_healthy
    environment:
      <<: *topo-env
      TABLET_ID: \"2001\"
      TABLET_PORT: \"16021\"
      GRPC_PORT: \"17021\"
      SHARD: \"${shard_2}\"${vttablet_env}
    volumes:
      - ./tablet/entrypoint.sh:/entrypoint.sh:ro
      - tablet-2001-socket:/vt/socket
    entrypoint: [\"/bin/bash\", \"/entrypoint.sh\"]
    healthcheck:
      test: [\"CMD\", \"curl\", \"-sf\", \"http://localhost:16021/debug/vars\"]
      interval: 5s
      timeout: 5s
      retries: 60"

    shard2_volumes="
  tablet-2001-data:
  tablet-2001-socket:"

    shard2_alloy="
  alloy-vttablet-2001:
    image: grafana/alloy:latest
    network_mode: \"service:mysql-2001\"
    depends_on:
      mysql-2001:
        condition: service_healthy
    environment:
      ALLOY_TARGET: \"localhost:16021\"
      ALLOY_SERVICE_NAME: \"vttablet\"
      ALLOY_WORKLOAD: \"${name}\"
      ALLOY_SHARDING: \"${sharding_label}\"
      ALLOY_VITESS_REF: \"${VITESS_REF}\"
    volumes:
      - ./alloy/vttablet.alloy:/etc/alloy/config.alloy:ro
    command: run /etc/alloy/config.alloy --server.http.listen-addr=0.0.0.0:12346"

    shard2_prom="
  prom-vttablet-2001:
    image: prom/prometheus:v2.51.2
    network_mode: \"service:mysql-2001\"
    depends_on:
      mysql-2001:
        condition: service_healthy
    volumes:
      - ./prometheus/vttablet-2001.yml:/etc/prometheus/prometheus.yml:ro
    command:
      - \"--enable-feature=agent\"
      - \"--config.file=/etc/prometheus/prometheus.yml\""

    vtctld_init_tablet_2001_dep="
      tablet-2001:
        condition: service_healthy"
  fi

  # Sysbench environment block
  local sysbench_env=""
  sysbench_env="      SYSBENCH_WORKLOAD: \"${sb_workload}\""
  [ -n "${sb_working_dir}" ] && sysbench_env="${sysbench_env}
      SYSBENCH_WORKING_DIR: \"${sb_working_dir}\""
  [ -n "${sb_extra}" ] && sysbench_env="${sysbench_env}
      SYSBENCH_EXTRA: \"${sb_extra}\""
  sysbench_env="${sysbench_env}
      THREADS: \"${sb_threads}\"
      TABLES: \"${sb_tables}\""
  if [ -n "${sb_scale}" ]; then
    sysbench_env="${sysbench_env}
      SCALE: \"${sb_scale}\""
  fi
  if [ -n "${sb_table_size}" ]; then
    sysbench_env="${sysbench_env}
      TABLE_SIZE: \"${sb_table_size}\""
  fi
  sysbench_env="${sysbench_env}
      WARMUP_TIME: \"${sb_warmup}\"
      RUN_TIME: \"${sb_runtime}\""

  # vtgate extra env
  local vtgate_env=""
  if [ -n "${vtgate_workload_mode}" ]; then
    vtgate_env="
      VTGATE_WORKLOAD_MODE: \"${vtgate_workload_mode}\""
  fi
  if [ -n "${vtgate_max_memory_rows}" ]; then
    vtgate_env="${vtgate_env}
      VTGATE_MAX_MEMORY_ROWS: \"${vtgate_max_memory_rows}\""
  fi

  # Prometheus agent volume mount paths are relative to --project-directory (bench/)
  # but we put the config in the benchmark subdir so use the full relative path
  local prom_vtgate_vol="./benchmarks/${name}/prometheus/vtgate.yml:/etc/prometheus/prometheus.yml:ro"
  local prom_vttablet1_vol="./benchmarks/${name}/prometheus/vttablet-1001.yml:/etc/prometheus/prometheus.yml:ro"
  local prom_vttablet2_vol="./benchmarks/${name}/prometheus/vttablet-2001.yml:/etc/prometheus/prometheus.yml:ro"

  # Update shard2_prom volume path
  if [ "${sharded}" = "true" ]; then
    shard2_prom="
  prom-vttablet-2001:
    image: prom/prometheus:v2.51.2
    network_mode: \"service:mysql-2001\"
    depends_on:
      mysql-2001:
        condition: service_healthy
    volumes:
      - ${prom_vttablet2_vol}
    command:
      - \"--enable-feature=agent\"
      - \"--config.file=/etc/prometheus/prometheus.yml\""
  fi

  cat > "${dir}/docker-compose.yml" <<EOF
# Auto-generated by generate.sh — do not edit manually.
# Benchmark: ${name} | Sharding: ${sharding_label} | Ref: ${VITESS_REF}

x-vitess-image: &vitess-image
  image: vitess-bench:${VITESS_REF}
  build:
    context: .
    dockerfile: vitess.Dockerfile
    args:
      VITESS_REF: ${VITESS_REF}

x-topo-env: &topo-env
  ETCDCTL_API: "3"

services:
  etcd:
    image: quay.io/coreos/etcd:v3.5.27
    environment:
      ETCD_LISTEN_CLIENT_URLS: "http://0.0.0.0:2379"
      ETCD_ADVERTISE_CLIENT_URLS: "http://etcd:2379"
      ETCD_DATA_DIR: /etcd-data
      ETCDCTL_API: "3"
    volumes:
      - etcd-data:/etcd-data
    networks:
      - bench
    healthcheck:
      test: ["CMD", "etcdctl", "--endpoints=http://localhost:2379", "endpoint", "health"]
      interval: 5s
      timeout: 5s
      retries: 12

  vtctld:
    <<: *vitess-image
    depends_on:
      etcd:
        condition: service_healthy
    environment:
      <<: *topo-env
    volumes:
      - ./vtctld/entrypoint.sh:/entrypoint.sh:ro
    entrypoint: ["/bin/bash", "/entrypoint.sh"]
    ports:
      - "15000:15000"
    networks:
      - bench
    healthcheck:
      test: ["CMD", "curl", "-sf", "http://localhost:15000/debug/health"]
      interval: 5s
      timeout: 5s
      retries: 24

  mysql-1001:
    <<: *vitess-image
    hostname: mysql-1001
    cpus: 4
    depends_on:
      vtctld:
        condition: service_healthy
    environment:
      TABLET_ID: "1001"
    volumes:
      - ./mysql/entrypoint.sh:/entrypoint.sh:ro
      - tablet-1001-data:/vt
      - tablet-1001-socket:/vt/socket
    entrypoint: ["/bin/bash", "/entrypoint.sh"]
    ports:
      - "9104:9104"
      - "16011:16011"
    networks:
      - bench
    security_opt:
      - apparmor=unconfined
    healthcheck:
      test: ["CMD-SHELL", "test -S /vt/socket/mysql1001.sock"]
      interval: 5s
      timeout: 5s
      retries: 60

  mysqld-exporter-1001:
    <<: *vitess-image
    network_mode: "service:mysql-1001"
    depends_on:
      mysql-1001:
        condition: service_healthy
    environment:
      TABLET_ID: "1001"
      MYSQLD_EXPORTER_PORT: "9104"
      DATA_SOURCE_NAME: "root@unix(/vt/socket/mysql1001.sock)/"
    volumes:
      - ./mysqld-exporter/entrypoint.sh:/entrypoint.sh:ro
      - tablet-1001-socket:/vt/socket
    entrypoint: ["/bin/bash", "/entrypoint.sh"]

  tablet-1001:
    <<: *vitess-image
    cpus: 4
    network_mode: "service:mysql-1001"
    depends_on:
      mysql-1001:
        condition: service_healthy
    environment:
      <<: *topo-env
      TABLET_ID: "1001"
      TABLET_PORT: "16011"
      GRPC_PORT: "17011"
      SHARD: "${shard_1}"${vttablet_env}
    volumes:
      - ./tablet/entrypoint.sh:/entrypoint.sh:ro
      - tablet-1001-socket:/vt/socket
    entrypoint: ["/bin/bash", "/entrypoint.sh"]
    healthcheck:
      test: ["CMD", "curl", "-sf", "http://localhost:16011/debug/vars"]
      interval: 5s
      timeout: 5s
      retries: 60
${shard2_services}
  vtctld-init:
    <<: *vitess-image
    depends_on:
      vtctld:
        condition: service_healthy
      tablet-1001:
        condition: service_healthy${vtctld_init_tablet_2001_dep}
    environment:
      <<: *topo-env
      SHARDS: "${shards}"
      VSCHEMA_FILE: "${vschema_file}"
    volumes:
      - ./vtctld-init/entrypoint.sh:/entrypoint.sh:ro
      - ./vschema:/vschema:ro
    entrypoint: ["/bin/bash", "/entrypoint.sh"]
    networks:
      - bench
    restart: "no"

  vtgate:
    <<: *vitess-image
    cpus: 4
    depends_on:
      vtctld-init:
        condition: service_completed_successfully
    environment:
      <<: *topo-env${vtgate_env}
    volumes:
      - ./vtgate/entrypoint.sh:/entrypoint.sh:ro
    entrypoint: ["/bin/bash", "/entrypoint.sh"]
    ports:
      - "13306:13306"
      - "15001:15001"
      - "15306:15306"
    networks:
      - bench
    healthcheck:
      test: ["CMD", "curl", "-sf", "http://localhost:15001/debug/health"]
      interval: 5s
      timeout: 5s
      retries: 24

  sysbench:
    image: sysbench-bench:latest
    build:
      context: .
      dockerfile: sysbench.Dockerfile
    depends_on:
      vtgate:
        condition: service_healthy
    environment:
${sysbench_env}
    networks:
      - bench
    profiles:
      - bench

  alloy-vtgate:
    image: grafana/alloy:latest
    environment:
      ALLOY_WORKLOAD: "${name}"
      ALLOY_SHARDING: "${sharding_label}"
      ALLOY_VITESS_REF: "${VITESS_REF}"
    volumes:
      - ./alloy/vtgate.alloy:/etc/alloy/config.alloy:ro
    command: run /etc/alloy/config.alloy --server.http.listen-addr=0.0.0.0:12345
    networks:
      - bench

  alloy-vttablet-1001:
    image: grafana/alloy:latest
    network_mode: "service:mysql-1001"
    depends_on:
      mysql-1001:
        condition: service_healthy
    environment:
      ALLOY_TARGET: "localhost:16011"
      ALLOY_SERVICE_NAME: "vttablet"
      ALLOY_WORKLOAD: "${name}"
      ALLOY_SHARDING: "${sharding_label}"
      ALLOY_VITESS_REF: "${VITESS_REF}"
    volumes:
      - ./alloy/vttablet.alloy:/etc/alloy/config.alloy:ro
    command: run /etc/alloy/config.alloy --server.http.listen-addr=0.0.0.0:12346
${shard2_alloy}
  prom-vtgate:
    image: prom/prometheus:v2.51.2
    volumes:
      - ${prom_vtgate_vol}
    command:
      - "--enable-feature=agent"
      - "--config.file=/etc/prometheus/prometheus.yml"
    networks:
      - bench

  prom-vttablet-1001:
    image: prom/prometheus:v2.51.2
    network_mode: "service:mysql-1001"
    depends_on:
      mysql-1001:
        condition: service_healthy
    volumes:
      - ${prom_vttablet1_vol}
    command:
      - "--enable-feature=agent"
      - "--config.file=/etc/prometheus/prometheus.yml"
${shard2_prom}
volumes:
  etcd-data:
  tablet-1001-data:
  tablet-1001-socket:${shard2_volumes}

networks:
  bench:
    external: true
EOF

  # --- Generate run.sh ---
  local shard_services="mysql-1001 tablet-1001 mysqld-exporter-1001"
  local alloy_services="alloy-vtgate alloy-vttablet-1001"
  local prom_services="prom-vtgate prom-vttablet-1001"

  if [ "${sharded}" = "true" ]; then
    shard_services="${shard_services} mysql-2001 tablet-2001 mysqld-exporter-2001"
    alloy_services="${alloy_services} alloy-vttablet-2001"
    prom_services="${prom_services} prom-vttablet-2001"
  fi

  cat > "${dir}/run.sh" <<RUNEOF
#!/usr/bin/env bash
# Auto-generated by generate.sh — do not edit manually.
# Benchmark: ${name} | Sharding: ${sharding_label} | Ref: ${VITESS_REF}
set -euo pipefail

SCRIPT_DIR="\$(cd "\$(dirname "\${BASH_SOURCE[0]}")" && pwd)"
BENCH_DIR="\$(cd "\${SCRIPT_DIR}/../.." && pwd)"

DC="docker compose --project-directory \${BENCH_DIR} -f \${SCRIPT_DIR}/docker-compose.yml -p bench"
MONITOR_DC="docker compose -f \${BENCH_DIR}/docker-compose.monitor.yml -p bench-monitor"

TEARDOWN=true
PPROF=false
for arg in "\$@"; do
  case "\$arg" in
    --no-teardown) TEARDOWN=false ;;
    --pprof)       PPROF=true ;;
  esac
done

docker network create bench 2>/dev/null || true

echo "==> Starting monitoring stack (Prometheus, Grafana, Pyroscope)..."
\${MONITOR_DC} up -d prometheus grafana pyroscope

echo ""
echo "==> Building Vitess image (ref=${VITESS_REF})..."
echo "    (This takes 20-30 min on first run; subsequent runs with same ref are instant.)"
\${DC} build vtctld sysbench

echo ""
echo "==> Bringing up cluster (etcd, vtctld, tablets, vtgate)..."
echo "    Benchmark: ${name}"
echo "    Sharding:  ${sharding_label}"
\${DC} up -d --wait etcd vtctld ${shard_services} vtctld-init vtgate

echo ""
echo "    Cluster is ready."
echo "    Prometheus: http://localhost:9090"
echo "    Grafana:    http://localhost:3000"
echo "    Pyroscope:  http://localhost:4040"
echo ""

echo "==> Preparing and warming up..."
\${DC} run --rm sysbench prepare-and-warmup

echo ""
echo "==> Starting sidecars (Alloy for profiling, Prometheus agents for metrics)..."
\${DC} up -d ${alloy_services} ${prom_services}

PPROF_PID=""
PROFILE_DIR=""
if \$PPROF; then
  PROFILE_DIR="\${SCRIPT_DIR}/profiles/\$(date +%Y%m%d_%H%M%S)"
  mkdir -p "\${PROFILE_DIR}"
  echo "==> pprof: collecting ${sb_runtime}s CPU profiles and allocation snapshots"
  echo "    Output: \${PROFILE_DIR}/"
  echo ""
  (
    curl -sf "http://localhost:15001/debug/pprof/profile?seconds=${sb_runtime}" \\
      -o "\${PROFILE_DIR}/vtgate-cpu.pprof" 2>/dev/null &
    curl -sf "http://localhost:16011/debug/pprof/profile?seconds=${sb_runtime}" \\
      -o "\${PROFILE_DIR}/vttablet-1001-cpu.pprof" 2>/dev/null &
RUNEOF

  if [ "${sharded}" = "true" ]; then
    cat >> "${dir}/run.sh" <<RUNEOF
    curl -sf "http://localhost:16021/debug/pprof/profile?seconds=${sb_runtime}" \\
      -o "\${PROFILE_DIR}/vttablet-2001-cpu.pprof" 2>/dev/null &
RUNEOF
  fi

  cat >> "${dir}/run.sh" <<RUNEOF
    wait
  ) &
  PPROF_PID=\$!
fi

echo "==> Running benchmark (${name}, threads=${sb_threads}, tables=${sb_tables}, time=${sb_runtime}s)..."
\${DC} run --rm sysbench run

echo ""
echo "==> Benchmark complete."

\${DC} stop ${alloy_services} ${prom_services}

if [ -n "\$PPROF_PID" ]; then
  echo "==> Waiting for CPU profile collection to finish..."
  wait "\$PPROF_PID" 2>/dev/null || true

  echo "==> Collecting allocation profiles..."
  curl -sf "http://localhost:15001/debug/pprof/heap" \\
    -o "\${PROFILE_DIR}/vtgate-heap.pprof" 2>/dev/null || true
  curl -sf "http://localhost:16011/debug/pprof/heap" \\
    -o "\${PROFILE_DIR}/vttablet-1001-heap.pprof" 2>/dev/null || true
RUNEOF

  if [ "${sharded}" = "true" ]; then
    cat >> "${dir}/run.sh" <<RUNEOF
  curl -sf "http://localhost:16021/debug/pprof/heap" \\
    -o "\${PROFILE_DIR}/vttablet-2001-heap.pprof" 2>/dev/null || true
RUNEOF
  fi

  cat >> "${dir}/run.sh" <<RUNEOF

  echo ""
  echo "    Profiles saved to: \${PROFILE_DIR}/"
  echo "    View with: go tool pprof -http=:8888 <profile.pprof>"
fi

echo ""
if \$TEARDOWN; then
  echo "==> Tearing down benchmark cluster..."
  \${DC} down -v
  echo "    Done. Monitoring stack is still running."
  echo "    To tear down monitoring: \${MONITOR_DC} down -v"
else
  echo "==> Cluster left running (--no-teardown). To shut down:"
  echo "    \${DC} down -v                  # benchmark only"
  echo "    \${MONITOR_DC} down -v   # monitoring"
fi
RUNEOF

  chmod +x "${dir}/run.sh"
  echo "  -> ${dir}/"
}

# --- Main ---
echo "Generating benchmarks (ref=${VITESS_REF})..."
echo ""

for bench in ${BENCHMARKS}; do
  generate_benchmark "${bench}"
done

echo ""
echo "Done. Run a benchmark with:"
echo "  ./bench/benchmarks/<name>/run.sh"
echo ""
echo "Available benchmarks:"
ls -1 "${OUTPUT_DIR}" 2>/dev/null | sed 's/^/  /'
