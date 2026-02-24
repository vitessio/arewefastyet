#!/usr/bin/env bash
set -euo pipefail

# SHARDS: space-separated list of "shard:tablet_id" pairs.
# Examples: "0:1001" (unsharded), "-80:1001 80-:2001" (sharded)
SHARDS="${SHARDS:-0:1001}"
VSCHEMA_FILE="${VSCHEMA_FILE:-}"

VTC="gosu vitess vtctldclient --server=vtctld:15999"

echo "Waiting for vtctld to be ready..."
for i in $(seq 1 120); do
  if ${VTC} GetKeyspaces >/dev/null 2>&1; then
    echo "vtctld is ready"
    break
  fi
  if [ "$i" -eq 120 ]; then
    echo "ERROR: vtctld not ready after 120s"
    exit 1
  fi
  sleep 1
done

echo "Waiting for keyspace 'main' to exist..."
for i in $(seq 1 60); do
  if ${VTC} GetKeyspace main >/dev/null 2>&1; then
    echo "Keyspace 'main' exists"
    break
  fi
  if [ "$i" -eq 60 ]; then
    echo "ERROR: keyspace 'main' not found after 60s"
    exit 1
  fi
  sleep 1
done

for entry in ${SHARDS}; do
  shard="${entry%%:*}"
  tablet_id="${entry##*:}"
  tablet_alias="local-$(printf '%010d' "${tablet_id}")"

  echo "Waiting for tablet ${tablet_alias} to be registered..."
  for i in $(seq 1 120); do
    if ${VTC} GetTablet "${tablet_alias}" >/dev/null 2>&1; then
      echo "Tablet ${tablet_alias} is registered"
      break
    fi
    if [ "$i" -eq 120 ]; then
      echo "ERROR: tablet ${tablet_alias} not registered after 120s"
      exit 1
    fi
    sleep 1
  done

  echo "Running InitShardPrimary for main/${shard} (tablet=${tablet_alias})..."
  ${VTC} InitShardPrimary \
    --force \
    "main/${shard}" \
    "${tablet_alias}"
done

echo "Applying VSchema..."
if [ -n "${VSCHEMA_FILE}" ] && [ -f "${VSCHEMA_FILE}" ]; then
  ${VTC} ApplyVSchema \
    --vschema="$(cat "${VSCHEMA_FILE}")" \
    main
else
  ${VTC} ApplyVSchema \
    --vschema='{"sharded":false}' \
    main
fi

echo "vtctld-init complete"
