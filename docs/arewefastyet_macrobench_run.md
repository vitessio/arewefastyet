## arewefastyet macrobench run

Run macro benchmarks and store the output in the mysql configuration provided.

### Synopsis

Run macro benchmarks using a fork of sysbench (https://github.com/planetscale/sysbench for OLTP and https://github.com/planetscale/sysbench-TPCC for TPCC)  and store the output in the mysql configuration provided.

```
arewefastyet macrobench run [flags]
```

### Examples

```
arewastyet macrobenchmark run --planetscale-db-database benchmark --planetscale-db-branch main --planetscale-db-org my-org --planetscale-db-service-token <token> --planetscale-db-service-token-name <token name>
```

### Options

```
  -h, --help                                       help for run
      --influx-database string                     Name of the database to use in InfluxDB.
      --influx-hostname string                     Hostname of InfluxDB.
      --influx-password string                     Password used to connect to InfluxDB.
      --influx-port string                         Port on which to InfluxDB listens. (default "8086")
      --influx-username string                     Username used to connect to InfluxDB.
      --macrobench-exec-uuid string                UUID of the parent execution, an empty string will set to NULL.
      --macrobench-git-ref string                  Git SHA referring to the macro benchmark.
      --macrobench-skip-steps strings              Slice of sysbench steps to skip.
      --macrobench-source string                   The source or origin of the macro benchmark trigger.
      --macrobench-sysbench-executable string      Path to the sysbench binary.
      --macrobench-type Type                       Type of macro benchmark.
      --macrobench-vtgate-planner-version string   Vtgate planner version running on Vitess
      --macrobench-working-directory string        Directory on which to execute sysbench.
      --macrobench-workload-path string            Path to the workload used by sysbench.
      --planetscale-db-branch string               PlanetscaleDB branch to use. (default "main")
      --planetscale-db-database string             PlanetscaleDB database name.
      --planetscale-db-host string                 Hostname of the PlanetscaleDB database.
      --planetscale-db-org string                  Name of the PlanetscaleDB organization.
      --planetscale-db-password string             Password used to authenticate to PlanetscaleDB.
      --planetscale-db-user string                 Username used to authenticate to PlanetscaleDB.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/arewefastyet/config.yaml)
```

### SEE ALSO

* [arewefastyet macrobench](arewefastyet_macrobench.md)	 - Top level command to manage macrobenchmarks

