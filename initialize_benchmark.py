from packet_vps import create_vps
from config import inventory_file
import uuid 
import json
import os
import yaml

def doesFileExists(filePathAndName):
    return os.path.exists(filePathAndName)

def init():
    vps = create_vps()
    data = {} 
    data['run'] = {
       'vps_id':vps[0],
       'ip_address':vps[1]
    }
     
    with open('config-lock.json', 'w') as outfile:
       json.dump(data, outfile)
    
    with open('ansible/'+inventory_file()) as f:
        data = yaml.load(f, Loader=yaml.FullLoader)

    for key in data:
        if isinstance(data[key],dict):
            print(1)
    

init()