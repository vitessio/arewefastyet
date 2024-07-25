## arewefastyet gen exec_metrics

For each execution, fetches the metrics from influxDB and store them to SQL if not already present.

```
arewefastyet gen exec_metrics [flags]
```

### Options

```
  -h, --help                                   help for exec_metrics
      --influx-database string                 Name of the database to use in InfluxDB.
      --influx-hostname string                 Hostname of InfluxDB.
      --influx-password string                 Password used to connect to InfluxDB.
      --influx-port string                     Port on which to InfluxDB listens. (default "8086")
      --influx-username string                 Username used to connect to InfluxDB.
      --planetscale-db-branch string           PlanetScaleDB branch to use. (default "main")
      --planetscale-db-database string         PlanetScaleDB database name.
      --planetscale-db-host string             Hostname of the PlanetScaleDB database.
      --planetscale-db-org string              Name of the PlanetScaleDB organization.
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

* [arewefastyet gen](arewefastyet_gen.md)	 - Generate things

