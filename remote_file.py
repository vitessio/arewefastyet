import pysftp
import os

def get_remote_oltp(hostname):
  # Hardcoded root
  myUsername = "root"
  #os.system('ssh-keygen -f "/root/.ssh/known_hosts" -R "'+hostname+'"')
  
  cnopts = pysftp.CnOpts()
  cnopts.hostkeys = None   

  with pysftp.Connection(host=hostname, username=myUsername, cnopts=cnopts) as sftp:
     print("Connection succesfully stablished ... ")
    
     # Define the file that you want to download from the remote directory
     remoteFilePath = '/tmp/oltp.json'

     # Define the local path where the file will be saved
     localFilePath = './report/oltp.json'

     sftp.get(remoteFilePath, localFilePath)
