# Command line interface

### Run sample command

Running the sample command requires:
- Completed the installation
- A proper `config.yml` file
- SSH access to the remote Packet servers.

```shell
clibench -c HEAD -s fork_terminal --config-file ./config/config.yaml -runo --ansible-dir ./ansible --tasks-scripts-dir ./scripts --tasks-reports-dir ./report
```

#### Source
To specify from which source the benchmark is called from. This is
added as a tag in mysql database where the benchmark run results are stored.

```
//Example of source name: Api_call,Webhook,test_call
--source <source name> or -s <source name>
```


### Run the tests

This command will run all the tests in the `test` directory.

```shell
python -m unittest discover -s test -v
```

### CLI flags
The `Env` column represents the environment variable name, and `Config name` represents the name found in the `config.yml` files.

| Flag | Env | Config name | Description |
| ---- | ----------- | ------- | ------- |
| `-web`   | _none_  | `web` | Runs the web server  |
| `--run-all`   |  _none_ | `tasks` | Run all the tasks  |
| `--run-tpcc` `-runt`  | _none_ | `tasks`  |  Run TPCC task |
| `--run-oltp` `-runo`  | _none_ | `tasks` | Run OLTP task |
| `--commit` `-c`   | _none_ | `commit` | Commit used to run the task(s)  |
| `--source`   | _none_ | `source` | Where is the task being run |
| `--tasks-scripts-dir`   | `BCLI_TASKS_SCRIPTS_DIR` | `tasks_scripts_dir` |  Directory where the task(s)'s scripts are  |
| `--tasks-reports-dir`   | `BCLI_TASKS_REPORTS_DIR` | `tasks_reports_dir` | Directory where the task(s)'s reports are  |
| `--ansible-dir`   | `BCLI_ANSIBLE_DIR` | `ansible_dir` | Ansible's directory  |
| `--inventory-file` `-invf`   | `BCLI_INVENTORY_FILE` | `inventory_file` | Inventory file used for the task(s)  |
| `--mysql-host`   | `BCLI_MYSQL_HOST` | `mysql_host` | Host of MySQL server |
| `--mysql-username`   | `BCLI_MYSQL_USER` | `mysql_username` | MySQL username  |
| `--mysql-password`   | `BCLI_MYSQL_PASSWORD` | `mysql_password` | MySQL password |
| `--mysql-database`   | `BCLI_MYSQL_DB` | `mysql_database` |  MySQL database |
| `--packet-token`   | `BCLI_PACKET_TOKEN` | `packet_token` |  Packet token |
| `--packet-project-id`   | `BCLI_PACKET_PROJECT_ID` | `packet_project_id` |  Packet project ID |
| `--api-key`   | `BCLI_API_KEY` | `api_key` | Web server API key |
| `--slack-api-token`   | `BCLI_SLACK_TOKEN` | `slack_api_token` | Slack API token |
| `--slack-channel`   | `BCLI_SLACK_CHANNEL` | `slack_channel` | Slack channel |
| `--config-file`   | `BCLI_CONFIG_FILE` | `config_file` | Path to configuration file |
| `--delete-benchmark`   | _none_ | `vps_id` | deletes a specific VPS given the VPS ID which can be found in Config-lock file |
| `--tasks-pprof`   | _none_ | `tasks_pprof` | Specify the profile configuration of vitess |


### Profiling

Profiling of vitess servers can be acheived through the `--tasks-pprof` flag.
The flag behaves similarly as the one in vitessio/vitess.

*Example:* `--tasks-pprof vtgate/cpu`, here we wants to profile vtgate's CPU.

The left-hand side of the flag represents the `targets` they can be `vtgate` or `vttablet`, and must be separated by a `/`.
We can thus specify multiple targets: `--tasks-pprof vttablet/vtgate/cpu`.

The right-hand side of the flag (after the rightmost `/`) is the arguments we will pass to the actual `-pprof` vitess flag. More information regarding the `-pprof` flag can be found [here](https://github.com/vitessio/vitess/blob/master/go/vt/servenv/pprof.go).
We can simply omit the `path=xx` part of the arguments as we set it ourselves within the benchmark configuration.

Once a task is completed, profiles will be copied from the remote server to the `tasks-reports-dir` directory.
