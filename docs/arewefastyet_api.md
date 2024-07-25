## arewefastyet api

Starts the api server of arewefastyet and the CRON service

```
arewefastyet api [flags]
```

### Options

```
      --gh-app-id int                            ID of the GitHub App
      --gh-installation-id int                   GitHub installation ID of this app
      --gh-port string                           Port on which to run the github app (default "8181")
      --gh-secret-key string                     Secret key used to authenticate
      --gh-webhook-secret string                 Secrets used to verify the webhooks
  -h, --help                                     help for api
      --planetscale-db-branch string             PlanetScaleDB branch to use. (default "main")
      --planetscale-db-database string           PlanetScaleDB database name.
      --planetscale-db-host string               Hostname of the PlanetScaleDB database.
      --planetscale-db-org string                Name of the PlanetScaleDB organization.
      --planetscale-db-password-read string      Password used to authenticate to the read-only servers of PlanetScaleDB.
      --planetscale-db-password-write string     Password used to authenticate to the write servers of PlanetScaleDB.
      --planetscale-db-user-read string          Username used to authenticate to the read-only servers of PlanetScaleDB.
      --planetscale-db-user-write string         Username used to authenticate to the write servers of PlanetScaleDB.
      --slack-channel string                     Slack channel on which to post messages
      --slack-token string                       Token used to authenticate Slack
      --web-benchmark-config-path string         Path to the configuration file folder for the benchmarks.
      --web-cron-nb-retry int                    Number of retries allowed for each cron job. (default 1)
      --web-cron-schedule string                 Execution CRON schedule defaults to every day at midnight. An empty string will result in no CRON. (default "@midnight")
      --web-cron-schedule-pull-requests string   Execution CRON schedule for pull requests benchmarks. An empty string will result in no CRON. Defaults to an execution every 5 minutes. (default "*/5 * * * *")
      --web-cron-schedule-tags string            Execution CRON schedule for tags/releases benchmarks. An empty string will result in no CRON. Defaults to an execution every minute. (default "*/1 * * * *")
      --web-mode string                          Specify the mode on which the server will run
      --web-port string                          Port used for the HTTP server (default "8080")
      --web-pr-label-trigger string              GitHub Pull Request label that will trigger the execution of new execution. (default "Benchmark me")
      --web-pr-label-trigger-planner-v3 string   GitHub Pull Request label that will trigger the execution of new execution using the V3 planner. (default "Benchmark me (V3)")
      --web-request-run-key string               Key to authenticate requests for custom benchmark runs.
      --web-source-exclude-filter strings        List of execution source to not execute. By default, all sources are ran.
      --web-source-filter strings                List of execution source that should be run. By default, all sources are ran.
      --web-vitess-path string                   Absolute path where the vitess directory is located or where it should be cloned (default "/")
```

### Options inherited from parent commands

```
      --config string    config file (default is $HOME/.config/arewefastyet/config.yaml)
      --secrets string   secrets file
```

### SEE ALSO

* [arewefastyet](arewefastyet.md)	 - Nightly Benchmarks Project

