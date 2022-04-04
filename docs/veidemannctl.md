## veidemannctl

Veidemann command line client

### Synopsis

A command line client for Veidemann which can manipulate configs and request status of the crawler.

```
veidemannctl [flags]
```

### Options

```
      --apiKey string               Api-key used for authentication instead of interactive logon trough IDP.
      --config string               config file (default is $HOME/.veidemannctl.yaml)
      --context string              The name of the veidemannconfig context to use.
  -c, --controllerAddress string    Address to the Controller service (default "localhost:50051")
  -d, --debug                       Turn on debugging
  -h, --help                        help for veidemannctl
      --serverNameOverride string   If set, it will override the virtual host name of authority (e.g. :authority header field) in requests.
```

### SEE ALSO

* [veidemannctl abort](veidemannctl_abort.md)	 - Abort one or more crawl executions
* [veidemannctl abortjobexecution](veidemannctl_abortjobexecution.md)	 - Abort one or more job executions
* [veidemannctl activeroles](veidemannctl_activeroles.md)	 - Get the active roles for the currently logged in user
* [veidemannctl completion](veidemannctl_completion.md)	 - Output bash completion code
* [veidemannctl config](veidemannctl_config.md)	 - Modify veidemannctl config files using subcommands
* [veidemannctl create](veidemannctl_create.md)	 - Create or update a config object
* [veidemannctl delete](veidemannctl_delete.md)	 - Delete a config object
* [veidemannctl get](veidemannctl_get.md)	 - Get the value(s) for an object type
* [veidemannctl import](veidemannctl_import.md)	 - Import data into Veidemann using subcommands
* [veidemannctl logconfig](veidemannctl_logconfig.md)	 - Configure logging
* [veidemannctl login](veidemannctl_login.md)	 - Initiate browser session for logging in to Veidemann
* [veidemannctl logout](veidemannctl_logout.md)	 - Log out of Veidemann
* [veidemannctl pause](veidemannctl_pause.md)	 - Pause crawler
* [veidemannctl report](veidemannctl_report.md)	 - Get report
* [veidemannctl run](veidemannctl_run.md)	 - Immediately run a crawl
* [veidemannctl script-parameters](veidemannctl_script-parameters.md)	 - Get the active script parameters for a Crawl Job
* [veidemannctl status](veidemannctl_status.md)	 - Get crawler status
* [veidemannctl unpause](veidemannctl_unpause.md)	 - Unpause crawler
* [veidemannctl update](veidemannctl_update.md)	 - Update the value(s) for an object type

