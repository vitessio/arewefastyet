## arewefastyet admin

Starts the admin application of arewefastyet

```
arewefastyet admin [flags]
```

### Options

```
      --admin-auth string                      The salt string to salt the GitHub Token
      --admin-gh-app-id string                 The ID of the GitHub App
      --admin-gh-app-secret string             The secret of the GitHub App
      --admin-mode string                      Specify the mode on which the server will run
      --admin-port string                      Port used for the HTTP server (default "8081")
  -h, --help                                   help for admin
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

* [arewefastyet](arewefastyet.md)	 - Nightly Benchmarks Project

