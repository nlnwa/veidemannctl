## veidemannctl config import-ca

Import file with trusted certificate chains for the idp and controller.

### Synopsis

Import file with trusted certificate chains for the idp and controller. These are in addition to the default certs configured for the OS.

```
veidemannctl config import-ca CA_CERT_FILE_NAME [flags]
```

### Options

```
  -h, --help   help for import-ca
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

