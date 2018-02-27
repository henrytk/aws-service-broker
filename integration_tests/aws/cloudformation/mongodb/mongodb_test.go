package mongodb_test

import (
	"context"
	"time"

	"github.com/henrytk/aws-service-broker/aws/cloudformation/mongodb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	DEFAULT_TIMEOUT = 15 * time.Minute
)

var _ = Describe("Mongodb", func() {

	It("Manages the lifecycle of a CloudFormation stack", func() {
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

		By("Updating the node instance type")
		_, err = mongoDBService.UpdateStack(
			context.Background(),
			instanceID,
			mongodb.InputParameters{NodeInstanceType: "m3.large"},
		)
		Expect(err).NotTo(HaveOccurred())

		Eventually(
			func() bool {
				completed, err := mongoDBService.UpdateStackCompleted(instanceID)
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
