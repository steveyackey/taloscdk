package taloscdk

import (
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type SecurityGroupProps struct {
	// Required.
	Vpc awsec2.IVpc

	// AllowTrafficFrom is the peer to allow ingress to the Kubernetes and Talos APIs.
	// Default: awsec2.Peer_AnyIpv4
	AllowTrafficFrom awsec2.IPeer
}

// NewSecurityGroup returns a security group that enables ingress to 6443, 50000, 50001,
// as well as all internal traffic within the security group.
// Requires a Vpc in the *SecurityGroupProps
func NewSecurityGroup(scope constructs.Construct, id *string, props *SecurityGroupProps) awsec2.SecurityGroup {
	if props == nil {
		props = &SecurityGroupProps{}
	}

	if props.Vpc == nil {
		panic("NewSecurityGroup() requires a Vpc in SecurityGroupProps")
	}

	// Default to allow from all
	if props.AllowTrafficFrom == nil {
		props.AllowTrafficFrom = awsec2.Peer_AnyIpv4()
	}

	sg := awsec2.NewSecurityGroup(scope, jsii.String("TalosSG"), &awsec2.SecurityGroupProps{
		Vpc:              props.Vpc,
		AllowAllOutbound: jsii.Bool(true),
		Description:      jsii.String("Talos Security Group"),
	})

	sg.AddIngressRule(
		props.AllowTrafficFrom,
		awsec2.NewPort(&awsec2.PortProps{
			Protocol:             awsec2.Protocol_TCP,
			FromPort:             jsii.Number(6443),
			ToPort:               jsii.Number(6443),
			StringRepresentation: jsii.String("6443")}),
		jsii.String("Kubernetes API"),
		jsii.Bool(false),
	)

	sg.AddIngressRule(
		props.AllowTrafficFrom,
		awsec2.NewPort(&awsec2.PortProps{
			Protocol:             awsec2.Protocol_TCP,
			FromPort:             jsii.Number(50000),
			ToPort:               jsii.Number(50001),
			StringRepresentation: jsii.String("50000-50001")}),
		jsii.String("Talos API"),
		jsii.Bool(false),
	)

	sg.AddIngressRule(
		sg,
		awsec2.Port_AllTraffic(),
		jsii.String("Allow all internal traffic between nodes"),
		jsii.Bool(false),
	)

	return sg
}
