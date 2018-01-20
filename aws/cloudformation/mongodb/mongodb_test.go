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
		mongoDBService        Service
	)

	BeforeEach(func() {
		fakeCloudFormationAPI = &fakes.FakeCloudFormationAPI{}
		mongoDBService = Service{Client: fakeCloudFormationAPI}
	})

	Describe("BuildStackTemplateParameters", func() {
		Describe("Mandatory parameters", func() {
			It("returns an error if bastion security group ID is empty", func() {
				_, err := mongoDBService.BuildStackTemplateParameters(InputParameters{"password", "", "keyPairName", "vpcId", "primary",
					"secondary0", "secondary1", "", "", "", "", "", "", ""})
				Expect(err).To(MatchError("Error building MongoDB parameters: bastion security group ID is empty"))
			})

			It("returns an error if key pair name is empty", func() {
				_, err := mongoDBService.BuildStackTemplateParameters(InputParameters{"password", "bastion", "", "vpcId", "primary",
					"secondary0", "secondary1", "", "", "", "", "", "", ""})
				Expect(err).To(MatchError("Error building MongoDB parameters: key pair name is empty"))
			})

			It("returns an error if VPC ID is empty", func() {
				_, err := mongoDBService.BuildStackTemplateParameters(InputParameters{"password", "bastion", "keyPairName", "", "primary",
					"secondary0", "secondary1", "", "", "", "", "", "", ""})
				Expect(err).To(MatchError("Error building MongoDB parameters: VPC ID is empty"))
			})

			It("returns an error if primary node subnet ID is empty", func() {
				_, err := mongoDBService.BuildStackTemplateParameters(InputParameters{"password", "bastion", "keyPairName", "vpcId", "",
					"secondary0", "secondary1", "", "", "", "", "", "", ""})
				Expect(err).To(MatchError("Error building MongoDB parameters: primary node subnet ID is empty"))
			})

			It("returns an error if secondary 0 node subnet ID is empty", func() {
				_, err := mongoDBService.BuildStackTemplateParameters(InputParameters{"password", "bastion", "keyPairName", "vpcId", "primary",
					"", "secondary1", "", "", "", "", "", "", ""})
				Expect(err).To(MatchError("Error building MongoDB parameters: secondary 0 node subnet ID is empty"))
			})

			It("returns an error if secondary 1 node subnet ID is empty", func() {
				_, err := mongoDBService.BuildStackTemplateParameters(InputParameters{"password", "bastion", "keyPairName", "vpcId", "primary",
					"secondary0", "", "", "", "", "", "", "", ""})
				Expect(err).To(MatchError("Error building MongoDB parameters: secondary 1 node subnet ID is empty"))
			})
		})

		Describe("Parameters with default values", func() {
			It("Adds all six optional parameters if non-empty", func() {
				parameters, err := mongoDBService.BuildStackTemplateParameters(InputParameters{"password", "bastion", "keyPairName", "vpcId", "primary",
					"secondary0", "secondary1", "3.2", "1", "1", "500", "io1", "200", "m4.xlarge"})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(parameters)).To(Equal(14))
			})
		})
	})

	Describe("Stack Creation", func() {
		It("should build valid input", func() {
			var parameters []*awscf.Parameter
			createStackInput := mongoDBService.BuildCreateStackInput("some-unique-id", parameters)
			err := createStackInput.Validate()
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
	})
})
