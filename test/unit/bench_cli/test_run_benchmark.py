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
from .context import configuration, run_benchmark, taskfac

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
    "tasks_pprof": None, "delete_benchmark": False
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
        self.task_factory = taskfac.TaskFactory()

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
        self.tmpdir = tempfile.mkdtemp()
        self.cfg = {
            "run_all": True,
            "tasks_reports_dir": self.tmpdir,
            "ansible_dir": self.tmpdir,
            "source": self.tmpdir,
            "tasks_pprof": "vttablet/mem",
            "inventory_file": "inv.yaml"
        }

    def tearDown(self) -> None:
        super().tearDown()

    def test_create_benchmark_runner_all_tasks(self):
        self.cfg["run_all"] = True
        self.cfg["run_oltp"] = False
        self.cfg["run_tpcc"] = False
        config = configuration.Config(self.cfg)
        benchmark_runner = run_benchmark.BenchmarkRunner(config)
        self.assertEqual(2, len(benchmark_runner.tasks))
        for i, task in enumerate(benchmark_runner.tasks):
            self.assertEqual(config.tasks[i], task.name())

    def test_create_benchmark_runner_oltp(self):
        self.cfg["run_all"] = False
        self.cfg["run_oltp"] = True
        self.cfg["run_tpcc"] = False
        config = configuration.Config(self.cfg)
        benchmark_runner = run_benchmark.BenchmarkRunner(config)
        self.assertEqual(1, len(benchmark_runner.tasks))
        self.assertEqual(config.tasks[0], benchmark_runner.tasks[0].name())

    def test_create_benchmark_runner_tpcc(self):
        self.cfg["run_all"] = False
        self.cfg["run_oltp"] = False
        self.cfg["run_tpcc"] = True
        config = configuration.Config(self.cfg)
        benchmark_runner = run_benchmark.BenchmarkRunner(config)
        self.assertEqual(1, len(benchmark_runner.tasks))
        self.assertEqual(config.tasks[0], benchmark_runner.tasks[0].name())


class TestCreationOfTaskCheckValues(unittest.TestCase):
    def test_create_task_check_values(self):
        task_factory = taskfac.TaskFactory()
        tcs = [
            {"source": "unit_test","inventory_file": "inv_file","task_name": "oltp", "tasks_reports_dir": "./report", "ansible_dir": "./ansible"},
            {"source": "unit_test", "inventory_file": "inv_file", "task_name": "tpcc", "tasks_reports_dir": "./report", "ansible_dir": "./ansible"}
        ]
        for tc in tcs:
            task = task_factory.create_task(tc["task_name"], tc["tasks_reports_dir"], tc["ansible_dir"], tc["inventory_file"], tc["source"], None)

            self.assertEqual(tc["task_name"], task.name())
            self.assertEqual(tc["tasks_reports_dir"], task.report_dir)
            self.assertEqual(tc["ansible_dir"], task.ansible_dir)
            self.assertEqual(tc["inventory_file"], task.ansible_inventory_file)
            self.assertEqual(tc["source"], task.source)

            expected_ansible_build_dir = os.path.join(tc["ansible_dir"], 'build')
            self.assertEqual(expected_ansible_build_dir, task.ansible_build_dir)

            expected_ansible_built_file = tc["inventory_file"].split('.')[0] + '-' + str(task.task_id) + '.yml'
            self.assertEqual(expected_ansible_built_file, task.ansible_built_inventory_file)
            self.assertEqual(os.path.join(expected_ansible_build_dir, expected_ansible_built_file), task.ansible_built_inventory_filepath)

            self.assertEqual(task.name().upper(), task.table_name())
            self.assertEqual(os.path.join(tc["tasks_reports_dir"], task.name() + "_v2.json"), task.report_path())
            self.assertEqual(os.path.join("./", task.name() + "_v2.json"), task.report_path("./"))

    def test_create_task_with_benchmark_runner_check_values(self):
        tcs = [
            {"name": "run_oltp"},
            {"name": "run_tpcc"}
        ]
        for tc in tcs:
            cfg = default_cfg_fields.copy()
            cfg[tc.get("name")] = True
            cfg.__delitem__("config_file")
            config = configuration.Config(cfg)
            benchmark_runner = run_benchmark.BenchmarkRunner(config)

            task = benchmark_runner.tasks[0]

            self.assertEqual(config.tasks[0], task.name())
            self.assertEqual(config.tasks_reports_dir, task.report_dir)
            self.assertEqual(config.ansible_dir, task.ansible_dir)
            self.assertEqual(config.get_inventory_file_path(), task.ansible_inventory_file)
            self.assertEqual(config.source, task.source)

            expected_ansible_build_dir = os.path.join(config.ansible_dir, 'build')
            self.assertEqual(expected_ansible_build_dir, task.ansible_build_dir)

            expected_ansible_built_file = config.inventory_file.split('.')[0] + '-' + str(task.task_id) + '.yml'
            self.assertEqual(expected_ansible_built_file, task.ansible_built_inventory_file)
            self.assertEqual(os.path.join(expected_ansible_build_dir, expected_ansible_built_file), task.ansible_built_inventory_filepath)

            self.assertEqual(task.name().upper(), task.table_name())
            self.assertEqual(os.path.join(config.tasks_reports_dir, task.name() + "_v2.json"), task.report_path())
            self.assertEqual(os.path.join("./", task.name() + "_v2.json"), task.report_path("./"))


