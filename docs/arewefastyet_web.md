## arewefastyet web

Starts the HTTP web server

```
arewefastyet web [flags]
```

### Options

```
      --db-database string          Database to use.
      --db-host string              Hostname of the database
      --db-password string          Password to authenticate the database.
      --db-user string              User used to connect to the database
  -h, --help                        help for web
      --influx-database string      Name of the database to use in InfluxDB.
      --influx-hostname string      Hostname of InfluxDB.
      --influx-password string      Password used to connect to InfluxDB.
      --influx-port string          Port on which to InfluxDB listens. (default "8086")
      --influx-username string      Username used to connect to InfluxDB.
      --web-api-key string          API key used to authenticate requests
      --web-mode string             Specify the mode on which the server will run
      --web-port string             Port used for the HTTP server (default "8080")
      --web-static-path string      Path to the static directory
      --web-template-path string    Path to the template directory
      --web-webhook-config string   Path to default config file used for Webhook
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/arewefastyet/config.yaml)
```

### SEE ALSO

* [arewefastyet](arewefastyet.md)	 - 

