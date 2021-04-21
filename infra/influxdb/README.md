## InfluxDB

The following command will create a new influxdb on your remote host. The host IP can be changed in `influxdb_inventory.yaml`.
```
ansible-playbook --extra-vars '{"admin_password": ADMIN_PASSWORD, "prometheus_password": PROMETHEUS_PASSWORD}' --ssh-common-args "-o StrictHostKeyChecking=no" --user root  --inventory influxdb_inventory.yaml influxdb.yaml 
```