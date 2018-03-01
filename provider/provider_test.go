package provider_test

import (
	"context"
	"encoding/json"
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
							"description": "No replicas. Disk: 400GB gp2. Instance: m3.large",
							"metadata": {},
							"cluster_replica_set_count": "1",
							"mongodb_version": "3.2",
							"mongodb_admin_username": "superadmin",
							"replica_shard_index": "1",
							"volume_size": "500",
							"volume_type": "io1",
							"iops": "300",
							"node_instance_type": "m3.large"
						},{
							"id": "uuid-3",
							"name": "enhanced",
							"description": "No replicas. Disk: 400GB gp2. Instance: m4.large",
							"metadata": {},
							"cluster_replica_set_count": "1",
							"mongodb_version": "3.2",
							"mongodb_admin_username": "superadmin",
							"replica_shard_index": "1",
							"volume_size": "500",
							"volume_type": "io1",
							"iops": "300",
							"node_instance_type": "m4.large"
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
						ParameterKey:     aws.String("MongoDBAdminPassword"),
						ParameterValue:   aws.String("08b4c2d4e74bba5478a634211027f7f4"),
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
						ParameterKey:     aws.String("MongoDBVersion"),
						ParameterValue:   aws.String("3.2"),
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
						ParameterValue:   aws.String("m3.large"),
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
				Expect(operationData).To(Equal(`{"type":"provision","service":"mongodb","stack_id":"id"}`))
			})
		})
	})

	Describe("Deprovision", func() {
		It("returns an error when it can't find the service", func() {
			deprovisionData := usbProvider.DeprovisionData{
				Service: brokerapi.Service{ID: "this-cannot-be-found"},
			}
			_, err := awsProvider.Deprovision(context.Background(), deprovisionData)
			Expect(err).To(MatchError("could not find service ID: this-cannot-be-found"))
		})

		Describe("Integration with the MongoDBService", func() {
			It("passes the correct parameters to AWS via the MongoDBService", func() {
				deprovisionData := usbProvider.DeprovisionData{
					InstanceID: "deleteme",
					Service:    brokerapi.Service{ID: "uuid-1"},
				}
				fakeCloudFormationAPI.DeleteStackReturns(
					&awscf.DeleteStackOutput{},
					nil,
				)
				_, err := awsProvider.Deprovision(context.Background(), deprovisionData)
				Expect(err).NotTo(HaveOccurred())

				expectedStackId := fakeMongoDBService.GenerateStackName(deprovisionData.InstanceID)
				deleteStackInput := fakeCloudFormationAPI.DeleteStackArgsForCall(0)
				Expect(deleteStackInput.StackName).To(Equal(aws.String(expectedStackId)))
			})

			It("returns an error if the AWS call fails", func() {
				deprovisionData := usbProvider.DeprovisionData{
					InstanceID: "deleteme",
					Service:    brokerapi.Service{ID: "uuid-1"},
				}
				fakeCloudFormationAPI.DeleteStackReturns(
					nil,
					errors.New("some-aws-api-error"),
				)
				_, err := awsProvider.Deprovision(context.Background(), deprovisionData)
				Expect(err).To(MatchError("some-aws-api-error"))
			})

			It("returns the correct values", func() {
				deprovisionData := usbProvider.DeprovisionData{
					InstanceID: "deleteme",
					Service:    brokerapi.Service{ID: "uuid-1"},
				}
				fakeCloudFormationAPI.DeleteStackReturns(
					&awscf.DeleteStackOutput{},
					nil,
				)
				operationData, err := awsProvider.Deprovision(context.Background(), deprovisionData)
				Expect(err).NotTo(HaveOccurred())
				Expect(operationData).To(Equal(`{"type":"deprovision","service":"mongodb","instance_id":"deleteme"}`))
			})
		})
	})

	Describe("Update", func() {
		It("errors if update parameters are sent", func() {
			updateData := usbProvider.UpdateData{
				Details: brokerapi.UpdateDetails{
					RawParameters: json.RawMessage(`{"field": "value"}`),
				},
			}
			_, err := awsProvider.Update(context.Background(), updateData)
			Expect(err).To(MatchError("update parameters are not supported"))
		})

		It("returns an error when it can't find the service", func() {
			updateData := usbProvider.UpdateData{
				Service: brokerapi.Service{ID: "this-cannot-be-found"},
			}
			_, err := awsProvider.Update(context.Background(), updateData)
			Expect(err).To(MatchError("could not find service ID: this-cannot-be-found"))
		})

		It("returns an error when it can't find the new plan", func() {
			updateData := usbProvider.UpdateData{
				Service: brokerapi.Service{ID: "uuid-1"},
				Plan:    brokerapi.ServicePlan{ID: "this-cannot-be-found"},
			}
			_, err := awsProvider.Update(context.Background(), updateData)
			Expect(err).To(MatchError("could not find plan ID: this-cannot-be-found"))
		})

		It("returns an error when it can't find the current plan", func() {
			updateData := usbProvider.UpdateData{
				Details: brokerapi.UpdateDetails{
					PreviousValues: brokerapi.PreviousValues{
						PlanID: "this-cannot-be-found",
					},
				},
				Service: brokerapi.Service{ID: "uuid-1"},
				Plan:    brokerapi.ServicePlan{ID: "uuid-3"},
			}
			_, err := awsProvider.Update(context.Background(), updateData)
			Expect(err).To(MatchError("could not find plan ID: this-cannot-be-found"))
		})

		Describe("Plans", func() {
			var (
				updateConfig *Config
			)

			BeforeEach(func() {
				updateConfig = &Config{
					Catalog: Catalog{
						Services: []Service{
							{
								Service: brokerapi.Service{ID: "uuid-1", Name: "mongodb"},
								Plans: []Plan{
									{
										ServicePlan:           brokerapi.ServicePlan{ID: "uuid-2"},
										MongoDBPlanParameters: MongoDBPlanParameters{},
									},
									{
										ServicePlan:           brokerapi.ServicePlan{ID: "uuid-3"},
										MongoDBPlanParameters: MongoDBPlanParameters{},
									},
								},
							},
						},
					},
				}
				fakeCloudFormationAPI.UpdateStackWithContextReturns(
					&awscf.UpdateStackOutput{StackId: aws.String("stack_id")},
					nil,
				)
			})

			It("allows you to update the node instance type", func() {
				updateConfig.Catalog.Services[0].Plans[0].MongoDBPlanParameters.NodeInstanceType = "m3.large"
				updateConfig.Catalog.Services[0].Plans[1].MongoDBPlanParameters.NodeInstanceType = "m4.large"
				awsProvider.Config = updateConfig
				updateData := usbProvider.UpdateData{
					Details: brokerapi.UpdateDetails{
						PreviousValues: brokerapi.PreviousValues{
							PlanID: "uuid-2",
						},
					},
					Service: brokerapi.Service{ID: "uuid-1"},
					Plan:    brokerapi.ServicePlan{ID: "uuid-3"},
				}
				_, err := awsProvider.Update(context.Background(), updateData)
				Expect(err).NotTo(HaveOccurred())
			})

			It("doesn't allow you to update the MongoDBAdminUsername", func() {
				updateConfig.Catalog.Services[0].Plans[0].MongoDBPlanParameters.MongoDBAdminUsername = "admin"
				updateConfig.Catalog.Services[0].Plans[1].MongoDBPlanParameters.MongoDBAdminUsername = "superadmin"
				awsProvider.Config = updateConfig
				updateData := usbProvider.UpdateData{
					Details: brokerapi.UpdateDetails{
						PreviousValues: brokerapi.PreviousValues{
							PlanID: "uuid-2",
						},
					},
					Service: brokerapi.Service{ID: "uuid-1"},
					Plan:    brokerapi.ServicePlan{ID: "uuid-3"},
				}
				_, err := awsProvider.Update(context.Background(), updateData)
				Expect(err).To(MatchError("updating MongoDB admin username is not supported"))
			})

			It("doesn't allow you to update the MongoDBVersion", func() {
				updateConfig.Catalog.Services[0].Plans[0].MongoDBPlanParameters.MongoDBVersion = "3.2"
				updateConfig.Catalog.Services[0].Plans[1].MongoDBPlanParameters.MongoDBVersion = "3.4"
				awsProvider.Config = updateConfig
				updateData := usbProvider.UpdateData{
					Details: brokerapi.UpdateDetails{
						PreviousValues: brokerapi.PreviousValues{
							PlanID: "uuid-2",
						},
					},
					Service: brokerapi.Service{ID: "uuid-1"},
					Plan:    brokerapi.ServicePlan{ID: "uuid-3"},
				}
				_, err := awsProvider.Update(context.Background(), updateData)
				Expect(err).To(MatchError("updating MongoDB version is not supported"))
			})

			It("doesn't allow you to update the ClusterReplicaSetCount", func() {
				updateConfig.Catalog.Services[0].Plans[0].MongoDBPlanParameters.ClusterReplicaSetCount = "1"
				updateConfig.Catalog.Services[0].Plans[1].MongoDBPlanParameters.ClusterReplicaSetCount = "3"
				awsProvider.Config = updateConfig
				updateData := usbProvider.UpdateData{
					Details: brokerapi.UpdateDetails{
						PreviousValues: brokerapi.PreviousValues{
							PlanID: "uuid-2",
						},
					},
					Service: brokerapi.Service{ID: "uuid-1"},
					Plan:    brokerapi.ServicePlan{ID: "uuid-3"},
				}
				_, err := awsProvider.Update(context.Background(), updateData)
				Expect(err).To(MatchError("updating cluster replica set count is not supported"))
			})

			It("doesn't allow you to update the ReplicaShardIndex", func() {
				updateConfig.Catalog.Services[0].Plans[0].MongoDBPlanParameters.ReplicaShardIndex = "0"
				updateConfig.Catalog.Services[0].Plans[1].MongoDBPlanParameters.ReplicaShardIndex = "1"
				awsProvider.Config = updateConfig
				updateData := usbProvider.UpdateData{
					Details: brokerapi.UpdateDetails{
						PreviousValues: brokerapi.PreviousValues{
							PlanID: "uuid-2",
						},
					},
					Service: brokerapi.Service{ID: "uuid-1"},
					Plan:    brokerapi.ServicePlan{ID: "uuid-3"},
				}
				_, err := awsProvider.Update(context.Background(), updateData)
				Expect(err).To(MatchError("updating replica shard index is not supported"))
			})

			It("doesn't allow you to update the VolumeSize", func() {
				updateConfig.Catalog.Services[0].Plans[0].MongoDBPlanParameters.VolumeSize = "400"
				updateConfig.Catalog.Services[0].Plans[1].MongoDBPlanParameters.VolumeSize = "500"
				awsProvider.Config = updateConfig
				updateData := usbProvider.UpdateData{
					Details: brokerapi.UpdateDetails{
						PreviousValues: brokerapi.PreviousValues{
							PlanID: "uuid-2",
						},
					},
					Service: brokerapi.Service{ID: "uuid-1"},
					Plan:    brokerapi.ServicePlan{ID: "uuid-3"},
				}
				_, err := awsProvider.Update(context.Background(), updateData)
				Expect(err).To(MatchError("updating volume size is not supported"))
			})

			It("doesn't allow you to update the VolumeType", func() {
				updateConfig.Catalog.Services[0].Plans[0].MongoDBPlanParameters.VolumeType = "gp2"
				updateConfig.Catalog.Services[0].Plans[1].MongoDBPlanParameters.VolumeType = "io1"
				awsProvider.Config = updateConfig
				updateData := usbProvider.UpdateData{
					Details: brokerapi.UpdateDetails{
						PreviousValues: brokerapi.PreviousValues{
							PlanID: "uuid-2",
						},
					},
					Service: brokerapi.Service{ID: "uuid-1"},
					Plan:    brokerapi.ServicePlan{ID: "uuid-3"},
				}
				_, err := awsProvider.Update(context.Background(), updateData)
				Expect(err).To(MatchError("updating volume type is not supported"))
			})

			It("doesn't allow you to update the Iops", func() {
				updateConfig.Catalog.Services[0].Plans[0].MongoDBPlanParameters.Iops = "100"
				updateConfig.Catalog.Services[0].Plans[1].MongoDBPlanParameters.Iops = "300"
				awsProvider.Config = updateConfig
				updateData := usbProvider.UpdateData{
					Details: brokerapi.UpdateDetails{
						PreviousValues: brokerapi.PreviousValues{
							PlanID: "uuid-2",
						},
					},
					Service: brokerapi.Service{ID: "uuid-1"},
					Plan:    brokerapi.ServicePlan{ID: "uuid-3"},
				}
				_, err := awsProvider.Update(context.Background(), updateData)
				Expect(err).To(MatchError("updating IOPS is not supported"))
			})
		})

		Describe("Integration with the MongoDBService", func() {
			It("passes the correct parameters to AWS via the MongoDBService", func() {
				updateData := usbProvider.UpdateData{
					Details: brokerapi.UpdateDetails{
						PreviousValues: brokerapi.PreviousValues{
							PlanID: "uuid-2",
						},
					},
					Service: brokerapi.Service{ID: "uuid-1"},
					Plan:    brokerapi.ServicePlan{ID: "uuid-3"},
				}
				fakeCloudFormationAPI.UpdateStackWithContextReturns(
					&awscf.UpdateStackOutput{StackId: aws.String("stack_id")},
					nil,
				)
				_, err := awsProvider.Update(context.Background(), updateData)
				Expect(err).NotTo(HaveOccurred())

				expectedParameters := []*awscf.Parameter{
					{
						ParameterKey:     aws.String("BastionSecurityGroupID"),
						ParameterValue:   nil,
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(true),
					},
					{
						ParameterKey:     aws.String("KeyPairName"),
						ParameterValue:   nil,
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(true),
					},
					{
						ParameterKey:     aws.String("VPC"),
						ParameterValue:   nil,
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(true),
					},
					{
						ParameterKey:     aws.String("PrimaryNodeSubnet"),
						ParameterValue:   nil,
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(true),
					},
					{
						ParameterKey:     aws.String("Secondary0NodeSubnet"),
						ParameterValue:   nil,
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(true),
					},
					{
						ParameterKey:     aws.String("Secondary1NodeSubnet"),
						ParameterValue:   nil,
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(true),
					},
					{
						ParameterKey:     aws.String("MongoDBAdminPassword"),
						ParameterValue:   nil,
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(true),
					},
					{
						ParameterKey:     aws.String("MongoDBAdminUsername"),
						ParameterValue:   nil,
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(true),
					},
					{
						ParameterKey:     aws.String("MongoDBVersion"),
						ParameterValue:   nil,
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(true),
					},
					{
						ParameterKey:     aws.String("ClusterReplicaSetCount"),
						ParameterValue:   nil,
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(true),
					},
					{
						ParameterKey:     aws.String("ReplicaShardIndex"),
						ParameterValue:   nil,
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(true),
					},
					{
						ParameterKey:     aws.String("VolumeSize"),
						ParameterValue:   nil,
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(true),
					},
					{
						ParameterKey:     aws.String("VolumeType"),
						ParameterValue:   nil,
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(true),
					},
					{
						ParameterKey:     aws.String("Iops"),
						ParameterValue:   nil,
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(true),
					},
					{
						ParameterKey:     aws.String("NodeInstanceType"),
						ParameterValue:   aws.String("m4.large"),
						ResolvedValue:    nil,
						UsePreviousValue: aws.Bool(false),
					},
				}
				Expect(fakeCloudFormationAPI.UpdateStackWithContextCallCount()).To(Equal(1))
				_, updateStackInput, _ := fakeCloudFormationAPI.UpdateStackWithContextArgsForCall(0)
				Expect(updateStackInput.Parameters).To(Equal(expectedParameters))
			})

			It("returns an error if the AWS call fails", func() {
				updateData := usbProvider.UpdateData{
					Details: brokerapi.UpdateDetails{
						PreviousValues: brokerapi.PreviousValues{
							PlanID: "uuid-2",
						},
					},
					Service: brokerapi.Service{ID: "uuid-1"},
					Plan:    brokerapi.ServicePlan{ID: "uuid-3"},
				}
				fakeCloudFormationAPI.UpdateStackWithContextReturns(
					nil,
					errors.New("some-aws-api-error"),
				)
				_, err := awsProvider.Update(context.Background(), updateData)
				Expect(err).To(MatchError("some-aws-api-error"))
			})

			It("returns the correct values", func() {
				updateData := usbProvider.UpdateData{
					Details: brokerapi.UpdateDetails{
						PreviousValues: brokerapi.PreviousValues{
							PlanID: "uuid-2",
						},
					},
					Service: brokerapi.Service{ID: "uuid-1"},
					Plan:    brokerapi.ServicePlan{ID: "uuid-3"},
				}
				fakeCloudFormationAPI.UpdateStackWithContextReturns(
					&awscf.UpdateStackOutput{StackId: aws.String("id")},
					nil,
				)
				operationData, err := awsProvider.Update(context.Background(), updateData)
				Expect(err).NotTo(HaveOccurred())
				Expect(operationData).To(Equal(`{"type":"update","service":"mongodb","stack_id":"id"}`))
			})
		})
	})

	Describe("LastOperation", func() {
		Describe("last operation data unmarshalling", func() {
			It("returns an error if the last operation type is unrecognised", func() {
				lastOperationData := usbProvider.LastOperationData{
					OperationData: `{"type": "restore", "service": "mongodb"}`,
				}
				_, _, err := awsProvider.LastOperation(context.Background(), lastOperationData)
				Expect(err).To(MatchError("unknown operation type 'restore'"))
			})

			It("returns an error if the last operation service isn't recognised", func() {
				lastOperationData := usbProvider.LastOperationData{
					OperationData: `{"type": "provision", "service": "BongoDB"}`,
				}
				_, _, err := awsProvider.LastOperation(context.Background(), lastOperationData)
				Expect(err).To(MatchError("unknown service 'BongoDB'"))
			})
		})

		Describe("Service `mongodb`", func() {
			Describe("provisioning", func() {
				It("makes the right calls and returns the right data when provision is complete", func() {
					lastOperationData := usbProvider.LastOperationData{
						InstanceID:    "id",
						OperationData: `{"type": "provision", "service": "mongodb", "stack_id": "id"}`,
					}
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus: aws.String(awscf.StackStatusCreateComplete),
								},
							},
						},
						nil,
					)
					state, description, err := awsProvider.LastOperation(context.Background(), lastOperationData)
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeCloudFormationAPI.DescribeStacksCallCount()).To(Equal(1))
					Expect(fakeCloudFormationAPI.DescribeStacksArgsForCall(0)).To(Equal(
						&awscf.DescribeStacksInput{
							StackName: aws.String(fakeMongoDBService.GenerateStackName("id")),
						},
					))
					Expect(state).To(Equal(brokerapi.Succeeded))
					Expect(description).To(Equal("provision succeeded"))
				})

				It("returns failure message when provision failed", func() {
					lastOperationData := usbProvider.LastOperationData{
						InstanceID:    "id",
						OperationData: `{"type": "provision", "service": "mongodb", "stack_id": "id"}`,
					}
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus: aws.String(awscf.StackStatusCreateFailed),
								},
							},
						},
						nil,
					)
					state, description, err := awsProvider.LastOperation(context.Background(), lastOperationData)
					Expect(err).NotTo(HaveOccurred())
					Expect(state).To(Equal(brokerapi.Failed))
					Expect(description).To(Equal("Final state of stack was not CREATE_COMPLETE. Got: CREATE_FAILED. Reason: no reason returned via the API"))
				})

				It("returns 'in progress' when provision failed", func() {
					lastOperationData := usbProvider.LastOperationData{
						InstanceID:    "id",
						OperationData: `{"type": "provision", "service": "mongodb", "stack_id": "id"}`,
					}
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus: aws.String(awscf.StackStatusCreateInProgress),
								},
							},
						},
						nil,
					)
					state, description, err := awsProvider.LastOperation(context.Background(), lastOperationData)
					Expect(err).NotTo(HaveOccurred())
					Expect(state).To(Equal(brokerapi.InProgress))
					Expect(description).To(Equal("provision in progress"))
				})
			})

			Describe("deprovisioning", func() {
				Context("the stack can no longer be retrieved from the AWS API", func() {
					It("makes the right calls and returns the right data when deprovision is complete", func() {
						lastOperationData := usbProvider.LastOperationData{
							InstanceID:    "id",
							OperationData: `{"type": "deprovision", "service": "mongodb", "stack_id": "id"}`,
						}
						fakeCloudFormationAPI.DescribeStacksReturns(
							&awscf.DescribeStacksOutput{},
							errors.New("Stack with id "+
								fakeMongoDBService.GenerateStackName(lastOperationData.InstanceID)+
								" does not exist",
							),
						)
						state, description, err := awsProvider.LastOperation(context.Background(), lastOperationData)
						Expect(err).NotTo(HaveOccurred())
						Expect(fakeCloudFormationAPI.DescribeStacksCallCount()).To(Equal(1))
						Expect(fakeCloudFormationAPI.DescribeStacksArgsForCall(0)).To(Equal(
							&awscf.DescribeStacksInput{
								StackName: aws.String(fakeMongoDBService.GenerateStackName("id")),
							},
						))
						Expect(state).To(Equal(brokerapi.Succeeded))
						Expect(description).To(Equal("deprovision succeeded"))
					})
				})
				Context("the AWS API returns an explicit completion message", func() {
					It("returns the right data", func() {
						lastOperationData := usbProvider.LastOperationData{
							InstanceID:    "id",
							OperationData: `{"type": "deprovision", "service": "mongodb", "stack_id": "id"}`,
						}
						fakeCloudFormationAPI.DescribeStacksReturns(
							&awscf.DescribeStacksOutput{
								Stacks: []*awscf.Stack{
									&awscf.Stack{
										StackStatus: aws.String(awscf.StackStatusDeleteComplete),
									},
								},
							},
							nil,
						)
						state, description, err := awsProvider.LastOperation(context.Background(), lastOperationData)
						Expect(err).NotTo(HaveOccurred())
						Expect(state).To(Equal(brokerapi.Succeeded))
						Expect(description).To(Equal("deprovision succeeded"))
					})
				})

				It("returns failure message when deprovision failed", func() {
					lastOperationData := usbProvider.LastOperationData{
						InstanceID:    "id",
						OperationData: `{"type": "deprovision", "service": "mongodb", "stack_id": "id"}`,
					}
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus: aws.String(awscf.StackStatusDeleteFailed),
								},
							},
						},
						nil,
					)
					state, description, err := awsProvider.LastOperation(context.Background(), lastOperationData)
					Expect(err).NotTo(HaveOccurred())
					Expect(state).To(Equal(brokerapi.Failed))
					Expect(description).To(Equal("Final state of stack was not DELETE_COMPLETE. Got: DELETE_FAILED. Reason: no reason returned via the API"))
				})

				It("returns 'in progress' when deprovision failed", func() {
					lastOperationData := usbProvider.LastOperationData{
						InstanceID:    "id",
						OperationData: `{"type": "deprovision", "service": "mongodb", "stack_id": "id"}`,
					}
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus: aws.String(awscf.StackStatusDeleteInProgress),
								},
							},
						},
						nil,
					)
					state, description, err := awsProvider.LastOperation(context.Background(), lastOperationData)
					Expect(err).NotTo(HaveOccurred())
					Expect(state).To(Equal(brokerapi.InProgress))
					Expect(description).To(Equal("deprovision in progress"))
				})
			})

			Describe("updating", func() {
				It("makes the right calls and returns the right data when update is complete", func() {
					lastOperationData := usbProvider.LastOperationData{
						InstanceID:    "id",
						OperationData: `{"type": "update", "service": "mongodb", "stack_id": "id"}`,
					}
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus: aws.String(awscf.StackStatusUpdateComplete),
								},
							},
						},
						nil,
					)
					state, description, err := awsProvider.LastOperation(context.Background(), lastOperationData)
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeCloudFormationAPI.DescribeStacksCallCount()).To(Equal(1))
					Expect(fakeCloudFormationAPI.DescribeStacksArgsForCall(0)).To(Equal(
						&awscf.DescribeStacksInput{
							StackName: aws.String(fakeMongoDBService.GenerateStackName("id")),
						},
					))
					Expect(state).To(Equal(brokerapi.Succeeded))
					Expect(description).To(Equal("update succeeded"))
				})

				It("returns failure message when update failed", func() {
					lastOperationData := usbProvider.LastOperationData{
						InstanceID:    "id",
						OperationData: `{"type": "update", "service": "mongodb", "stack_id": "id"}`,
					}
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus: aws.String(awscf.StackStatusUpdateRollbackComplete),
								},
							},
						},
						nil,
					)
					state, description, err := awsProvider.LastOperation(context.Background(), lastOperationData)
					Expect(err).NotTo(HaveOccurred())
					Expect(state).To(Equal(brokerapi.Failed))
					Expect(description).To(Equal("Final state of stack was not UPDATE_COMPLETE. Got: UPDATE_ROLLBACK_COMPLETE. Reason: no reason returned via the API"))
				})

				It("returns 'in progress' when update failed", func() {
					lastOperationData := usbProvider.LastOperationData{
						InstanceID:    "id",
						OperationData: `{"type": "update", "service": "mongodb", "stack_id": "id"}`,
					}
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus: aws.String(awscf.StackStatusUpdateInProgress),
								},
							},
						},
						nil,
					)
					state, description, err := awsProvider.LastOperation(context.Background(), lastOperationData)
					Expect(err).NotTo(HaveOccurred())
					Expect(state).To(Equal(brokerapi.InProgress))
					Expect(description).To(Equal("update in progress"))
				})
			})
		})
	})
})
