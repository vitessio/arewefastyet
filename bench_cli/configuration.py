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
#   - returns data from config.yaml file
# -------------------------------------------------------------------------------------------------------------------------------------

import yaml
import os

import bench_cli.connection as connection

def create_cfg(web, tasks, commit, source, inventory_file, mysql_host, mysql_username,
               mysql_password, mysql_database, packet_token, packet_project_id,
               api_key, slack_api_token, slack_channel, config_file, ansible_dir,
               tasks_scripts_dir, tasks_reports_dir,delete_vps,running_benchmarks):
    return {
        "web": web, "tasks": tasks, "commit": commit, "source": source,
        "inventory_file":inventory_file, "mysql_host": mysql_host,
        "mysql_username": mysql_username, "mysql_password": mysql_password,
        "mysql_database": mysql_database, "packet_token": packet_token,
        "packet_project_id": packet_project_id, "api_key": api_key,
        "slack_api_token": slack_api_token, "slack_channel": slack_channel,
        "config_file": config_file, "ansible_dir": ansible_dir,
        "tasks_scripts_dir": tasks_scripts_dir, "tasks_reports_dir": tasks_reports_dir,
        "delete_vps": delete_vps,"running_benchmarks": running_benchmarks
    }

class Config:
    def __init__(self, cfg):
        self.__load_config(cfg)
        if self.config_file:
            cfg = {**cfg, **self.__read_from_file__()}
            self.__load_config(cfg)

    def __read_from_file__(self):
        with open(self.config_file) as f:
            return yaml.load(f, Loader=yaml.FullLoader)

    def __load_config(self, cfg) -> None:
        self.web: bool = cfg["web"]
        self.tasks: [str] = cfg["tasks"]
        self.commit: str = cfg["commit"]
        self.source: str = cfg["source"]
        self.inventory_file: str = cfg["inventory_file"]
        self.mysql_host: str = cfg["mysql_host"]
        self.mysql_username: str = cfg["mysql_username"]
        self.mysql_password: str = cfg["mysql_password"]
        self.mysql_database: str = cfg["mysql_database"]
        self.packet_token: str = cfg["packet_token"]
        self.packet_project_id: str = cfg["packet_project_id"]
        self.api_key: str = cfg["api_key"]
        self.slack_api_token: str = cfg["slack_api_token"]
        self.slack_channel: str = cfg["slack_channel"]
        self.config_file: str = cfg["config_file"]
        self.ansible_dir: str = cfg["ansible_dir"]
        self.tasks_scripts_dir: str = cfg["tasks_scripts_dir"]
        self.tasks_reports_dir: str = cfg["tasks_reports_dir"]
        self.delete_vps: str = cfg["delete_vps"]
        self.running_benchmarks: bool = cfg["running_benchmarks"] 

    def get_inventory_file_path(self) -> str:
        """
        Build the inventory file path from the given ansible directory and inventory file.
        @return: str
        """
        return os.path.join(self.ansible_dir, self.inventory_file)

    def unsafe_dump(self, echo=True) -> str:
        """
        Dumps the configuration data.
        @param echo: If True, prints the dump
        @return: str
        """
        attrs = vars(self)
        dumpstr = '\n'.join("%s: %s" % item for item in attrs.items())
        if echo:
            print(dumpstr)
        return dumpstr

    def valid_to_run(self) -> bool:
        """
        Check if the configuration allows us to run tests.
        @return: bool

        @todo: Returns False if the inventory file cannot be resolved
        """
        if not self.commit or not self.source or not self.inventory_file:
            # TODO: throw error instead
            return False
        return True

    def mysql_connect(self):
        """
        Connect to the mysql.
        """
        return connection.connectdb(self.mysql_host, self.mysql_username, self.mysql_password, self.mysql_database)
