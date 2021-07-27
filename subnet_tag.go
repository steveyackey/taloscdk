package taloscdk

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/jsii-runtime-go"
)

// TagSubnets is used to tag all subnets within a vpc with the appropriate ELB role.
// It is used to determine which subnets in a VPC can be used within an ELB.
// Ref: https://github.com/aws/aws-cdk/blob/6f2384ddc180e944c9564a543351b8df2f75c1a7/packages/%40aws-cdk/aws-eks/lib/cluster.ts#L1499-L1513
func TagSubnets(vpc awsec2.IVpc) {
	private := vpc.PrivateSubnets()
	public := vpc.PublicSubnets()

	for _, s := range *private {
		awscdk.Tags_Of(s).Add(jsii.String("kubernetes.io/role/internal-elb"), jsii.String("1"), nil)
	}

	for _, s := range *public {
		awscdk.Tags_Of(s).Add(jsii.String("kubernetes.io/role/elb"), jsii.String("1"), nil)
	}
}
