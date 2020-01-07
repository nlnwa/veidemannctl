## veidemannctl report

Get log report

### Synopsis

Request a report.

```
veidemannctl report [flags]
```

### Options

```
  -h, --help   help for report
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
* [veidemannctl report crawlexecution](veidemannctl_report_crawlexecution.md)	 - A brief description of your command
* [veidemannctl report crawllog](veidemannctl_report_crawllog.md)	 - A brief description of your command
* [veidemannctl report jobexecution](veidemannctl_report_jobexecution.md)	 - A brief description of your command
* [veidemannctl report pagelog](veidemannctl_report_pagelog.md)	 - A brief description of your command
* [veidemannctl report query](veidemannctl_report_query.md)	 - Run a databse query

