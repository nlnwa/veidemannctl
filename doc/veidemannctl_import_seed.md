## veidemannctl import seed

Import seeds

### Synopsis

Import seeds

```
veidemannctl import seed [flags]
```

### Options

```
      --checkuri               Check the uri for liveness and follow 301
      --checkuri-timeout int   Timeout in ms when checking uri for liveness. (default 500)
      --crawljob-id string     Set crawlJob ID for new seeds.
  -b, --db-directory string    Directory for storing state db (default "/tmp/veidemannctl")
      --dry-run                Run the import without writing anything to Veidemann
  -e, --errorfile string       File to write errors to.
  -f, --filename string        Filename or directory to read from. If input is a directory, all files ending in .yaml or .json will be tried. An input of '-' will read from stdin.
  -h, --help                   help for seed
  -r, --reset-db               Clean state db
      --toplevel               Convert URI to toplevel by removing path.
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

