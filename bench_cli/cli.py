#!/usr/bin/env python

# ------------------------------------------------------------------------------------------------------------------------------------
# Copyright 2021 The Vitess Authors.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#    http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# demonstrates to:
#   - Creates Cli to interact with the benchmarks
#      1. --runall: runs OLTP and TPCC
#      2. --commit: specify commit hash or branch name
#      3. --run-tpcc: Runs TPCC
#      4. --run-oltp: Runs OLTP
#      5. --inventory-file: Set inventory file to run
# -------------------------------------------------------------------------------------------------------------------------------------

import click
import uuid
import configuration
import run_benchmark
import sys

@click.command()
@click.option("--run-all",                  is_flag=True, help="runs OLTP and TPCC")
@click.option("--run-tpcc", "-runt",        is_flag=True, help="Runs TPCC")
@click.option("--run-oltp", "-runo",        is_flag=True, help="Runs OLTP")
@click.option("--commit", "-c",             help="Specify commit hash or branch name ")
@click.option("--source", "-s",             help="Mention the source from where the cli is called")
@click.option("--tasks-scripts-dir",        help="Path to tasks scripts directory", envvar="BCLI_TASKS_SCRIPTS_DIR")
@click.option("--tasks-reports-dir",        help="Path to tasks reports directory", envvar="BCLI_TASKS_REPORTS_DIR")
@click.option("--ansible-dir",              help="Path to the Ansible directory", envvar="BCLI_ANSIBLE_DIR")
@click.option("--inventory-file", "-invf",  help="Mention inventory file to call", envvar="BCLI_INVENTORY_FILE")
@click.option("--mysql-host",               help="MySQL server hostname", envvar="BCLI_MYSQL_HOST")
@click.option("--mysql-username",           help="MySQL server username", envvar="BCLI_MYSQL_USER")
@click.option("--mysql-password",           help="MySQL server password", envvar="BCLI_MYSQL_PASSWORD")
@click.option("--mysql-database",           help="MySQL database to use", envvar="BCLI_MYSQL_DB")
@click.option("--packet-token",             help="Token used to authenticate Packet", envvar="BCLI_PACKET_TOKEN")
@click.option("--packet-project-id",        help="Packet project ID", envvar="BCLI_PACKET_PROJECT_ID")
@click.option("--api-key",                  help="API key", envvar="BCLI_API_KEY")
@click.option("--slack-api-token",          help="Slack API token", envvar="BCLI_SLACK_TOKEN")
@click.option("--slack-channel",            help="Slack channel", envvar="BCLI_SLACK_CHANNEL")
@click.option("--config-file",              help="Configuration file path", envvar="BCLI_CONFIG_FILE")
def main(run_all, run_tpcc, run_oltp, commit, source, inventory_file, mysql_host, mysql_username,
         mysql_password, mysql_database, packet_token, packet_project_id, api_key, slack_api_token,
         slack_channel, config_file, ansible_dir, tasks_scripts_dir, tasks_reports_dir):

    cfg = configuration.Config(configuration.create_cfg(
        run_to_task_array(run_all, run_oltp, run_tpcc), commit, source, inventory_file, mysql_host,
        mysql_username, mysql_password, mysql_database, packet_token, packet_project_id,
        api_key, slack_api_token, slack_channel, config_file, ansible_dir, tasks_scripts_dir,
        tasks_reports_dir)
    )

    cfg.dump()

    if cfg.valid_to_run() and len(cfg.tasks) > 0:
        run_id = str(uuid.uuid4())
        run_benchmark.run_tasks(cfg, run_id)
    else:
        ctx = click.get_current_context()
        click.echo(ctx.get_help())
        sys.exit(1)

def run_to_task_array(all, oltp, tpcc):
    tasks = []
    if oltp or all:
        tasks.append("oltp")
    if tpcc or all:
        tasks.append("tpcc")
    return tasks

if __name__ == "__main__":
    main()