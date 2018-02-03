package mongodb_test

import (
	"github.com/henrytk/aws-service-broker/aws/cloudformation/mongodb"
	"github.com/henrytk/aws-service-broker/integration_tests/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"

	"testing"
)

var (
	region                 = "eu-west-1"
	mongoDBService         *mongodb.Service
	instanceID             string
	vpc                    *helpers.Vpc
	err                    error
	id                     string
	ok                     bool
	keyPairName            string
	primaryNodeSubnetId    string
	secondary0NodeSubnetId string
	secondary1NodeSubnetId string
	mongoDBAdminPassword   string
	vpcId                  string
	bastionSecurityGroupId string
)

func TestMongodb(t *testing.T) {
	BeforeSuite(func() {
		instanceID = uuid.NewV4().String()
		vpc = helpers.SetupVpc(region, instanceID)
		vpcId = *vpc.VpcId
		primaryNodeSubnetId = *vpc.Subnets[3].SubnetId
		secondary0NodeSubnetId = *vpc.Subnets[4].SubnetId
		secondary1NodeSubnetId = *vpc.Subnets[5].SubnetId
		bastionSecurityGroupId = *vpc.SecurityGroups[0].GroupId
		mongoDBAdminPassword = "volunteer-pilot"
		keyPairName = vpc.KeyPairName

		mongoDBService, err = mongodb.NewService(region)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		helpers.DestroyVpc(vpc)
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongodb Suite")
}
