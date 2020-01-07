## veidemannctl config use-context

Sets the current-context

### Synopsis

Sets the current-context

Examples:
  # Use the context for the prod cluster
  veidemannctl config use-context prod


```
veidemannctl config use-context CONTEXT_NAME [flags]
```

### Options

```
  -h, --help   help for use-context
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

* [veidemannctl config](veidemannctl_config.md)	 - Modify veidemannctl config files using subcommands

