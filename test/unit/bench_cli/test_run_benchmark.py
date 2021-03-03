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
from .context import configuration, run_benchmark

sample_inv_file = '''
---
all:
  vars:
    vitess_git_version: "HEAD"
    cell: local
    keyspace: main
    # OLTP
    oltp_table_size: 10000000
    oltp_threads: 100
    oltp_preparation_time: 30
    oltp_warmup_time: 10
    oltp_test_time: 900
    oltp_number_tables: 50
    # TPCC
    tpcc_warehouses: 1000
    tpcc_threads: 300
    tpcc_load_threads: 25
    tpcc_preparation_time: 900
    tpcc_ensure_time: 900
    tpcc_warmup_time: 90
    tpcc_test_time: 900
    tpcc_number_tables: 1
  hosts:
    52.250.110.51:
      storage_device:
        device: nvme0n1
        partition: nvme0n1p1
  children:
    sysbench:
      hosts:
        52.250.110.51:
    prometheus:
      hosts:
        52.250.110.51:
    etcd:
      hosts:
        52.250.110.51:
    vtctld:
      hosts:
        52.250.110.51:
    vtgate:
      vars:
        vtgate_query_cache_size: 1000
        vtgate_max_goproc: 6
      hosts:
        52.250.110.51:
          gateways:
            - id: 1
              port: 15001
              mysql_port: 13306
              grpc_port: 15306
            - id: 2
              port: 15002
              mysql_port: 13307
              grpc_port: 15307
            - id: 3
              port: 15003
              mysql_port: 13308
              grpc_port: 15308
            - id: 4
              port: 15004
              mysql_port: 13309
              grpc_port: 15309
            - id: 5
              port: 15005
              mysql_port: 13310
              grpc_port: 15310
            - id: 6
              port: 15006
              mysql_port: 13311
              grpc_port: 15311
    vttablet:
      vars:
        vitess_memory_ratio: 0.6
        vttablet_query_cache_size: 10000
        vttablet_max_goproc: 24
      hosts:
        52.250.110.51:
          tablets:
            - id: 1001
              keyspace: main
              shard: -80
              pool_size: 500
              transaction_cap: 2000
              port: 16001
              grpc_port: 17001
              mysql_port: 18001
              mysqld_exporter_port: 9104
            - id: 2001
              keyspace: main
              shard: 80-
              pool_size: 500
              transaction_cap: 2000
              port: 16002
              grpc_port: 17002
              mysql_port: 18002
              mysqld_exporter_port: 9105
'''

default_cfg_fields = {
    "web": False, "tasks": ['oltp'], "commit": "HEAD", "source": "testing",
    "inventory_file": "inventory_file", "mysql_host": "localhost",
    "mysql_username": "root", "mysql_password": "password",
    "mysql_database": "main", "packet_token": "AB12345",
    "packet_project_id": "AABB11", "api_key": "123-ABC-456-EFG",
    "slack_api_token": "slack-token", "slack_channel": "general",
    "config_file": "./config", "ansible_dir": "./ansible",
    "tasks_scripts_dir": "./scripts", "tasks_reports_dir": "./reports",
    "tasks_pprof": None
}

def data_to_tmp_yaml(prefix, suffix, data):
    tmpdir = tempfile.mkdtemp()
    f, tmpcfg = tempfile.mkstemp(suffix, prefix, tmpdir, text=True)
    os.write(f, yaml.dump(data).encode())
    os.close(f)
    return tmpdir, tmpcfg

class TestTaskFactoryCreateProperTaskType(unittest.TestCase):
    def setUp(self) -> None:
        super().setUp()
        self.tmpdir = tempfile.mkdtemp()
        self.task_factory = run_benchmark.TaskFactory()

    def tearDown(self) -> None:
        super().tearDown()
        shutil.rmtree(self.tmpdir)

    def test_create_tpcc(self):
        task = self.task_factory.create_task("tpcc", self.tmpdir, self.tmpdir, "inv_file", "unit_test", None)
        self.assertEqual("tpcc", task.name())

    def test_create_oltp(self):
        task = self.task_factory.create_task("oltp", self.tmpdir, self.tmpdir, "inv_file", "unit_test", None)
        self.assertEqual("oltp", task.name())

