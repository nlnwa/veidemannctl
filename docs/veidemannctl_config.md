## veidemannctl config

Modify veidemannctl config files using subcommands

### Synopsis

Modify veidemannctl config files using subcommands

```
veidemannctl config [flags]
```

### Options

```
  -h, --help   help for config
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
* [veidemannctl config create-context](veidemannctl_config_create-context.md)	 - Creates a new context
* [veidemannctl config current-context](veidemannctl_config_current-context.md)	 - Displays the current-context
* [veidemannctl config get-address](veidemannctl_config_get-address.md)	 - Displays Veidemann controller service address
* [veidemannctl config get-apikey](veidemannctl_config_get-apikey.md)	 - Displays Veidemann authentication api-key
* [veidemannctl config get-server-name-override](veidemannctl_config_get-server-name-override.md)	 - Displays the server name override
* [veidemannctl config import-ca](veidemannctl_config_import-ca.md)	 - Import file with trusted certificate chains for the idp and controller.
* [veidemannctl config list-contexts](veidemannctl_config_list-contexts.md)	 - Displays the known contexts
* [veidemannctl config set-address](veidemannctl_config_set-address.md)	 - Sets the address to Veidemann controller service
* [veidemannctl config set-apikey](veidemannctl_config_set-apikey.md)	 - Sets the api-key to use for authentication
* [veidemannctl config set-server-name-override](veidemannctl_config_set-server-name-override.md)	 - Sets the server name override
* [veidemannctl config use-context](veidemannctl_config_use-context.md)	 - Sets the current-context

