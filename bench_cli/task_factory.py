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

import bench_cli.task as task
import bench_cli.task_oltp as oltp
import bench_cli.task_tpcc as tpcc


class TaskFactory:
    def __init__(self):
        pass

    def create_task(self, task_name, report_dir, ansible_dir, inventory_file, source, pprof) -> task.Task:
        """
        Create a task children based on the given task_name.
        The task created can either be "oltp" (OLTP) or "tpcc" (TPCC).

        @param: task_name
        @param: report_dir: Path to task's report directory
        @param: ansible_dir: Path to the Ansible directory to use
        @param: inventory_file: Filename of the inventory to use
        @param: source: The task's source
        @param: pprof: The pprof configuration of the task
        """
        if task_name == "oltp":
            return oltp.OLTP(report_dir, ansible_dir, inventory_file, source, pprof)
        elif task_name == "tpcc":
            return tpcc.TPCC(report_dir, ansible_dir, inventory_file, source, pprof)
