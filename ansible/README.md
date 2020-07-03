## Ansible Provisioner for Vitess

### Testing individual roles

Most roles in the playbook have been configured to be tested with molecule. These can be run by issuing

`molecule verify`

### Configure the system

Create an inventory with the following structure

```yaml
# Standard Ansible Inventory ref. https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html#inheriting-variable-values-group-variables-for-groups-of-groups
all:
  vars:
    # Default to HEAD
    vitess_git_version: "<THE VERSION YOU WANT>"
  hosts:
    # All Hosts must be listed here
    <host_ip>:
       # optional storage device
       # If configured we will attempt to partition and mount this
       # To serve as your mysql storage
       storage_device:
         device:
         partition:
  children:
    sysbench:
    prometheus:
    etcd:
    vtctld:
    vtgate:
      hosts:
        <host_ip>:
          gateways:
            - id:
              port:
              mysql_port
              grpc_port:
    vttablet:
      hosts:
        <host_ip>:
          tablets:
            - id:
              keyspace:
              shard:
```

### Run the Scripts

Given a configured inventory. Running a full provision and test can be done with the following command

`ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i test-inventory.yml full.yml -u root -e provision=True -e clean=True`

This will run delete any existing deployment and then run provision.yml, and configure.yml. A similar command is

`ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i test-inventory.yml full.yml -u root -e update=True`

This will update an existing deployment and the run the tests. It will not clean and re provision the cluster.

### Run Tests

`sysbench --luajit-cmd=off --threads=50 --time=300 --mysql-db=main --mysql-host=127.0.0.1 --mysql-port=3306 --db-ps-mode=disable --db-driver=mysql --report-interval=10 --auto-inc=off --tables=50 --table_size=5000000 --range_selects=0 --rand-type=uniform oltp_read_write prepare`

`sysbench --luajit-cmd=off --threads=50 --time=300 --mysql-db=main --mysql-host=127.0.0.1 --mysql-port=3306 --db-ps-mode=disable --db-driver=mysql --report-interval=10 --auto-inc=off --tables=50 --table_size=5000000  --range_selects=0 --rand-type=uniform oltp_read_write run`

### Get Performance Data

## Vtgate

`go tool pprof -seconds=120 -http ':8080' http://139.178.85.73:15001/debug/pprof/profile`

## Vttablet

`go tool pprof -seconds=120 -http ':8080' http://139.178.85.73:16001/debug/pprof/profile`
