## arewefastyet microbench run

Run micro benchmarks from the <root dir> on <pkg>, and outputs to <output file>.

### Synopsis

Runs all the micro benchmarks from the <root dir> on <pkg>, and parses the output and saves it to mysql if the configuration is provided. 
The output can also be outputted to <output file>.

```
arewefastyet microbench run [root dir] <pkg> <output file> [flags]
```

### Examples

```
arewastyet microbenchmark run ~/vitess ./go/vt/sqlparser output.txt
```

### Options

```
  -h, --help                                       help for run
      --microbench-exec-uuid string                UUID of the parent execution, an empty string will set to NULL.
      --planetscale-db-branch string               PlanetscaleDB branch to use. (default "main")
      --planetscale-db-database string             PlanetscaleDB database name.
      --planetscale-db-org string                  Name of the PlanetscaleDB organization.
      --planetscale-db-service-token string        PlanetscaleDB service token value.
      --planetscale-db-service-token-name string   PlanetscaleDB service token name.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/arewefastyet/config.yaml)
```

### SEE ALSO

* [arewefastyet microbench](arewefastyet_microbench.md)	 - Top level command to manage microbenchmarks

