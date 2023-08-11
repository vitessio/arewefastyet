## arewefastyet completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	arewefastyet completion fish | source

To load completions for every new session, execute once:

	arewefastyet completion fish > ~/.config/fish/completions/arewefastyet.fish

You will need to start a new shell for this setup to take effect.


```
arewefastyet completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --config string    config file (default is $HOME/.config/arewefastyet/config.yaml)
      --secrets string   secrets file
```

### SEE ALSO

* [arewefastyet completion](arewefastyet_completion.md)	 - Generate the autocompletion script for the specified shell

