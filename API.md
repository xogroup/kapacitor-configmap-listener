# 1.0.0 API Reference

<!-- TOC -->

- [1.0.0 API Reference](#100-api-reference)
    - [Commands](#commands)
        - [Required Flags](#required-flags)
        - [Optional Flags](#optional-flags)
    - [ConfigMap](#configmap)
        - [Contract](#contract)
        - [Replacement Tokens](#replacement-tokens)
    - [TICK](#tick)
        - [`autoscaling`](#autoscaling)

<!-- /TOC -->

## Commands
Flags need to be formatted as `-<key1>=<value1>,-<key2>=<value2>,...-<keyn>=<valuen>`.

### Required Flags

* `-kapacitorurl` - set the url to the `kapacitord` server.  Needs to include the schema, host and port values.  eg: `http://localhost:9092`

### Optional Flags

* `-incluster` - configure the context of where this controller is to run.  If it is inside a Kubernetes cluster, safe defaults will be used to fetch API data from environment variables.
* `-cleansubscriptions` - enable influx subscription cleanup
* `-influxurl` - set the url to the `influxd` server.  Needs to include the schema, host and port values.  eg: `http://localhost:8086`
* `-influxusername` - set the username for influx connection
* `-influxpassword` - set the password for the influx connection
* `-influxssl` - set the http schema through this property
* `-influxunsafessl` - skip ssl verification
* `-kubeconfig` - path to the `kubectl` configuration file.
* `-prefixname` - a custom prefix string to use for filtering `ConfigMaps` coming from the Kubernetes API
* `-loglevel` - a value of `0-5` for capturing `panic -> debug` log messages.

## ConfigMap
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
  # namespace deployment resides
  namespace: >-
    { "type" : "string", "value" : "{deployment-namespace}" }
  # name of deployment to scale for
  deploymentName: >-
    { "type" : "string", "value" : "{deployment-name}" }
  # name of release
  releaseName: >-
    { "type" : "string", "value" : "{release-name}" }
  # minimum amount of replicas running >-
  minReplicaCount: >-
    { "type" : "int", "value" : 1 }    
  # maximum amount of replicas running >-
  maxReplicaCount: >-
    { "type" : "int", "value" : 20 }    
```

### Replacement Tokens
* `{config-map-name}` - name appending the `kapacitor-hpa-rule-` prefix.
* `{release-name}` - release name issued from `helm` if given, or something arbitrary associated to the application.
* `{deployment-namespace}` - namespace the deployment resides in
* `{deployment-name}` - must be the same name registered with the `ReplicateSet`.  This is the key used to target the `replicaCount` change within the TICK script.

## TICK

### `autoscaling`
This script is combined with the `ConfigMap` to scale up/down a replicate set based on collected data available from InfluxDB.  We can assume all measurements are available based on any `telegraf` forwarded telemetry along with your own custom data points pushed.

```
// Database
var database string
// Retention policy for database
var retentionPolicy string
// Dataset collected within the retention policy
var measurement string
// Optional where filter
var whereFilter = lambda: TRUE
// Field data to use for the processing
var field string
// The time scale to calculate the average against
var movingAverageCount = 60
// Namespace the deployment lives in
var namespace = 'default'
// Deployment this is scaling for
var deploymentName = 'placeholder'
// Threshold for triggering
var target = 10.0
// Time interval per scaling up
var scalingCooldown = 2m
// Time interval per scaling down
var descalingCooldown = 5m
// Minimum replica count to maintain regardless of needs
var minReplicaCount = 1
// Maximum replica count to stop scaling at regardless of needs
var maxReplicaCount = 20
	
stream
	|from()
		.database(database)
		.retentionPolicy(retentionPolicy)
		.measurement(measurement)
		.where(whereFilter)
		.truncate(1s)
	|movingAverage(field, movingAverageCount)
		.as('averageResource')	
	|eval(lambda: int(ceil("averageResource" / float(target))))
		.as('replicaCount')
	|k8sAutoscale()
		// We are scaling a deployment.
		.kind('deployments')
		// The namespace of the deployment
		.namespace(namespace)
		// The name of the replicaset to scale is found in the 'replicaset' tag.
		.resourceName(deploymentName)
		// Set the cool down timer values.
		.increaseCooldown(scalingCooldown)
		.decreaseCooldown(descalingCooldown)
		// The minimum amount of replica to have regardless of averages
		.min(minReplicaCount)
		// The maximum amount of replica to have regardless of averages
		.max(maxReplicaCount)
		// Compute the desired number of replicas based on the
		// avg_point_per_second and target values.
		.replicas(lambda: "replicaCount")
```