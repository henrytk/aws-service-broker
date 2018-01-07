package provider_test

import (
	"context"

	"github.com/henrytk/aws-service-broker/aws/cloudformation/fakes"
	"github.com/henrytk/aws-service-broker/aws/cloudformation/mongodb"
	. "github.com/henrytk/aws-service-broker/provider"
	usbProvider "github.com/henrytk/universal-service-broker/provider"
	"github.com/pivotal-cf/brokerapi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provider", func() {
	var (
		rawConfig             []byte
		config                *Config
		fakeCloudFormationAPI *fakes.FakeCloudFormationAPI
		fakeMongoDBService    mongodb.Service
		awsProvider           AWSProvider
	)

	BeforeEach(func() {
		rawConfig = []byte(`
			{
				"basic_auth_username": "username",
				"basic_auth_password": "password",
				"log_level": "info",
				"aws_config": {
					"region": "eu-west-1"
				},
				"catalog": {
					"services": [{
						"id": "uuid-1",
						"name": "mongodb",
						"description": "MongoDB clusters via AWS CloudFormation",
						"bindable": true,
						"requires": [],
						"metadata": {},
						"bastion_security_group_id": "sg-xxxxxx",
						"key_pair_name": "key_pair_name",
						"vpc_id": "vpc-xxxxxx",
						"primary_node_subnet_id": "subnet-xxxxxx",
						"secondary_0_node_subnet_id": "subnet-xxxxxx",
						"secondary_1_node_subnet_id": "subnet-xxxxxx",
						"plans": [{
							"id": "uuid-2",
							"name": "basic",
							"description": "No replicas. Disk: 400GB gp2. Instance: m4.large",
							"metadata": {},
							"cluster_replica_set_count": "1"
						},{
							"id": "uuid-2",
							"name": "replica-set-3",
							"description": "A replica set of 3 instances. Disk: 400GB gp2. Instance: m4.large",
							"metadata": {},
							"cluster_replica_set_count": "3"
						}]
					}]
				}
			}
		`)
		var err error
		config, err = DecodeConfig(rawConfig)
		Expect(err).NotTo(HaveOccurred())
		fakeCloudFormationAPI = &fakes.FakeCloudFormationAPI{}
		fakeMongoDBService = mongodb.Service{Client: fakeCloudFormationAPI}
		awsProvider = AWSProvider{Config: config, MongoDBService: fakeMongoDBService}
	})

	Describe("Provision", func() {
		It("returns an error when it can't find the service", func() {
			provisionData := usbProvider.ProvisionData{
				Service: brokerapi.Service{ID: "this-cannot-be-found"},
			}
			_, _, err := awsProvider.Provision(context.Background(), provisionData)
			Expect(err).To(MatchError("could not find service ID: this-cannot-be-found"))
		})
		It("returns an error when it can't find the plan", func() {
			provisionData := usbProvider.ProvisionData{
				Service: brokerapi.Service{ID: "uuid-1"},
				Plan:    brokerapi.ServicePlan{ID: "this-cannot-be-found"},
			}
			_, _, err := awsProvider.Provision(context.Background(), provisionData)
			Expect(err).To(MatchError("could not find plan ID: this-cannot-be-found"))
		})
	})
})
