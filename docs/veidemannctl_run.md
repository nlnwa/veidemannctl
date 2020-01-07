## veidemannctl run

Immediately run a crawl

### Synopsis

Run a crawl. If seedId is submitted only this seed will be run using the configuration
from the submitted jobId. This will run even if the seed is not configured to use the jobId.
If seedId is not submitted then all the seeds wich are configured to use the submitted jobId will be crawled.

```
veidemannctl run jobId [seedId] [flags]
```

### Options

```
  -h, --help   help for run
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

