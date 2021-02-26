import os
import sys
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../../../')))

import bench_cli.cli as cli
import bench_cli.configuration as configuration
import bench_cli.run_benchmark as run_benchmark