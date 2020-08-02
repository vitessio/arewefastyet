import packet
from config import packet_token , packet_project_id
import time


def auth_packet():
    token = packet_token()
    return packet.Manager(auth_token=token)

def create_vps():
    project_id = packet_project_id()
    manager = auth_packet()

    device = manager.create_device(project_id=project_id,
                               hostname="benchmark-1",
                               plan="m2.xlarge.x86", facility='ams1',
                               operating_system='centos_8')
    
    while True:
            if manager.get_device(device.id).state == "active":
                break
            time.sleep(2)
    
    ips = manager.list_device_ips(device.id)

    return [device.id,ips[0].address]

def delete_vps(id):
    project_id = packet_project_id()
    manager = auth_packet()

    device = manager.get_device(id)
    device.delete()



