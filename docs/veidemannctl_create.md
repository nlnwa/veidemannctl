## veidemannctl create

Create or update a config object

### Synopsis

Create or update a config object

```
veidemannctl create [flags]
```

### Options

```
  -f, --filename string   Filename or directory to read from. If input is a directory, all files ending in .yaml or .json will be tried. An input of '-' will read from stdin.
  -h, --help              help for create
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

