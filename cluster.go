package taloscdk

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsautoscaling"
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	awselbv2 "github.com/aws/aws-cdk-go/awscdk/awselasticloadbalancingv2"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type ControlPlaneProps struct {
	// ClusterName is used for tagging all resources with kubernetes.io/cluster/<name>=owned
	// Default: talos
	ClusterName *string

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
	// Default: NLB DNS name.
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

	// MinInstances to use with the autoscaling group
	// Default: jsii.Number(1)
	MinInstances *float64

	// MaxInstances to use with the autoscaling group
	// Default: jsii.Number(1)
	MaxInstances *float64

	// IAMRole used when launching the instance.
	// If planning to create AWS load balancers, it's best to use
	// taloscdk.NewControlPlaneIAMRole() or taloscdk.NewWorkerIAMRole()
	// Default: NewControlPlaneIAMRole()
	IAMRole awsiam.Role

	// DesiredCapacity of the autoscaling group
	// Best practice: leave it nil. If you set a value, it will always reset the number of
	// nodes to this number each time you run `cdk deploy`
	// Default: nil
	DesiredCapacity *float64 // leave nil if using any autoscaling features, otherwise it will be replaced each `cdk deploy`

	// InternetFacingNLB determines whether or not the control plane NLB should be
	// created in public subnets (or left in the private subnets)
	// Default: jsii.Bool(true)
	InternetFacingNLB *bool
}

type ControlPlane struct {
	constructs.Construct
	SecurityGroup awsec2.SecurityGroup
	Vpc           awsec2.IVpc
	ASG           awsautoscaling.AutoScalingGroup
	NLB           awselbv2.NetworkLoadBalancer
	IAMRole       awsiam.Role
}

type WorkerASGProps struct {
	// ClusterName is used for tagging all resources with kubernetes.io/cluster/<name>=owned
	// Default: talos
	ClusterName *string

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
	// Default: NLB DNS name.
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
	// Vpc is required and stack will panic if not given.
	// awsec2.NewVpc(), awsec2.Vpc_FromLookup() will return a usable VPC
	Vpc awsec2.IVpc

	// Subnets to allow the instance to be deployed into
	// Default: &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PUBLIC}
	SubnetSelection *awsec2.SubnetSelection

	// MinInstances to use with the autoscaling group
	// Default: jsii.Number(1)
	MinInstances *float64

	// MaxInstances to use with the autoscaling group
	// Default: jsii.Number(1)
	MaxInstances *float64

	// IAMRole used when launching the instance.
	// If planning to create AWS load balancers, it's best to use
	// taloscdk.NewControlPlaneIAMRole() or taloscdk.NewWorkerIAMRole()
	// Default: NewWorkerIAMRole()
	IAMRole awsiam.Role

	// DesiredCapacity of the autoscaling group
	// Best practice: leave it nil. If you set a value, it will always reset the number of
	// nodes to this number each time you run `cdk deploy`
	// Default: nil
	DesiredCapacity *float64 // leave nil if using any autoscaling features, otherwise it will be replaced each `cdk deploy`
}

