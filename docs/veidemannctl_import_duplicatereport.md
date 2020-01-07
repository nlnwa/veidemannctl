## veidemannctl import duplicatereport

List duplicated seeds in Veidemann

### Synopsis

List duplicated seeds in Veidemann

```
veidemannctl import duplicatereport [kind] [flags]
```

### Options

```
  -b, --db-directory string   Directory for storing state db (default "/tmp/veidemannctl")
  -h, --help                  help for duplicatereport
  -o, --outFile string        File to write output.
  -r, --reset-db              Clean state db
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

* [veidemannctl import](veidemannctl_import.md)	 - Import data into Veidemann using subcommands

