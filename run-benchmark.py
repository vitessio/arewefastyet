import time
import subprocess
import os
from config import inventory_file
import uuid 

def tasks():
   print('------------- Initialize VPS ------------------')
   run_id = uuid.uuid1()
   os.system('python initialize_benchmark.py ')
   print('------------- Running Benchamrk ------------------')
   os.system('./run '+ inventory_file())
   print('------------- Adding results to the database ------------------')
   os.system('python report.py ' + run_id)

tasks()
