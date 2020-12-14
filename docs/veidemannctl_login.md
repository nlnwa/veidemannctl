## veidemannctl login

Initiate browser session for logging in to Veidemann

```
veidemannctl login [flags]
```

### Options

```
  -h, --help     help for login
  -m, --manual   Manually copy and paste login url and code. Use this to log in from a remote terminal.
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

