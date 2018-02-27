package mongodb

import (
	"github.com/aws/aws-sdk-go/aws"
	awscf "github.com/aws/aws-sdk-go/service/cloudformation"
)

func (s *Service) BuildUpdateStackParameters(p InputParameters) []*awscf.Parameter {
	var parameters []*awscf.Parameter

	value, usePreviousValue := updateParameterValue(p.BastionSecurityGroupId)
	parameters = append(parameters, &awscf.Parameter{
		ParameterKey:     aws.String(string(bastionSecurityGroupIdSPK)),
		ParameterValue:   value,
		UsePreviousValue: usePreviousValue,
	})

	value, usePreviousValue = updateParameterValue(p.KeyPairName)
	parameters = append(parameters, &awscf.Parameter{
		ParameterKey:     aws.String(string(keyPairNameSPK)),
		ParameterValue:   value,
		UsePreviousValue: usePreviousValue,
	})

	value, usePreviousValue = updateParameterValue(p.VpcId)
	parameters = append(parameters, &awscf.Parameter{
		ParameterKey:     aws.String(string(vpcIdSPK)),
		ParameterValue:   value,
		UsePreviousValue: usePreviousValue,
	})

	value, usePreviousValue = updateParameterValue(p.PrimaryNodeSubnetId)
	parameters = append(parameters, &awscf.Parameter{
		ParameterKey:     aws.String(string(primaryNodeSubnetIdSPK)),
		ParameterValue:   value,
		UsePreviousValue: usePreviousValue,
	})

	value, usePreviousValue = updateParameterValue(p.Secondary0NodeSubnetId)
	parameters = append(parameters, &awscf.Parameter{
		ParameterKey:     aws.String(string(secondary0NodeSubnetIdSPK)),
		ParameterValue:   value,
		UsePreviousValue: usePreviousValue,
	})

	value, usePreviousValue = updateParameterValue(p.Secondary1NodeSubnetId)
	parameters = append(parameters, &awscf.Parameter{
		ParameterKey:     aws.String(string(secondary1NodeSubnetIdSPK)),
		ParameterValue:   value,
		UsePreviousValue: usePreviousValue,
	})

	value, usePreviousValue = updateParameterValue(p.MongoDBAdminPassword)
	parameters = append(parameters, &awscf.Parameter{
		ParameterKey:     aws.String(string(mongoDBAdminPasswordSPK)),
		ParameterValue:   value,
		UsePreviousValue: usePreviousValue,
	})

	value, usePreviousValue = updateParameterValue(p.MongoDBAdminUsername)
	parameters = append(parameters, &awscf.Parameter{
		ParameterKey:     aws.String(string(mongoDBAdminUsernameSPK)),
		ParameterValue:   value,
		UsePreviousValue: usePreviousValue,
	})

	value, usePreviousValue = updateParameterValue(p.MongoDBVersion)
	parameters = append(parameters, &awscf.Parameter{
		ParameterKey:     aws.String(string(mongoDBVersionSPK)),
		ParameterValue:   value,
		UsePreviousValue: usePreviousValue,
	})

	value, usePreviousValue = updateParameterValue(p.ClusterReplicaSetCount)
	parameters = append(parameters, &awscf.Parameter{
		ParameterKey:     aws.String(string(clusterReplicaSetCountSPK)),
		ParameterValue:   value,
		UsePreviousValue: usePreviousValue,
	})

	value, usePreviousValue = updateParameterValue(p.ReplicaShardIndex)
	parameters = append(parameters, &awscf.Parameter{
		ParameterKey:     aws.String(string(replicaShardIndexSPK)),
		ParameterValue:   value,
		UsePreviousValue: usePreviousValue,
	})

	value, usePreviousValue = updateParameterValue(p.VolumeSize)
	parameters = append(parameters, &awscf.Parameter{
		ParameterKey:     aws.String(string(volumeSizeSPK)),
		ParameterValue:   value,
		UsePreviousValue: usePreviousValue,
	})

	value, usePreviousValue = updateParameterValue(p.VolumeType)
	parameters = append(parameters, &awscf.Parameter{
		ParameterKey:     aws.String(string(volumeTypeSPK)),
		ParameterValue:   value,
		UsePreviousValue: usePreviousValue,
	})

	value, usePreviousValue = updateParameterValue(p.Iops)
	parameters = append(parameters, &awscf.Parameter{
		ParameterKey:     aws.String(string(iopsSPK)),
		ParameterValue:   value,
		UsePreviousValue: usePreviousValue,
	})

	value, usePreviousValue = updateParameterValue(p.NodeInstanceType)
	parameters = append(parameters, &awscf.Parameter{
		ParameterKey:     aws.String(string(nodeInstanceTypeSPK)),
		ParameterValue:   value,
		UsePreviousValue: usePreviousValue,
	})
	return parameters
}

func updateParameterValue(input string) (*string, *bool) {
	if input == "" {
		return nil, aws.Bool(true)
	}
	return aws.String(input), aws.Bool(false)
}
