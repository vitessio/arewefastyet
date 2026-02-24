#!/usr/bin/env bash
set -euo pipefail

TABLET_ID="${TABLET_ID:?TABLET_ID must be set}"

MYSQL_SOCKET="/vt/socket/mysql${TABLET_ID}.sock"
MYSQLCTLD_SOCKET="/vt/socket/mysqlctld${TABLET_ID}.sock"

mkdir -p /vt/socket
chown -R vitess:vitess /vt

echo "Starting mysqlctld for tablet ${TABLET_ID}..."
exec gosu vitess mysqlctld \
  --tablet-uid="${TABLET_ID}" \
  --mysql-socket="${MYSQL_SOCKET}" \
  --socket-file="${MYSQLCTLD_SOCKET}" \
  --db-charset=utf8mb4
