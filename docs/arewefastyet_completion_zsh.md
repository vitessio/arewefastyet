## arewefastyet completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(arewefastyet completion zsh)

To load completions for every new session, execute once:

#### Linux:

	arewefastyet completion zsh > "${fpath[1]}/_arewefastyet"

#### macOS:

	arewefastyet completion zsh > $(brew --prefix)/share/zsh/site-functions/_arewefastyet

You will need to start a new shell for this setup to take effect.


```
arewefastyet completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --config string    config file (default is $HOME/.config/arewefastyet/config.yaml)
      --secrets string   secrets file
```

### SEE ALSO

* [arewefastyet completion](arewefastyet_completion.md)	 - Generate the autocompletion script for the specified shell

