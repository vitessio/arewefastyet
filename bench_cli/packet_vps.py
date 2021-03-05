# --------------------------------------------------------------------------------------------------------------------------------
# Copyright 2021 The Vitess Authors.
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
import time
import uuid

# ---------------------------------------------------- Authenticate token ---------------------------------------------------------

def auth_packet(packet_token):
    return packet.Manager(auth_token=packet_token)

# ----------------------------------------------------------------------------------------------------------------------------------
# ---------------------------------------------------- Creates packet server -------------------------------------------------------

def create_vps(packet_token, packet_project_id, run_id: uuid.UUID):
    manager = auth_packet(packet_token)
    device = manager.create_device(project_id=packet_project_id,
                                   hostname="benchmark-" + run_id.__str__(),
                                   plan="m2.xlarge.x86", facility='ams1',
                                   operating_system='centos_8')

    while True:
        if manager.get_device(device.id).state == "active":
            break
        time.sleep(2)

    ips = manager.list_device_ips(device.id)
    return [device.id, ips[0].address]

# ------------------------------------------------------------------------------------------------------------------------------------
# ---------------------------------------------------- Deletes packet server ---------------------------------------------------------

def delete_vps(packet_token, device_id):
    manager = auth_packet(packet_token)

    device = manager.get_device(device_id)
    device.delete()

# ------------------------------------------------------------------------------------------------------------------------------------