class TestCreateBenchmarkRunner(unittest.TestCase):
    def setUp(self) -> None:
        super().setUp()
        cfg_file_data = default_cfg_fields.copy()
        cfg_file_data.__delitem__("config_file")
        cfg_file_data.__delitem__("tasks")
        self.tmpdir, self.tmpcfg = data_to_tmp_yaml("config", ".yaml", cfg_file_data)
        cp_cfg_fields = default_cfg_fields.copy()
        cp_cfg_fields["config_file"] = self.tmpcfg
        cp_cfg_fields["tasks"] = ['oltp', 'tpcc']
        self.cfg = configuration.create_cfg(**cp_cfg_fields)
        self.config = configuration.Config(self.cfg)

    def tearDown(self) -> None:
        super().tearDown()
        shutil.rmtree(self.tmpdir)

    def test_create_benchmark_runner_all_tasks(self):
        benchmark_runner = run_benchmark.BenchmarkRunner(self.config)
        self.assertEqual(2, len(benchmark_runner.tasks))
        for i, task in enumerate(benchmark_runner.tasks):
            self.assertEqual(self.cfg["tasks"][i], task.name())

    def test_create_benchmark_runner_oltp(self):
        self.cfg["tasks"] = ['oltp']
        self.config.tasks = self.cfg["tasks"]
        benchmark_runner = run_benchmark.BenchmarkRunner(self.config)
        self.assertEqual(1, len(benchmark_runner.tasks))
        self.assertEqual(self.cfg["tasks"][0], benchmark_runner.tasks[0].name())

    def test_create_benchmark_runner_tpcc(self):
        self.cfg["tasks"] = ['tpcc']
        self.config.tasks = self.cfg["tasks"]
        benchmark_runner = run_benchmark.BenchmarkRunner(self.config)
        self.assertEqual(1, len(benchmark_runner.tasks))
        self.assertEqual(self.cfg["tasks"][0], benchmark_runner.tasks[0].name())

class TestCreationOfTaskCheckValues(unittest.TestCase):
    def setUp(self) -> None:
        super().setUp()
        cfg_file_data = default_cfg_fields.copy()
        cfg_file_data.__delitem__("config_file")
        cfg_file_data.__delitem__("tasks_reports_dir")
        cfg_file_data.__delitem__("ansible_dir")
        self.tmpdir, self.tmpcfg = data_to_tmp_yaml("config", ".yaml", cfg_file_data)
        cp_cfg_fields = default_cfg_fields.copy()
        cp_cfg_fields["config_file"] = self.tmpcfg
        cp_cfg_fields["tasks_reports_dir"] = self.tmpdir
        cp_cfg_fields["ansible_dir"] = self.tmpdir
        cp_cfg_fields["tasks"] = ['oltp', 'tpcc']
        self.cfg = configuration.create_cfg(**cp_cfg_fields)
        self.config = configuration.Config(self.cfg)
        self.task_factory = run_benchmark.TaskFactory()

    def tearDown(self) -> None:
        super().tearDown()
        shutil.rmtree(self.tmpdir)

    def test_create_task_check_values(self):
        tcs = [
            {"source": "unit_test","inventory_file": "inv_file","task_name": "oltp"},
            {"source": "unit_test", "inventory_file": "inv_file", "task_name": "tpcc"}
        ]
        for tc in tcs:
            task = self.task_factory.create_task(tc["task_name"], self.tmpdir, self.tmpdir, tc["inventory_file"], tc["source"], None)

            self.assertEqual(tc["task_name"], task.name())
            self.assertEqual(self.tmpdir, task.report_dir)
            self.assertEqual(self.tmpdir, task.ansible_dir)
            self.assertEqual(tc["inventory_file"], task.ansible_inventory_file)
            self.assertEqual(tc["source"], task.source)

            expected_ansible_build_dir = os.path.join(self.tmpdir, 'build')
            self.assertEqual(expected_ansible_build_dir, task.ansible_build_dir)

            expected_ansible_built_file = tc["inventory_file"].split('.')[0] + '-' + str(task.task_id) + '.yml'
            self.assertEqual(expected_ansible_built_file, task.ansible_built_inventory_file)
            self.assertEqual(os.path.join(expected_ansible_build_dir, expected_ansible_built_file), task.ansible_built_inventory_filepath)

            self.assertEqual(task.name().upper(), task.table_name())
            self.assertEqual(os.path.join(self.tmpdir, task.name() + "_v2.json"), task.report_path())
            self.assertEqual(os.path.join("./", task.name() + "_v2.json"), task.report_path("./"))

    def test_create_task_with_benchmark_runner_check_values(self):
        tcs = [
            {"tasks": ["oltp"]},
            {"tasks": ["tpcc"]}
        ]
        for tc in tcs:
            self.cfg["tasks"] = tc["tasks"]
            self.config.tasks = tc["tasks"]
            benchmark_runner = run_benchmark.BenchmarkRunner(self.config)

            task = benchmark_runner.tasks[0]

            self.assertEqual(self.cfg["tasks"][0], task.name())
            self.assertEqual(self.tmpdir, task.report_dir)
            self.assertEqual(self.tmpdir, task.ansible_dir)
            self.assertEqual(self.config.get_inventory_file_path(), task.ansible_inventory_file)
            self.assertEqual(self.cfg["source"], task.source)

            expected_ansible_build_dir = os.path.join(self.tmpdir, 'build')
            self.assertEqual(expected_ansible_build_dir, task.ansible_build_dir)

            expected_ansible_built_file = self.cfg["inventory_file"].split('.')[0] + '-' + str(task.task_id) + '.yml'
            self.assertEqual(expected_ansible_built_file, task.ansible_built_inventory_file)
            self.assertEqual(os.path.join(expected_ansible_build_dir, expected_ansible_built_file), task.ansible_built_inventory_filepath)

            self.assertEqual(task.name().upper(), task.table_name())
            self.assertEqual(os.path.join(self.tmpdir, task.name() + "_v2.json"), task.report_path())
            self.assertEqual(os.path.join("./", task.name() + "_v2.json"), task.report_path("./"))

