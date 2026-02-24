#!/usr/bin/env bash
set -euo pipefail

TABLET_ID="${TABLET_ID:?TABLET_ID must be set}"
MYSQLD_EXPORTER_PORT="${MYSQLD_EXPORTER_PORT:-9104}"

MYSQL_SOCKET="/vt/socket/mysql${TABLET_ID}.sock"

echo "Waiting for MySQL socket..."
while [ ! -S "${MYSQL_SOCKET}" ]; do sleep 2; done

echo "Starting mysqld_exporter on port ${MYSQLD_EXPORTER_PORT}..."
exec gosu vitess mysqld_exporter \
  --web.listen-address=":${MYSQLD_EXPORTER_PORT}" \
  --no-collect.slave_status
