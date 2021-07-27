module single-node

go 1.16

replace github.com/steveyackey/taloscdk => ../../

require (
	github.com/aws/aws-cdk-go/awscdk v1.114.0-devpreview
	github.com/aws/constructs-go/constructs/v3 v3.3.97
	github.com/aws/jsii-runtime-go v1.31.0
	github.com/steveyackey/taloscdk v0.0.0-20210718000409-a46c3799927f
	github.com/stretchr/testify v1.7.0

	// for testing
	github.com/tidwall/gjson v1.7.4
)
