# AreWeFastYet
Uses sysbench to run a benchmark on vitess. 
The following codebase:
1. Creates a VPS
2. Runs the ansibles on the VPS
3. reads results from the VPS and stores in the database
4. Kills the VPS

# Installation steps

### Requirements :
1. CentOS 8 or Ubuntu 
2. Python 3
3. Mysql Server
4. Packet API key

### Install python 3.7.x and set up virtual environment
Ubuntu: https://linuxize.com/post/how-to-install-python-3-7-on-ubuntu-18-04/
CentOS: https://tecadmin.net/install-python-3-7-on-centos-8/

#### Install virtual environment
Ubuntu: https://www.liquidweb.com/kb/how-to-setup-a-python-virtual-environment-on-ubuntu-18-04/
CentOS: https://www.liquidweb.com/kb/how-to-setup-a-python-virtual-environment-on-centos/

```
virtualenv -python=python3.7 benchmark
// activate virtual environment
source benchmark/bin/activate
```

### Install python libraries 
```
pip install -r requirements.txt
```
### Install Ansible 
https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html

### Configure Ansible
https://github.com/vitessio/arewefastyet/blob/modify-ansible/ansible/README.md
#### ansible dependencies
```
cd ansible
ansible-galaxy install cloudalchemy.node-exporter

```

### Create SSH key for ansible or use exsisting
https://docs.github.com/en/enterprise/2.15/user/articles/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent

### Build Vitess
https://vitess.io/docs/contributing/build-on-ubuntu/

### Change files to executables
```
chmod +x run-benchmark
chmod +x scheduler
```
### Create file config.yaml 
```
mysql_host: <mysql hostname>
mysql_username: <mysql username>
mysql_password: <mysql password>
mysql_database: <mysql database>
inventory_file: <inventory filename>
packet_token : <packet token>
packet_project_id : <packet project id>
api_key: <api key you want for the flask server>
slack_api_token: <slack_token>
slack_channel: <channel name>
```
Ex : 
```
mysql_host: localhost:3306
mysql_username: vitess
mysql_password: vitess123
mysql_database: vitess_benchmark
inventory_file: packet-inventory.yml
packet_token: bgRy8otJVWUmtpDDadSdSDSfgsAtY1xnRNg
packet_project_id: dba22084-f8c7-4aaf-9e0a-weSASFDd
api_key: db084-f8c7-4aaf-9e0a-waeasSd
slack_api_token: xoxb-12862423802725-128sdffsddSD6946-lzieR3PQXsdfsd2TmmFlpcQeb
slack_channel: benchmark
```
Inventory file from ansible directory

### Run Scheduler
```
python3 scheduler.py <time> & 
```
### Example run for API
```
curl -X GET 'http://127.0.0.1:5000/allresults' -H 'api-key:b084-f8c7-4aaf-9e0a-waeasSd'
```
Api routes : https://github.com/vitessio/arewefastyet/blob/master/server.py




