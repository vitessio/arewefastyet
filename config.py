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
    return data["all"]["vars"]["vitess_git_version"]
