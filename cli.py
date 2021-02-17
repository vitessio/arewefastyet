# ------------------------------------------------------------------------------------------------------------------------------------
# Copyright 2020 The Vitess Authors.
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
# -------------------------------------------------------------------------------------------------------------------------------------

import click
import uuid
from run_benchmark import run_tasks


# ----------------------------------------------- CLI flags --------------------------------------------------------

@click.command()
@click.option("--runall", is_flag=True, help="runs OLTP and TPCC")
@click.option("--commit",  "-c", help="specify commit")
@click.option("--source",  "-s", help="mention the source from where the cli is called")
@click.option("--run-tpcc", "-runt",is_flag=True, help="Runs TPCC")
@click.option("--run-oltp", "-runo",is_flag=True, help="Runs OLTP")
@click.option("--inventory-file", "-invf", help="Inventory File path")

# -------------------------------------------------------------------------------------------------------------------
# -------------------------------------- Actions when flags are called ----------------------------------------------

def main(runall,commit,source,run_tpcc,run_oltp,inventory_file):

    if runall and commit and source:
           run_id = str(uuid.uuid4())
           # TODO: add to CLI flags
           tasks = ["oltp", "tpcc"]

           run_tasks(commit, run_id, source, tasks)
    elif run_tpcc:
        print("run tpcc")

    # Display --help information when no flag is called
    else:
        ctx = click.get_current_context()
        click.echo(ctx.get_help())

# -------------------------------------------------------------------------------------------------------------------

if __name__ == "__main__":
    main()
