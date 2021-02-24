# Makefile

### Targets

| Name | Action |
| -----| ------ |
| `virtual_env`   | Check if virtualenv is present, if not, create it |
| `install`   | Call `virtual_env`, then install all depedencies (ansible & pip) |
| `oltp`   | Call `virtual_env`, then run an OLTP benchmark |
| `tpcc`   | Call `virtual_env`, then run a TPCC benchmark |
| `molecule_converge_all`   | Call `virtual_env`, then create a molecule and converge all the roles |

### Variables

| Name | Description | Default |
| ---- | ----------- | ------- |
| `PY_VERSION`   | Python version used | 3.7 |
| `VIRTUALENV_PATH`   | Path to virtualenv folder | benchmark |
| `ANSIBLE_PATH`   | Path to ansible folder | ansible |
| `RUN_COMMIT`   | Vitess commit used to run the benchmark | HEAD |
| `RUN_INVENTORY_FILE`   | Path to the ansible inventory file for benchmark| koz-inventory-unsharded-test.yml |
