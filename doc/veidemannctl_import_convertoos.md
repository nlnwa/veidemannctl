## veidemannctl import convertoos

Convert Out of Scope file(s) to seed import file

### Synopsis

Convert Out of Scope file(s) to seed import file

```
veidemannctl import convertoos [flags]
```

### Options

```
      --checkuri               Check the uri for liveness and follow 301 (default true)
      --checkuri-timeout int   Timeout in ms when checking uri for liveness. (default 2000)
  -b, --db-directory string    Directory for storing state db (default "/tmp/veidemannctl")
  -e, --errorfile string       File to write errors to.
  -f, --filename string        Filename or directory to read from. If input is a directory, all files ending in .yaml or .json will be tried. An input of '-' will read from stdin. (required)
  -h, --help                   help for convertoos
  -o, --outfile string         File to write result to. (required)
  -r, --reset-db               Clean state db
      --toplevel               Convert URI to toplevel by removing path. (default true)
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

