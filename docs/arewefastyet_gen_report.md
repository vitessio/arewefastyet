## arewefastyet gen report

Generate comparison between two sha commits of Vitess

```
arewefastyet gen report [flags]
```

### Options

```
      --compare-from string   SHA for Vitess that we want to compare from
      --compare-to string     SHA for Vitess that we want to compare to
      --db-database string    Database to use.
      --db-host string        Hostname of the database
      --db-password string    Password to authenticate the database.
      --db-user string        User used to connect to the database
  -h, --help                  help for report
      --report-file string    File created that stores the report. (default "./report.pdf")
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/arewefastyet/config.yaml)
```

### SEE ALSO

* [arewefastyet gen](arewefastyet_gen.md)	 - Generate things

