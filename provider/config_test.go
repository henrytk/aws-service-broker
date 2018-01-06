package provider_test

import (
	"encoding/json"

	. "github.com/henrytk/aws-service-broker/provider"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var (
		rawConfig json.RawMessage
	)

	Context("when there is no Catalog defined", func() {
		BeforeEach(func() {
			rawConfig = json.RawMessage(`{}`)
		})

		It("returns an error", func() {
			_, err := DecodeConfig(rawConfig)
			Expect(err).To(MatchError("Error decoding config: no catalog found"))
		})
	})

	Context("when there are no services configured", func() {
		BeforeEach(func() {
			rawConfig = json.RawMessage(`
				{
					"catalog": {
						"services": []
					}
				}
			`)
		})

		It("returns an error", func() {
			_, err := DecodeConfig(rawConfig)
			Expect(err).To(MatchError("Error decoding config: at least one service must be configured"))
		})
	})

	Context("when the service name is not recognised", func() {
		BeforeEach(func() {
			rawConfig = json.RawMessage(`
				{
					"catalog": {
						"services": [
							{
								"name": "mangoDB"
							}
						]
					}
				}
			`)
		})

		It("returns an error", func() {
			_, err := DecodeConfig(rawConfig)
			Expect(err).To(MatchError("Error decoding config: service name mangoDB not recognised"))
		})
	})

	Context("when a service has no plans", func() {
		BeforeEach(func() {
			rawConfig = json.RawMessage(`
				{
					"catalog": {
						"services": [
							{
								"name": "mongodb",
								"plans": []
							}
						]
					}
				}
			`)
		})

		It("returns an error", func() {
			_, err := DecodeConfig(rawConfig)
			Expect(err).To(MatchError("Error decoding config: at least one plan must be configured for service mongodb"))
		})
	})

	Context("when given valid config", func() {
		BeforeEach(func() {
			rawConfig = json.RawMessage(`
				{
					"catalog": {
						"services": [
							{
								"name": "mongodb",
								"description": "MongoDB clusters",
								"bastion_security_group_id": "irrelevant",
								"key_pair_name": "key_pair_name",
								"vpc_id": "irrelevant",
								"primary_node_subnet_id": "irrelevant",
								"secondary_0_node_subnet_id": "irrelevant",
								"secondary_1_node_subnet_id": "irrelevant",
								"plans": [
									{
										"id": "1",
										"description": "No replicas",
										"cluster_replica_set_count": "1"
									},
									{
										"id": "2",
										"description": "Replica set of 3",
										"cluster_replica_set_count": "3"
									}
								]
							}
						]
					}
				}
			`)
		})

		It("decodes both catalog data and provider-specific data into one structure", func() {
			config, err := DecodeConfig(rawConfig)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(config.Catalog.Services)).To(Equal(1))
			service := config.Catalog.Services[0]
			Expect(service.Description).To(Equal("MongoDB clusters"))
			Expect(service.KeyPairName).To(Equal("key_pair_name"))
			Expect(len(service.Plans)).To(Equal(2))
			plan1 := service.Plans[0]
			Expect(plan1.Description).To(Equal("No replicas"))
			Expect(plan1.ClusterReplicaSetCount).To(Equal("1"))
			plan2 := service.Plans[1]
			Expect(plan2.Description).To(Equal("Replica set of 3"))
			Expect(plan2.ClusterReplicaSetCount).To(Equal("3"))
		})
	})
})
