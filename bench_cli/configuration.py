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

from typing import Dict

import bench_cli.connection as connection


class Config:
    def __init__(self, cfg: Dict[str, any]):
        self.config_file: str = None

        # Web related
        self.web: bool = None
        self.web_api_key: str = None

        # Task related
        self.tasks: [str] = None
        self.tasks_run_oltp: bool = None
        self.tasks_run_tpcc: bool = None
        self.tasks_run_all: bool = None
        self.tasks_commit: str = None
        self.tasks_source: str = None
        self.tasks_inventory_file: str = None
        self.tasks_ansible_dir: str = None
        self.tasks_scripts_dir: str = None
        self.tasks_reports_dir: str = None
        self.tasks_pprof: str = None
        self.tasks_pprof_options = None
        self.tasks_upload_to_aws = None

        # MySQL related
        self.mysql_host: str = None
        self.mysql_username: str = None
        self.mysql_password: str = None
        self.mysql_database: str = None

        # Packet related
        self.packet_token: str = None
        self.packet_project_id: str = None

        # Slack related
        self.slack_api_token: str = None
        self.slack_channel: str = None

        # Other commands
        self.delete_benchmark: str = None

        if cfg.get("config_file") is not None:
            self.__load_config(self.__read_from_file(cfg.get("config_file")))
        self.__load_config(cfg)
        self.__parse()

    def __read_from_file(self, file: str):
        with open(file) as f:
            return yaml.load(f, Loader=yaml.SafeLoader)

    def __load_config(self, cfg) -> None:
        for key in cfg:
            if cfg[key] is not None:
                setattr(self, key, cfg[key])

    def __parse(self):
        self.__parse_profiling()
        self.__parse_run_task()

    def __parse_profiling(self):
        self.tasks_pprof_options = None
        if self.tasks_pprof is None:
            return
        splits = self.tasks_pprof.split("/")
        pprof_targets = splits[:len(splits) - 1]
        pprof_args = splits[-1]
        if len(pprof_targets) == 0 or ("vttablet" not in pprof_targets and "vtgate" not in pprof_targets):
            raise AttributeError("profiling needs a target (vttablet, vtgate)")
        # TODO: check pprof_args based on vitess check
        self.tasks_pprof_options = {
            "targets": pprof_targets,
            "args": pprof_args,
        }

    def __parse_run_task(self):
        self.tasks = []
        if self.tasks_run_oltp or self.tasks_run_all:
            self.tasks.append("oltp")
        if self.tasks_run_tpcc or self.tasks_run_all:
            self.tasks.append("tpcc")

    def get_inventory_file_path(self) -> str:
        """
        Build the inventory file path from the given ansible directory and inventory file.
        @return: str
        """
        return os.path.join(self.tasks_ansible_dir, self.tasks_inventory_file)

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
        """
        if not self.tasks_commit or not self.tasks_source or not self.tasks_inventory_file:
            return False
        return True

    def mysql_connect(self):
        """
        Connect to the mysql.
        """
        return connection.connectdb(self.mysql_host, self.mysql_username, self.mysql_password, self.mysql_database)
