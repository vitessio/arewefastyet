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
#   - Authenticate token  
#   - create a packet server (m2.xlarge.x86)
#   - deletes a packet server
# --------------------------------------------------------------------------------------------------------------------------------

import packet
from config import packet_token , packet_project_id
import time

# ---------------------------------------------------- Authenticate token ---------------------------------------------------------

def auth_packet():
    token = packet_token()
    return packet.Manager(auth_token=token)

# ----------------------------------------------------------------------------------------------------------------------------------
# ---------------------------------------------------- Creates packet server -------------------------------------------------------

def create_vps(id):
    project_id = packet_project_id()
    manager = auth_packet()

    device = manager.create_device(project_id=project_id,
                               hostname="benchmark-" + id,
                               plan="m2.xlarge.x86", facility='ams1',
                               operating_system='centos_8')
    
    while True:
            if manager.get_device(device.id).state == "active":
                break
            time.sleep(2)
    
    ips = manager.list_device_ips(device.id)

    return [device.id,ips[0].address]

# ------------------------------------------------------------------------------------------------------------------------------------
# ---------------------------------------------------- Deletes packet server ---------------------------------------------------------

def delete_vps(id):
    project_id = packet_project_id()
    manager = auth_packet()

    device = manager.get_device(id)
    device.delete()

# ------------------------------------------------------------------------------------------------------------------------------------



