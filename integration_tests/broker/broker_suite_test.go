package broker_test

import (
	"github.com/henrytk/aws-service-broker/integration_tests/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"

	"testing"
)

var (
	region                 = "eu-west-1"
	testID                 string
	vpc                    *helpers.Vpc
	id                     string
	keyPairName            string
	primaryNodeSubnetId    string
	secondary0NodeSubnetId string
	secondary1NodeSubnetId string
	vpcId                  string
	bastionSecurityGroupId string
)

func TestBroker(t *testing.T) {
	BeforeSuite(func() {
		testID = uuid.NewV4().String()
		vpc = helpers.SetupVpc(region, testID)
		vpcId = *vpc.VpcId
		primaryNodeSubnetId = *vpc.Subnets[3].SubnetId
		secondary0NodeSubnetId = *vpc.Subnets[4].SubnetId
		secondary1NodeSubnetId = *vpc.Subnets[5].SubnetId
		bastionSecurityGroupId = *vpc.SecurityGroups[0].GroupId
		keyPairName = vpc.KeyPairName
	})

	AfterSuite(func() {
		helpers.DestroyVpc(vpc)
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Broker Suite")
}
