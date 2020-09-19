# -------------------------------------------------------------------------------------------------------------------------------------------------
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
#   - Run API server (must have api key in header eg: curl -X GET 'http://127.0.0.1:5000/allresults' -H 'api-key:b0wewer')
#       - /run [GET] paramter [commit=<commit hash>] - run benchmark and notify result on slack channel
#       - /run_scheduler [GET] paramter [time=<server time>] - run benchmark on specified time everyday and notify result
#                                                              on slack channel
#       - /servertime - returns server time
#       - /allresults - returns JSON of all benchmark results in the database
#       - /filter_result [GET] paramters [date=<reverse order for mysql>,commit=<commit hash>&commit=<commit hash>&...n,test_no=<int>]
#                     - filters and returns result based on argument given
#       - /webhook [POST] - triggers benchmark run on current HEAD (Called from github)        
#    - Web App routes
#       - / [GET] - returns home page
#       - /search_compare [GET] - returns search page
#       - /request_benchmark [GET] - returns request run for benchmark
#
#   Future code fix: Normalize code (Reduce Code duplication)
# -----------------------------------------------------------------------------------------------------------------------------------------------------

from flask import Flask ,request ,jsonify, render_template, Response
import os
import datetime
import uuid 
from multiprocessing import Process
from connection import connectdb
from config import mysql_connect,api_key,slack_api_token,slack_channel
from slack import WebClient
from slack.errors import SlackApiError
import ssl


app = Flask(__name__)

# ----------------------------------------------------------------- Render Home page -------------------------------------------------------------------

@app.route('/')
def home():
    data_oltp = graph_data('oltp')
    data_tpcc = graph_data('tpcc')
    return render_template("index.html", data_oltp=data_oltp, data_tpcc=data_tpcc)

# -------------------------------------------------------------------------------------------------------------------------------------------------------
# ------------------------------------------------------------- Render Search_Compare -------------------------------------------------------------------

@app.route('/search_compare')
def search_compare():
    searchcommit = request.args.get('search_commit')
    compare_commit_1 = request.args.get('compare_commit_1')
    compare_commit_2 = request.args.get('compare_commit_2')

    search_result = []
    search_result_tpcc = []
    compare_result_1 = []
    compare_result_tpcc_1 = []
    compare_result_2 = []
    compare_result_tpcc_2 = []
    
    # flag for empty result
    search_flag_empty = "No"
    search_flag_tpcc_empty = "No"
    compare_1_flag_empty = "No"
    compare_1_flag_tpcc_empty = "No"
    compare_2_flag_empty = "No"
    compare_2_flag_tpcc_empty = "No"

    if searchcommit != None:
       search_result = search_commit(searchcommit,'oltp')
       search_result_tpcc = search_commit(searchcommit,'tpcc')
       if search_result == None:
           search_flag_empty = "Yes"
           search_result = []
        
       if search_result_tpcc == None:
           search_flag_tpcc_empty = "Yes"
           search_result_tpcc = []
    
    if compare_commit_1 != None:
       compare_result_1 = search_commit(compare_commit_1,'oltp')
       compare_result_tpcc_1 = search_commit(compare_commit_1,'tpcc')
       if compare_result_1 == None:
           compare_1_flag_empty = "Yes"
           compare_result_1 = []
       if compare_result_tpcc_1 == None:
           compare_1_flag_tpcc_empty = "Yes"
           compare_result_tpcc_1 = []

    if compare_commit_2 != None:
       compare_result_2 = search_commit(compare_commit_2,'oltp')
       compare_result_tpcc_2 = search_commit(compare_commit_2,'tpcc')
       if compare_result_2 == None:
           compare_2_flag_empty = "Yes"
           compare_result_2 = []
       if compare_result_tpcc_2 == None:
           compare_2_flag_tpcc_empty = "Yes"
           compare_result_tpcc_2 = []

    
    # returns: search result, search result flag, compare result 1, compare result 1 flag, compare result 2, compare result flag 2
    return render_template("search_compare.html",search_result=search_result,search_result_tpcc=search_result_tpcc,search_flag_empty=search_flag_empty,
    search_flag_tpcc_empty=search_flag_tpcc_empty,compare_result_1=compare_result_1,compare_result_tpcc_1=compare_result_tpcc_1,
    compare_1_flag_empty=compare_1_flag_empty,compare_1_flag_tpcc_empty=compare_1_flag_tpcc_empty,
    compare_result_2=compare_result_2,compare_result_tpcc_2=compare_result_tpcc_2,compare_2_flag_empty=compare_2_flag_empty,compare_2_flag_tpcc_empty=compare_2_flag_empty)

