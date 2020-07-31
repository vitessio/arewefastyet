from connection import connectdb
import yaml

def read_config():
    with open('config.yaml') as f:
      return yaml.load(f, Loader=yaml.FullLoader)

def mysql_connect():
    data = read_config()
    return connectdb(data["mysql_host"],data["mysql_username"],data["mysql_password"],"vitess_benchmark")

def vitess_git_version():
    data = read_config()

    with open('ansible/'+data["inventory_file"]) as f:
        data = yaml.load(f, Loader=yaml.FullLoader)
        print(data)
    return data["all"]["vars"]["vitess_git_version"]

def packet_token():
    data = read_config()
    return data["packet_token"]

def packet_project_id():
    data = read_config()
    return data["packet_project_id"]

def inventory_file():
    data = read_config()
    return data["inventory_file"]