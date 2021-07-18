package taloscdk

import (
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type SingleNodeProps struct {
	NodeName       *string
	TalosVersion   *string
	SecurityGroup  *awsec2.SecurityGroup
	Role           TalosRole
	SingleNodeOnly *bool
}

func NewSingleNode(scope constructs.Construct, id *string, props *SingleNodeProps) constructs.Construct {
	construct := constructs.NewConstruct(scope, jsii.String("TalosNode"), &constructs.ConstructOptions{})
	NewSecurityGroup(construct, jsii.String("TalosSG"), &SecurityGroupProps{})
	return construct
}
