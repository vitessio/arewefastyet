#!/usr/bin/env bash
set -euo pipefail

TABLET_ID="${TABLET_ID:?TABLET_ID must be set}"
TABLET_PORT="${TABLET_PORT:?TABLET_PORT must be set}"
GRPC_PORT="${GRPC_PORT:?GRPC_PORT must be set}"
SHARD="${SHARD:-0}"
VTTABLET_EXTRA_FLAGS="${VTTABLET_EXTRA_FLAGS:-}"

MYSQL_SOCKET="/vt/socket/mysql${TABLET_ID}.sock"
MYSQLCTLD_SOCKET="/vt/socket/mysqlctld${TABLET_ID}.sock"

echo "Starting vttablet for tablet ${TABLET_ID}..."
exec gosu vitess vttablet \
  --tablet-hostname="${HOSTNAME}" \
  --tablet-path="local-${TABLET_ID}" \
  --init-keyspace=main \
  --init-shard="${SHARD}" \
  --init-db-name-override=main \
  --init-tablet-type=replica \
  --port="${TABLET_PORT}" \
  --grpc-port="${GRPC_PORT}" \
  --mysqlctl-socket="${MYSQLCTLD_SOCKET}" \
  --db-socket="${MYSQL_SOCKET}" \
  --topo-implementation=etcd2 \
  --topo-global-server-address=etcd:2379 \
  --topo-global-root=/vitess/main \
  --queryserver-config-pool-size=500 \
  --queryserver-config-transaction-cap=2000 \
  --queryserver-config-query-timeout=300s \
  --service-map=grpc-queryservice,grpc-tabletmanager,grpc-updatestream \
  --pprof-http \
  ${VTTABLET_EXTRA_FLAGS}
