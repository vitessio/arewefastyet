from packet_vps import create_vps
from config import inventory_file
import json
import os
import yaml
import sys

def doesFileExists(filePathAndName):
    return os.path.exists(filePathAndName)

def init():
    vps = create_vps(sys.argv[1])
    
    if doesFileExists('config-lock.json'):
      with open('config-lock.json') as json_file:
          data = json.load(json_file)
         
      data['run'].append({
        'run_id':sys.argv[1],
        'vps_id':vps[0],
        'ip_address':vps[1]
      })
     
      with open('config-lock.json', 'w') as outfile:
        json.dump(data, outfile)
    
    else:
       data = {}
       data['run'] = []
       data['run'].append({
        'run_id':sys.argv[1],
        'vps_id':vps[0],
        'ip_address':vps[1]
       })
       with open('config-lock.json', 'w') as outfile:
        json.dump(data, outfile)

    with open('ansible/'+inventory_file()) as f:
        data = yaml.load(f, Loader=yaml.FullLoader)
    
    # Changes ip address with new ip address
    data = recursive_dict(data,vps[1])

    print(data)
    
    with open('ansible/'+inventory_file() , 'w') as f:
      yaml.dump(data,f)
    

def recursive_dict(data,ip):
     for k, v in data.items():
        if isinstance(data[k],dict) and k == "hosts":
            data[k] = recursive_dict_ip(data[k],ip)
        elif isinstance(data[k],dict):
            data[k] = recursive_dict(data[k],ip)
     return data

def recursive_dict_ip(data,ip):
    for k, v in data.items():
        old_key = k
        data[ip] = data.pop(old_key)
    return data


init()