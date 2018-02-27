package mongodb

import "github.com/aws/aws-sdk-go/aws"

var (
	capabilities = []*string{aws.String("CAPABILITY_IAM")}
)

type StackParameterKey string

var (
	bastionSecurityGroupIdSPK StackParameterKey = "BastionSecurityGroupID"
	keyPairNameSPK            StackParameterKey = "KeyPairName"
	vpcIdSPK                  StackParameterKey = "VPC"
	primaryNodeSubnetIdSPK    StackParameterKey = "PrimaryNodeSubnet"
	secondary0NodeSubnetIdSPK StackParameterKey = "Secondary0NodeSubnet"
	secondary1NodeSubnetIdSPK StackParameterKey = "Secondary1NodeSubnet"
	mongoDBAdminPasswordSPK   StackParameterKey = "MongoDBAdminPassword"
	mongoDBAdminUsernameSPK   StackParameterKey = "MongoDBAdminUsername"
	mongoDBVersionSPK         StackParameterKey = "MongoDBVersion"
	clusterReplicaSetCountSPK StackParameterKey = "ClusterReplicaSetCount"
	replicaShardIndexSPK      StackParameterKey = "ReplicaShardIndex"
	volumeSizeSPK             StackParameterKey = "VolumeSize"
	volumeTypeSPK             StackParameterKey = "VolumeType"
	iopsSPK                   StackParameterKey = "Iops"
	nodeInstanceTypeSPK       StackParameterKey = "NodeInstanceType"
)

type InputParameters struct {
	BastionSecurityGroupId string
	KeyPairName            string
	VpcId                  string
	PrimaryNodeSubnetId    string
	Secondary0NodeSubnetId string
	Secondary1NodeSubnetId string
	MongoDBAdminPassword   string
	MongoDBAdminUsername   string
	MongoDBVersion         string
	ClusterReplicaSetCount string
	ReplicaShardIndex      string
	VolumeSize             string
	VolumeType             string
	Iops                   string
	NodeInstanceType       string
}
