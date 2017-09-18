# kapacitor-configmap-listener
A shuttle application to move `ConfigMaps` into Kapacitor.  This uses the informer framework for watch notification against the Kubernetes API. 

https://travis-ci.org/xogroup/kapacitor-configmap-listener.svg

Lead Maintainer: [Lam Chan](https://github.com/lamchakchan)

## Introduction
Kapacitor is a data calculator designed around [TICK script](https://docs.influxdata.com/kapacitor/v1.3/tick/).  These scripts can be used to create alerts and events based off of custom thresholds.  Currently, there is no absolute way to set up the TICK environment.  Kubernetes allows for easy boot strapping of containers into a hosting environment.  But there is no easy way to feed Kapacitor TICK scripts to scale containers per deployment.  Also, the subscription binding mechanism between Kapacitor and InfluxDB is less than desirable with the lack of orphaned subscription cleanup.  This controller is a solution to solve these sets of problem

## What does it do?
In short, this controller listens to the Kubernetes State for `ConfigMap` changes.  Any changes will be replicated to Kapacitor.

### Operations

* Delete all Kapacitor subscriptions
  * Any subscriptions that are still used will automatically regenerate
* On start, collects all filtered `ConfigMaps` from Kubernetes
  * Filtered applied is base on `prefix` name matching on the `ConfigMap` name property
* The initial list of `ConfigMaps` are upserted to Kapacitor as task
* All future `ConfigMap` changes are listened on and translated to the appropriate create, update, delete task comands to Kapacitor.

## Installation
This installation guide assumes `go` and [`glide`](https://github.com/Masterminds/glide) is installed

```
go get github.com/xogroup/kapacitor-configmap-listener
cd $GOPATH/src/github.com/xogroup/kapacitor-configmap-listener
glide install
go build
```

Or using Docker
```
docker pull xogroup/kapacitor-configmap-listener
```

## Usage

```
kapacitor-configmap-listener -kapacitorurl http://xyz.com:9092

# or

KAPACITOR_URL=http://xyz.com:9092 kapacitor-configmap-listener
```

Or using Docker
```
docker run -d xogroup/kapacitor-configmap-listener -kapacitorurl http://xyz.com:9092

# or

docker run -d -e KAPACITOR_URL=http://xyz.com:9092 xogroup/kapacitor-configmap-listener
```

## Dependencies

* [client-go@v4.0.0](https://github.com/kubernetes/client-go)
* [kapacitor-client@v1.0.0](https://github.com/influxdata/kapacitor/tree/master/client/v1)
* [logrus@1.0.3](https://github.com/sirupsen/logrus)

## Documentation

### Diagrams

Look at the [Diagrams](DIAGRAM.md).

### API

See the [API Reference](API.md).

### Examples

Check out the [Examples](Example.md).

## Contributing

We love community and contributions! Please check out our [guidelines](.github/CONTRIBUTING.md) before making any PRs.

## Setting up for development

1. Clone this project and `cd` into the project directory
2. Run the commands below

```
glide install
go build
```