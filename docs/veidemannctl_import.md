## veidemannctl import

Import data into Veidemann using subcommands

```
veidemannctl import [flags]
```

### Options

```
  -h, --help   help for import
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
* [veidemannctl import convertoos](veidemannctl_import_convertoos.md)	 - Convert Out of Scope file(s) to seed import file
* [veidemannctl import duplicatereport](veidemannctl_import_duplicatereport.md)	 - List duplicated seeds in Veidemann
* [veidemannctl import seed](veidemannctl_import_seed.md)	 - Import seeds

