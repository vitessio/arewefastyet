#!/usr/bin/env bash

# shellcheck disable=SC2046
supervisord_file="$(dirname $(realpath $0))/supervisord.conf"

cp $supervisord_file /etc/supervisor/supervisord.conf

supervisorctl reload