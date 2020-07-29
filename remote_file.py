import pysftp


def get_remote_oltp():
  myHostname = "147.75.80.149"
  myUsername = "root"

  with pysftp.Connection(host=myHostname, username=myUsername) as sftp:
     print("Connection succesfully stablished ... ")
    
     # Define the file that you want to download from the remote directory
     remoteFilePath = '/tmp/oltp.json'

     # Define the local path where the file will be saved
     localFilePath = './report/oltp.json'

     sftp.get(remoteFilePath, localFilePath)