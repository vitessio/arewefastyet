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
import shutil
import json
import yaml
import pysftp
import abc

import bench_cli.reporting as reporting
import bench_cli.configuration as configuration
import bench_cli.packet_vps as packet_vps
import bench_cli.get_head_hash as get_head_hash

class Task:
    def __init__(self, report_dir: str, ansible_dir: str, inventory_file: str, source: str,
                 create_build_dir: bool = True):
        self.task_id = uuid.uuid4()
        self.device_id = 0
        self.device_ip = ""
        self.source = source
        self.commit_hash = ""
        self.report = None
        self.report_dir = report_dir
        self.ansible_dir = ansible_dir
        self.ansible_build_dir = os.path.join(self.ansible_dir, "build")
        self.ansible_inventory_file = inventory_file
        if self.ansible_inventory_file.find('.') is None:
            self.ansible_inventory_file += '.yml'
        self.ansible_built_inventory_file = self.__build_built_inventory_filename()
        self.ansible_built_inventory_filepath = os.path.join(self.ansible_build_dir, self.ansible_built_inventory_file)

        if create_build_dir and (
                not os.path.exists(self.ansible_build_dir) or not os.path.isdir(self.ansible_build_dir)):
            os.mkdir(self.ansible_build_dir)

    def __build_built_inventory_filename(self):
        splits = self.ansible_inventory_file.split('.')
        if len(splits) == 1:
            splits.append('.yml')
        return os.path.basename(splits[0] + "-" + str(self.task_id) + splits[1])

    def append_state_to_file(self, filepath: str):
        curr_state = self.get_state()
        dump = {'run': []}
        if os.path.exists(filepath):
            with open(filepath, 'r') as f:
                dump = {**dump, **json.load(f)}
        dump['run'].append(curr_state)
        with open(filepath, 'w') as outfile:
            json.dump(dump, outfile)

    def get_state(self):
        return {
            'task_name': self.name(),
            'run_id': self.task_id.__str__(),
            'source': self.source,
            'commit': self.commit_hash,
            'vps_id': self.device_id,
            'ip_address': self.device_ip,
            'inventory_file': os.path.basename(self.ansible_inventory_file)
        }

    def create_device(self, packet_token, packet_project_id):
        vps = packet_vps.create_vps(packet_token, packet_project_id, self.task_id)
        self.device_id = vps[0]
        self.device_ip = vps[1]
        self.append_state_to_file("config-lock.json")

    def delete_device(self, packet_token):
        packet_vps.delete_vps(packet_token, self.device_id)

    def build_ansible_inventory(self, commit_hash: str):
        shutil.copy2(self.ansible_inventory_file, self.ansible_built_inventory_filepath)
        with open(self.ansible_inventory_file, 'r') as invf:
            invdata = yaml.load(invf, Loader=yaml.FullLoader)
        self.__recursive_dict(invdata)
        # TODO: handle any commit
        if commit_hash == 'HEAD':
            self.commit_hash = get_head_hash.head_commit_hash()
        invdata["all"]["vars"]["vitess_git_version"] = self.commit_hash
        with open(self.ansible_built_inventory_filepath, 'w') as builtf:
            yaml.dump(invdata, builtf)
        pass

    def __recursive_dict(self, inventory_data):
        for k, _ in inventory_data.items():
            if isinstance(inventory_data[k], dict) and k == "hosts":
                inventory_data[k] = self.__recursive_dict_ip(inventory_data[k])
            elif isinstance(inventory_data[k], dict):
                inventory_data[k] = self.__recursive_dict(inventory_data[k])
        return inventory_data

    def __recursive_dict_ip(self, data):
        for k, _ in data.items():
            old_key = k
            data[self.device_ip] = data.pop(old_key)
        return data

    def __get_remote_task_report(self, echo=False):
        username = "root"
        cnopts = pysftp.CnOpts()
        cnopts.hostkeys = None
        remote_report_file = self.report_path() + '.bench_report'
        with pysftp.Connection(host=self.device_ip, username=username, cnopts=cnopts) as sftp:
            if echo:
                print("Connection succesfully stablished ... ")
            remote_file_path = '/tmp/' + self.name() + '.json'
            sftp.get(remote_file_path, remote_report_file)
        with open(remote_report_file, 'r') as f:
            task_report = json.load(f)
        return task_report

    def save_report(self):
        task_report = self.__get_remote_task_report()
        if len(task_report) == 0:
            return
        self.report = {**self.get_state(), 'results': task_report[0]}
        with open(self.report_path(), 'w') as f:
            json.dump(self.report, f)

    def clean_up(self):
        os.remove(self.ansible_built_inventory_file)

    @abc.abstractmethod
    def run(self, script_path: str):
        pass

    @abc.abstractmethod
    def name(self) -> str:
        pass

    @abc.abstractmethod
    def report_path(self, base: str = None) -> str:
        pass

    @abc.abstractmethod
    def table_name(self) -> str:
        pass


class OLTP(Task):
    def name(self) -> str:
        return 'oltp'

    def run(self, script_path: str):
        # TODO: Use Ansible Python API
        os.system(os.path.join(script_path, "run-oltp") + ' ' + self.ansible_built_inventory_filepath + ' ' + self.ansible_dir)

    def report_path(self, base: str = None) -> str:
        if base is not None:
            return os.path.join(base, "oltp_v2.json")
        return os.path.join(self.report_dir, "oltp_v2.json")

    def table_name(self) -> str:
        return "OLTP"


class TPCC(Task):
    def name(self) -> str:
        return 'tpcc'

    def run(self, script_path: str):
        # TODO: Use Ansible Python API
        os.system(os.path.join(script_path, "run-tpcc") + ' ' + self.ansible_built_inventory_filepath + ' ' + self.ansible_dir)

    def report_path(self, base: str = None) -> str:
        if base is not None:
            return os.path.join(base, "tpcc_v2.json")
        return os.path.join(self.report_dir, "tpcc_v2.json")

    def table_name(self) -> str:
        return "TPCC"


class TaskFactory:
    def __init__(self):
        pass

    def create_task(self, task_name, report_dir, ansible_dir, inventory_file, source) -> Task:
        if task_name == "oltp":
            return OLTP(report_dir, ansible_dir, inventory_file, source)
        elif task_name == "tpcc":
            return TPCC(report_dir, ansible_dir, inventory_file, source)


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
            tasks.append(task_factory.create_task(task_name, self.config.tasks_reports_dir,
                                                  self.config.ansible_dir,
                                                  self.config.get_inventory_file_path(),
                                                  self.config.source)
                         )
        return tasks

    def run(self):
        for task in self.tasks:
            task.create_device(self.config.packet_token, self.config.packet_project_id)
            task.build_ansible_inventory(self.config.commit)
            task.run(self.config.tasks_scripts_dir)
            task.save_report()
            reporting.save_to_mysql(self.config, task.report, task.table_name())
            reporting.send_slack_message(self.config.slack_api_token, self.config.slack_channel, task.report_path())