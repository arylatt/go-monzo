## monzo completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(monzo completion zsh); compdef _monzo monzo

To load completions for every new session, execute once:

#### Linux:

	monzo completion zsh > "${fpath[1]}/_monzo"

#### macOS:

	monzo completion zsh > $(brew --prefix)/share/zsh/site-functions/_monzo

You will need to start a new shell for this setup to take effect.


```
monzo completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [monzo completion](monzo_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 8-Jan-2023
