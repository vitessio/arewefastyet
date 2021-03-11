import click
import bench_cli.configuration as configuration
import bench_cli.server.server as server
import bench_cli.configuration as configuration
import bench_cli.run_benchmark as run_benchmark
import bench_cli.packet_vps as vps
import bench_cli.cli as cli

@click.command()
@click.option("--web",                              is_flag=True, help="Only runs the web UI")
@click.option("--web-api-key",                      help="API key", envvar="BCLI_API_KEY")
def web(*arg, **kwargs):
    cli.cfg.set_config(dict(locals().items()).get("kwargs"))
    if cli.cfg.web is True:
        server.main(cli.cfg)