class TestBuildAnsibleInventoryFile(unittest.TestCase):
    def setup_inventory(self, filepath, inv_data):
        f = open(filepath, "w+")
        f.write(inv_data)
        f.close()

    def test_build_ansible_inventory_created(self):
        tmpdir = tempfile.mkdtemp()
        inventory_yml = "inventory.yml"
        config = configuration.Config({"source": "unit_test", "inventory_file": inventory_yml, "run_oltp": True, "tasks_reports_dir": tmpdir, "ansible_dir": tmpdir})
        self.setup_inventory(config.get_inventory_file_path(), sample_inv_file)
        benchmark_runner = run_benchmark.BenchmarkRunner(config)
        task = benchmark_runner.tasks[0]
        task.build_ansible_inventory('HEAD')

        exptected_path = os.path.join(config.ansible_dir, "build", inventory_yml.split('.')[0] + '-' + task.task_id.__str__() + ".yml")
        self.assertEqual(exptected_path, task.ansible_built_inventory_filepath)
        self.assertTrue(os.path.exists(exptected_path))

    def test_build_ansible_inventory_pprof(self):
        tmpdir = tempfile.mkdtemp()
        inventory_yml = "inventory.yml"
        config = configuration.Config({"tasks_pprof": "vtgate/cpu", "source": "unit_test", "inventory_file": inventory_yml, "run_oltp": True, "tasks_reports_dir": tmpdir, "ansible_dir": tmpdir})
        self.setup_inventory(config.get_inventory_file_path(), sample_inv_file)
        benchmark_runner = run_benchmark.BenchmarkRunner(config)
        task = benchmark_runner.tasks[0]
        task.build_ansible_inventory('HEAD')

        invf = open(task.ansible_built_inventory_filepath, 'r')
        invdata = yaml.load(invf, Loader=yaml.FullLoader)
        invf.close()

        self.assertEqual(["vtgate"], invdata["all"]["vars"]["pprof_targets"])
        self.assertEqual("cpu", invdata["all"]["vars"]["pprof_args"])

    def test_build_ansible_inventory_commit_is_pr(self):
        tmpdir = tempfile.mkdtemp()
        inventory_yml = "inventory.yml"
        commit = "1"  # represents pull request #1

        config = configuration.Config({"tasks_pprof": "vtgate/cpu", "source": "unit_test", "inventory_file": inventory_yml, "run_oltp": True, "tasks_reports_dir": tmpdir, "ansible_dir": tmpdir})
        self.setup_inventory(config.get_inventory_file_path(), sample_inv_file)
        benchmark_runner = run_benchmark.BenchmarkRunner(config)
        task = benchmark_runner.tasks[0]
        task.build_ansible_inventory(commit)

        invf = open(task.ansible_built_inventory_filepath, 'r')
        invdata = yaml.load(invf, Loader=yaml.FullLoader)
        invf.close()

        self.assertEqual(1, invdata["all"]["vars"]["vitess_git_version_pr_nb"])
        self.assertEqual("pull/1/head:1", invdata["all"]["vars"]["vitess_git_version_fetch_pr"])
        self.assertEqual("21a4f62c614f19f6717e6161ec049628aa119f52", invdata["all"]["vars"]["vitess_git_version"])  # SHA taken from https://github.com/vitessio/vitess/pull/1/commits/21a4f62c614f19f6717e6161ec049628aa119f52


if __name__ == '__main__':
    unittest.main()