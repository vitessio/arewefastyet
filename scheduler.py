import schedule
import time
import subprocess
import os
import yaml
import uuid 
import sys

def job():

    print("--------------- Init benchmark ----------------", end='\n')

    run_id = uuid.uuid4()
    os.system('python initialize_benchmark.py '+ str(run_id))

    print("--------------- Starting ansible ----------------", end='\n')

    #To avoid segmentation fault 

    with open('config.yaml') as f:
      data = yaml.load(f, Loader=yaml.FullLoader)

    os.system('./run '+ data["inventory_file"])

    print('------------- Adding results to the database ------------------', end='\n')

    os.system('python report.py ' + str(run_id))

    print("---------------------------------------------------------", end='\n')


# Runs everyday at <specified time>
schedule.every().day.at(sys.argv[1]).do(job)

while True:
    schedule.run_pending()
    time.sleep(1)