# -------------------------------------------------------------------------------------------------------------------------------------------------------
# ----------------------------------------------------------- Render request_benchmark ------------------------------------------------------------------

@app.route('/request_benchmark')
def request_benchmark():
    name = request.args.get('name')
    commit_hash = request.args.get('commit_hash')
    email_id = request.args.get('email_id')
    
    Message = ""
    status = ""
    
    # Check if all arguments have a value and then send a slack message
    if name != None and commit_hash != None and email_id != None:
        ssl._create_default_https_context = ssl._create_unverified_context

        client = WebClient(slack_api_token())

        try:
          response = client.chat_postMessage(
            channel='#' + slack_channel(),
            text=""" Request Benchmark run 
            Name: """ + name + """
            Commit hash: """ + commit_hash + """
            Email ID: """ + email_id + """ """)

          #assert response["message"]["text"] == """ Request Benchmark run 
          #  Name: """ + name + """
          #  Commit hash: """ + commit_hash + """
          #  Email ID: """ + email_id + """ """
          Message = "Sent Succesfully"
          status = "success"

        except SlackApiError as e:
        # You will get a SlackApiError if "ok" is False
             assert e.response["ok"] is False
             assert e.response["error"]  # str like 'invalid_auth', 'channel_not_found'
             Message = f"Got an error: {e.response['error']}"
             status = "warning"

    
    return render_template("request_benchmark.html",message=Message,status=status)

# -------------------------------------------------------------------------------------------------------------------------------------------------------
# ---------------------------------------------------------------- runs benchmark -----------------------------------------------------------------------

@app.route('/run')
def run_benchmark():
    key = request.headers.get('api-key')

    if key == None:
        return "please add api key in header"

    if key != api_key():
        return "wrong api key"

    commit = request.args.get('commit')
    run_id = uuid.uuid4()

    os.system('./run-benchmark ' + commit + ' ' + str(run_id) + ' api_call &')
    #os.system('python run-benchmark.py ' + commit + ' ' + str(run_id) + ' api_call' + ' &')
    
    return 'Result will be updated on mysql database and you will be notified on slack <br> Run_id: ' + str(run_id)
        
# --------------------------------------------------------------------------------------------------------------------------------------------------------
# -------------------------------------------------------- runs benchmark based on time given ----------------------------------------------------------

@app.route('/run_scheduler')
def nightly_bechmark():
    key = request.headers.get('api-key')

    if key == None:
        return "please add api key in header"

    if key != api_key():
        return "wrong api key"

    time = request.args.get('time')
    os.system('./scheduler ' + time + ' &')
    #os.system('python scheduler.py ' + time + ' &')
    return 'benchmark will at server time ' + time + '. Result will be updated on mysql database and you will be notified on slack'

@app.route('/servertime')
def server_time():
    return str(datetime.datetime.now())

# --------------------------------------------------------------------------------------------------------------------------------------------------------------
# ----------------------------------------------------- Returns all information in the database -----------------------------------------------------------------

