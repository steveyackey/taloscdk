# taloscdk

AWS CDK constructs in Go for deploying Talos-based Kubernetes clusters. 

For more examples on utilizing the constructs, check out the /examples directory. 

## Getting Started
If you'd like a quick walkthrough on getting started, visit my [blog](https://www.steveyackey.com/post/taloscdk/).

## Goal
The goal of this construct library is to simplify the deployment of Kubernetes clusters running Talos, and supporting the needed policies to successfully run the aws-controller-manager for creating AWS loadbalancers via Kubernetes manifests.

## Using the Constructs
To use the constructs in your own stacks, run:
```
go get github.com/steveyackey/taloscdk
```

## Requirements
- [Go >= v1.16](https://golang.org/dl/)
- [CDK >= v1.114](https://docs.aws.amazon.com/cdk/latest/guide/getting_started.html#getting_started_install)
- [talosctl](https://www.talos.dev/docs/v0.11/introduction/quickstart/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
