import schedule
import time
import subprocess
import os
import yaml
import uuid 

def job():

    run_id = uuid.uuid4()
    os.system('python initialize_benchmark.py '+ str(run_id))

    print("--------------- Starting ansible ----------------")

    #To avoid segmentation fault 

    with open('config.yaml') as f:
      data = yaml.load(f, Loader=yaml.FullLoader)

    os.system('./run '+ data["inventory_file"])

    print('------------- Adding results to the database ------------------')
    
    os.system('python report.py ' + str(run_id))


    print("---------------------------------------------------------")


# Runs everyday at <specified time>
schedule.every().day.at("14:50").do(job)

while True:
    schedule.run_pending()
    time.sleep(1)
