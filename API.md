# 1.0.0 API Reference

<!-- TOC -->

- [1.0.0 API Reference](#100-api-reference)
    - [Commands](#commands)
        - [Required Flags](#required-flags)
        - [Optional Flags](#optional-flags)
    - [ConfigMap Contract](#configmap-contract)
        - [Contract](#contract)
        - [Replacement Tokens](#replacement-tokens)

<!-- /TOC -->

## Commands

### Required Flags

* `-kapacitorurl` - set the url to the `kapacitord` server.  Needs to include the schema, host and port values.  eg: `http://localhost:9092`

### Optional Flags

* `-incluster` - configure the context of where this controller is to run.  If it is inside a Kubernetes cluster, safe defaults will be used to fetch API data from environment variables.
* `-kubeconfig` - path to the `kubectl` configuration file.
* `-prefixname` - a custom prefix string to use for filtering `ConfigMaps` coming from the Kubernetes API
* `-loglevel` - a value of `0-5` for capturing `panic -> debug` log messages.

## ConfigMap Contract
The configuration is loosely based off of the `kapacitor define <task> -vars vars.json` file contract located [here](https://docs.influxdata.com/kapacitor/v1.3/guides/template_tasks/).

### Contract
```
kind: ConfigMap
apiVersion: v1
metadata:
  name: kapacitor-hpa-rule-{config-map-name}
  namespace: eng
data:
  # default tick template
  template: >-
    { "type" : "string", "value" : "autoscaling" }
  # target is the desired number of request per second per host
  target: >-
    { "type" : "float", "value" : 11 }
  # only one scaling event will be triggered by this time interval
  # https://docs.influxdata.com/kapacitor/v1.3/tick/syntax/#durations
  scalingCooldown: >-
    { "type" : "duration", "value" : "1m0s" }
  # only one descaling even twill be triggered by this time interval
  # https://docs.influxdata.com/kapacitor/v1.3/tick/syntax/#durations
  descalingCooldown: >-
    { "type" : "duration", "value" : "2m0s" }
  # database
  database: >-
    { "type" : "string", "value" : "telegraf" }
  # retention policy for database
  retentionPolicy: >-
    { "type" : "string", "value" : "2-weeks" }
  # dataset collected within the retention policy
  measurement: >-
    { "type" : "string", "value" : "docker_container_cpu" }
  # filter dataset for only preproduction data
  where_filter: >-
    { "type" : "lambda", "value" : "\"environment\" == 'preproduction'" }
  # data to be used for comparison with target
  field: >-
    { "type" : "string", "value" : "usage_percent" }
  # name of deployment to scale for
  deploymentName: >-
    { "type" : "string", "value" : "{deployment-name}" }
  # name of release
  releaseName: >-
    { "type" : "string", "value" : "{release-name}" }
```

### Replacement Tokens
* `{config-map-name}` - name appending the `kapacitor-hpa-rule-` prefix
* `{release-name}` - release name issued from `helm` if given, or something arbitrary associated to the application.
* `{deployment-name}` - must be the same name registered with the `ReplicateSet`.  This is the key used to target the `replicaCount` change within the TICK script.