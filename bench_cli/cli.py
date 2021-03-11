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
#  - https://github.com/vitessio/arewefastyet/blob/master/docs/cli.md
# -------------------------------------------------------------------------------------------------------------------------------------

import click
from bench_cli.commands import benchmark, tasks, web
import bench_cli.configuration as configuration

cfg = None

@click.group(help="CLI to configure, run and control benchmark runs")
@click.option("--tasks-ansible-dir",  help="Path to the Ansible directory", envvar="BCLI_TASKS_ANSIBLE_DIR")
@click.option("--config-file",        help="Configuration file path", envvar="BCLI_CONFIG_FILE")
@click.option("--mysql-host",                       help="MySQL server hostname", envvar="BCLI_MYSQL_HOST")
@click.option("--mysql-username",                   help="MySQL server username", envvar="BCLI_MYSQL_USER")
@click.option("--mysql-password",                   help="MySQL server password", envvar="BCLI_MYSQL_PASSWORD")
@click.option("--mysql-database",                   help="MySQL database to use", envvar="BCLI_MYSQL_DB")
@click.option("--packet-token",                     help="Token used to authenticate Packet", envvar="BCLI_PACKET_TOKEN")
@click.option("--packet-project-id",                help="Packet project ID", envvar="BCLI_PACKET_PROJECT_ID")
@click.option("--slack-api-token",                  help="Slack API token", envvar="BCLI_SLACK_TOKEN")
def main(*args, **kwargs):
    global cfg
    cfg = configuration.Config(dict(locals().items()).get("kwargs"))
    pass

main.add_command(benchmark.benchmark)
main.add_command(tasks.tasks)
main.add_command(web.web)

if __name__ == '__main__':
    main()
'''
import click
import sys

import bench_cli.server.server as server
import bench_cli.configuration as configuration
import bench_cli.run_benchmark as run_benchmark
import bench_cli.packet_vps as vps

@click.command()
@click.option("--web",                              is_flag=True, help="Only runs the web UI")
@click.option("--web-api-key",                      help="API key", envvar="BCLI_API_KEY")
@click.option("--tasks-run-all",                    is_flag=True, help="runs OLTP and TPCC")
@click.option("--tasks-run-tpcc", "-tpcc",          is_flag=True, help="Runs TPCC")
@click.option("--tasks-run-oltp", "-oltp",          is_flag=True, help="Runs OLTP")
@click.option("--tasks-commit", "-c",               help="Specify commit hash or branch name ")
@click.option("--tasks-source", "-s",               help="Mention the source from where the cli is called")
@click.option("--tasks-scripts-dir",                help="Path to tasks scripts directory", envvar="BCLI_TASKS_SCRIPTS_DIR")
@click.option("--tasks-reports-dir",                help="Path to tasks reports directory", envvar="BCLI_TASKS_REPORTS_DIR")
@click.option("--tasks-pprof",                      help="Profiling option for Vitess")
@click.option("--tasks-upload-to-aws", "-aws",      is_flag=True, help="Upload the task report to AWS S3")
@click.option("--tasks-ansible-dir",                help="Path to the Ansible directory", envvar="BCLI_TASKS_ANSIBLE_DIR")
@click.option("--tasks-inventory-file", "-invf",    help="Mention inventory file to call", envvar="BCLI_TASKS_INVENTORY_FILE")
@click.option("--config-file",                      help="Configuration file path", envvar="BCLI_CONFIG_FILE")
@click.option("--mysql-host",                       help="MySQL server hostname", envvar="BCLI_MYSQL_HOST")
@click.option("--mysql-username",                   help="MySQL server username", envvar="BCLI_MYSQL_USER")
@click.option("--mysql-password",                   help="MySQL server password", envvar="BCLI_MYSQL_PASSWORD")
@click.option("--mysql-database",                   help="MySQL database to use", envvar="BCLI_MYSQL_DB")
@click.option("--packet-token",                     help="Token used to authenticate Packet", envvar="BCLI_PACKET_TOKEN")
@click.option("--packet-project-id",                help="Packet project ID", envvar="BCLI_PACKET_PROJECT_ID")
@click.option("--slack-api-token",                  help="Slack API token", envvar="BCLI_SLACK_TOKEN")
@click.option("--slack-channel",                    help="Slack channel", envvar="BCLI_SLACK_CHANNEL")
@click.option("--delete-benchmark", "-d",           help="Delete VPS")
def main(*args, **kwargs):
    cfg = configuration.Config(dict(locals().items()).get("kwargs"))

    cfg.unsafe_dump()

    if cfg.delete_benchmark is not None:
        delete_benchmark_procedure(cfg)
    elif cfg.web is True:
        server.main(cfg)
    elif cfg.valid_to_run() and len(cfg.tasks) > 0:
        benchmark_runner = run_benchmark.BenchmarkRunner(cfg, echo=True)
        benchmark_runner.run()
    else:
        ctx = click.get_current_context()
        click.echo(ctx.get_help())
        sys.exit(1)


def delete_benchmark_procedure(cfg: configuration.Config):
    vps.delete_vps(cfg.packet_token, cfg.delete_benchmark)



def run_to_task_array(all, oltp, tpcc) -> [str]:
    """
    Transforms the tasks given through CLI flags into an array of string.
    @param all: All tasks
    @param oltp: OLTP task
    @param tpcc: TPCC task
    @return: [str]
    """
    tasks = []
    if oltp or all:
        tasks.append("oltp")
    if tpcc or all:
        tasks.append("tpcc")
    return tasks




if __name__ == "__main__":
    main()  # pylint: disable=E1120
'''
