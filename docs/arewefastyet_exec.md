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
arewefastyet exec --exec-git-ref 4a70d3d226113282554b393a97f893d133486b94  --planetscale-db-database benchmark --planetscale-db-branch main --planetscale-db-org my-org --planetscale-db-service-token <token> --planetscale-db-service-token-name <token name>
--exec-source config_micro_remote --ansible-inventory-file microbench_inventory.yml --ansible-playbook-file microbench.yml --ansible-root-directory ./ansible/
--equinix-instance-type m2.xlarge.x86 --equinix-token tok --equinix-project-id id

```

### Options

```
      --ansible-inventory-file string        Inventory file used by Ansible
      --ansible-playbook-file string         Playbook file used by Ansible
      --ansible-root-directory string        Root directory of Ansible
      --exec-git-ref string                  Git reference on which the benchmarks will run.
      --exec-go-version string               Defines the golang version that will be used by this execution. (default "1.17")
      --exec-pull-nb int                     Defines the number of the pull request against which to execute.
      --exec-root-dir string                 Path to the root directory of exec.
      --exec-schema string                   Path to the VSchema for this benchmark.
      --exec-server-address string           The IP address of the server on which the benchmark will be executed.
      --exec-source string                   Name of the source that triggered the execution.
      --exec-type string                     Defines the execution type (oltp, tpcc, micro).
      --exec-vtgate-planner-version string   Defines the vtgate planner version to use. Valid values are: V3, Gen4, Gen4Greedy and Gen4Fallback. (default "V3")
  -h, --help                                 help for exec
      --planetscale-db-branch string         PlanetscaleDB branch to use. (default "main")
      --planetscale-db-database string       PlanetscaleDB database name.
      --planetscale-db-host string           Hostname of the PlanetscaleDB database.
      --planetscale-db-org string            Name of the PlanetscaleDB organization.
      --planetscale-db-password string       Password used to authenticate to PlanetscaleDB.
      --planetscale-db-user string           Username used to authenticate to PlanetscaleDB.
      --stats-remote-db-database string      Name of the stats remote database.
      --stats-remote-db-host string          Hostname of the stats remote database.
      --stats-remote-db-password string      Password to authenticate the stats remote database.
      --stats-remote-db-port string          Port of the stats remote database.
      --stats-remote-db-user string          User used to connect to the stats remote database
```

### Options inherited from parent commands

```
      --config string    config file (default is $HOME/.config/arewefastyet/config.yaml)
      --secrets string   secrets file
```

### SEE ALSO

* [arewefastyet](arewefastyet.md)	 - Nightly Benchmarks Project

