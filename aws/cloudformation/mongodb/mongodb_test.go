package mongodb_test

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	awscf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/henrytk/aws-service-broker/aws/cloudformation/fakes"
	. "github.com/henrytk/aws-service-broker/aws/cloudformation/mongodb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mongodb", func() {
	var (
		fakeCloudFormationAPI *fakes.FakeCloudFormationAPI
		mongoDBService        *Service
		inputParameters       InputParameters
	)

	BeforeEach(func() {
		fakeCloudFormationAPI = &fakes.FakeCloudFormationAPI{}
		mongoDBService = &Service{Client: fakeCloudFormationAPI}
		inputParameters = InputParameters{
			BastionSecurityGroupId: "bastion",
			KeyPairName:            "keypairname",
			VpcId:                  "vpc-id",
			PrimaryNodeSubnetId:    "primary",
			Secondary0NodeSubnetId: "secondary0",
			Secondary1NodeSubnetId: "secondary1",
			MongoDBAdminPassword:   "password",
			MongoDBAdminUsername:   "admin",
			MongoDBVersion:         "3.4",
			ClusterReplicaSetCount: "1",
			ReplicaShardIndex:      "0",
			VolumeSize:             "400",
			VolumeType:             "io1",
			Iops:                   "200",
			NodeInstanceType:       "m4.xlarge",
		}
	})

	Describe("BuildUpdateStackParameters", func() {
		It("creates a parameter for each input field", func() {
			parameters := mongoDBService.BuildUpdateStackParameters(inputParameters)
			Expect(len(parameters)).To(Equal(15))

			By("overriding their previous value")
			for _, v := range parameters {
				Expect(v.ParameterValue).NotTo(BeNil())
				Expect(*v.UsePreviousValue).To(BeFalse())
			}
		})

		It("uses the previous value when a parameter isn't provided", func() {
			inputParameters.BastionSecurityGroupId = ""
			parameters := mongoDBService.BuildUpdateStackParameters(inputParameters)
			Expect(len(parameters)).To(Equal(15))
			Expect(*parameters[0].ParameterKey).To(Equal("BastionSecurityGroupID"))
			Expect(parameters[0].ParameterValue).To(BeNil())
			Expect(*parameters[0].UsePreviousValue).To(BeTrue())
		})
	})

	Describe("BuildCreateStackParameters", func() {
		Describe("Mandatory parameters", func() {
			It("returns an error if bastion security group ID is empty", func() {
				inputParameters.BastionSecurityGroupId = ""
				_, err := mongoDBService.BuildCreateStackParameters(inputParameters)
				Expect(err).To(MatchError("Error building MongoDB parameters: bastion security group ID is empty"))
			})

			It("returns an error if key pair name is empty", func() {
				inputParameters.KeyPairName = ""
				_, err := mongoDBService.BuildCreateStackParameters(inputParameters)
				Expect(err).To(MatchError("Error building MongoDB parameters: key pair name is empty"))
			})

			It("returns an error if VPC ID is empty", func() {
				inputParameters.VpcId = ""
				_, err := mongoDBService.BuildCreateStackParameters(inputParameters)
				Expect(err).To(MatchError("Error building MongoDB parameters: VPC ID is empty"))
			})

			It("returns an error if primary node subnet ID is empty", func() {
				inputParameters.PrimaryNodeSubnetId = ""
				_, err := mongoDBService.BuildCreateStackParameters(inputParameters)
				Expect(err).To(MatchError("Error building MongoDB parameters: primary node subnet ID is empty"))
			})

			It("returns an error if secondary 0 node subnet ID is empty", func() {
				inputParameters.Secondary0NodeSubnetId = ""
				_, err := mongoDBService.BuildCreateStackParameters(inputParameters)
				Expect(err).To(MatchError("Error building MongoDB parameters: secondary 0 node subnet ID is empty"))
			})

			It("returns an error if secondary 1 node subnet ID is empty", func() {
				inputParameters.Secondary1NodeSubnetId = ""
				_, err := mongoDBService.BuildCreateStackParameters(inputParameters)
				Expect(err).To(MatchError("Error building MongoDB parameters: secondary 1 node subnet ID is empty"))
			})

			It("returns an error if MongoDB admin password is empty", func() {
				inputParameters.MongoDBAdminPassword = ""
				_, err := mongoDBService.BuildCreateStackParameters(inputParameters)
				Expect(err).To(MatchError("Error building MongoDB parameters: MongoDB admin password is empty"))
			})
		})

		Describe("Parameters with default values", func() {
			It("Adds all six optional parameters if non-empty", func() {
				parameters, err := mongoDBService.BuildCreateStackParameters(inputParameters)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(parameters)).To(Equal(15))
			})
		})
	})

	Describe("BuildCreateStackInput", func() {
		It("should build valid input", func() {
			var parameters []*awscf.Parameter
			createStackInput := mongoDBService.BuildCreateStackInput("some-unique-id", parameters)
			err := createStackInput.Validate()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("BuildUpdateStackInput", func() {
		It("should build valid input", func() {
			var parameters []*awscf.Parameter
			updateStackInput := mongoDBService.BuildUpdateStackInput("some-unique-id", parameters)
			err := updateStackInput.Validate()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Getting stack information", func() {
		Describe("GetStackState", func() {
			Context("when stack has been created successfully", func() {
				It("returns the state with no error", func() {
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus: aws.String(awscf.StackStatusCreateComplete),
								},
							},
						}, nil,
					)
					state, reason, err := mongoDBService.GetStackState("irrelevant")
					Expect(err).NotTo(HaveOccurred())
					Expect(state).To(Equal(awscf.StackStatusCreateComplete))
					Expect(reason).To(Equal("no reason returned via the API"))
				})
			})

			Context("when there is an error getting the stack information", func() {
				It("returns no state information and an error", func() {
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{},
						}, errors.New("Error calling DescribeStacks"),
					)
					state, reason, err := mongoDBService.GetStackState("irrelevant")
					Expect(err).To(MatchError("Error calling DescribeStacks"))
					Expect(state).To(BeEmpty())
					Expect(reason).To(BeEmpty())
				})
			})

			Context("when multiple stacks are returned", func() {
				It("returns no state information and an error", func() {
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus: aws.String(awscf.StackStatusCreateComplete),
								},
								&awscf.Stack{
									StackStatus: aws.String(awscf.StackStatusCreateComplete),
								},
							},
						}, nil,
					)
					state, reason, err := mongoDBService.GetStackState("irrelevant")
					Expect(err).To(MatchError("Error checking stack state: number of stacks was not 1"))
					Expect(state).To(BeEmpty())
					Expect(reason).To(BeEmpty())
				})
			})

			Context("when stack has failed to create", func() {
				It("returns the state with a reason and no error", func() {
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus:       aws.String(awscf.StackStatusCreateFailed),
									StackStatusReason: aws.String("some reason for failure"),
								},
							},
						}, nil,
					)
					state, reason, err := mongoDBService.GetStackState("irrelevant")
					Expect(err).NotTo(HaveOccurred())
					Expect(state).To(Equal(awscf.StackStatusCreateFailed))
					Expect(reason).To(Equal("some reason for failure"))
				})
			})
		})

		Describe("CreateStackCompleted", func() {
			Context("when failing to get stack information", func() {
				It("returns false with an error", func() {
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{},
						}, errors.New("Error calling DescribeStacks"),
					)
					completed, err := mongoDBService.CreateStackCompleted("irrelevant")
					Expect(err).To(MatchError("Error calling DescribeStacks"))
					Expect(completed).To(BeFalse())
				})
			})

			Context("when stack has been created successfully", func() {
				It("returns true with no error", func() {
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus: aws.String(awscf.StackStatusCreateComplete),
								},
							},
						}, nil,
					)
					completed, err := mongoDBService.CreateStackCompleted("irrelevant")
					Expect(err).NotTo(HaveOccurred())
					Expect(completed).To(BeTrue())
				})
			})

			Context("when stack creation fails", func() {
				It("returns false and an error", func() {
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus:       aws.String(awscf.StackStatusCreateFailed),
									StackStatusReason: aws.String("something went wrong"),
								},
							},
						}, nil,
					)
					completed, err := mongoDBService.CreateStackCompleted("irrelevant")
					Expect(err).To(MatchError("Final state of stack was not CREATE_COMPLETE. Got: CREATE_FAILED. Reason: something went wrong"))
					Expect(completed).To(BeTrue())
				})
			})

			Context("when stack creation is still in progress", func() {
				It("returns false and no error", func() {
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus: aws.String(awscf.StackStatusCreateInProgress),
								},
							},
						}, nil,
					)
					completed, err := mongoDBService.CreateStackCompleted("irrelevant")
					Expect(err).NotTo(HaveOccurred())
					Expect(completed).To(BeFalse())
				})
			})
		})

		Describe("DeleteStackCompleted", func() {
			Context("when failing to get stack information", func() {
				Context("if it is due to the stack not existing", func() {
					It("assumes the deletion is complete", func() {
						fakeCloudFormationAPI.DescribeStacksReturns(
							&awscf.DescribeStacksOutput{
								Stacks: []*awscf.Stack{},
							}, errors.New("ValidationError: Stack with id "+mongoDBService.GenerateStackName("irrelevant")+" does not exist"),
						)
						completed, err := mongoDBService.DeleteStackCompleted("irrelevant")
						Expect(err).NotTo(HaveOccurred())
						Expect(completed).To(BeTrue())
					})
				})

				Context("if it is due to some other error", func() {
					It("doesn't consider it complete and returns the error", func() {
						fakeCloudFormationAPI.DescribeStacksReturns(
							&awscf.DescribeStacksOutput{
								Stacks: []*awscf.Stack{},
							}, errors.New("Error calling DescribeStacks"),
						)
						completed, err := mongoDBService.DeleteStackCompleted("irrelevant")
						Expect(err).To(MatchError("Error calling DescribeStacks"))
						Expect(completed).To(BeFalse())
					})
				})
			})

			Context("when stack has been deleted successfully", func() {
				It("returns true with no error", func() {
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus: aws.String(awscf.StackStatusDeleteComplete),
								},
							},
						}, nil,
					)
					completed, err := mongoDBService.DeleteStackCompleted("irrelevant")
					Expect(err).NotTo(HaveOccurred())
					Expect(completed).To(BeTrue())
				})
			})

			Context("when stack deletion fails", func() {
				It("returns false and an error", func() {
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus:       aws.String(awscf.StackStatusDeleteFailed),
									StackStatusReason: aws.String("something went wrong"),
								},
							},
						}, nil,
					)
					completed, err := mongoDBService.DeleteStackCompleted("irrelevant")
					Expect(err).To(MatchError("Final state of stack was not DELETE_COMPLETE. Got: DELETE_FAILED. Reason: something went wrong"))
					Expect(completed).To(BeTrue())
				})
			})

			Context("when stack deletion is still in progress", func() {
				It("returns false and no error", func() {
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus: aws.String(awscf.StackStatusDeleteInProgress),
								},
							},
						}, nil,
					)
					completed, err := mongoDBService.DeleteStackCompleted("irrelevant")
					Expect(err).NotTo(HaveOccurred())
					Expect(completed).To(BeFalse())
				})
			})
		})

		Describe("UpdateStackCompleted", func() {
			Context("when failing to get stack information", func() {
				It("returns false with an error", func() {
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{},
						}, errors.New("Error calling DescribeStacks"),
					)
					completed, err := mongoDBService.UpdateStackCompleted("irrelevant")
					Expect(err).To(MatchError("Error calling DescribeStacks"))
					Expect(completed).To(BeFalse())
				})
			})

			Context("when stack has been updated successfully", func() {
				It("returns true with no error", func() {
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus: aws.String(awscf.StackStatusUpdateComplete),
								},
							},
						}, nil,
					)
					completed, err := mongoDBService.UpdateStackCompleted("irrelevant")
					Expect(err).NotTo(HaveOccurred())
					Expect(completed).To(BeTrue())
				})
			})

			Context("when update stack fails", func() {
				It("returns false and an error", func() {
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus:       aws.String(awscf.StackStatusUpdateRollbackComplete),
									StackStatusReason: aws.String("something went wrong"),
								},
							},
						}, nil,
					)
					completed, err := mongoDBService.UpdateStackCompleted("irrelevant")
					Expect(err).To(MatchError("Final state of stack was not UPDATE_COMPLETE. Got: UPDATE_ROLLBACK_COMPLETE. Reason: something went wrong"))
					Expect(completed).To(BeTrue())
				})
			})

			Context("when update stack is still in progress", func() {
				It("returns false and no error", func() {
					fakeCloudFormationAPI.DescribeStacksReturns(
						&awscf.DescribeStacksOutput{
							Stacks: []*awscf.Stack{
								&awscf.Stack{
									StackStatus: aws.String(awscf.StackStatusUpdateInProgress),
								},
							},
						}, nil,
					)
					completed, err := mongoDBService.UpdateStackCompleted("irrelevant")
					Expect(err).NotTo(HaveOccurred())
					Expect(completed).To(BeFalse())
				})
			})
		})
	})
})
