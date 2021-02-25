import unittest
from .context import configuration

# default_cfg = configuration.create_cfg(
#     web=False, tasks=[], commit="HEAD", source="testing",
#     inventory_file="inventory_file", mysql_host="localhost", mysql_username="root",
#     mysql_password="password", mysql_database="main", packet_token="12345",
#     packet_project_id="AABB11", api_key="123-ABC-456-EFG", slack_api_token="slack-token",
#     slack_channel="general", config_file="./config", ansible_dir="./ansible",
#     tasks_scripts_dir="./scripts", tasks_reports_dir="./reports"
# )

default_cfg_fields = {
    "web": False, "tasks": [], "commit": "HEAD", "source": "testing",
    "inventory_file": "inventory_file", "mysql_host": "localhost",
    "mysql_username": "root", "mysql_password": "password",
    "mysql_database": "main", "packet_token": "12345",
    "packet_project_id": "AABB11", "api_key": "123-ABC-456-EFG",
    "slack_api_token": "slack-token", "slack_channel": "general",
    "config_file": "./config", "ansible_dir": "./ansible",
    "tasks_scripts_dir": "./scripts", "tasks_reports_dir": "./reports"
}

class TestCreateConfig(unittest.TestCase):
    def test_create_default_config(self):
        cfg = configuration.create_cfg(default_cfg_fields["web"], default_cfg_fields["tasks"], default_cfg_fields["commit"],
                                       default_cfg_fields["source"], default_cfg_fields["inventory_file"], default_cfg_fields["mysql_host"],
                                       default_cfg_fields["mysql_username"], default_cfg_fields["mysql_password"], default_cfg_fields["mysql_database"],
                                       default_cfg_fields["packet_token"], default_cfg_fields["packet_project_id"], default_cfg_fields["api_key"],
                                       default_cfg_fields["slack_api_token"], default_cfg_fields["slack_channel"], default_cfg_fields["config_file"],
                                       default_cfg_fields["ansible_dir"], default_cfg_fields["tasks_scripts_dir"], default_cfg_fields["tasks_reports_dir"])
        assert cfg.__eq__(default_cfg_fields)

if __name__ == '__main__':
    unittest.main()