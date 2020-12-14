## veidemannctl script-parameters

Get the active script parameters for a Crawl Job

### Synopsis

Get the active script parameters for a Crawl Job

Examples:
  # See active script parameters for a Crawl Job
  veidemannctl script-parameters 5604f0cc-315d-4091-8d6e-1b17a7eb990b

  # See active script parameters for a Crawl Job and eventual overrides from Seed and Entity
  veidemannctl script-parameters 5604f0cc-315d-4091-8d6e-1b17a7eb990b 9f89ca44-afe0-4f8f-808f-9df1a0fe64c9


```
veidemannctl script-parameters CRAWLJOB_CONFIG_ID [SEED_ID] [flags]
```

### Options

```
  -h, --help   help for script-parameters
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

