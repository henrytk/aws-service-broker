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
		It("returns an error", func() {
			rawConfig = json.RawMessage(`{}`)
			_, err := DecodeConfig(rawConfig)
			Expect(err).To(MatchError("Config error: no catalog found"))
		})
	})

	Context("when there are no services configured", func() {
		It("returns an error", func() {
			rawConfig = json.RawMessage(`
				{
					"catalog": {
						"services": []
					}
				}
			`)
			_, err := DecodeConfig(rawConfig)
			Expect(err).To(MatchError("Config error: at least one service must be configured"))
		})
	})

	Context("when the service name is not recognised", func() {
		It("returns an error", func() {
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
			_, err := DecodeConfig(rawConfig)
			Expect(err).To(MatchError("Config error: service name mangoDB not recognised"))
		})
	})

	Context("when a service has no plans", func() {
		It("returns an error", func() {
			rawConfig = json.RawMessage(`
				{
					"catalog": {
						"services": [
							{
								"name": "mongodb",
								"bastion_security_group_id": "irrelevant",
								"key_pair_name": "key_pair_name",
								"vpc_id": "irrelevant",
								"primary_node_subnet_id": "irrelevant",
								"secondary_0_node_subnet_id": "irrelevant",
								"secondary_1_node_subnet_id": "irrelevant",
								"plans": []
							}
						]
					}
				}
			`)
			_, err := DecodeConfig(rawConfig)
			Expect(err).To(MatchError("Config error: at least one plan must be configured for service mongodb"))
		})
	})

	Context("when given valid config", func() {
		It("decodes both catalog data and provider-specific data into one structure", func() {
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

	Describe("Mandatory parameters", func() {
		It("returns an error if bastion security group ID is empty", func() {
			rawConfig = json.RawMessage(`
				{
					"catalog": {
						"services": [
							{
								"name": "mongodb",
								"description": "MongoDB clusters",
								"bastion_security_group_id": "",
								"key_pair_name": "non-empty",
								"vpc_id": "non-empty",
								"primary_node_subnet_id": "non-empty",
								"secondary_0_node_subnet_id": "non-empty",
								"secondary_1_node_subnet_id": "non-empty",
								"plans": []
							}
						]
					}
				}
			`)
			_, err := DecodeConfig(rawConfig)
			Expect(err).To(MatchError("Config error: must provide bastion security group ID"))
		})

		It("returns an error if key pair name is empty", func() {
			rawConfig = json.RawMessage(`
				{
					"catalog": {
						"services": [
							{
								"name": "mongodb",
								"description": "MongoDB clusters",
								"bastion_security_group_id": "non-empty",
								"key_pair_name": "",
								"vpc_id": "non-empty",
								"primary_node_subnet_id": "non-empty",
								"secondary_0_node_subnet_id": "non-empty",
								"secondary_1_node_subnet_id": "non-empty",
								"plans": []
							}
						]
					}
				}
			`)
			_, err := DecodeConfig(rawConfig)
			Expect(err).To(MatchError("Config error: must provide key pair name"))
		})

		It("returns an error if VPC ID is empty", func() {
			rawConfig = json.RawMessage(`
				{
					"catalog": {
						"services": [
							{
								"name": "mongodb",
								"description": "MongoDB clusters",
								"bastion_security_group_id": "non-empty",
								"key_pair_name": "non-empty",
								"vpc_id": "",
								"primary_node_subnet_id": "non-empty",
								"secondary_0_node_subnet_id": "non-empty",
								"secondary_1_node_subnet_id": "non-empty",
								"plans": []
							}
						]
					}
				}
			`)
			_, err := DecodeConfig(rawConfig)
			Expect(err).To(MatchError("Config error: must provide VPC ID"))
		})

		It("returns an error if primary node subnet ID is empty", func() {
			rawConfig = json.RawMessage(`
				{
					"catalog": {
						"services": [
							{
								"name": "mongodb",
								"description": "MongoDB clusters",
								"bastion_security_group_id": "non-empty",
								"key_pair_name": "non-empty",
								"vpc_id": "non-empty",
								"primary_node_subnet_id": "",
								"secondary_0_node_subnet_id": "non-empty",
								"secondary_1_node_subnet_id": "non-empty",
								"plans": []
							}
						]
					}
				}
			`)
			_, err := DecodeConfig(rawConfig)
			Expect(err).To(MatchError("Config error: must provide primary node subnet ID"))
		})

		It("returns an error if secondary node 0 subnet ID is empty", func() {
			rawConfig = json.RawMessage(`
				{
					"catalog": {
						"services": [
							{
								"name": "mongodb",
								"description": "MongoDB clusters",
								"bastion_security_group_id": "non-empty",
								"key_pair_name": "non-empty",
								"vpc_id": "non-empty",
								"primary_node_subnet_id": "non-empty",
								"secondary_0_node_subnet_id": "",
								"secondary_1_node_subnet_id": "non-empty",
								"plans": []
							}
						]
					}
				}
			`)
			_, err := DecodeConfig(rawConfig)
			Expect(err).To(MatchError("Config error: must provide secondary 0 node subnet ID"))
		})

		It("returns an error if secondary node 1 subnet ID is empty", func() {
			rawConfig = json.RawMessage(`
				{
					"catalog": {
						"services": [
							{
								"name": "mongodb",
								"description": "MongoDB clusters",
								"bastion_security_group_id": "non-empty",
								"key_pair_name": "non-empty",
								"vpc_id": "non-empty",
								"primary_node_subnet_id": "non-empty",
								"secondary_0_node_subnet_id": "non-empty",
								"secondary_1_node_subnet_id": "",
								"plans": []
							}
						]
					}
				}
			`)
			_, err := DecodeConfig(rawConfig)
			Expect(err).To(MatchError("Config error: must provide secondary 1 node subnet ID"))
		})
	})
})
