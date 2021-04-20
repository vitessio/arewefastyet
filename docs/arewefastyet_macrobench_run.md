## arewefastyet macrobench run



```
arewefastyet macrobench run [flags]
```

### Options

```
      --db-database string                      Database to use.
      --db-host string                          Hostname of the database
      --db-password string                      Password to authenticate the database.
      --db-user string                          User used to connect to the database
  -h, --help                                    help for run
      --macrobench-exec-uuid string             UUID of the parent execution, an empty string will set to NULL.
      --macrobench-git-ref string               Git SHA referring to the macro benchmark.
      --macrobench-skip-steps strings           Slice of sysbench steps to skip.
      --macrobench-source string                The source or origin of the macro benchmark trigger.
      --macrobench-sysbench-executable string   Path to the sysbench binary.
      --macrobench-type MacroBenchmarkType      Type of macro benchmark.
      --macrobench-working-directory string     Directory on which to execute sysbench.
      --macrobench-workload-path string         Path to the workload used by sysbench.
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/arewefastyet/config.yaml)
```

### SEE ALSO

* [arewefastyet macrobench](arewefastyet_macrobench.md)	 - 

