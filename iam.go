package taloscdk

import (
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

// NewControlPlaneIAMRole returns a new awsiam.Role with minimum permissions
// to utilize the aws-controller-manager for creating ELBs from your cluster.
// Returns a role with an inline policy created via taloscdk.NewControlPlaneIAMPolicyDocument()
func NewControlPlaneIAMRole(scope constructs.Construct, id *string) awsiam.Role {
	return awsiam.NewRole(scope, id, &awsiam.RoleProps{
		InlinePolicies: &map[string]awsiam.PolicyDocument{
			"ControlPlanePolicy": NewControlPlaneIAMPolicyDocument(scope, jsii.String("ControlPlanePolicy")),
		},
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), nil),
	})
}

// NewWorkerIAMRole returns a new awsiam.Role with minimum permissions
// to utilize the aws-controller-manager for creating ELBs from your cluster.
// Returns a role with an inline policy created via taloscdk.NewWorkerIAMPolicyDocument()
func NewWorkerIAMRole(scope constructs.Construct, id *string) awsiam.Role {
	return awsiam.NewRole(scope, id, &awsiam.RoleProps{
		InlinePolicies: &map[string]awsiam.PolicyDocument{
			"WorkerPolicy": NewWorkerIAMPolicyDocument(scope, jsii.String("WorkerPolicy")),
		},
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), nil),
	})
}

func NewControlPlaneIAMPolicyDocument(scope constructs.Construct, id *string) awsiam.PolicyDocument {
	return awsiam.NewPolicyDocument(&awsiam.PolicyDocumentProps{
		Statements: &[]awsiam.PolicyStatement{
			//awsiam.PolicyStatement_FromJson(controlPlanePolicy),
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Effect: awsiam.Effect_ALLOW,
				Actions: &[]*string{
					jsii.String("autoscaling:DescribeAutoScalingGroups"),
					jsii.String("autoscaling:DescribeLaunchConfigurations"),
					jsii.String("autoscaling:DescribeTags"),
					jsii.String("ec2:DescribeInstances"),
					jsii.String("ec2:DescribeRegions"),
					jsii.String("ec2:DescribeRouteTables"),
					jsii.String("ec2:DescribeSecurityGroups"),
					jsii.String("ec2:DescribeSubnets"),
					jsii.String("ec2:DescribeVolumes"),
					jsii.String("ec2:CreateSecurityGroup"),
					jsii.String("ec2:CreateTags"),
					jsii.String("ec2:CreateVolume"),
					jsii.String("ec2:ModifyInstanceAttribute"),
					jsii.String("ec2:ModifyVolume"),
					jsii.String("ec2:AttachVolume"),
					jsii.String("ec2:AuthorizeSecurityGroupIngress"),
					jsii.String("ec2:CreateRoute"),
					jsii.String("ec2:DeleteRoute"),
					jsii.String("ec2:DeleteSecurityGroup"),
					jsii.String("ec2:DeleteVolume"),
					jsii.String("ec2:DetachVolume"),
					jsii.String("ec2:RevokeSecurityGroupIngress"),
					jsii.String("ec2:DescribeVpcs"),
					jsii.String("elasticloadbalancing:AddTags"),
					jsii.String("elasticloadbalancing:AttachLoadBalancerToSubnets"),
					jsii.String("elasticloadbalancing:ApplySecurityGroupsToLoadBalancer"),
					jsii.String("elasticloadbalancing:CreateLoadBalancer"),
					jsii.String("elasticloadbalancing:CreateLoadBalancerPolicy"),
					jsii.String("elasticloadbalancing:CreateLoadBalancerListeners"),
					jsii.String("elasticloadbalancing:ConfigureHealthCheck"),
					jsii.String("elasticloadbalancing:DeleteLoadBalancer"),
					jsii.String("elasticloadbalancing:DeleteLoadBalancerListeners"),
					jsii.String("elasticloadbalancing:DescribeLoadBalancers"),
					jsii.String("elasticloadbalancing:DescribeLoadBalancerAttributes"),
					jsii.String("elasticloadbalancing:DetachLoadBalancerFromSubnets"),
					jsii.String("elasticloadbalancing:DeregisterInstancesFromLoadBalancer"),
					jsii.String("elasticloadbalancing:ModifyLoadBalancerAttributes"),
					jsii.String("elasticloadbalancing:RegisterInstancesWithLoadBalancer"),
					jsii.String("elasticloadbalancing:SetLoadBalancerPoliciesForBackendServer"),
					jsii.String("elasticloadbalancing:AddTags"),
					jsii.String("elasticloadbalancing:CreateListener"),
					jsii.String("elasticloadbalancing:CreateTargetGroup"),
					jsii.String("elasticloadbalancing:DeleteListener"),
					jsii.String("elasticloadbalancing:DeleteTargetGroup"),
					jsii.String("elasticloadbalancing:DescribeListeners"),
					jsii.String("elasticloadbalancing:DescribeLoadBalancerPolicies"),
					jsii.String("elasticloadbalancing:DescribeTargetGroups"),
					jsii.String("elasticloadbalancing:DescribeTargetHealth"),
					jsii.String("elasticloadbalancing:ModifyListener"),
					jsii.String("elasticloadbalancing:ModifyTargetGroup"),
					jsii.String("elasticloadbalancing:RegisterTargets"),
					jsii.String("elasticloadbalancing:DeregisterTargets"),
					jsii.String("elasticloadbalancing:SetLoadBalancerPoliciesOfListener"),
					jsii.String("iam:CreateServiceLinkedRole"),
					jsii.String("kms:DescribeKey"),
				},
				Resources: jsii.Strings("*"),
			}),
		},
	})
}

func NewWorkerIAMPolicyDocument(scope constructs.Construct, id *string) awsiam.PolicyDocument {
	policy := awsiam.NewPolicyDocument(&awsiam.PolicyDocumentProps{
		Statements: &[]awsiam.PolicyStatement{
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Effect: awsiam.Effect_ALLOW,
				Actions: &[]*string{
					jsii.String("ec2:DescribeInstances"),
					jsii.String("ec2:DescribeRegions"),
					jsii.String("ecr:GetAuthorizationToken"),
					jsii.String("ecr:BatchCheckLayerAvailability"),
					jsii.String("ecr:GetDownloadUrlForLayer"),
					jsii.String("ecr:GetRepositoryPolicy"),
					jsii.String("ecr:DescribeRepositories"),
					jsii.String("ecr:ListImages"),
					jsii.String("ecr:BatchGetImage"),
				},
				Resources: jsii.Strings("*"),
			}),
		},
	})

	return policy
}
