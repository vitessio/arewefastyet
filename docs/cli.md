# Command line interface

## Describes the cli commands to use:
The implementation can be found on cli.py:
1. When running any OLTP or TPCC the source and commit hash must be provided

### Sample run
```
./cli --run-all -c HEAD -s api_call
```

#### Run both OLTP and TPCC benchmarks
```
--run-all
```

#### Run TPCC
```
--run-tpcc or -runt
```

#### Run OLTP
```
--run-oltp or -runo
```

#### Provide Inventory file to run
If not provided the inventory file provided in config.yaml will be called by1 default.
```
--inventory-file=<inventory file name> or -invf=<inventory file name>
```

#### To Specify commit hash or branch name
```
--commit <commit hash or branch name> or -c <commit hash or branch name>
```

#### Source
To specify from which source the benchmark is called from. This is
added as a tag in mysql database where the benchmark run results are stored.

```
//Example of source name: Api_call,Webhook,test_call
--source <source name> or -s <source name>
```
