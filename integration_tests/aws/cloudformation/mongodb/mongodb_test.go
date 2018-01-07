package mongodb_test

import (
	"time"

	"github.com/henrytk/aws-service-broker/aws/cloudformation/mongodb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
)

var (
	DEFAULT_TIMEOUT = 15 * time.Minute
)

var _ = Describe("Mongodb", func() {
	var (
		instanceID string
	)

	It("Manages the lifecycle of a CloudFormation stack", func() {
		instanceID = uuid.NewV4().String()
		By("Creating a stack")
		_, err := mongoDBService.CreateStack(
			instanceID,
			mongodb.InputParameters{
				KeyPairName:            keyPairName,
				PrimaryNodeSubnetId:    primaryNodeSubnetId,
				Secondary0NodeSubnetId: secondary0NodeSubnetId,
				Secondary1NodeSubnetId: secondary1NodeSubnetId,
				MongoDBAdminPassword:   mongoDBAdminPassword,
				VpcId:                  vpcId,
				BastionSecurityGroupId: bastionSecurityGroupId,
			},
		)
		Expect(err).NotTo(HaveOccurred())

		By("Polling for creation completion")
		Eventually(
			func() bool {
				completed, err := mongoDBService.CreateStackCompleted(instanceID)
				Expect(err).NotTo(HaveOccurred())
				return completed
			},
			DEFAULT_TIMEOUT,
			30*time.Second,
		).Should(BeTrue())

		By("Deleting the stack")
		err = mongoDBService.DeleteStack(instanceID)
		Expect(err).NotTo(HaveOccurred())

		By("Polling for deletion completion")
		Eventually(
			func() bool {
				completed, err := mongoDBService.DeleteStackCompleted(instanceID)
				Expect(err).NotTo(HaveOccurred())
				return completed
			},
			DEFAULT_TIMEOUT,
			30*time.Second,
		).Should(BeTrue())
	})
})
