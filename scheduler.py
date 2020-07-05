import schedule
import time
import subprocess

def job():
    print("--------------- Starting ansible ----------------")
    rc = subprocess.call("./run-ansible")
    print("--------------- Ending ansible -----------------")

# Runs everyday at <specified time>
schedule.every().day.at("20:58").do(job)

while True:
    schedule.run_pending()
    time.sleep(1)
