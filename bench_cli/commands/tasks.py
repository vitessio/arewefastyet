import click
import bench_cli.configuration as configuration
import bench_cli.run_benchmark as run_benchmark
import bench_cli.packet_vps as vps
import bench_cli.cli as cli

@click.command()
@click.option("--delete-benchmark", "-d",           help="Delete VPS")
def tasks(*arg, **kwargs):
    cli.cfg.set_config(dict(locals().items()).get("kwargs"))
    if cli.cfg.delete_benchmark is not None:
        delete_benchmark_procedure(cli.cfg)

def delete_benchmark_procedure(cfg: configuration.Config):
    vps.delete_vps(cli.cfg.packet_token, cli.cfg.delete_benchmark)
