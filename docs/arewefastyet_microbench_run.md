## arewefastyet microbench run

Run micro benchmarks from the <root dir> on <pkg>, and outputs to <output path>.

```
arewefastyet microbench run [root dir] <pkg> <output path> [flags]
```

### Options

```
      --db-database string            Database to use.
      --db-host string                Hostname of the database
      --db-password string            Password to authenticate the database.
      --db-user string                User used to connect to the database
  -h, --help                          help for run
      --microbench-exec-uuid string   UUID of the parent execution, an empty string will set to NULL.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/arewefastyet/config.yaml)
```

### SEE ALSO

* [arewefastyet microbench](arewefastyet_microbench.md)	 - 

