from flask import Flask ,request ,jsonify
import os
import datetime
from multiprocessing import Process
from connection import connectdb
from config import mysql_connect

app = Flask(__name__)

@app.route('/')
def run_benchmark():
    os.system('python run-benchmark.py')
    return ''

@app.route('/run_scheduler')
def nightly_bechmark():
    time = request.args.get('time')
    heavy_process = Process(  # Create a daemonic process with heavy scheduler
        target=scheduler(time),
        daemon=True
    )
    heavy_process.start()
    return time

def scheduler(time):
    process = os.system('python scheduler.py ' + time)
    print("Process finished")

@app.route('/servertime')
def server_time():
    return str(datetime.datetime.now())

# Returns all information in the database 
@app.route('/allresults')
def all_results():
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
    
   