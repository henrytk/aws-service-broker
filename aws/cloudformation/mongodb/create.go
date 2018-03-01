package mongodb

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	awscf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/henrytk/aws-service-broker/aws/cloudformation/templates"
)

var (
	timeoutInMinutes int64 = 15
)

func (s *Service) CreateStack(id string, inputParameters InputParameters) (*awscf.CreateStackOutput, error) {
	parameters, err := s.BuildCreateStackParameters(inputParameters)
	if err != nil {
		return nil, err
	}
	createStackInput := s.BuildCreateStackInput(id, parameters)
	return s.Client.CreateStack(createStackInput)
}

func (s *Service) BuildCreateStackParameters(p InputParameters) ([]*awscf.Parameter, error) {
	usePreviousValue := false
	var parameters []*awscf.Parameter

	if p.BastionSecurityGroupId == "" {
		return parameters, errors.New("Error building MongoDB parameters: bastion security group ID is empty")
	} else {
		parameters = append(parameters, &awscf.Parameter{
			ParameterKey:     aws.String(string(bastionSecurityGroupIdSPK)),
			ParameterValue:   aws.String(p.BastionSecurityGroupId),
			UsePreviousValue: aws.Bool(usePreviousValue),
		})
	}

	if p.KeyPairName == "" {
		return parameters, errors.New("Error building MongoDB parameters: key pair name is empty")
	} else {
		parameters = append(parameters, &awscf.Parameter{
			ParameterKey:     aws.String(string(keyPairNameSPK)),
			ParameterValue:   aws.String(p.KeyPairName),
			UsePreviousValue: aws.Bool(usePreviousValue),
		})
	}

	if p.VpcId == "" {
		return parameters, errors.New("Error building MongoDB parameters: VPC ID is empty")
	} else {
		parameters = append(parameters, &awscf.Parameter{
			ParameterKey:     aws.String(string(vpcIdSPK)),
			ParameterValue:   aws.String(p.VpcId),
			UsePreviousValue: aws.Bool(usePreviousValue),
		})
	}

	if p.PrimaryNodeSubnetId == "" {
		return parameters, errors.New("Error building MongoDB parameters: primary node subnet ID is empty")
	} else {
		parameters = append(parameters, &awscf.Parameter{
			ParameterKey:     aws.String(string(primaryNodeSubnetIdSPK)),
			ParameterValue:   aws.String(p.PrimaryNodeSubnetId),
			UsePreviousValue: aws.Bool(usePreviousValue),
		})
	}

	if p.Secondary0NodeSubnetId == "" {
		return parameters, errors.New("Error building MongoDB parameters: secondary 0 node subnet ID is empty")
	} else {
		parameters = append(parameters, &awscf.Parameter{
			ParameterKey:     aws.String(string(secondary0NodeSubnetIdSPK)),
			ParameterValue:   aws.String(p.Secondary0NodeSubnetId),
			UsePreviousValue: aws.Bool(usePreviousValue),
		})
	}

	if p.Secondary1NodeSubnetId == "" {
		return parameters, errors.New("Error building MongoDB parameters: secondary 1 node subnet ID is empty")
	} else {
		parameters = append(parameters, &awscf.Parameter{
			ParameterKey:     aws.String(string(secondary1NodeSubnetIdSPK)),
			ParameterValue:   aws.String(p.Secondary1NodeSubnetId),
			UsePreviousValue: aws.Bool(usePreviousValue),
		})
	}

	if p.MongoDBAdminPassword == "" {
		return parameters, errors.New("Error building MongoDB parameters: MongoDB admin password is empty")
	} else {
		parameters = append(parameters, &awscf.Parameter{
			ParameterKey:     aws.String(string(mongoDBAdminPasswordSPK)),
			ParameterValue:   aws.String(p.MongoDBAdminPassword),
			UsePreviousValue: aws.Bool(usePreviousValue),
		})
	}

	if p.MongoDBAdminUsername != "" {
		parameters = append(parameters, &awscf.Parameter{
			ParameterKey:     aws.String(string(mongoDBAdminUsernameSPK)),
			ParameterValue:   aws.String(p.MongoDBAdminUsername),
			UsePreviousValue: aws.Bool(usePreviousValue),
		})
	}

	if p.MongoDBVersion != "" {
		parameters = append(parameters, &awscf.Parameter{
			ParameterKey:     aws.String(string(mongoDBVersionSPK)),
			ParameterValue:   aws.String(p.MongoDBVersion),
			UsePreviousValue: aws.Bool(usePreviousValue),
		})
	}

	if p.ClusterReplicaSetCount != "" {
		parameters = append(parameters, &awscf.Parameter{
			ParameterKey:     aws.String(string(clusterReplicaSetCountSPK)),
			ParameterValue:   aws.String(p.ClusterReplicaSetCount),
			UsePreviousValue: aws.Bool(usePreviousValue),
		})
	}

	if p.ReplicaShardIndex != "" {
		parameters = append(parameters, &awscf.Parameter{
			ParameterKey:     aws.String(string(replicaShardIndexSPK)),
			ParameterValue:   aws.String(p.ReplicaShardIndex),
			UsePreviousValue: aws.Bool(usePreviousValue),
		})
	}

	if p.VolumeSize != "" {
		parameters = append(parameters, &awscf.Parameter{
			ParameterKey:     aws.String(string(volumeSizeSPK)),
			ParameterValue:   aws.String(p.VolumeSize),
			UsePreviousValue: aws.Bool(usePreviousValue),
		})
	}

	if p.VolumeType != "" {
		parameters = append(parameters, &awscf.Parameter{
			ParameterKey:     aws.String(string(volumeTypeSPK)),
			ParameterValue:   aws.String(p.VolumeType),
			UsePreviousValue: aws.Bool(usePreviousValue),
		})
	}

	if p.Iops != "" {
		parameters = append(parameters, &awscf.Parameter{
			ParameterKey:     aws.String(string(iopsSPK)),
			ParameterValue:   aws.String(p.Iops),
			UsePreviousValue: aws.Bool(usePreviousValue),
		})
	}

	if p.NodeInstanceType != "" {
		parameters = append(parameters, &awscf.Parameter{
			ParameterKey:     aws.String(string(nodeInstanceTypeSPK)),
			ParameterValue:   aws.String(p.NodeInstanceType),
			UsePreviousValue: aws.Bool(usePreviousValue),
		})
	}

	return parameters, nil
}

func (s *Service) BuildCreateStackInput(id string, parameters []*awscf.Parameter) *awscf.CreateStackInput {
	stackName := s.GenerateStackName(id)
	mongoDBStackTemplate := string(templates.MongoDBStack)
	return &awscf.CreateStackInput{
		Capabilities:       capabilities,
		ClientRequestToken: aws.String("create-" + stackName),
		Parameters:         parameters,
		StackName:          aws.String(stackName),
		TemplateBody:       aws.String(mongoDBStackTemplate),
		TimeoutInMinutes:   aws.Int64(timeoutInMinutes),
	}
}
