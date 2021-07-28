package taloscdk

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type SingleNodeProps struct {
	// ClusterName is used for tagging all resources with kubernetes.io/cluster/<name>=owned
	// Default: talos
	ClusterName *string

	// NodeName is used for naming your EC2 instances
	// Default: jsii.String("talos")
	NodeName *string

	// MachineImageName is used for searching AMI by name and supports * wildcard.
	// Be sure to select an arch that matches your instance type.
	// It's typically easiest to use a wildcard for the region so that it works cross-region.
	// Format: talos-<Version>-<AWSRegion>-<arch>
	// Default: talos-v0.11.2-*-amd64
	MachineImageName *string

	// MachineImageAMI is used to get the image from an AMI.
	// Talos AMIs can be found in the docs: https://www.talos.dev/docs/v0.11/cloud-platforms/aws/ (sub v0.11 for current version)
	// Example: {"us-east-1": jsii.String("ami-0fdb2f5cb915076a3")}  (us-east-1 amd64 v0.11 image)
	// Defaults to using MachineImageName
	MachineImageAMI *map[string]*string

	// TalosNodeConfig is a *string of the controlplane.yaml or join.yaml you've generated with
	// `talosctl gen config <clusterName> <endpoint>`
	// To load a node config use taloscdk.LoadConfig("<yourConfig>")
	//
	// Example:
	// config, err := taloscdk.LoadConfig("cluster-config/controlplane.yaml")
	// if err != nil {
	// 	panic("Could not load talos config")
	// }
	// TalosNodeConfig is required
	TalosNodeConfig *string

	// TransformConfig sets whether or not to change the endpoint in our TalosNodeConfig to
	// the OverwriteValue
	// Default: jsii.Bool(true)
	TransformConfig *bool

	// EndpointToOverwrite  is the <endpoint> you used when running `talosctl gen config <clusterName> https://<endpoint>:6443`
	// This will overwrite the <endpoint> in your config, while keeping https:// and the port (:6443).
	// For example: in https://talos.cluster:6443, if you overwrite "talos.cluster", it would become https://YourOverwriteValue:6443
	// By default, the OverwriteValue does not include protocl or port.
	EndpointToOverwrite *string

	// OverwriteValue to replace EndpointToOverwrite
	// Default: EIP. Can use GetEIPAddress() to get from another node.
	OverwriteValue *string

	// InstanceType is used to determine the size/arch of the instance.
	// Default: t3.small (amd64). Meets min specs: https://www.talos.dev/docs/v0.11/introduction/system-requirements/
	InstanceType awsec2.InstanceType

	// SecurityGroup for the instance.
	// To create a security group to use with multiple images, you can use:
	// taloscdk.NewSecutiyGroup()
	// Default: Generates a new security group, opening ports: 6443, 50000, 50001 to the any peer
	SecurityGroup awsec2.SecurityGroup

	// Vpc selects the AWS VPC to deploy your instance into.
	// Default for NewSingleNode(): Default VPC
	Vpc awsec2.IVpc

	// Subnets to allow the instance to be deployed into
	// Default: &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PUBLIC}
	SubnetSelection *awsec2.SubnetSelection

	// CreateEIP enables an ElasticIP to be created and allocated to your instance.
	// This is generally used as the cluster endpoint in a single node cluster.
	// Default: jsii.Bool("true")
	CreateEIP *bool

	// IAMRole used when launching the instance.
	// If planning to create AWS load balancers, it's best to use
	// taloscdk.NewControlPlaneIAMRole() or taloscdk.NewWorkerIAMRole()
	// Default: NewControlPlaneIAMRole()
	IAMRole awsiam.Role
}

type SingleNode struct {
	constructs.Construct
	// SecurityGroup used or created by NewSingleNode()
	SecurityGroup awsec2.SecurityGroup

	// VPC of the node
	Vpc awsec2.IVpc

	// EIP (if allocated/assigned)
	EIP awsec2.CfnEIP
}

func (s *SingleNode) GetEIPAddress() *string {
	return s.EIP.Ref()
}

