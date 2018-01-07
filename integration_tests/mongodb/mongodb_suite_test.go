package mongodb_test

import (
	"os"
	"strings"

	"github.com/henrytk/aws-service-broker/aws/cloudformation/mongodb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"

	"testing"
)

var (
	region                 = "eu-west-1"
	mongoDBService         mongodb.Service
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
		mongoDBService, err = mongodb.NewService(region)
		Expect(err).NotTo(HaveOccurred())
		id = "test" + strings.Replace(uuid.NewV4().String(), "-", "", -1)
		assertEnvVar(&keyPairName, "ASB_KEY_PAIR")
		assertEnvVar(&primaryNodeSubnetId, "ASB_PRIMARY_NODE")
		assertEnvVar(&secondary0NodeSubnetId, "ASB_SECONDARY_0_NODE")
		assertEnvVar(&secondary1NodeSubnetId, "ASB_SECONDARY_1_NODE")
		assertEnvVar(&mongoDBAdminPassword, "ASB_MONGODB_ADMIN_PASSWORD")
		assertEnvVar(&vpcId, "ASB_VPC_ID")
		assertEnvVar(&bastionSecurityGroupId, "ASB_BASTION_SECURITY_GROUP")
	})
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongodb Suite")
}

func assertEnvVar(parameter *string, key string) {
	*parameter, ok = os.LookupEnv(key)
	Expect(ok).To(BeTrue(), "key "+key+" not set")
	Expect(*parameter).NotTo(BeEmpty(), "for key "+key)
}
