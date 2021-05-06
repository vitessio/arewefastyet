## Monitoring Server

Creates and configures the monitoring server composed of a Grafana frontend, an InfluxDB server, and a Prometheus server.

### Configure

The host IP can be configured using the command:
```
sed -i.bak 's/HOST_IP_0/${MY_HOST_IP}/g' ./inventory.yaml
```

### Run

> The host's IP can be modified in `./inventory.yaml`

```
ansible-playbook --extra-vars '{"grafana_password": ${GF_PASS}, "influxdb_admin_password": ${INF_PASS},  "influxdb_prometheus_password": ${INF_PROM_PASS}, "prometheus_password": ${PROM_PASS}}' --ssh-common-args "-o StrictHostKeyChecking=no" --user root --inventory inventory.yaml playbook.yaml
```