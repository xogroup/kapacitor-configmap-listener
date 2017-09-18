# Diagrams
<!-- TOC -->

- [Diagrams](#diagrams)
    - [Architecture of Ecosystem](#architecture-of-ecosystem)

<!-- /TOC -->
## Architecture of Ecosystem

![Image of Ecosystem](/image/architecture.png)

In an ideal world, we would want a single Kapacitor setup to capture and calculate all metrics with.  While it is easier to couple a Kapacitor instance with a deployment, the scale of deployments in Kubernetes makes the management of it unsustainable.  For each Kapacitor instance, a port has to be designated for inbound InfluxDB subscriptions.  Also, there is a parity of InfluxDB subscription per Kapacitor instance which will eventually result in memory management issues.  

The `kapacitor-configmap-listener` (KCL) supports the ideal world scenario for running a single Kapacitor instance for the entire Kubernetes cluster workload.  Assuming you are using `helm` packaging or something similar, a `ConfigMap` can be created along side the deployment with Kapacitor specified `Task` details.  The KCL will capture the `ConfigMap` and translate the to Kapacitor `Task`.  Kapacitor will then run data through these `Task` to create an action with.  A packaged action with KCL is a TICK script used for autoscaling a `ReplicaSet`.