// NewControlPlane creates a new NLB and control plane backed by an autoscaling group
func NewControlPlane(scope constructs.Construct, id *string, props *ControlPlaneProps) ControlPlane {
	construct := awscdk.NewConstruct(scope, jsii.String(*id))
	if props.ClusterName == nil {
		props.ClusterName = jsii.String("talos")
	}

	if props.Vpc == nil {
		panic("Vpc is required")
	}

	if props.MachineImageName == nil && props.MachineImageAMI == nil {
		props.MachineImageName = jsii.String("talos-v0.11.2-*-amd64")
	}

	if props.SecurityGroup == nil {
		props.SecurityGroup = NewSecurityGroup(construct, jsii.String("SG"), &SecurityGroupProps{
			Vpc: props.Vpc,
		})
	}

	if props.SubnetSelection == nil {
		props.SubnetSelection = &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PUBLIC}
	}

	if props.InstanceType == nil {
		props.InstanceType = awsec2.InstanceType_Of(awsec2.InstanceClass_BURSTABLE3, awsec2.InstanceSize_SMALL)
	}

	if props.IAMRole == nil {
		props.IAMRole = NewControlPlaneIAMRole(construct, jsii.String("Role"))
	}

	nlb := awselbv2.NewNetworkLoadBalancer(construct, jsii.String("CP-NLB"), &awselbv2.NetworkLoadBalancerProps{
		Vpc:            props.Vpc,
		InternetFacing: jsii.Bool(true),
	})

	if props.OverwriteValue == nil {
		props.OverwriteValue = nlb.LoadBalancerDnsName()
	}

	if props.EndpointToOverwrite == nil && props.TransformConfig != nil {
		panic("Requested config transform but missing EndpointToOverwrite.")
	}

	if *props.TransformConfig {
		props.TalosNodeConfig = TransformConfig(props.TalosNodeConfig, *props.EndpointToOverwrite, *props.OverwriteValue)
	}

	if props.InternetFacingNLB == nil {
		props.InternetFacingNLB = jsii.Bool(true)
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

	TagSubnets(props.Vpc)

	cpAsg := awsautoscaling.NewAutoScalingGroup(construct, jsii.String("TalosCP"), &awsautoscaling.AutoScalingGroupProps{
		AllowAllOutbound: jsii.Bool(true),
		DesiredCapacity:  props.DesiredCapacity,
		MinCapacity:      props.MinInstances,
		MaxCapacity:      props.MaxInstances,
		VpcSubnets:       props.SubnetSelection,
		Vpc:              props.Vpc,
		InstanceType:     props.InstanceType,
		MachineImage:     image,
		Role:             props.IAMRole,
		SecurityGroup:    props.SecurityGroup,
	})

	targets := awselbv2.NewNetworkTargetGroup(construct, jsii.String("targetgroup-6443"), &awselbv2.NetworkTargetGroupProps{
		Port: jsii.Number(6443),
		HealthCheck: &awselbv2.HealthCheck{
			Enabled:  jsii.Bool(true),
			Port:     jsii.String("6443"),
			Protocol: awselbv2.Protocol_TCP,
		},
		TargetType: awselbv2.TargetType_INSTANCE,
		Vpc:        props.Vpc,
	})

	cpAsg.AttachToNetworkTargetGroup(targets)

	nlb.AddListener(jsii.String("talos-cp-listener-6443"), &awselbv2.BaseNetworkListenerProps{
		Port:                jsii.Number(6443),
		Protocol:            awselbv2.Protocol_TCP,
		DefaultTargetGroups: &[]awselbv2.INetworkTargetGroup{targets},
	})

	awscdk.Tags_Of(construct).Add(jsii.String(fmt.Sprintf("kubernetes.io/cluster/%s", *props.ClusterName)), jsii.String("owned"), &awscdk.TagProps{ApplyToLaunchedInstances: jsii.Bool(true)})

	return ControlPlane{Construct: construct, SecurityGroup: props.SecurityGroup, Vpc: props.Vpc, ASG: cpAsg, NLB: nlb, IAMRole: props.IAMRole}
}

func NewWorkerASG(scope constructs.Construct, id *string, props *WorkerASGProps) awsautoscaling.AutoScalingGroup {
	construct := awscdk.NewConstruct(scope, jsii.String(*id))

	if props.ClusterName == nil {
		props.ClusterName = jsii.String("talos")
	}

	if props.Vpc == nil {
		panic("Vpc is required")
	}

	if props.MachineImageName == nil && props.MachineImageAMI == nil {
		props.MachineImageName = jsii.String("talos-v0.11.2-*-amd64")
	}

	if props.SecurityGroup == nil {
		props.SecurityGroup = NewSecurityGroup(construct, jsii.String("SG"), &SecurityGroupProps{
			Vpc: props.Vpc,
		})
	}

	if props.SubnetSelection == nil {
		props.SubnetSelection = &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PUBLIC}
	}

	if props.InstanceType == nil {
		props.InstanceType = awsec2.InstanceType_Of(awsec2.InstanceClass_BURSTABLE3, awsec2.InstanceSize_SMALL)
	}

	if props.IAMRole == nil {
		props.IAMRole = NewWorkerIAMRole(construct, jsii.String("Role"))
	}

	nlb := awselbv2.NewNetworkLoadBalancer(construct, jsii.String("CP-NLB"), &awselbv2.NetworkLoadBalancerProps{
		Vpc:            props.Vpc,
		InternetFacing: jsii.Bool(true),
	})

	if props.OverwriteValue == nil {
		props.OverwriteValue = nlb.LoadBalancerDnsName()
	}

	if props.EndpointToOverwrite == nil && props.TransformConfig != nil {
		panic("Requested config transform but missing EndpointToOverwrite.")
	}

	if *props.TransformConfig {
		props.TalosNodeConfig = TransformConfig(props.TalosNodeConfig, *props.EndpointToOverwrite, *props.OverwriteValue)
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
	TagSubnets(props.Vpc)

	asg := awsautoscaling.NewAutoScalingGroup(construct, jsii.String("WorkerASG"), &awsautoscaling.AutoScalingGroupProps{
		AllowAllOutbound: jsii.Bool(true),
		DesiredCapacity:  props.DesiredCapacity,
		MinCapacity:      props.MinInstances,
		MaxCapacity:      props.MaxInstances,
		VpcSubnets:       props.SubnetSelection,
		Vpc:              props.Vpc,
		InstanceType:     props.InstanceType,
		MachineImage:     image,
		Role:             props.IAMRole,
		SecurityGroup:    props.SecurityGroup,
	})

	awscdk.Tags_Of(construct).Add(jsii.String(fmt.Sprintf("kubernetes.io/cluster/%s", *props.ClusterName)), jsii.String("owned"), &awscdk.TagProps{ApplyToLaunchedInstances: jsii.Bool(true)})

	return asg
}
