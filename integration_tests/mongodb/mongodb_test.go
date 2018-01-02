package mongodb_test

import (
	"time"

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
			id,
			keyPairName,
			primaryNodeSubnetId,
			secondary0NodeSubnetId,
			secondary1NodeSubnetId,
			mongoDBAdminPassword,
			vpcId,
			bastionSecurityGroupId,
		)
		Expect(err).NotTo(HaveOccurred())

		By("Polling for creation completion")
		Eventually(
			func() bool {
				completed, err := mongoDBService.CreateStackCompleted(id)
				Expect(err).NotTo(HaveOccurred())
				return completed
			},
			DEFAULT_TIMEOUT,
			30*time.Second,
		).Should(BeTrue())

		By("Deleting the stack")
		err = mongoDBService.DeleteStack(id)
		Expect(err).NotTo(HaveOccurred())

		By("Polling for deletion completion")
		Eventually(
			func() bool {
				completed, err := mongoDBService.DeleteStackCompleted(id)
				Expect(err).NotTo(HaveOccurred())
				return completed
			},
			DEFAULT_TIMEOUT,
			30*time.Second,
		).Should(BeTrue())
	})
})
