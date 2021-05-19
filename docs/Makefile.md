# Makefile

### Targets

| Name | Action |
| -----| ------ |
| `all`   | Call `build` |
| `virtual_env`   | Check if virtualenv is present, if not, create it |
| `install`   | Call `virtual_env`, then install all dependencies (ansible & pip) |
| `build`   | Build the CLI binary |
| `gendoc`   | Generate the CLI documentation |
| `test`   | Start Golang unit tests |
| `molecule_converge_all`   | Call `virtual_env`, then create a molecule and converge all the roles |
| `molecule_dryrun_all`   | Call `virtual_env`, then converge and execute all tasks except oltp and tpcc |

### Variables

| Name | Description | Default |
| ---- | ----------- | ------- |
| `PY_VERSION`   | Python version used | 3.7 |
| `VIRTUALENV_PATH`   | Path to virtualenv folder | benchmark |
| `BIN_NAME`   | Name of the CLI binary | arewefastyetcli |
