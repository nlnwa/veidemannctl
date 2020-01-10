## veidemannctl report pagelog

View page log

### Synopsis

View page log.

```
veidemannctl report pagelog [flags]
```

### Options

```
  -f, --filename string   File name to write to
  -q, --filter string     Filter objects by field (i.e. meta.description=foo
  -h, --help              help for pagelog
  -o, --output string     Output format (table|wide|json|yaml|template|template-file) (default "table")
  -p, --page int32        The page number
  -s, --pagesize int32    Number of objects to get (default 10)
  -t, --template string   A Go template used to format the output
  -w, --watch             Get a continous stream of changes
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

* [veidemannctl report](veidemannctl_report.md)	 - Get log report

