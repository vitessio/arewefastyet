## arewefastyet gen report

Generate comparison between two sha commits of Vitess

### Synopsis

Command to generate a pdf-report that compares the microbenchmark and macrobenchmark runs between the two sha commits of Vitess provided

```
arewefastyet gen report [flags]
```

### Examples

```
arewefastyet gen report --compare-from sha1 --compare-to sha2 --report-file report.pdf
```

### Options

```
      --compare-from string      SHA for Vitess that we want to compare from
      --compare-to string        SHA for Vitess that we want to compare to
      --db-database string       Database to use.
      --db-host string           Hostname of the database
      --db-password string       Password to authenticate the database.
      --db-user string           User used to connect to the database
  -h, --help                     help for report
      --influx-database string   Name of the database to use in InfluxDB.
      --influx-hostname string   Hostname of InfluxDB.
      --influx-password string   Password used to connect to InfluxDB.
      --influx-port string       Port on which to InfluxDB listens. (default "8086")
      --influx-username string   Username used to connect to InfluxDB.
      --report-file string       File created that stores the report. (default "./report.pdf")
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/arewefastyet/config.yaml)
```

### SEE ALSO

* [arewefastyet gen](arewefastyet_gen.md)	 - Generate things

