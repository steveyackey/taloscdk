package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/steveyackey/taloscdk"
)

type SingleNodeStackProps struct {
	awscdk.StackProps
}

func NewSingleNodeStack(scope constructs.Construct, id string, props *SingleNodeStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	config, err := taloscdk.LoadConfig("./controlplane.yaml")
	if err != nil {
		panic("Could not load talos config")
	}

	cp := taloscdk.NewSingleNode(stack, jsii.String("TalsoSingleNodeCluster"), &taloscdk.SingleNodeProps{
		TalosNodeConfig:     config,
		TransformConfig:     jsii.Bool(true),
		EndpointToOverwrite: jsii.String("talos.cluster"),
	})

	// Optional Second Node.
	// Creating a second node will allow you to use the aws-controller-manager for
	// loadbalancers without needing to remove the master node label on the node
	// (along with the nodeAffinity from the aws-controller-manager).

	// workerConfig, err := taloscdk.LoadConfig("./join.yaml")
	// if err != nil {
	// 	panic("Could not load talos config")
	// }

	// taloscdk.NewSingleNode(stack, jsii.String("TalosWorker"), &taloscdk.SingleNodeProps{
	// 	TalosNodeConfig:     workerConfig,
	// 	TransformConfig:     jsii.Bool(true),
	// 	EndpointToOverwrite: jsii.String("talos.cluster"),
	// 	OverwriteValue:      cp.GetEIPAddress(),
	// 	Vpc: awsec2.Vpc_FromLookup(stack, jsii.String("VPC"), &awsec2.VpcLookupOptions{
	// 		IsDefault: jsii.Bool(true),
	// 	}),
	// 	SecurityGroup: cp.SecurityGroup,
	// 	CreateEIP:     jsii.Bool(false),
	// 	IAMRole:       taloscdk.NewWorkerIAMRole(stack, jsii.String("WorkerRole")),
	// })

	// Output the EIP. You may need to clean this up manually when destroying this stack.
	awscdk.NewCfnOutput(stack, jsii.String("TalosSingleNodeClusterEndpoint"), &awscdk.CfnOutputProps{
		Value:       cp.GetEIPAddress(),
		Description: jsii.String("Use this IP address in your talosconfig as the endpoint and node."),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewSingleNodeStack(app, "SingleNodeStack", &SingleNodeStackProps{
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
