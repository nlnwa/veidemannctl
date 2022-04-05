## veidemannctl report jobexecution

Get current status for job executions

### Synopsis

Get current status for job executions.

```
veidemannctl report jobexecution [flags]
```

### Options

```
      --desc              Order descending
  -f, --filename string   File name to write to
  -q, --filter strings    Filter objects by field (i.e. meta.description=foo
      --from string       From start time
  -h, --help              help for jobexecution
      --order-by string   Order by path
  -o, --output string     Output format (table|wide|json|yaml|template|template-file) (default "table")
  -p, --page int32        The page number
  -s, --pagesize int32    Number of objects to get (default 10)
      --state strings     Filter objects by state(s)
  -t, --template string   A Go template used to format the output
      --to string         To start time
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

* [veidemannctl report](veidemannctl_report.md)	 - Get report

