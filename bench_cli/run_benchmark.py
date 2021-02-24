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
#   - Creates Run ID and runs Benchmark
#
# Arguments: python run-benchmark.py <commit hash> <run id> <source>
# -------------------------------------------------------------------------------------------------------------------------------------

import os
import configuration
import uuid
from pathlib import Path
from initialize_benchmark import init
from reporting import add_oltp, add_tpcc
import abc

class Task:
    def __init__(self):
        self.task_id = uuid.uuid4()

    def run(self):
        pass

    def save_results(self):
        pass

    @abc.abstractmethod
    def name(self) -> str:
        pass

class OLTP(Task):
    def name(self) -> str:
        return 'oltp'

class TPCC(Task):
    def name(self) -> str:
        return 'tpcc'


class TaskFactory:
    def __init__(self):
        pass

    def create_task(self, task_name) -> Task:
        if task_name == "oltp":
            return OLTP()
        elif task_name == "tpcc":
            return TPCC()

# ------------------------------------------------------ Runs benchmark tasks ---------------------------------------------------------

class BenchmarkRunner:
    def __init__(self, config: configuration.Config, echo=False):
        self.runner_id = uuid.uuid4()
        self.config = config
        self.tasks = self.__instantiate_tasks__()
        if echo:
            print('Runner ' + self.runner_id.__str__() + ' created.')

    def __instantiate_tasks__(self) -> [Task]:
        tasks = []
        task_factory = TaskFactory()
        for task_name in self.config.tasks:
            tasks.append(task_factory.create_task(task_name))
        return tasks

    def run(self):
        for task in self.tasks:
            print(task.name())
            # init(cfg, run_id)
            #
            #
            # os.system(cfg.tasks_scripts_dir + '/' + task_info['run_script'] + ' ' + Path(
            #     cfg.get_inventory_file_path()).stem + '-' + str(run_id) + '.yml')
            #
            # task_info['save_results'](cfg, run_id)
