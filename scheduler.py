import schedule
import time
import subprocess
import os
import yaml
import uuid 
import sys

def job():
    os.system('python run-benchmark.py')


# Runs everyday at <specified time>
schedule.every().day.at(sys.argv[1]).do(job)

while True:
    schedule.run_pending()
    time.sleep(1)
