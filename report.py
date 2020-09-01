# --------------------------------------------------------------------------------------------------------------------------------
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
#   - OLTP results to a database 
#   - Sends the inventory file used and oltp.json as a slack message
#   - Deletes the packet baremetal server used to run the Ansibles 
#
#  Arguments: python report.py <run id> <source> <oltp or tpcc>
# --------------------------------------------------------------------------------------------------------------------------------

import datetime
from connection import connectdb
from config import mysql_connect,vitess_git_version,slack_api_token,slack_channel,inventory_file
from remote_file import get_remote_oltp
from packet_vps import delete_vps
from slack import WebClient
from slack.errors import SlackApiError
from pathlib import Path
import os
import json
import sys
import ssl


# ------------------------------------------------------ IP and Project ID --------------------------------------------------------

def get_ip_and_project_id(run_id):
  with open('config-lock.json') as json_file:
    data = json.load(json_file)
  
  print(data)
  for i in data['run']:
    if i['run_id'] == run_id:
      l = [i['ip_address'],i['vps_id']]
      data['run'].remove(i)
      with open('config-lock.json', 'w') as outfile:
       json.dump(data, outfile)
      return l
  return None

# ----------------------------------------------------------------------------------------------------------------------------------
# ---------------------------------------------------- Send Slack Message ----------------------------------------------------------

def send_slack_message():
    ssl._create_default_https_context = ssl._create_unverified_context
   
    client = WebClient(slack_api_token())

    # Upload OLTP file to slack
    try:
       filepath="./report/oltp.json"
       response = client.files_upload(
         channels='#'+slack_channel(),
         file=filepath)
       assert response["file"]  # the uploaded file 
    except SlackApiError as e:
    # You will get a SlackApiError if "ok" is False
       assert e.response["ok"] is False
       assert e.response["error"]  # str like 'invalid_auth', 'channel_not_found'
       print(f"Got an error: {e.response['error']}")

# -------------------------------------------------------------------------------------------------------------------------------------
# -------------------------------------------- Adds run_id , vitess_version , source --------------------------------------------------

def add_data_oltp_report(run_id,source):
    with open('report/oltp.json') as f:
        data = json.load(f)
    data[0]["run_id"] = run_id
    data[0]["source"] = source
    data[0]["commit_id"] = vitess_git_version(Path('./ansible/' + inventory_file()).stem + '-' + run_id + '.yml')
    with open('report/oltp.json', 'w') as outfile:
       json.dump(data, outfile)
    
# --------------------------------------------------------------------------------------------------------------------------------------
# -------------------------------------------------- Remove Inventory file -------------------------------------------------------------

def remove_inventory_file(id):
    os.remove('./ansible/' + Path('./ansible/' + inventory_file()).stem + '-' + id + '.yml')

# --------------------------------------------------------------------------------------------------------------------------------------
# -------------------------------------- Main function for report and add OLTP to database ---------------------------------------------

def add_oltp():
    # Read the argument for the run id and source
    run_id = sys.argv[1]
    source = sys.argv[2]

    config_lock = get_ip_and_project_id(run_id)

    get_remote_oltp(config_lock[0])
    
    # Add run_id , source , vitess_git_version to oltp.json
    add_data_oltp_report(run_id,source)

    # Send report file
    send_slack_message()

    # local variable db connection object
    conn = mysql_connect()

    # source (https://www.w3schools.com/python/python_mysql_insert.asp)
    mycursor = conn.cursor()

    # current date and time
    now = datetime.datetime.now()

    format = '%Y-%m-%d %H:%M:%S'

    mysql_timestamp = now.strftime(format)


    # Sets data varaible to null
    data = None

    benchmark = "INSERT INTO benchmark(commit,Datetime,source) values(%s,%s,%s)"
    mycursor.execute(benchmark, (vitess_git_version(Path('./ansible/' + inventory_file()).stem + '-' + run_id + '.yml'),mysql_timestamp),source)
    conn.commit()

    mycursor.execute("select *from benchmark ORDER BY test_no DESC LIMIT 1;")
    result = mycursor.fetchall()
    test_no = result[0][0]


    ## TODO: replace with calling ssh server and get remote file
    with open('report/oltp.json') as f:
      data = json.load(f)

    # Inserting for oltp
    oltp = "INSERT INTO OLTP(time,threads,test_no,tps,latency,errors,reconnects) values(%s,%s,%s,%s,%s,%s,%s)"
    mycursor.execute(oltp,(data[0]["time"],data[0]["threads"],test_no,data[0]["tps"],data[0]["latency"],data[0]["errors"],data[0]["reconnects"]))
    conn.commit()

    #get oltp_no
    mycursor.execute("select OLTP_no from OLTP where test_no = %s ORDER BY OLTP_no DESC LIMIT 1;",(test_no,))
    result = mycursor.fetchall()
    oltp_no = result[0][0]

    # Inserting for oltp_qps
    oltp_qps = "INSERT INTO qps(OLTP_no,total_qps,reads_qps,writes_qps,other_qps) values(%s,%s,%s,%s,%s)"
    mycursor.execute(oltp_qps,(oltp_no,data[0]["qps"]["total"],data[0]["qps"]["reads"],data[0]["qps"]["writes"],data[0]["qps"]["other"]))
    conn.commit()

    # Delete vps instance
    delete_vps(config_lock[1])
    
    # remove inventory file
    remove_inventory_file(run_id)

    return test_no

# -------------------------------------------------------------------------------------------------------------------------------------------------

if sys.argv[3] == "oltp":
   add_oltp()
