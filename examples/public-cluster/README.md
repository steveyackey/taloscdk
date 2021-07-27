# public-cluster

**NOTICE**: Go support for CDK is still in Developer Preview. This implies that APIs may
change while we address early feedback from the community. We would love to hear
about your experience through GitHub issues.

## Useful commands

 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized CloudFormation template
 * `go test`         run unit tests

## Notes

This CDK app will deploy a public cluster with:
  - An autoscaling group for the control plane (starting with 1 node).
  - An NLB for the control plane
  - An autoscaling group for the worker nodes (starting with 1 node)
  The instructions below will also install the aws-controller-manager for creating loadbalancers via AWS. All resources are created in the default VPC and in public subnets. 

## Talos Setup

Generate configs:
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
Replace the endpoint in your talosconfig with the control plane node's IP.

talosctl --talosconfig=talosconfig dmesg -n <ip>
talosctl --talosconfig=talosconfig bootstrap -n <ip>
talosctl --talosconfig=talosconfig dmesg -n <ip> -f # until node is ready and etcd bootstrap is complete
talosctl --talosconfig=talosconfig kubeconfig -n <ip> .  # omit the dot to merge with your current kubeconfig
```

You can now use your new cluster with the newly pulled kubeconfig.
```
kubectl --kubeconfig=./kubeconfig get nodes
```