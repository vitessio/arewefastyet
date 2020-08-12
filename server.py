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
#   - Run API server to (must have api key in header eg: curl -X GET 'http://127.0.0.1:5000/allresults' -H 'api-key:b0wewer')
#       - /run - run benchmark and notify result on slack channel
#       - /run_scheduler [GET] paramter [time=<server time>] - run benchmark on specified time everyday and notify result
#                                                              on slack channel
#       - /servertime - returns server time
#       - /allresults - returns JSON of all benchmark results in the database
#       - /filter_result [GET] paramters [date=<reverse order for mysql>,commit=<commit hash>&commit=<commit hash>&...n,test_no=<int>]
#                     - filters and returns result based on argument given
# -----------------------------------------------------------------------------------------------------------------------------------------------------

from flask import Flask ,request ,jsonify
import os
import datetime
from multiprocessing import Process
from connection import connectdb
from config import mysql_connect,api_key

app = Flask(__name__)

# ---------------------------------------------------------------- runs benchmark -----------------------------------------------------------------------

@app.route('/run')
def run_benchmark():
    key = request.headers.get('api-key')

    if key == None:
        return "please add api key in header"

    if key != api_key():
        return "wrong api key"
    
    os.system('python run-benchmark.py ' + '&')
    
    return 'Result will be updated on mysql database and you will be notified on slack'
        
# --------------------------------------------------------------------------------------------------------------------------------------------------------
# ---------------------------------------------------------- runs benchmark based on time given ----------------------------------------------------------

@app.route('/run_scheduler')
def nightly_bechmark():
    key = request.headers.get('api-key')

    if key == None:
        return "please add api key in header"

    if key != api_key():
        return "wrong api key"

    time = request.args.get('time')
    os.system('python scheduler.py ' + time + ' &')
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
        'datetime':benchmark[i][2],
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
        'datetime':benchmark[i][2],
        'oltp':oltp
        })
        
    

    return jsonify(data)

# -----------------------------------------------------------------------------------------------------------------------------------------------------------------