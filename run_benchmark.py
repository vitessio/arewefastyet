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
#   - Creates Run ID and runs Benchmark
#
# Arguments: python run-benchmark.py <commit hash> <run id> <source>
# -------------------------------------------------------------------------------------------------------------------------------------

import time
import subprocess
import os
from config import inventory_file
from pathlib import Path
import sys
from initialize_benchmark import init
from report import add_oltp, add_tpcc

# ------------------------------------------------------ Runs benchmark tasks ---------------------------------------------------------

def init_task(name, script, save_results):
   return {
      "name": name,
      "run_script": script,
      "save_results": save_results
   }

tasks_list = {
   "oltp": init_task("oltp", "run-oltp", add_oltp),
   "tpcc": init_task("tpcc", "run-tpcc", add_tpcc)
}

def print_step(task, step):
   print('-------------', task, '-', step, '-------------')

def create_task(task):
   return tasks_list.get(task)

def run_tasks(commit, run_id, source, tasks):
   for task in tasks:
      task_info = create_task(task)

      print_step(task_info['name'], 'Initialize VPS')
      init(run_id, commit)

      print_step(task_info['name'], 'Running Benchmark')
      os.system('./' + task_info['run_script'] + ' ' + Path('./ansible/build/' + inventory_file()).stem + '-' + str(run_id) + '.yml')

      print_step(task_info['name'], 'Saving Results')
      task_info['save_results'](run_id, source)
