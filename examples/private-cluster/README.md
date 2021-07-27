# private-cluster

**NOTICE**: Go support for CDK is still in Developer Preview. This implies that APIs may
change while we address early feedback from the community. We would love to hear
about your experience through GitHub issues.

## Useful commands

 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized CloudFormation template
 * `go test`         run unit tests

## Notes

This CDK app will deploy a private cluster with:
  - An autoscaling group for the control plane (starting with 3 nodes in a private subnet).
  - An NLB for the control plane (only accessible via private subnets and the bastion host)
  - An autoscaling group for the worker nodes (starting with 3 nodes in a private subnet)
  - A bastion host for using SSM Session Manager to remotely access your private cluster
  - A new VPC with private and public subnets

  The instructions below will also install the aws-controller-manager for creating loadbalancers via AWS. All resources are created in the new VPC and in private subnets. This cluster is able to create public loadblancers to access services within the cluster as well as private loadbalancers. 

## Talos Setup

### Generate configs
```bash
talosctl gen config talos https://talos.cluster:6443 \
    --with-examples=false --with-docs=false \
    --config-patch='[{"op":"replace", "path":"/machine/kubelet", "value": {"registerWithFQDN": true}},
        {"op":"replace", "path":"/cluster/externalCloudProvider", "value": {
            "enabled": true,
            "manifests": [
                "https://raw.githubusercontent.com/kubernetes/cloud-provider-aws/v1.20.0-alpha.0/manifests/rbac.yaml", 
                "https://raw.githubusercontent.com/kubernetes/cloud-provider-aws/v1.20.0-alpha.0/manifests/aws-cloud-controller-manager-daemonset.yaml"
            ]
        }}]'
```

Deploy CDK:
```
cdk deploy
```

### Update Talosconfig
Replace the endpoint in your talosconfig with one of your control plane node's IP addresses.
Connect to the EC2 via SSM Session Manager, and create paste your talosconfig to ~/.talos/config

Install [kubectl](https://kubernetes.io/docs/tasks/tools/) and [talosctl](https://www.talos.dev/docs/v0.11/introduction/quickstart/) on the bastion host and then continue with the steps below: 

```
talosctl dmesg -n <ip>
talosctl bootstrap -n <ip>
talosctl dmesg -n <ip> -f # until node is ready and etcd bootstrap is complete
talosctl kubeconfig -n <ip>
```

You can now use your new cluster with the newly pulled kubeconfig.
```
kubectl --kubeconfig=./kubeconfig get nodes
```