## veidemannctl report query

Run a databse query

### Synopsis

Run a databse query. The query should be a java script string like the ones used by RethinkDb javascript driver.

```
veidemannctl report query [queryString|file] args... [flags]
```

### Options

```
  -f, --filename string   File name to write to
  -h, --help              help for query
  -o, --output string     Output format (json|yaml|template|template-file) (default "json")
  -p, --page int32        The page number
  -s, --pagesize int32    Number of objects to get (default 10)
  -t, --template string   A Go template used to format the output
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

