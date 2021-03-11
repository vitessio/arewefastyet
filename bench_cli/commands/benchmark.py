import click
import bench_cli.configuration as configuration
import bench_cli.run_benchmark as run_benchmark
import bench_cli.packet_vps as vps
import bench_cli.cli as cli

@click.command()
@click.option("--tasks-run-all",                    is_flag=True, help="runs OLTP and TPCC")
@click.option("--tasks-run-tpcc", "-tpcc",          is_flag=True, help="Runs TPCC")
@click.option("--tasks-run-oltp", "-oltp",          is_flag=True, help="Runs OLTP")
@click.option("--tasks-commit", "-c",               help="Specify commit hash or branch name ")
@click.option("--tasks-source", "-s",               help="Mention the source from where the cli is called")
@click.option("--tasks-scripts-dir",                help="Path to tasks scripts directory", envvar="BCLI_TASKS_SCRIPTS_DIR")
@click.option("--tasks-reports-dir",                help="Path to tasks reports directory", envvar="BCLI_TASKS_REPORTS_DIR")
@click.option("--tasks-pprof",                      help="Profiling option for Vitess")
@click.option("--tasks-upload-to-aws", "-aws",      is_flag=True, help="Upload the task report to AWS S3")
@click.option("--tasks-inventory-file", "-invf",    help="Mention inventory file to call", envvar="BCLI_TASKS_INVENTORY_FILE")
def benchmark(*arg, **kwargs):
    cli.cfg.set_config(dict(locals().items()).get("kwargs"))
    if cli.cfg.valid_to_run() and len(cli.cfg.tasks) > 0:
        benchmark_runner = run_benchmark.BenchmarkRunner(cli.cfg, echo=True)
        benchmark_runner.run()


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
    #cli.cfg.unsafe_dump()

    #if cfg.valid_to_run() and len(cfg.tasks) > 0:
    #    benchmark_runner = run_benchmark.BenchmarkRunner(cfg, echo=True)
    #    benchmark_runner.run()
    #else:
    #    print(f'Building this repo using default method...')