// NewSingleNode creates a new EC2 instance that runs Talos.
// Required SingleNodeProps:
//     TalosNodeConfig, EndpointToOverwrite (if TransformConfig==true)
func NewSingleNode(scope constructs.Construct, id *string, props *SingleNodeProps) SingleNode {
	construct := awscdk.NewConstruct(scope, jsii.String(*id))

	if props.ClusterName == nil {
		props.ClusterName = jsii.String("talos")
	}
	if props.NodeName == nil {
		props.NodeName = jsii.String("talos-node")
	}
	if props.Vpc == nil {
		props.Vpc = awsec2.Vpc_FromLookup(construct, jsii.String("VPC"), &awsec2.VpcLookupOptions{
			IsDefault: jsii.Bool(true),
		})
	}

	if props.TalosNodeConfig == nil {
		panic("TalosNodeConfig cannot be nil. taloscdk.LoadConfig() can be used to load the needed file.")
	}

	if props.SecurityGroup == nil {
		props.SecurityGroup = NewSecurityGroup(construct, jsii.String("SG"), &SecurityGroupProps{
			Vpc: props.Vpc,
		})
	}

	if props.InstanceType == nil {
		props.InstanceType = awsec2.InstanceType_Of(awsec2.InstanceClass_BURSTABLE3, awsec2.InstanceSize_SMALL)
	}

	if props.SubnetSelection == nil {
		props.SubnetSelection = &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PUBLIC}
	}

	if props.MachineImageName == nil && props.MachineImageAMI == nil {
		props.MachineImageName = jsii.String("talos-v0.11.2-*-amd64")
	}

	if props.CreateEIP == nil {
		props.CreateEIP = jsii.Bool(true)
	}

	var eip awsec2.CfnEIP
	if *props.CreateEIP {
		eip = awsec2.NewCfnEIP(construct, jsii.String("EIP"), &awsec2.CfnEIPProps{})
		eip.ApplyRemovalPolicy(awscdk.RemovalPolicy_DESTROY, nil)
	}

	if props.OverwriteValue == nil {
		props.OverwriteValue = eip.Ref()
	}

	if props.EndpointToOverwrite == nil && props.TransformConfig != nil {
		panic("Requested config transform but missing EndpointToOverwrite.")
	}

	if *props.TransformConfig {
		props.TalosNodeConfig = TransformConfig(props.TalosNodeConfig, *props.EndpointToOverwrite, *props.OverwriteValue)
	}

	if props.IAMRole == nil {
		props.IAMRole = NewControlPlaneIAMRole(construct, jsii.String("Role"))
	}

	var image awsec2.IMachineImage

	if props.MachineImageAMI != nil {
		image = awsec2.NewGenericLinuxImage(props.MachineImageAMI, &awsec2.GenericLinuxImageProps{
			UserData: awsec2.UserData_Custom(props.TalosNodeConfig),
		})
	} else {
		image = awsec2.NewLookupMachineImage(&awsec2.LookupMachineImageProps{
			Name:     props.MachineImageName,
			Owners:   jsii.Strings("540036508848"),
			UserData: awsec2.UserData_Custom(props.TalosNodeConfig),
		})
	}

	instance := awsec2.NewInstance(construct, jsii.String("Instance"), &awsec2.InstanceProps{
		InstanceName:  props.NodeName,
		InstanceType:  props.InstanceType,
		MachineImage:  image,
		Vpc:           props.Vpc,
		SecurityGroup: props.SecurityGroup,
		VpcSubnets:    props.SubnetSelection,
		Role:          props.IAMRole,
	})

	if *props.CreateEIP {
		awsec2.NewCfnEIPAssociation(construct, jsii.String("EIPAssoc"), &awsec2.CfnEIPAssociationProps{InstanceId: instance.InstanceId(), Eip: eip.Ref()})
	}

	awscdk.Tags_Of(construct).Add(jsii.String(fmt.Sprintf("kubernetes.io/cluster/%s", *props.ClusterName)), jsii.String("owned"), nil)
	TagSubnets(props.Vpc)

	return SingleNode{Construct: construct, SecurityGroup: props.SecurityGroup, Vpc: props.Vpc, EIP: eip}
}
