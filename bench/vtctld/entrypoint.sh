#!/usr/bin/env bash
set -euo pipefail

echo "Waiting for etcd to be ready..."
for i in $(seq 1 60); do
  if etcdctl --endpoints=http://etcd:2379 endpoint health 2>/dev/null; then
    echo "etcd is ready"
    break
  fi
  if [ "$i" -eq 60 ]; then
    echo "ERROR: etcd not ready after 60s"
    exit 1
  fi
  sleep 1
done

echo "Adding cell info..."
# Runs directly against etcd (no vtctld needed). Idempotent — ignore error if already exists.
gosu vitess vtctl \
  --topo-implementation=etcd2 \
  --topo-global-server-address=etcd:2379 \
  --topo-global-root=/vitess/main \
  AddCellInfo \
  --root=/vitess/main/local \
  --server_address=etcd:2379 \
  local || true

echo "Starting vtctld..."
exec gosu vitess vtctld \
  --cell=local \
  --service-map=grpc-vtctl,grpc-vtctld \
  --port=15000 \
  --grpc-port=15999 \
  --topo-implementation=etcd2 \
  --topo-global-server-address=etcd:2379 \
  --topo-global-root=/vitess/main
