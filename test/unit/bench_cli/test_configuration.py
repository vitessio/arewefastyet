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
import tempfile
import shutil
import os
import yaml
from .context import configuration

default_cfg_fields = {
    "web": False, "tasks": [], "commit": "HEAD", "source": "testing",
    "inventory_file": "inventory_file", "mysql_host": "localhost",
    "mysql_username": "root", "mysql_password": "password",
    "mysql_database": "main", "packet_token": "AB12345",
    "packet_project_id": "AABB11", "api_key": "123-ABC-456-EFG",
    "slack_api_token": "slack-token", "slack_channel": "general",
    "config_file": "./config", "ansible_dir": "./ansible",
    "tasks_scripts_dir": "./scripts", "tasks_reports_dir": "./reports"
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

def init_tmp_config_dir(cfg_fields):
    tmpdir = tempfile.mkdtemp()
    f, tmpcfg = tempfile.mkstemp(".yaml", "config", tmpdir, text=True)
    os.write(f, yaml.dump(cfg_fields).encode())
    os.close(f)
    return tmpdir, tmpcfg

class TestCreateConfig(unittest.TestCase):
    def test_create_default_cfg_dict(self):
        cfg = configuration.create_cfg(**default_cfg_fields)
        self.assertEqual(cfg, default_cfg_fields)

    def test_create_incorrect_cfg_dict(self):
        with self.assertRaises(TypeError) as ctx:
            configuration.create_cfg()
        self.assertTrue('create_cfg() missing' in ctx.exception.__str__())
        self.assertTrue('required positional arguments' in ctx.exception.__str__())

    def test_create_config(self):
        cfg_file_data = default_cfg_fields.copy()
        cfg_file_data.__delitem__("config_file")
        tmpdir, tmpcfg = init_tmp_config_dir(cfg_file_data)

        cp_cfg_fields = default_cfg_fields.copy()
        cp_cfg_fields["config_file"] = tmpcfg
        cfg = configuration.create_cfg(**cp_cfg_fields)

        config = configuration.Config(cfg)
        for key in cfg:
            self.assertEqual(cfg[key], config.__getattribute__(key))
        shutil.rmtree(tmpdir)

class TestValidToRun(unittest.TestCase):
    def setUp(self) -> None:
        super().setUp()

        cfg_file_data = default_cfg_fields.copy()
        cfg_file_data.__delitem__("config_file")
        self.tmpdir, self.tmpcfg = init_tmp_config_dir(cfg_file_data)

        cp_cfg_fields = default_cfg_fields.copy()
        cp_cfg_fields["config_file"] = self.tmpcfg
        self.cfg = configuration.create_cfg(**cp_cfg_fields)

        self.config = configuration.Config(self.cfg)

    def tearDown(self) -> None:
        super().tearDown()
        shutil.rmtree(self.tmpdir)

    def test_valid_to_run_true(self):
        valid = self.config.valid_to_run()
        self.assertTrue(valid, "should be valid")

    def test_valid_to_run_no_commit(self):
        self.config.commit = None
        valid = self.config.valid_to_run()
        self.assertFalse(valid, "should be invalid")

    def test_valid_to_run_no_source(self):
        self.config.source = None
        valid = self.config.valid_to_run()
        self.assertFalse(valid, "should be invalid")

    def test_valid_to_run_no_inventory_file(self):
        self.config.inventory_file = None
        valid = self.config.valid_to_run()
        self.assertFalse(valid, "should be invalid")

class TestInventoryFile(unittest.TestCase):
    def setUp(self) -> None:
        super().setUp()
        cfg_file_data = default_cfg_fields.copy()
        cfg_file_data.__delitem__("config_file")
        self.tmpdir, self.tmpcfg = init_tmp_config_dir(cfg_file_data)

        cp_cfg_fields = default_cfg_fields.copy()
        cp_cfg_fields["config_file"] = self.tmpcfg
        self.cfg = configuration.create_cfg(**cp_cfg_fields)

        self.config = configuration.Config(self.cfg)

    def tearDown(self) -> None:
        super().tearDown()
        shutil.rmtree(self.tmpdir)

    def test_get_inventory_file_path(self):
        expected_path = os.path.join("./ansible", "inventory_file")
        path = self.config.get_inventory_file_path()
        self.assertEqual(expected_path, path, "invalid path")

if __name__ == '__main__':
    unittest.main()