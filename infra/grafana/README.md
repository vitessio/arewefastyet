## Grafana

### Requirement

Please install the [`grafana`](https://github.com/cloudalchemy/ansible-grafana) role first.
```
ansible-galaxy install cloudalchemy.grafana
```

### Run

The following command will spin up a Grafana server on your host. The host IP can be changed in `grafana_inventory.yaml`.
```
ansible-playbook --extra-vars '{"grafana_password": GRAFANA_PASSWORD}' --ssh-common-args "-o StrictHostKeyChecking=no" --user root --inventory grafana_inventory.yaml grafana.yaml 
```