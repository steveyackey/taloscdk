package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/steveyackey/taloscdk"
)

type PublicClusterStackProps struct {
	awscdk.StackProps
}

func NewPublicClusterStack(scope constructs.Construct, id string, props *PublicClusterStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Look up the default VPC
	vpc := awsec2.Vpc_FromLookup(stack, jsii.String("VPC"), &awsec2.VpcLookupOptions{
		IsDefault: jsii.Bool(true),
	})

	// Load controlplane.yaml
	config, err := taloscdk.LoadConfig("./controlplane.yaml")
	if err != nil {
		panic("Could not load talos config")
	}

	// Create a new control plane Autoscaling Group
	cp := taloscdk.NewControlPlane(stack, jsii.String("TalosCP"), &taloscdk.ControlPlaneProps{
		TalosNodeConfig:     config,
		TransformConfig:     jsii.Bool(true),
		EndpointToOverwrite: jsii.String("talos.cluster"),
		Vpc:                 vpc,
		SubnetSelection:     &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PUBLIC},
	})

	// Load the join.yaml worker config
	workerConfig, err := taloscdk.LoadConfig("./join.yaml")
	if err != nil {
		panic("Could not load talos config")
	}

	// Create a new ASG for worker nodes
	// For an example of creating a separate worker security group, see the
	// private-cluster example.
	taloscdk.NewWorkerASG(stack, jsii.String("WorkerASG"), &taloscdk.WorkerASGProps{
		TalosNodeConfig:     workerConfig,
		TransformConfig:     jsii.Bool(true),
		EndpointToOverwrite: jsii.String("talos.cluster"),
		OverwriteValue:      cp.NLB.LoadBalancerDnsName(),
		Vpc:                 vpc,
		SecurityGroup:       cp.SecurityGroup,
		SubnetSelection:     &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PUBLIC},
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewPublicClusterStack(app, "PublicClusterStack", &PublicClusterStackProps{
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
	// return nil

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
