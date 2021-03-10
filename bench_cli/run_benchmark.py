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
import uuid

import bench_cli.reporting as reporting
import bench_cli.configuration as configuration
import bench_cli.task as task
import bench_cli.task_factory as taskfac


class BenchmarkRunner:
    def __init__(self, config: configuration.Config, echo=False):
        self.runner_id = uuid.uuid4()
        self.config = config
        self.tasks = self.__instantiate_tasks()
        if echo:
            print('Runner ' + self.runner_id.__str__() + ' created.')

    def __instantiate_tasks(self) -> [task.Task]:
        tasks = []
        task_factory = taskfac.TaskFactory()
        for task_name in self.config.tasks:
            tasks.append(task_factory.create_task(task_name, self.config.tasks_reports_dir,
                                                  self.config.tasks_ansible_dir,
                                                  self.config.get_inventory_file_path(),
                                                  self.config.tasks_source,
                                                  self.config.tasks_pprof_options)
                         )
        return tasks

    def run(self):
        """
        Run the BenchmarkRunner's tasks one by one.
        """
        for task in self.tasks:
            task.create_device(self.config.packet_token, self.config.packet_project_id)
            task.create_task_data_directory()
            task.build_ansible_inventory(self.config.tasks_commit)
            task.run()
            task.save_report()
            task.download_remote_pprof_folder()

            report_url = None
            if self.config.tasks_upload_to_aws:
                report_url = task.upload_report_to_aws()
            reporting.save_to_mysql(self.config, task.report, task.table_name())
            reporting.send_slack_message(self.config.slack_api_token, self.config.slack_channel,
                                         task_id=task.task_id.__str__(),
                                         report_url=report_url,
                                         filename=task.report_path())
