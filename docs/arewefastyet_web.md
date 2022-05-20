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
  -h, --help                                     help for web
      --planetscale-db-branch string             PlanetscaleDB branch to use. (default "main")
      --planetscale-db-database string           PlanetscaleDB database name.
      --planetscale-db-host string               Hostname of the PlanetscaleDB database.
      --planetscale-db-org string                Name of the PlanetscaleDB organization.
      --planetscale-db-password string           Password used to authenticate to PlanetscaleDB.
      --planetscale-db-user string               Username used to authenticate to PlanetscaleDB.
      --slack-channel string                     Slack channel on which to post messages
      --slack-token string                       Token used to authenticate Slack
      --web-cron-nb-retry int                    Number of retries allowed for each cron job. (default 1)
      --web-cron-schedule string                 Execution CRON schedule defaults to every day at midnight. An empty string will result in no CRON. (default "@midnight")
      --web-cron-schedule-pull-requests string   Execution CRON schedule for pull requests benchmarks. An empty string will result in no CRON. Defaults to an execution every 5 minutes. (default "*/5 * * * *")
      --web-cron-schedule-tags string            Execution CRON schedule for tags/releases benchmarks. An empty string will result in no CRON. Defaults to an execution every minute. (default "*/1 * * * *")
      --web-macrobench-oltp-config string        Path to the configuration file used to execute OLTP macrobenchmark.
      --web-macrobench-tpcc-config string        Path to the configuration file used to execute TPCC macrobenchmark.
      --web-microbench-config string             Path to the configuration file used to execute microbenchmark.
      --web-mode string                          Specify the mode on which the server will run
      --web-port string                          Port used for the HTTP server (default "8080")
      --web-pr-label-trigger string              GitHub Pull Request label that will trigger the execution of new execution. (default "Benchmark me")
      --web-pr-label-trigger-planner-v3 string   GitHub Pull Request label that will trigger the execution of new execution using the V3 planner. (default "Benchmark me (V3)")
      --web-static-path string                   Path to the static directory
      --web-template-path string                 Path to the template directory
      --web-vitess-path string                   Absolute path where the vitess directory is located or where it should be cloned (default "/")
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/arewefastyet/config.yaml)
```

### SEE ALSO

* [arewefastyet](arewefastyet.md)	 - Nightly Benchmarks Project

