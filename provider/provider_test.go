package provider_test

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	awscf "github.com/aws/aws-sdk-go/service/cloudformation"
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
		fakeMongoDBService    *mongodb.Service
		awsProvider           *AWSProvider
	)

	BeforeEach(func() {
		rawConfig = []byte(`
			{
				"basic_auth_username": "username",
				"basic_auth_password": "password",
				"log_level": "info",
				"secret": "pocket-dialer",
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
							"cluster_replica_set_count": "1",
							"mongodb_version": "3.2",
							"mongodb_admin_username": "superadmin",
							"replica_shard_index": "1",
							"volume_size": "500",
							"volume_type": "io1",
							"iops": "300",
							"node_instance_type": "m4.xxlarge"
						},{
							"id": "uuid-2",
							"name": "replica-set-3",
							"description": "A replica set of 3 instances. Disk: 400GB gp2. Instance: m4.large",
							"metadata": {},
							"cluster_replica_set_count": "3",
							"mongodb_version": "3.2",
							"mongodb_admin_username": "superadmin",
							"replica_shard_index": "1",
							"volume_size": "500",
							"volume_type": "io1",
							"iops": "300",
							"node_instance_type": "m4.xxlarge"
						}]
					}]
				}
			}
		`)
		var err error
		config, err = DecodeConfig(rawConfig)
		Expect(err).NotTo(HaveOccurred())
		fakeCloudFormationAPI = &fakes.FakeCloudFormationAPI{}
		fakeMongoDBService = &mongodb.Service{Client: fakeCloudFormationAPI}
		awsProvider = &AWSProvider{Config: config, MongoDBService: fakeMongoDBService}
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

		Describe("Integration with the MongoDBService", func() {
			It("passes the correct parameters to AWS via the MongoDBService", func() {
				provisionData := usbProvider.ProvisionData{
					Service: brokerapi.Service{ID: "uuid-1"},
					Plan:    brokerapi.ServicePlan{ID: "uuid-2"},
				}
				fakeCloudFormationAPI.CreateStackReturns(
					&awscf.CreateStackOutput{StackId: aws.String("id")},
					nil,
				)
				_, _, err := awsProvider.Provision(context.Background(), provisionData)
				Expect(err).NotTo(HaveOccurred())

				expectedParameters := []*awscf.Parameter{
					{
						ParameterKey:     aws.String("MongoDBAdminPassword"),
						ParameterValue:   aws.String("CLTC1OdLulR4pjQhECf39A=="),
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(false),
					},
					{
						ParameterKey:     aws.String("BastionSecurityGroupID"),
						ParameterValue:   aws.String("sg-xxxxxx"),
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(false),
					},
					{
						ParameterKey:     aws.String("KeyPairName"),
						ParameterValue:   aws.String("key_pair_name"),
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(false),
					},
					{
						ParameterKey:     aws.String("VPC"),
						ParameterValue:   aws.String("vpc-xxxxxx"),
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(false),
					},
					{
						ParameterKey:     aws.String("PrimaryNodeSubnet"),
						ParameterValue:   aws.String("subnet-xxxxxx"),
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(false),
					},
					{
						ParameterKey:     aws.String("Secondary0NodeSubnet"),
						ParameterValue:   aws.String("subnet-xxxxxx"),
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(false),
					},
					{
						ParameterKey:     aws.String("Secondary1NodeSubnet"),
						ParameterValue:   aws.String("subnet-xxxxxx"),
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(false),
					},
					{
						ParameterKey:     aws.String("MongoDBVersion"),
						ParameterValue:   aws.String("3.2"),
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(false),
					},
					{
						ParameterKey:     aws.String("MongoDBAdminUsername"),
						ParameterValue:   aws.String("superadmin"),
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(false),
					},
					{
						ParameterKey:     aws.String("ClusterReplicaSetCount"),
						ParameterValue:   aws.String("1"),
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(false),
					},
					{
						ParameterKey:     aws.String("ReplicaShardIndex"),
						ParameterValue:   aws.String("1"),
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(false),
					},
					{
						ParameterKey:     aws.String("VolumeSize"),
						ParameterValue:   aws.String("500"),
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(false),
					},
					{
						ParameterKey:     aws.String("VolumeType"),
						ParameterValue:   aws.String("io1"),
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(false),
					},
					{
						ParameterKey:     aws.String("Iops"),
						ParameterValue:   aws.String("300"),
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(false),
					},
					{
						ParameterKey:     aws.String("NodeInstanceType"),
						ParameterValue:   aws.String("m4.xxlarge"),
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(false),
					},
				}
				createStackInput := fakeCloudFormationAPI.CreateStackArgsForCall(0)
				Expect(createStackInput.Parameters).To(Equal(expectedParameters))
			})

			It("returns an error if the AWS call fails", func() {
				provisionData := usbProvider.ProvisionData{
					Service: brokerapi.Service{ID: "uuid-1"},
					Plan:    brokerapi.ServicePlan{ID: "uuid-2"},
				}
				fakeCloudFormationAPI.CreateStackReturns(
					nil,
					errors.New("some-aws-api-error"),
				)
				_, _, err := awsProvider.Provision(context.Background(), provisionData)
				Expect(err).To(MatchError("some-aws-api-error"))
			})

			It("returns the correct values", func() {
				provisionData := usbProvider.ProvisionData{
					Service: brokerapi.Service{ID: "uuid-1"},
					Plan:    brokerapi.ServicePlan{ID: "uuid-2"},
				}
				fakeCloudFormationAPI.CreateStackReturns(
					&awscf.CreateStackOutput{StackId: aws.String("id")},
					nil,
				)
				dashboardURL, operationData, err := awsProvider.Provision(context.Background(), provisionData)
				Expect(err).NotTo(HaveOccurred())
				Expect(dashboardURL).To(BeEmpty())
				Expect(operationData).To(Equal(`{"type":"provision","stack_id":"id"}`))
			})
		})
	})
})
