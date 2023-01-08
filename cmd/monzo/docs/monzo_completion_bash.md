## monzo completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(monzo completion bash)

To load completions for every new session, execute once:

#### Linux:

	monzo completion bash > /etc/bash_completion.d/monzo

#### macOS:

	monzo completion bash > $(brew --prefix)/etc/bash_completion.d/monzo

You will need to start a new shell for this setup to take effect.


```
monzo completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [monzo completion](monzo_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 8-Jan-2023