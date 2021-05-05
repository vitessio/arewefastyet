## Monitoring Server

Creates and configures the monitoring server composed of a Grafana frontend, an InfluxDB server, and a Prometheus server.

### Run

> The host's IP can be modified in `./inventory.yaml`

```
ansible-playbook --extra-vars '{"grafana_password": ${GF_PASS}, "influxdb_admin_password": ${INF_PASS},  "influxdb_prometheus_password": ${INF_PROM_PASS}, "prometheus_password": ${PROM_PASS}}' --ssh-common-args "-o StrictHostKeyChecking=no" --user root --inventory inventory.yaml playbook.yaml
```