package mongodb

import (
	"github.com/aws/aws-sdk-go/aws"
	awscf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/henrytk/aws-service-broker/aws/cloudformation/templates"
)

type StackParameterKey string

var (
	timeoutInMinutes int64 = 15

	capabilities = []*string{aws.String("CAPABILITY_IAM")}

	usePreviousValue                            = false
	keyPairNameSPK            StackParameterKey = "KeyPairName"
	primaryNodeSubnetIdSPK    StackParameterKey = "PrimaryNodeSubnet"
	secondary0NodeSubnetIdSPK StackParameterKey = "Secondary0NodeSubnet"
	secondary1NodeSubnetIdSPK StackParameterKey = "Secondary1NodeSubnet"
	mongoDBAdminPasswordSPK   StackParameterKey = "MongoDBAdminPassword"
	vpcIdSPK                  StackParameterKey = "VPC"
	bastionSecurityGroupIdSPK StackParameterKey = "BastionSecurityGroupID"
)

func (s Service) CreateStack(
	id,
	keyPairName,
	primaryNodeSubnetId,
	secondary0NodeSubnetId,
	secondary1NodeSubnetId,
	mongoDBAdminPassword,
	vpcId,
	bastionSecurityGroupId string,
) (*awscf.CreateStackOutput, error) {
	parameters := buildParameters(
		keyPairName,
		primaryNodeSubnetId,
		secondary0NodeSubnetId,
		secondary1NodeSubnetId,
		mongoDBAdminPassword,
		vpcId,
		bastionSecurityGroupId,
	)
	createStackInput := BuildCreateStackInput(id, parameters)
	return s.Client.CreateStack(&createStackInput)
}

func buildParameters(
	keyPairName,
	primaryNodeSubnetId,
	secondary0NodeSubnetId,
	secondary1NodeSubnetId,
	mongoDBAdminPassword,
	vpcId,
	bastionSecurityGroupId string,
) []*awscf.Parameter {
	return []*awscf.Parameter{
		&awscf.Parameter{
			ParameterKey:     aws.String(string(keyPairNameSPK)),
			ParameterValue:   aws.String(keyPairName),
			UsePreviousValue: aws.Bool(usePreviousValue),
		},
		&awscf.Parameter{
			ParameterKey:     aws.String(string(primaryNodeSubnetIdSPK)),
			ParameterValue:   aws.String(primaryNodeSubnetId),
			UsePreviousValue: aws.Bool(usePreviousValue),
		},
		&awscf.Parameter{
			ParameterKey:     aws.String(string(secondary0NodeSubnetIdSPK)),
			ParameterValue:   aws.String(secondary0NodeSubnetId),
			UsePreviousValue: aws.Bool(usePreviousValue),
		},
		&awscf.Parameter{
			ParameterKey:     aws.String(string(secondary1NodeSubnetIdSPK)),
			ParameterValue:   aws.String(secondary1NodeSubnetId),
			UsePreviousValue: aws.Bool(usePreviousValue),
		},
		&awscf.Parameter{
			ParameterKey:     aws.String(string(mongoDBAdminPasswordSPK)),
			ParameterValue:   aws.String(mongoDBAdminPassword),
			UsePreviousValue: aws.Bool(usePreviousValue),
		},
		&awscf.Parameter{
			ParameterKey:     aws.String(string(vpcIdSPK)),
			ParameterValue:   aws.String(vpcId),
			UsePreviousValue: aws.Bool(usePreviousValue),
		},
		&awscf.Parameter{
			ParameterKey:     aws.String(string(bastionSecurityGroupIdSPK)),
			ParameterValue:   aws.String(bastionSecurityGroupId),
			UsePreviousValue: aws.Bool(usePreviousValue),
		},
	}
}

func BuildCreateStackInput(id string, parameters []*awscf.Parameter) awscf.CreateStackInput {
	mongoDBStack := string(templates.MongoDBStack)
	return awscf.CreateStackInput{
		Capabilities:       capabilities,
		ClientRequestToken: aws.String("create-" + id),
		Parameters:         parameters,
		StackName:          aws.String(id),
		TemplateBody:       aws.String(mongoDBStack),
		TimeoutInMinutes:   &timeoutInMinutes,
	}
}
