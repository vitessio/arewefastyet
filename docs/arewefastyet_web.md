## arewefastyet web

Starts the HTTP web server

### Synopsis

web command starts the HTTP web server, the credentials for which are provided as command line arguments or the configuration file. 
It uses a MySQL and InfluxDB instances to read the metrics that will be displayed. It has an interface for macrobenchmark and microbenchmark runs.

```
arewefastyet web [flags]
```

### Examples

```
arewefastyet web --db-database benchmark --db-host localhost --db-password <db-password> --db-user <db-username>  
--influx-database benchmark-influx --influx-hostname localhost --influx-password <influx-password>
--influx-port <influx-port> --influx-username <influx-username> --web-api-key <web-api-key>
--web-mode production --web-port <web-port> --web-static-path ./server/static --web-template-path ./server/template
--web-webhook-config ./config.yaml
```

### Options

```
      --db-database string                  Database to use.
      --db-host string                      Hostname of the database
      --db-password string                  Password to authenticate the database.
      --db-user string                      User used to connect to the database
  -h, --help                                help for web
      --influx-database string              Name of the database to use in InfluxDB.
      --influx-hostname string              Hostname of InfluxDB.
      --influx-password string              Password used to connect to InfluxDB.
      --influx-port string                  Port on which to InfluxDB listens. (default "8086")
      --influx-username string              Username used to connect to InfluxDB.
      --web-cron-nb-retry int               Number of retries allowed for each cron job.
      --web-cron-schedule string            Execution CRON schedule defaults to every day at midnight. An empty string will result in no CRON. (default "@midnight")
      --web-macrobench-oltp-config string   Path to the configuration file used to execute OLTP macrobenchmark.
      --web-macrobench-tpcc-config string   Path to the configuration file used to execute TPCC macrobenchmark.
      --web-microbench-config string        Path to the configuration file used to execute microbenchmark.
      --web-mode string                     Specify the mode on which the server will run
      --web-port string                     Port used for the HTTP server (default "8080")
      --web-pr-label-trigger string         GitHub Pull Request label that will trigger the execution of new execution. (default "Benchmark me")
      --web-static-path string              Path to the static directory
      --web-template-path string            Path to the template directory
      --web-vitess-path string              Absolute path where the vitess directory is located or where it should be cloned (default "/")
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/arewefastyet/config.yaml)
```

### SEE ALSO

* [arewefastyet](arewefastyet.md)	 - Nightly Benchmarks Project

