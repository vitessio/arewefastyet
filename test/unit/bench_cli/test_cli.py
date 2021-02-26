import unittest
from .context import cli

class TestRunToTaskArray(unittest.TestCase):
    def test_all_run(self):
        result = ["oltp", "tpcc"]
        self.assertEqual(result, cli.run_to_task_array(all=True, oltp=False, tpcc=False))
        self.assertEqual(result, cli.run_to_task_array(all=True, oltp=True, tpcc=False))
        self.assertEqual(result, cli.run_to_task_array(all=True, oltp=False, tpcc=True))
        self.assertEqual(result, cli.run_to_task_array(all=True, oltp=True, tpcc=True))
        self.assertEqual(result, cli.run_to_task_array(all=False, oltp=True, tpcc=True))

    def test_no_run(self):
        result = []
        self.assertEqual(result, cli.run_to_task_array(all=False, oltp=False, tpcc=False))

    def test_oltp_run(self):
        result = ["oltp"]
        self.assertEqual(result, cli.run_to_task_array(all=False, oltp=True, tpcc=False))

    def test_tpcc_run(self):
        result = ["tpcc"]
        self.assertEqual(result, cli.run_to_task_array(all=False, oltp=False, tpcc=True))

if __name__ == '__main__':
    unittest.main()