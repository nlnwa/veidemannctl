## veidemannctl logconfig

Configure logging

### Synopsis

Configure logging.

### Options

```
  -h, --help   help for logconfig
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
* [veidemannctl logconfig delete](veidemannctl_logconfig_delete.md)	 - Delete a logger
* [veidemannctl logconfig list](veidemannctl_logconfig_list.md)	 - List configured loggers
* [veidemannctl logconfig set](veidemannctl_logconfig_set.md)	 - Configure logger