@app.route('/allresults')
def all_results():

    key = request.headers.get('api-key')

    if key == None:
        return "please add api key in header"

    if key != api_key():
        return "wrong api key"

    conn = mysql_connect()
    mycursor = conn.cursor()
    
    # Basic run info
    sql = "SELECT * FROM benchmark;"
    mycursor.execute(sql)

    benchmark = mycursor.fetchall()
    data = {}
    data['benchmark'] = [] 

    for i in range(len(benchmark)):
        oltp = []
        # Oltp information 

        sql = "SELECT * FROM OLTP where test_no = %s;"
        adr = (benchmark[i][0], )
        mycursor.execute(sql,adr)
        
        oltp_result = mycursor.fetchall()

        for j in range(len(oltp_result)):
           qps = []
           sql = "SELECT * FROM qps where OLTP_no = %s;"
           adr = (oltp_result[j][0], )
           mycursor.execute(sql,adr)

           qps_result = mycursor.fetchall()

           for k in range(len(qps_result)):
               qps.append({
                   'qps_no': qps_result[k][0],
                   'TPCC_no': qps_result[k][1],
                   'total_qps': str(qps_result[k][2]),
                   'reads_qps': str(qps_result[k][3]),
                   'writes_qps':str(qps_result[k][4]),
                   'other_qps': str(qps_result[k][5]),
                   'OLTP_no': qps_result[k][6]
               })

           oltp.append({
             'oltp_no': oltp_result[j][0],
             'test_no': oltp_result[j][1],
             'tps': str(oltp_result[j][2]),
             'latency': str(oltp_result[j][3]),
             'errors': str(oltp_result[j][4]),
             'reconnects': str(oltp_result[j][5]),
             'time': oltp_result[j][6],
             'threads': oltp_result[j][7],
             'qps': qps
           })

        data['benchmark'].append({
        'id':benchmark[i][0],
        'commit':benchmark[i][1],
        'datetime':str(benchmark[i][2]),
        'oltp':oltp
        })
        
    

    return jsonify(data)

# ----------------------------------------------------------------------------------------------------------------------------------------------------------------------
# ----------------------------------------------------------- Returns all information in the database ------------------------------------------------------------------

@app.route('/filter_results')
def filter_results():
    
    key = request.headers.get('api-key')

    if key == None:
        return "please add api key in header"

    if key != api_key():
        return "wrong api key"

    date = request.args.get('date')
    commit = request.args.getlist('commit')
    test_no = request.args.get('test_no') 

    conn = mysql_connect()
    mycursor = conn.cursor()
    
    commit = tuple(commit)
    
    if test_no != None:
        # Basic run info
        sql = "SELECT * FROM benchmark where test_no=%s;"
        adr = (test_no, )
        mycursor.execute(sql,adr)
    else:
       if date != None and commit != None:
           sql = 'SELECT * FROM benchmark where DateTime BETWEEN %s AND %s AND commit IN ("' + '","'.join(map(str, commit)) + '")'
           adr = (date + ' 00:00:00', date + ' 23:59:59',)
           mycursor.execute(sql,adr)
       elif date != None:
           sql = "SELECT * FROM benchmark where DateTime BETWEEN %s AND %s;"
           adr = (date + ' 00:00:00', date + ' 23:59:59',)
           mycursor.execute(sql,adr)
       elif commit != None:
           sql = 'SELECT * FROM benchmark where commit IN ("' + '","'.join(map(str, commit)) + '")'
           print(sql)
           mycursor.execute(sql)
       else:
           return 'use /allresults to view all results'
    
    
        
    benchmark = mycursor.fetchall()
    data = {}
    data['benchmark'] = [] 

    for i in range(len(benchmark)):
        oltp = []
        # Oltp information 

        sql = "SELECT * FROM OLTP where test_no = %s;"
        adr = (benchmark[i][0], )
        mycursor.execute(sql,adr)
        
        oltp_result = mycursor.fetchall()

        for j in range(len(oltp_result)):
           qps = []
           sql = "SELECT * FROM qps where OLTP_no = %s;"
           adr = (oltp_result[j][0], )
           mycursor.execute(sql,adr)

           qps_result = mycursor.fetchall()

           for k in range(len(qps_result)):
               qps.append({
                   'qps_no': qps_result[k][0],
                   'TPCC_no': qps_result[k][1],
                   'total_qps': str(qps_result[k][2]),
                   'reads_qps': str(qps_result[k][3]),
                   'writes_qps':str(qps_result[k][4]),
                   'other_qps': str(qps_result[k][5]),
                   'OLTP_no': qps_result[k][6]
               })

           oltp.append({
             'oltp_no': oltp_result[j][0],
             'test_no': oltp_result[j][1],
             'tps': str(oltp_result[j][2]),
             'latency': str(oltp_result[j][3]),
             'errors': str(oltp_result[j][4]),
             'reconnects': str(oltp_result[j][5]),
             'time': oltp_result[j][6],
             'threads': oltp_result[j][7],
             'qps': qps
           })

        data['benchmark'].append({
        'id':benchmark[i][0],
        'commit':benchmark[i][1],
        'datetime':str(benchmark[i][2]),
        'oltp':oltp
        })
        
    

    return jsonify(data)

