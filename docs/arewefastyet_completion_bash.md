## arewefastyet completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(arewefastyet completion bash)

To load completions for every new session, execute once:

#### Linux:

	arewefastyet completion bash > /etc/bash_completion.d/arewefastyet

#### macOS:

	arewefastyet completion bash > $(brew --prefix)/etc/bash_completion.d/arewefastyet

You will need to start a new shell for this setup to take effect.


```
arewefastyet completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --config string    config file (default is $HOME/.config/arewefastyet/config.yaml)
      --secrets string   secrets file
```

### SEE ALSO

* [arewefastyet completion](arewefastyet_completion.md)	 - Generate the autocompletion script for the specified shell