class TestBuildAnsibleInventoryFile(unittest.TestCase):
    def setUp(self) -> None:
        super().setUp()
        self.tmpdir = None
        self.tmpcfg = None

    def tearDown(self) -> None:
        super().tearDown()
        if self.tmpdir is not None:
            shutil.rmtree(self.tmpdir)

    def setup(self, overrides):
        cfg_file_data = default_cfg_fields.copy()
        for key in overrides:
            cfg_file_data.__delitem__(key)
        self.tmpdir, self.tmpcfg = data_to_tmp_yaml("config", ".yaml", cfg_file_data)
        cp_cfg_fields = default_cfg_fields.copy()
        for key in overrides:
            if overrides[key] == "__tmpdir":
                cp_cfg_fields[key] = self.tmpdir
            elif overrides[key] == "__tmpcfg":
                cp_cfg_fields[key] = self.tmpcfg
            else:
                cp_cfg_fields[key] = overrides[key]
        cfg = configuration.create_cfg(**cp_cfg_fields)
        self.config = configuration.Config(cfg)
        path = self.config.get_inventory_file_path()
        f = open(path, "w+")
        f.write(sample_inv_file)
        f.close()
        self.benchmark_runner = run_benchmark.BenchmarkRunner(self.config)

    def test_build_ansible_inventory_created(self):
        inventory_yml = "inventory.yml"
        cf = {
            "config_file": "__tmpcfg",
            "tasks_reports_dir": "__tmpdir",
            "ansible_dir": "__tmpdir",
            "inventory_file": inventory_yml
        }
        self.setup(cf)
        self.setup(cf)
        task = self.benchmark_runner.tasks[0]
        task.build_ansible_inventory('HEAD')

        exptectedPath = os.path.join(self.tmpdir, "build", inventory_yml.split('.')[0] + '-' + task.task_id.__str__() + ".yml")
        self.assertEqual(exptectedPath, task.ansible_built_inventory_filepath)
        self.assertTrue(os.path.exists(exptectedPath))

    def test_build_ansible_inventory_pprof(self):
        inventory_yml = "inventory.yml"
        cf = {
            "config_file": "__tmpcfg",
            "tasks_reports_dir": "__tmpdir",
            "ansible_dir": "__tmpdir",
            "tasks_pprof": "vtgate/cpu",
            "inventory_file": inventory_yml
        }
        self.setup(cf)
        task = self.benchmark_runner.tasks[0]
        task.build_ansible_inventory('HEAD')

        invf = open(task.ansible_built_inventory_filepath, 'r')
        invdata = yaml.load(invf, Loader=yaml.FullLoader)
        invf.close()
        self.assertEqual(["vtgate"], invdata["all"]["vars"]["pprof_targets"])
        self.assertEqual("cpu", invdata["all"]["vars"]["pprof_args"])


if __name__ == '__main__':
    unittest.main()