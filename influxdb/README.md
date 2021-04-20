## InfluxDB

### Guide

Followed [this](https://www.influxdata.com/blog/deploying-influxdb-with-ansible/) guide to setup the playbook.

### Run

```
ansible-playbook --ssh-common-args "-o StrictHostKeyChecking=no" --user root  --inventory influxdb_inventory.yaml influxdb.yaml 
```