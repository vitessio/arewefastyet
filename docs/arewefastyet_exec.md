## arewefastyet exec

Execute a task

### Synopsis

Execute a task based on the given terraform and ansible configuration.
It handles the creation, configuration, and cleanup of the infrastructure.

```
arewefastyet exec [flags]
```

### Examples

```
arewefastyet exec --exec-git-ref 4a70d3d226113282554b393a97f893d133486b94  --db-database benchmark --db-host localhost --db-password <db-password> --db-user <db-username>
--exec-source config_micro_remote --ansible-inventory-files microbench_inventory.yml --ansible-playbook-files microbench.yml --ansible-root-directory ./ansible/
--equinix-instance-type m2.xlarge.x86 --equinix-token tok --equinix-project-id id

```

### Options

```
      --ansible-inventory-files strings      List of inventory files used by Ansible
      --ansible-playbook-files strings       List of playbook files used by Ansible
      --ansible-root-directory string        Root directory of Ansible
      --db-database string                   Database to use.
      --db-host string                       Hostname of the database
      --db-password string                   Password to authenticate the database.
      --db-user string                       User used to connect to the database
      --equinix-instance-type string         Instance type to use for the creation of a new node
      --equinix-project-id string            Project ID to use for Equinix Metal
      --equinix-token string                 Auth Token for Equinix Metal
      --exec-git-ref string                  Git reference on which the benchmarks will run.
      --exec-pull-nb int                     Defines the number of the pull request against which to execute.
      --exec-root-dir string                 Path to the root directory of exec.
      --exec-source string                   Name of the source that triggered the execution.
      --exec-type string                     Defines the execution type (oltp, tpcc, micro).
      --exec-vtgate-planner-version string   Defines the vtgate planner version to use. Valid values are: V3, Gen4, Gen4Greedy and Gen4Fallback. (default "V3")
  -h, --help                                 help for exec
      --infra-path string                    Path to the infra directory
      --stats-remote-db-database string      Name of the stats remote database.
      --stats-remote-db-host string          Hostname of the stats remote database.
      --stats-remote-db-password string      Password to authenticate the stats remote database.
      --stats-remote-db-port string          Port of the stats remote database.
      --stats-remote-db-user string          User used to connect to the stats remote database
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/arewefastyet/config.yaml)
```

### SEE ALSO

* [arewefastyet](arewefastyet.md)	 - Nightly Benchmarks Project

