## veidemannctl completion

Output bash completion code

### Synopsis

Output bash completion code. The shell code must be evalutated to provide
interactive completion of veidemannctl commands.  This can be done by sourcing it from the .bash _profile.

Example:
  ## Load the kubectl completion code for bash into the current shell
  source <(veidemannctl completion)


```
veidemannctl completion [flags]
```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --apiKey string               Api-key used for authentication instead of interactive logon trough IDP.
      --config string               config file (default is $HOME/.veidemannctl.yaml)
      --context string              The name of the veidemannconfig context to use.
  -c, --controllerAddress string    Address to the Controller service (default "localhost:50051")
  -d, --debug                       Turn on debugging
      --serverNameOverride string   If set, it will override the virtual host name of authority (e.g. :authority header field) in requests.
```

### SEE ALSO

* [veidemannctl](veidemannctl.md)	 - Veidemann command line client