# -----------------------------------------------------------------------------------------------------------------------------------------------------------------
# ----------------------------------------------------------- Returns results for the last 7 days ----------------------------------------------------------------

def graph_data(Type):
    conn = mysql_connect()
    mycursor = conn.cursor()

    sql = "SELECT test_no,commit,datetime FROM benchmark WHERE DateTime BETWEEN DATE(NOW()) - INTERVAL 7 DAY AND DATE(NOW()) AND source IN('scheduler','webhook') ORDER BY DateTime DESC;"
    mycursor.execute(sql)
        
    benchmark = mycursor.fetchall()
    data = {}
    data['benchmark'] = [] 

    print(len(benchmark))

    for i in range(len(benchmark)):
        oltp_tpcc = []
        # Oltp information 
        if Type == 'oltp':
          sql = "SELECT * FROM OLTP where test_no = %s;"
        elif Type == 'tpcc':
          sql = "SELECT * FROM TPCC where test_no = %s;"  

        adr = (benchmark[i][0], )
        mycursor.execute(sql,adr)
        
        oltp_tpcc_result = mycursor.fetchall()

        for j in range(len(oltp_tpcc_result)):
           qps = []
           if Type == 'oltp':
              sql = "SELECT * FROM qps where OLTP_no = %s;"
           elif Type == 'tpcc':
              sql = "SELECT * FROM qps where TPCC_no = %s;"
        
           adr = (oltp_tpcc_result[j][0], )
           mycursor.execute(sql,adr)

           qps_result = mycursor.fetchall()

           for k in range(len(qps_result)):

               if Type == 'oltp':
                qps.append({
                    'qps_no': qps_result[k][0],
                    'total_qps': str(qps_result[k][2]),
                    'reads_qps': str(qps_result[k][3]),
                    'writes_qps':str(qps_result[k][4]), 
                    'other_qps': str(qps_result[k][5]),
                    'OLTP_no': qps_result[k][6]
                })

               elif Type == 'tpcc':
                 qps.append({
                     'qps_no': qps_result[k][0],
                     'total_qps': str(qps_result[k][2]),
                     'reads_qps': str(qps_result[k][3]),
                     'writes_qps':str(qps_result[k][4]), 
                     'other_qps': str(qps_result[k][5]),
                     'TPCC_no': qps_result[k][6]
                 })
           if Type == 'oltp':
             oltp_tpcc.append({
               'oltp_no': oltp_tpcc_result[j][0],
               'test_no': oltp_tpcc_result[j][1],
               'tps': str(oltp_tpcc_result[j][2]),
               'latency': str(oltp_tpcc_result[j][3]),
               'errors': str(oltp_tpcc_result[j][4]),
               'reconnects': str(oltp_tpcc_result[j][5]),
               'time': oltp_tpcc_result[j][6],
               'threads': str(oltp_tpcc_result[j][7]),
               'qps': qps
             })

           elif Type == 'tpcc':
              oltp_tpcc.append({
               'tpcc_no': oltp_tpcc_result[j][0],
               'test_no': oltp_tpcc_result[j][1],
               'tps': str(oltp_tpcc_result[j][2]),
               'latency': str(oltp_tpcc_result[j][3]),
               'errors': str(oltp_tpcc_result[j][4]),
               'reconnects': str(oltp_tpcc_result[j][5]),
               'time': oltp_tpcc_result[j][6],
               'threads': str(oltp_tpcc_result[j][7]),
               'qps': qps
             })

        if Type == 'oltp':
         data['benchmark'].append({
         'id':benchmark[i][0],
         'commit':benchmark[i][1],
         'oltp':oltp_tpcc,
         'datetime':str(benchmark[i][2])
         })
        
        elif Type == 'tpcc':
         data['benchmark'].append({
         'id':benchmark[i][0],
         'commit':benchmark[i][1],
         'tpcc':oltp_tpcc,
         'datetime':str(benchmark[i][2])
         })
         
    return data

