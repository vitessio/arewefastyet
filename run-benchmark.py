import time
import subprocess
import os
from config import inventory_file
import uuid 

def tasks():
   print('------------- Initialize VPS ------------------')
   run_id = uuid.uuid4()
   os.system('python initialize_benchmark.py '+ str(run_id))
   print('------------- Running Benchamrk ------------------')
   os.system('./run '+ inventory_file())
   print('------------- Adding results to the database ------------------')
   os.system('python report.py ' + str(run_id))

tasks()
