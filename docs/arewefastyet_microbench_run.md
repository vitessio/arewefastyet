## arewefastyet microbench run

Run micro benchmarks from the <root dir> on <pkg>, and outputs to <output file>.

### Synopsis

Runs all the micro benchmarks from the <root dir> on <pkg>, and parses the output and saves it to mysql if the configuration is provided. 
The output can also be outputted to <output file>.

```
arewefastyet microbench run [root dir] <pkg> <output file> [flags]
```

### Options

```
  -h, --help                                   help for run
      --microbench-exec-uuid string            UUID of the parent execution, an empty string will set to NULL.
      --microbench-run-profile                 Run goproc profiling for each micro-benchmark.
      --planetscale-db-branch string           PlanetscaleDB branch to use. (default "main")
      --planetscale-db-database string         PlanetscaleDB database name.
      --planetscale-db-host string             Hostname of the PlanetscaleDB database.
      --planetscale-db-org string              Name of the PlanetscaleDB organization.
      --planetscale-db-password-read string    Password used to authenticate to the read-only servers of PlanetScaleDB.
      --planetscale-db-password-write string   Password used to authenticate to the write servers of PlanetScaleDB.
      --planetscale-db-user-read string        Username used to authenticate to the read-only servers of PlanetScaleDB.
      --planetscale-db-user-write string       Username used to authenticate to the write servers of PlanetScaleDB.
```

### Options inherited from parent commands

```
      --config string    config file (default is $HOME/.config/arewefastyet/config.yaml)
      --secrets string   secrets file
```

### SEE ALSO

* [arewefastyet microbench](arewefastyet_microbench.md)	 - Top level command to manage microbenchmarks