# -----------------------------------------------------------------------------------------------------------------------------------------------------------------
# ------------------------------------------------------------ Returns results based on commit hash ---------------------------------------------------------------

def search_commit(commit,Type):
    conn = mysql_connect()
    mycursor = conn.cursor()

    sql = "SELECT * FROM benchmark where commit=%s;"
    adr = (commit, )
    mycursor.execute(sql,adr)
        
    benchmark = mycursor.fetchall()
    data = {}
    data['benchmark'] = [] 

    for i in range(len(benchmark)):
        oltp_tpcc = []
        # Oltp information 
        if Type == 'oltp':
          sql = "SELECT * FROM OLTP where test_no = %s;"
        elif Type == 'tpcc':
          sql = "SELECT * FROM TPCC where test_no = %s;"  

        adr = (benchmark[i][0], )
        mycursor.execute(sql,adr)
        
        oltp_tpcc_result = mycursor.fetchall()

        for j in range(len(oltp_tpcc_result)):
           qps = []
           if Type == 'oltp':
              sql = "SELECT * FROM qps where OLTP_no = %s;"
           elif Type == 'tpcc':
              sql = "SELECT * FROM qps where TPCC_no = %s;"
        
           adr = (oltp_tpcc_result[j][0], )
           mycursor.execute(sql,adr)

           qps_result = mycursor.fetchall()

           for k in range(len(qps_result)):

               if Type == 'oltp':
                qps.append({
                    'qps_no': qps_result[k][0],
                    'total_qps': str(qps_result[k][2]),
                    'reads_qps': str(qps_result[k][3]),
                    'writes_qps':str(qps_result[k][4]), 
                    'other_qps': str(qps_result[k][5]),
                    'OLTP_no': qps_result[k][6]
                })

               elif Type == 'tpcc':
                 qps.append({
                     'qps_no': qps_result[k][0],
                     'total_qps': str(qps_result[k][2]),
                     'reads_qps': str(qps_result[k][3]),
                     'writes_qps':str(qps_result[k][4]), 
                     'other_qps': str(qps_result[k][5]),
                     'TPCC_no': qps_result[k][6]
                 })
           if Type == 'oltp':
             oltp_tpcc.append({
               'oltp_no': oltp_tpcc_result[j][0],
               'test_no': oltp_tpcc_result[j][1],
               'tps': str(oltp_tpcc_result[j][2]),
               'latency': str(oltp_tpcc_result[j][3]),
               'errors': str(oltp_tpcc_result[j][4]),
               'reconnects': str(oltp_tpcc_result[j][5]),
               'time': oltp_tpcc_result[j][6],
               'threads': str(oltp_tpcc_result[j][7]),
               'qps': qps
             })

           elif Type == 'tpcc':
              oltp_tpcc.append({
               'tpcc_no': oltp_tpcc_result[j][0],
               'test_no': oltp_tpcc_result[j][1],
               'tps': str(oltp_tpcc_result[j][2]),
               'latency': str(oltp_tpcc_result[j][3]),
               'errors': str(oltp_tpcc_result[j][4]),
               'reconnects': str(oltp_tpcc_result[j][5]),
               'time': oltp_tpcc_result[j][6],
               'threads': str(oltp_tpcc_result[j][7]),
               'qps': qps
             })

        if Type == 'oltp':
         data['benchmark'].append({
         'id':benchmark[i][0],
         'commit':benchmark[i][1],
         'oltp':oltp_tpcc,
         'datetime':str(benchmark[i][2])
         })
        
        elif Type == 'tpcc':
         data['benchmark'].append({
         'id':benchmark[i][0],
         'commit':benchmark[i][1],
         'tpcc':oltp_tpcc,
         'datetime':str(benchmark[i][2])
         })
         
        return data

# ----------------------------------------------------------------------------------------------------------------------------------------------------------------
# ------------------------------------------------------------ Triggers benchmark on every push ------------------------------------------------------------------

@app.route('/webhook', methods=['POST'])
def respond():

    if request.json["ref"] == "refs/heads/master":
       commit = 'HEAD'
       run_id = uuid.uuid4()
       os.system('./run-benchmark ' + commit + ' ' + str(run_id) + ' webhook &')

    return Response(status=200)

# ----------------------------------------------------------------------------------------------------------------------------------------------------------------