# 1.0.0 API Reference

<!-- TOC -->

- [1.0.0 API Reference](#100-api-reference)
    - [Commands](#commands)
        - [Required Flags](#required-flags)
        - [Optional Flags](#optional-flags)

<!-- /TOC -->

## Commands

### Required Flags

* `-kapacitorurl` - set the url to the `kapacitord` server.  Needs to include the schema, host and port values.  eg: `http://localhost:9092`

### Optional Flags

* `-incluster` - configure the context of where this controller is to run.  If it is inside a Kubernetes cluster, safe defaults will be used to fetch API data from environment variables.
* `-kubeconfig` - path to the `kubectl` configuration file.
* `-prefixname` - a custom prefix string to use for filtering `ConfigMaps` coming from the Kubernetes API
* `-loglevel` - a value of `0-5` for capturing `panic -> debug` log messages.