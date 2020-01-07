## veidemannctl delete

Delete a config object

### Synopsis

Delete a config object.

Valid object types include:
  * browserConfig
  * browserScript
  * collection
  * crawlConfig
  * crawlEntity
  * crawlHostGroupConfig
  * crawlJob
  * crawlScheduleConfig
  * politenessConfig
  * roleMapping
  * seed

Examples:
  #Delete a seed.
  veidemannctl delete seed 407a9600-4f25-4f17-8cff-ee1b8ee950f6

```
veidemannctl delete <kind> [id] ... [flags]
```

### Options

```
      --dry-run         Set to false to execute delete (default true)
  -q, --filter string   Delete objects by field (i.e. meta.description=foo)
  -h, --help            help for delete
  -l, --label string    Delete objects by label (<type>:<value> | <value>)
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

