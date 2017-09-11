# kapacitor-configmap-listener

<!-- TOC -->

- [kapacitor-configmap-listener](#kapacitor-configmap-listener)
    - [Description](#description)
    - [How to install](#how-to-install)
    - [How to build](#how-to-build)
    - [How to upgrade dependencies](#how-to-upgrade-dependencies)

<!-- /TOC -->

## Description
A shuttle application to move configmaps into Kapacitor.

## How to install
The installation process requires the [`glide`](https://github.com/Masterminds/glide) package management tool.

Run `glide install` to bring down all dependencies followed by `glide up -v` to remove nested vendor dependencies.

## How to build

Run `go build -o kcl`.

## How to upgrade dependencies

Run `glide up` - updates dependencies based on [`semver`](http://semver.org/)