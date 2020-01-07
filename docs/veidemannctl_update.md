## veidemannctl update

Update the value(s) for an object type

### Synopsis

Update one or many objects with new values.

```
veidemannctl update [object_type] [flags]
```

### Examples

```
# Add CrawlJob for a seed.
veidemannctl update seed -n "http://www.gwpda.org/" -u seed.jobRef+=crawlJob:e46863ae-d076-46ca-8be3-8a8ef72e709

# Replace all configured CrawlJobs for a seed with a new one.
veidemannctl update seed -n "http://www.gwpda.org/" -u seed.jobRef=crawlJob:e46863ae-d076-46ca-8be3-8a8ef72e709
```

### Options

```
  -q, --filter string        Filter objects by field (i.e. meta.description=foo)
  -h, --help                 help for update
  -l, --label string         List objects by label (<type>:<value> | <value>)
  -s, --limit int32          Limit the number of objects to update. 0 = no limit
  -n, --name string          List objects by name (accepts regular expressions)
  -u, --updateField string   Filter objects by field (i.e. meta.description=foo)
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

