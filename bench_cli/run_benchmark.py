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
import packet_vps
import uuid
import shutil
import json
import yaml
import get_head_hash
from pathlib import Path
from initialize_benchmark import init
from reporting import add_oltp, add_tpcc
import abc

class Task:
    def __init__(self, ansible_dir: str, inventory_file: str, create_build_dir: bool = True):
        self.task_id = uuid.uuid4()
        self.device_id = 0
        self.device_ip = ""
        self.commit_hash = ""
        self.ansible_dir = ansible_dir
        self.ansible_build_dir = os.path.join(self.ansible_dir, "build")
        self.ansible_inventory_file = inventory_file + '.yml'
        self.ansible_built_inventory_file = self.__build_built_inventory_filename__(inventory_file)
        self.ansible_built_inventory_filepath = os.path.join(self.ansible_build_dir, self.ansible_built_inventory_file)

        if create_build_dir and (not os.path.exists(self.ansible_build_dir) or not os.path.isdir(self.ansible_build_dir)):
            os.mkdir(self.ansible_build_dir)

    def __build_built_inventory_filename__(self, inventory_file: str):
        return os.path.basename(inventory_file + "-" + str(self.task_id) + ".yml")

    def append_state(self, lock_json: str = "config-lock.json"):
        dump = {'run': []}
        if os.path.exists(lock_json):
            with open(lock_json) as lock_json_file:
                dump = json.load(lock_json_file)
        dump['run'].append({
            'task_name': self.name(),
            'task_id': self.task_id,
            'vps_id': self.device_ip,
            'ip_address': self.device_ip,
            'inventory_file': self.ansible_inventory_file
        })
        with open(lock_json, 'w') as outfile:
            json.dump(dump, outfile)

    def create_device(self, packet_token, packet_project_id):
        vps = packet_vps.create_vps(packet_token, packet_project_id, self.task_id)
        self.device_id = vps[0]
        self.device_ip = vps[1]
        self.append_state()

    def build_ansible_inventory(self, commit_hash: str):
        shutil.copy2(self.ansible_inventory_file, self.ansible_built_inventory_filepath)
        with open(self.ansible_inventory_file, 'r') as invf:
            invdata = yaml.load(invf, Loader=yaml.FullLoader)
        self.__recursive_dict__(invdata)
        # TODO: handle any commit
        if commit_hash == 'HEAD':
            self.commit_hash = get_head_hash.head_commit_hash()
        invdata["all"]["vars"]["vitess_git_version"] = self.commit_hash
        with open(self.ansible_built_inventory_filepath, 'w') as builtf:
            yaml.dump(invdata, builtf)
        pass

    def __recursive_dict__(self, inventory_data):
        for k, _ in inventory_data.items():
            if isinstance(inventory_data[k], dict) and k == "hosts":
                inventory_data[k] = self.__recursive_dict_ip__(inventory_data[k])
            elif isinstance(inventory_data[k], dict):
                inventory_data[k] = self.__recursive_dict__(inventory_data[k])
        return inventory_data

    def __recursive_dict_ip__(self, data):
        for k, _ in data.items():
            old_key = k
            data[self.device_ip] = data.pop(old_key)
        return data

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

    def create_task(self, task_name, ansible_dir, inventory_file) -> Task:
        if task_name == "oltp":
            return OLTP(ansible_dir, inventory_file)
        elif task_name == "tpcc":
            return TPCC(ansible_dir, inventory_file)

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
            tasks.append(task_factory.create_task(task_name, self.config.ansible_dir, self.config.get_inventory_file_path()))
        return tasks

    def run(self):
        for task in self.tasks:
            print(task.name())
            task.create_device(self.config.packet_token, self.config.packet_project_id)
            task.build_ansible_inventory(self.config.commit)
            #
            # os.system(cfg.tasks_scripts_dir + '/' + task_info['run_script'] + ' ' + Path(
            #     cfg.get_inventory_file_path()).stem + '-' + str(run_id) + '.yml')
            #
            # task_info['save_results'](cfg, run_id)
