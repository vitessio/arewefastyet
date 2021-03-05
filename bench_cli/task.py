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
import uuid
import shutil
import json
from typing import Optional

import yaml
import abc

import bench_cli.packet_vps as packet_vps
import bench_cli.get_head_hash as get_head_hash
import bench_cli.get_from_remote as get_from_remote
import bench_cli.zip as zip
import bench_cli.aws as aws
import bench_cli.github_api as ghapi


class Task:
    def __init__(self, report_dir: str, ansible_dir: str, inventory_file: str, source: str, pprof,
                 create_build_dir: bool = True):
        self.task_id = uuid.uuid4()
        self.device_id = 0
        self.device_ip = ""
        self.source = source
        self.commit_hash = ""
        self.report = None
        self.pprof = pprof
        self.report_dir = report_dir
        self.ansible_dir = ansible_dir
        self.ansible_build_dir = os.path.join(self.ansible_dir, "build")
        self.ansible_inventory_file = inventory_file
        self.ansible_built_inventory_file = self.__build_built_inventory_filename()
        self.ansible_built_inventory_filepath = os.path.join(self.ansible_build_dir, self.ansible_built_inventory_file)

        if create_build_dir and (
                not os.path.exists(self.ansible_build_dir) or not os.path.isdir(self.ansible_build_dir)):
            os.mkdir(self.ansible_build_dir)

    def __build_built_inventory_filename(self):
        splits = os.path.basename(self.ansible_inventory_file).split('.')
        return splits[0] + "-" + str(self.task_id) + '.yml'

    def create_task_data_directory(self):
        dir = os.path.join(self.report_dir, self.name() + "-"+ self.task_id.__str__()[:8])
        if os.path.exists(dir) is False:
            os.mkdir(dir)
        self.report_dir = dir

    def append_state_to_file(self, filepath: str):
        """
        Append the task's state to the given file.
        @param filepath: Filepath used to append the state
        """
        curr_state = self.get_state()
        dump = {'run': []}
        if os.path.exists(filepath):
            with open(filepath, 'r') as f:
                dump = {**dump, **json.load(f)}
        dump['run'].append(curr_state)
        with open(filepath, 'w') as outfile:
            json.dump(dump, outfile)

    def get_state(self):
        """
        Returns the task's state.

        @todo: optimize
        """
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
        """
        Creates the task's Packet Device.
        @param: packet_token: The packet API token
        @param: packet_project_id: The packet project ID
        """
        vps = packet_vps.create_vps(packet_token, packet_project_id, self.task_id)
        self.device_id = vps[0]
        self.device_ip = vps[1]
        self.append_state_to_file("config-lock.json")

    def delete_device(self, packet_token):
        """
        Delete the task's Packet Device.
        @param: packet_token: The packet API token
        """
        packet_vps.delete_vps(packet_token, self.device_id)

    def build_ansible_inventory(self, commit_hash: str):
        """
        Create a new Ansible inventory file and build it by appending the task's configuration.
        The file will be used to run the task on Ansible.

        @param: commit_hash: Commit hash that will be used to run the task
        """
        shutil.copy2(self.ansible_inventory_file, self.ansible_built_inventory_filepath)
        with open(self.ansible_inventory_file, 'r') as invf:
            invdata = yaml.load(invf, Loader=yaml.FullLoader)
        self.__recursive_dict(invdata)

        self.commit_hash, is_pr = ghapi.resolve_ref(commit_hash)
        if self.commit_hash is None:
            self.commit_hash = commit_hash
        invdata["all"]["vars"]["vitess_git_version"] = self.commit_hash
        if is_pr:
            invdata["all"]["vars"]["vitess_git_version_pr_nb"] = int(commit_hash)
            invdata["all"]["vars"]["vitess_git_version_fetch_pr"] = "pull/{0}/head:{0}".format(commit_hash)

        if self.pprof:
            invdata["all"]["vars"]["pprof_targets"] = self.pprof["targets"]
            invdata["all"]["vars"]["pprof_args"] = self.pprof["args"]

        with open(self.ansible_built_inventory_filepath, 'w') as builtf:
            yaml.dump(invdata, builtf)

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

    def download_remote_report(self):
        src = '/tmp/' + self.name() + '.json'
        dest = self.report_path() + '.bench_report'
        get_from_remote.get_from_remote(self.device_ip, "root", src, dest)
        with open(dest, 'r') as f:
            task_report = json.load(f)
        return task_report

    def download_remote_pprof_folder(self):
        src_dir = '/tmp/pprof'
        dest_dir = os.path.join(self.report_dir, "pprof")
        get_from_remote.get_from_remote(self.device_ip, "root", src_dir, dest_dir, is_directory=True, create_dest=True)

    def upload_report_to_aws(self) -> Optional[str]:
        zip.zipdir(self.report_dir, "zip", self.report_dir)
        return aws.upload_file(self.report_dir+".zip", object_name=os.path.basename(self.report_dir+".zip"))

    def save_report(self):
        """
        Save the task's state to a report file.
        """
        task_report = self.download_remote_report()
        if len(task_report) == 0:
            return
        self.report = {**self.get_state(), 'results': task_report[0]}
        with open(self.report_path(), 'w') as f:
            json.dump(self.report, f)

    def clean_up(self):
        """
        Removes the ansible_built_inventory_file of the file system.
        """
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
