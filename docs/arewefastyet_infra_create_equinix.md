## arewefastyet infra create equinix

Create an Equinix Metal instance

### Synopsis

Command used to create a new equinix metal instance based on terraform configuration

```
arewefastyet infra create equinix [flags]
```

### Examples

```
arewefastyet infra create equinix --infra-path ./infra --equinix-instance-type m2.xlarge.x86 --equinix-token tok --equinix-project-id id
```

### Options

```
      --equinix-instance-type string   Instance type to use for the creation of a new node
      --equinix-project-id string      Project ID to use for Equinix Metal
      --equinix-token string           Auth Token for Equinix Metal
  -h, --help                           help for equinix
```

### Options inherited from parent commands

```
      --config string       config file (default is $HOME/.config/arewefastyet/config.yaml)
      --infra-path string   Path to the infra directory
```

### SEE ALSO

* [arewefastyet infra create](arewefastyet_infra_create.md)	 - Create a new instance

