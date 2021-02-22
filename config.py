# ------------------------------------------------------------------------------------------------------------------------------------
# Copyright 2020 The Vitess Authors.
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

from connection import connectdb
import yaml

# ------------------------------------------------- Reads config.yaml file ------------------------------------------------------------

def read_config():
    with open('config.yaml') as f:
      return yaml.load(f, Loader=yaml.SafeLoader)

# -------------------------------------------------------------------------------------------------------------------------------------
# -------------------------------------- Reads Mysql parameters and return connection object ------------------------------------------

def mysql_connect():
    data = read_config()
    return connectdb(data["mysql_host"],data["mysql_username"],data["mysql_password"],data["mysql_database"])

# -------------------------------------------------------------------------------------------------------------------------------------
# -------------------------------------- Returns Vitess commit hash used in inventory file --------------------------------------------

def vitess_git_version(inventory_file):
    data = read_config()

    with open(inventory_file) as f:
        data = yaml.load(f, Loader=yaml.SafeLoader)
        print(data)
    return data["all"]["vars"]["vitess_git_version"]

# -------------------------------------------------------------------------------------------------------------------------------------
# ------------------------------------------------------ Returns Packet token ---------------------------------------------------------

def packet_token():
    data = read_config()
    return data["packet_token"]

# -------------------------------------------------------------------------------------------------------------------------------------
# ---------------------------------------------------- Returns Packet project ID ------------------------------------------------------

def packet_project_id():
    data = read_config()
    return data["packet_project_id"]

# -------------------------------------------------------------------------------------------------------------------------------------
# ---------------------------------------------------- Returns Inventory file name ----------------------------------------------------


def inventory_file_default():
    data = read_config()
    global inventory_file
    inventory_file = str(data["inventory_file"])

def set_inventory_file(file):
    global inventory_file
    inventory_file = file

def get_inventory_file():
    return inventory_file

# Sets to default in config file
if 'inventory_file' not in globals():
    inventory_file_default()
    print(get_inventory_file())


# -------------------------------------------------------------------------------------------------------------------------------------
# -------------------------------------------------- Returns API key for flask server -------------------------------------------------

def api_key():
    data = read_config()
    return data["api_key"]

# -------------------------------------------------------------------------------------------------------------------------------------
# ----------------------------------------------------- Returns Slack API token -------------------------------------------------------

def slack_api_token():
    data = read_config()
    return data["slack_api_token"]

# -------------------------------------------------------------------------------------------------------------------------------------
# ---------------------------------------------------- Returns Slack channel name -----------------------------------------------------

def slack_channel():
    data = read_config()
    return data["slack_channel"]

# -------------------------------------------------------------------------------------------------------------------------------------
