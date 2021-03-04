# --------------------------------------------------------------------------------------------------------------------------------
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
#   - OLTP results to a database
#   - Sends the inventory file used and oltp.json as a slack message
#   - Deletes the packet baremetal server used to run the Ansibles
#
#  Arguments: python report.py <run id> <source> <oltp or tpcc>
# --------------------------------------------------------------------------------------------------------------------------------

import datetime
from slack import WebClient
from slack.errors import SlackApiError
import ssl

import bench_cli.configuration as configuration

# ----------------------------------------------------------------------------------------------------------------------------------
# ---------------------------------------------------- Send Slack Message ----------------------------------------------------------


def send_slack_message(slack_api_token: str, slack_channel: str, report_path: str):
    ssl._create_default_https_context = ssl._create_unverified_context

    client = WebClient(slack_api_token)

    # Upload OLTP file to slack
    try:
       response = client.files_upload(channels='#'+slack_channel, file=report_path)
       assert response["file"]  # the uploaded file
    except SlackApiError as e:
    # You will get a SlackApiError if "ok" is False
       assert e.response["ok"] is False
       assert e.response["error"]  # str like 'invalid_auth', 'channel_not_found'
       print(f"Got an error: {e.response['error']}")


# --------------------------------------------------------------------------------------------------------------------------------------
# -------------------------------------- Main function for report and add OLTP to database ---------------------------------------------

def save_to_mysql(cfg: configuration.Config, report, table_name: str):
    conn = cfg.mysql_connect()

    # source (https://www.w3schools.com/python/python_mysql_insert.asp)
    mycursor = conn.cursor()

    # current date and time
    now = datetime.datetime.now()

    format = '%Y-%m-%d %H:%M:%S'

    mysql_timestamp = now.strftime(format)

    benchmark = "INSERT INTO benchmark(commit,Datetime,source) values(%s,%s,%s)"
    mycursor.execute(benchmark, (cfg.commit, mysql_timestamp, cfg.source))
    conn.commit()

    mycursor.execute("select * from benchmark ORDER BY test_no DESC LIMIT 1;")
    result = mycursor.fetchall()
    test_no = result[0][0]

    # Inserting for table name
    sql_insert = "INSERT INTO "+table_name+" (time,threads,test_no,tps,latency,errors,reconnects) values(%s,%s,%s,%s,%s,%s,%s)"
    mycursor.execute(sql_insert, (report['results']["time"], report['results']["threads"], test_no, report['results']["tps"], report['results']["latency"], report['results']["errors"], report['results']["reconnects"]))
    conn.commit()

    # Get {{table_name}}_no
    mycursor.execute("select "+table_name+"_no from "+table_name+" where test_no = %s ORDER BY "+table_name+"_no DESC LIMIT 1;", (test_no,))
    result = mycursor.fetchall()
    task_id_res = result[0][0]

    # Inserting for {{table_name}}_qps
    insert_qps = "INSERT INTO qps("+table_name+"_no,total_qps,reads_qps,writes_qps,other_qps) values(%s,%s,%s,%s,%s)"
    mycursor.execute(insert_qps, (task_id_res, report['results']["qps"]["total"], report['results']["qps"]["reads"], report['results']["qps"]["writes"], report['results']["qps"]["other"]))
    conn.commit()