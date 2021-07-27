package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/steveyackey/taloscdk"
)

type PrivateClusterStackProps struct {
	awscdk.StackProps
}

func NewPrivateClusterStack(scope constructs.Construct, id string, props *PrivateClusterStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Create a new VPC and tag it
	vpc := awsec2.NewVpc(stack, jsii.String("TalosVPC"), nil)

	// Load controlplane.yaml
	config, err := taloscdk.LoadConfig("./controlplane.yaml")
	if err != nil {
		panic("Could not load talos config")
	}

	// Create a new control plane Autoscaling Group with an NLB only available to private subnets
	cp := taloscdk.NewControlPlane(stack, jsii.String("TalosCP"), &taloscdk.ControlPlaneProps{
		ClusterName:         jsii.String("talos"),
		TalosNodeConfig:     config,
		TransformConfig:     jsii.Bool(true),
		EndpointToOverwrite: jsii.String("talos.cluster"),
		Vpc:                 vpc,
		SubnetSelection:     &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PRIVATE},
		InternetFacingNLB:   jsii.Bool(false),
		MinInstances:        jsii.Number(3), // If you don't include Max, it defaults MaxInstances to equal MinInstances
	})

	// Load the join.yaml worker config
	workerConfig, err := taloscdk.LoadConfig("./join.yaml")
	if err != nil {
		panic("Could not load talos config")
	}

	// Create a new SG for the worker ASG. Could use cp.SecurityGroup instead, but this
	// prevents having to expose port 6443 on the worker nodes.
	workerSG := awsec2.NewSecurityGroup(stack, jsii.String("TalosSG"), &awsec2.SecurityGroupProps{
		Vpc:              vpc,
		AllowAllOutbound: jsii.Bool(true),
		Description:      jsii.String("Talos Security Group"),
	})

	workerSG.AddIngressRule(
		workerSG,
		awsec2.Port_AllTraffic(),
		jsii.String("Allow all internal traffic between worker nodes"),
		jsii.Bool(false),
	)

	workerSG.AddIngressRule(
		cp.SecurityGroup,
		awsec2.Port_AllTraffic(),
		jsii.String("Allow all internal traffic between worker and control plane nodes"),
		jsii.Bool(false),
	)

	cp.SecurityGroup.AddIngressRule(
		workerSG,
		awsec2.Port_AllTraffic(),
		jsii.String("Allow all internal traffic between worker and control plane nodes"),
		jsii.Bool(false),
	)

	taloscdk.NewWorkerASG(stack, jsii.String("WorkerASG"), &taloscdk.WorkerASGProps{
		TalosNodeConfig:     workerConfig,
		TransformConfig:     jsii.Bool(true),
		EndpointToOverwrite: jsii.String("talos.cluster"),
		OverwriteValue:      cp.NLB.LoadBalancerDnsName(),
		Vpc:                 vpc,
		SecurityGroup:       workerSG,
		SubnetSelection:     &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PRIVATE},
		MinInstances:        jsii.Number(3),
		MaxInstances:        jsii.Number(5),
	})

	// New Bastion host for using SSM Session Manager.
	// Use this instance to bootstrap the cluster and run kubectl commands.
	awsec2.NewBastionHostLinux(stack, jsii.String("BastionHost"), &awsec2.BastionHostLinuxProps{
		Vpc: vpc,
	})
	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewPrivateClusterStack(app, "PrivateClusterStack", &PrivateClusterStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	//return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
