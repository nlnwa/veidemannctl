## veidemannctl config set-server-name-override

Sets the server name override

### Synopsis

Sets the server name override.

Use this when there is a mismatch between exposed server name for the cluster and the certificate. The use is a security
risk and is only recommended for testing.

Examples:
  # Sets the server name override to test.local
  veidemannctl config set-server-name-override test.local


```
veidemannctl config set-server-name-override HOST [flags]
```

### Options

```
  -h, --help   help for set-server-name-override
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

