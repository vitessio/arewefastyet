# ------------------------------------------------------------------------------------------------------------------------------------
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
#   - create a copy of inventory file for the run 
#   - calls create_vps to create packet server 
#   - add run information to config-lock.json file 
#   - changes the ip address in the copy of the inventory file 
#
# Arguments: python initialize_benchmark.py <run id> <commit hash>
# -------------------------------------------------------------------------------------------------------------------------------------

from packet_vps import create_vps
from config import inventory_file
from pathlib import Path
import json
import os
import yaml
import sys
import shutil
from get_head_hash import head_commit_hash

# ------------------------------------------------- Checks if file exists or not --------------------------------------------------------

def doesFileExists(filePathAndName):
    return os.path.exists(filePathAndName)

# ----------------------------------------------------------------------------------------------------------------------------------------
# ------------------------------------------------- Initializes benchmark process --------------------------------------------------------

def init(run_id, commit_hash):
    vps = create_vps(run_id)

    # make sure build folder is present
    if Path('./ansible/build').exists() == False:
        os.mkdir('./ansible/build')
    # create copy of inventory file
    shutil.copy2('./ansible/'+inventory_file(), './ansible/build/' + Path('./ansible/' + inventory_file()).stem + '-' + run_id + '.yml')
    
    if doesFileExists('config-lock.json'):
      with open('config-lock.json') as json_file:
          data = json.load(json_file)
         
      data['run'].append({
        'run_id':run_id,
        'vps_id':vps[0],
        'ip_address':vps[1]
      })
     
      with open('config-lock.json', 'w') as outfile:
        json.dump(data, outfile)
    
    else:
       data = {}
       data['run'] = []
       data['run'].append({
        'run_id':run_id,
        'vps_id':vps[0],
        'ip_address':vps[1]
       })
       with open('config-lock.json', 'w') as outfile:
        json.dump(data, outfile)

    with open('ansible/'+inventory_file()) as f:
        data = yaml.load(f, Loader=yaml.FullLoader)
    
    # Changes ip address with new ip address
    data = recursive_dict(data,vps[1])

    # if HEAD then get commit hash
    if commit_hash == 'HEAD':
       commit_hash = head_commit_hash()

    data["all"]["vars"]["vitess_git_version"] = commit_hash

    print(data)
    
    with open('ansible/build/' + Path('./ansible/' + inventory_file()).stem + '-' + run_id + '.yml' , 'w') as f:
      yaml.dump(data,f)
    
# ----------------------------------------------------------------------------------------------------------------------------------------
# ------------------------------------------------- Changes IP recursively ---------------------------------------------------------------

def recursive_dict(data,ip):
     for k, _ in data.items():
        if isinstance(data[k],dict) and k == "hosts":
            data[k] = recursive_dict_ip(data[k],ip)
        elif isinstance(data[k],dict):
            data[k] = recursive_dict(data[k],ip)
     return data

def recursive_dict_ip(data,ip):
    for k, _ in data.items():
        old_key = k
        data[ip] = data.pop(old_key)
    return data
