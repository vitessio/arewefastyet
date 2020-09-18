# --------------------------------------------------------------------------------------------------------------------------------
# Copyright 2020 The Vitess Authors.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#    http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# demonstrates to:
#   - Getting OLTP and TPCC file from a remote server using a SSH server
# --------------------------------------------------------------------------------------------------------------------------------

import pysftp
import os

# ------------------------------------------------ get OLTP file from remote server -----------------------------------------------

def get_remote_oltp(hostname,Type):
  # Hardcoded root
  myUsername = "root"
  #os.system('ssh-keygen -f "/root/.ssh/known_hosts" -R "'+hostname+'"')
  
  cnopts = pysftp.CnOpts()
  cnopts.hostkeys = None   
  

  with pysftp.Connection(host=hostname, username=myUsername, cnopts=cnopts) as sftp:
     print("Connection succesfully stablished ... ")
    
     if Type == "oltp":
       # Define the file that you want to download from the remote directory
       remoteFilePath = '/tmp/oltp.json'

       # Define the local path where the file will be saved
       localFilePath = './report/oltp.json'

       sftp.get(remoteFilePath, localFilePath)

     elif Type == "tpcc":
       # Define the file that you want to download from the remote directory
       remoteFilePath = '/tmp/tpcc.json'

       # Define the local path where the file will be saved
       localFilePath = './report/tpcc.json'

       sftp.get(remoteFilePath, localFilePath)

# -----------------------------------------------------------------------------------------------------------------------------------