# kapacitor-configmap-listener

<!-- TOC -->

- [kapacitor-configmap-listener](#kapacitor-configmap-listener)
    - [Description](#description)
    - [How to use](#how-to-use)
    - [Dependencies](#dependencies)
    - [Development and Maintenance](#development-and-maintenance)
        - [How to install](#how-to-install)
        - [How to build](#how-to-build)
        - [How to upgrade dependencies](#how-to-upgrade-dependencies)

<!-- /TOC -->

## Description
A shuttle application to move `ConfigMaps` into Kapacitor.  This uses the informer framework for watch notification against the API. 

## How to use

## Dependencies

* [client-go@v4.0.0](https://github.com/kubernetes/client-go)
* [kapacitor-client@v1.0.0](https://github.com/influxdata/kapacitor/tree/master/client/v1)

## Development and Maintenance

### How to install
The installation process requires the [`glide`](https://github.com/Masterminds/glide) package management tool.

Run `glide install` to bring down all dependencies followed by `glide up -v` to remove nested vendor dependencies.

### How to build

Run `go build`.

### How to upgrade dependencies

Run `glide up` - updates dependencies based on [`semver`](http://semver.org/)