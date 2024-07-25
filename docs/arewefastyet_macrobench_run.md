## arewefastyet macrobench run

Run macro benchmarks and store the output in the mysql configuration provided.

### Synopsis

Run macro benchmarks using a fork of sysbench (https://github.com/planetscale/sysbench for OLTP and https://github.com/planetscale/sysbench-TPCC for TPCC)  and store the output in the mysql configuration provided.

```
arewefastyet macrobench run [flags]
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
      --macrobench-skip-steps string               Slice of sysbench steps to skip.
      --macrobench-sysbench-executable string      Path to the sysbench binary.
      --macrobench-type Type                       Type of macro benchmark.
      --macrobench-vtgate-planner-version string   Vtgate planner version running on Vitess
      --macrobench-vtgate-web-ports strings        List of the web port for each VTGate.
      --macrobench-working-directory string        Directory on which to execute sysbench.
      --macrobench-workload-path string            Path to the workload used by sysbench.
      --planetscale-db-branch string               PlanetScaleDB branch to use. (default "main")
      --planetscale-db-database string             PlanetScaleDB database name.
      --planetscale-db-host string                 Hostname of the PlanetScaleDB database.
      --planetscale-db-org string                  Name of the PlanetScaleDB organization.
      --planetscale-db-password-read string        Password used to authenticate to the read-only servers of PlanetScaleDB.
      --planetscale-db-password-write string       Password used to authenticate to the write servers of PlanetScaleDB.
      --planetscale-db-user-read string            Username used to authenticate to the read-only servers of PlanetScaleDB.
      --planetscale-db-user-write string           Username used to authenticate to the write servers of PlanetScaleDB.
```

### Options inherited from parent commands

```
      --config string    config file (default is $HOME/.config/arewefastyet/config.yaml)
      --secrets string   secrets file
```

### SEE ALSO

* [arewefastyet macrobench](arewefastyet_macrobench.md)	 - Top level command to manage macrobenchmarks

