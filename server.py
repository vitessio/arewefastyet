from flask import Flask ,request
import os
import datetime
from multiprocessing import Process
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