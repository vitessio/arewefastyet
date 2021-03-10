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

import os
import ansible_runner
import shutil
import tempfile

import bench_cli.task as task

class TPCC(task.Task):
    def name(self) -> str:
        """
        Returns the task's name
        """
        return 'tpcc'

    def run(self, config: ansible_runner.RunnerConfig = None):
        """
        Runs the task.
        """
        tmpdir = tempfile.mkdtemp()
        ssh_priv_key = open(os.path.expanduser('~/.ssh/id_rsa')).read()

        runner = ansible_runner.run(
            ident=self.task_id,
            private_data_dir=tmpdir,
            project_dir=self.ansible_dir,
            artifact_dir=os.path.abspath(os.path.join(self.ansible_dir, "artifacts")),
            playbook=os.path.abspath(os.path.join(self.ansible_dir, "full.yml")),
            inventory=[os.path.abspath(self.ansible_built_inventory_filepath)],
            ssh_key=ssh_priv_key,
            extravars=dict({"provision": True, "clean": True, "tpcc": "true"}),
            envvars=dict({"OBJC_DISABLE_INITIALIZE_FORK_SAFETY": "YES"}),
            cmdline="-u root",
        )
        if runner.status == "failed" or runner.rc is not 0:
            raise RuntimeError("task execution failed, ansible finished with {0}".format(runner.rc))
        shutil.rmtree(tmpdir)

    def report_path(self, base: str = None) -> str:
        """
        Returns the path of the task report directory.

        @param: base: Folder to use as base for the report directory
        """
        if base is not None:
            return os.path.join(base, "tpcc_v2.json")
        return os.path.join(self.report_dir, "tpcc_v2.json")

    def table_name(self) -> str:
        """
        Returns the task's table name
        """
        return "TPCC"
