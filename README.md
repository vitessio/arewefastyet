# AreWeFastYet
Uses sysbench to benchmark vitess and also includes a schduler runs once a day

# Installation steps

### Requirements :
1. Ubuntu or CentOS 8 
2. Python 3
3. Mysql Server 

### Install python libraries 
```
pip3 install -r requirement.txt
```
### Install Ansible 
https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html

### Configure Ansible
https://github.com/vitessio/arewefastyet/blob/modify-ansible/ansible/README.md

### Create SSH key for ansible or use exsisting
https://docs.github.com/en/enterprise/2.15/user/articles/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent


### Create file config.yaml 
Ex : 
```
mysql_host: localhost:3306
mysql_username: vitess
mysql_password: vitess123
inventory_file: packet-inventory.yml
```
Inventory file from ansible directory

### Run Scheduler
```
python3 scheduler.py & > process.txt 
```




