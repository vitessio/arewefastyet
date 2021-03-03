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

import unittest
import os

from .context import configuration

default_cfg_fields = {
    "web": False, "tasks": [], "commit": "HEAD", "source": "testing",
    "inventory_file": "inventory_file", "mysql_host": "localhost",
    "mysql_username": "root", "mysql_password": "password",
    "mysql_database": "main", "packet_token": "AB12345",
    "packet_project_id": "AABB11", "api_key": "123-ABC-456-EFG",
    "slack_api_token": "slack-token", "slack_channel": "general",
    "config_file": "./config", "ansible_dir": "./ansible",
    "tasks_scripts_dir": "./scripts", "tasks_reports_dir": "./reports",
    "tasks_pprof": None, "delete_benchmark": False
}

default_cfg_file_yaml = "" \
                        "mysql_host: localhost\n" \
                        "mysql_username: root\n"\
                        "mysql_password: password\n"\
                        "mysql_database: main\n"\
                        "inventory_file: inventory_file\n"\
                        "packet_token : AB12345\n"\
                        "packet_project_id : AABB11\n"\
                        "api_key: 123-ABC-456-EFG\n"\
                        "slack_api_token: slack-token\n"\
                        "slack_channel: general\n"


class TestValidToRun(unittest.TestCase):
    def test_valid_to_run_true(self):
        config = configuration.Config({"commit": "HEAD", "source": "test", "inventory_file": "file"})
        valid = config.valid_to_run()
        self.assertTrue(valid, "should be valid")

    def test_valid_to_run_no_commit(self):
        config = configuration.Config({"source": "test", "inventory_file": "file"})
        valid = config.valid_to_run()
        self.assertFalse(valid, "should be invalid")

    def test_valid_to_run_no_source(self):
        config = configuration.Config({"commit": "HEAD", "inventory_file": "file"})
        valid = config.valid_to_run()
        self.assertFalse(valid, "should be invalid")

    def test_valid_to_run_no_inventory_file(self):
        config = configuration.Config({"commit": "HEAD", "source": "test"})
        valid = config.valid_to_run()
        self.assertFalse(valid, "should be invalid")


class TestInventoryFile(unittest.TestCase):
    def test_get_inventory_file_path(self):
        ansible_dir = "./ansible"
        inventory_file = "inventory_file"
        config = configuration.Config({"ansible_dir": ansible_dir, "inventory_file": inventory_file})

        expected_path = os.path.join(ansible_dir, inventory_file)
        path = config.get_inventory_file_path()
        self.assertEqual(expected_path, path, "invalid path")


class TestParseProfilingInformation(unittest.TestCase):
    def test_vtgate_cpu_pprof(self):
        config = configuration.Config({"tasks_pprof": "vtgate/cpu"})
        self.assertEqual(["vtgate"], config.tasks_pprof_options["targets"])
        self.assertEqual("cpu", config.tasks_pprof_options["args"])

    def test_vttablet_cpu_pprof(self):
        config = configuration.Config({"tasks_pprof": "vttablet/cpu"})
        self.assertEqual(["vttablet"], config.tasks_pprof_options["targets"])
        self.assertEqual("cpu", config.tasks_pprof_options["args"])

    def test_vttablet_vtgate_cpu_pprof(self):
        config = configuration.Config({"tasks_pprof": "vttablet/vtgate/cpu"})
        self.assertEqual(["vttablet", "vtgate"], config.tasks_pprof_options["targets"])
        self.assertEqual("cpu", config.tasks_pprof_options["args"])

    def test_incorrect_tasks_pprof(self):
        with self.assertRaises(AttributeError) as ctx:
            configuration.Config({"tasks_pprof": "incorrect"})
        self.assertTrue('profiling needs a target' in ctx.exception.__str__())


if __name__ == '__main__':
    unittest.main()
