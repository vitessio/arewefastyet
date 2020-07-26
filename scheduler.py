import schedule
import time
import subprocess
import os

def job():
    print("--------------- Starting ansible partial ----------------")

    #TODO : Changed to the ansible bash script
    os.system('./run ')

    print("--------------- Adding reports to MySql -----------------")

    # Not calling method directly due to segmentation fault
    os.system('python3 report.py')


    print("---------------------------------------------------------")


# Runs everyday at <specified time>
schedule.every().day.at("14:50").do(job)

while True:
    schedule.run_pending()
    time.sleep(1)
