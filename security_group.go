// package construct contains custom CDK constructs for creating Talos clusters
package taloscdk

import (
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type SecurityGroupProps struct {
	Vpc              awsec2.IVpc
	AllowTrafficFrom awsec2.IPeer
}

func NewSecurityGroup(scope constructs.Construct, id *string, props *SecurityGroupProps) awsec2.SecurityGroup {

	if props.Vpc == nil {
		props.Vpc = awsec2.Vpc_FromLookup(scope, jsii.String("vpc"), &awsec2.VpcLookupOptions{
			IsDefault: jsii.Bool(true),
		})
	}

	// Default to allow from all
	if props.AllowTrafficFrom == nil {
		props.AllowTrafficFrom = awsec2.Peer_AnyIpv4()
	}

	cpSecurityGroup := awsec2.NewSecurityGroup(scope, id, &awsec2.SecurityGroupProps{
		Vpc:               props.Vpc,
		AllowAllOutbound:  jsii.Bool(true),
		Description:       jsii.String("Talos Control Plane Security Group"),
		SecurityGroupName: id,
	})

	cpSecurityGroup.AddIngressRule(
		props.AllowTrafficFrom,
		awsec2.NewPort(&awsec2.PortProps{
			Protocol:             awsec2.Protocol_TCP,
			FromPort:             jsii.Number(6443),
			ToPort:               jsii.Number(6443),
			StringRepresentation: jsii.String("6443")}),
		jsii.String("Kubernetes API"),
		jsii.Bool(false),
	)

	cpSecurityGroup.AddIngressRule(
		props.AllowTrafficFrom,
		awsec2.NewPort(&awsec2.PortProps{
			Protocol:             awsec2.Protocol_TCP,
			FromPort:             jsii.Number(50000),
			ToPort:               jsii.Number(50001),
			StringRepresentation: jsii.String("50000-50001")}),
		jsii.String("Talos API"),
		jsii.Bool(false),
	)

	cpSecurityGroup.AddIngressRule(
		cpSecurityGroup,
		awsec2.Port_AllTraffic(),
		jsii.String("Allow all internal traffic"),
		jsii.Bool(false),
	)
}
