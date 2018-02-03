package mongodb_test

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/henrytk/aws-service-broker/aws/cloudformation/mongodb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	region                 = "eu-west-1"
	awsSession             *session.Session
	ec2Service             *ec2.EC2
	mongoDBService         *mongodb.Service
	err                    error
	id                     string
	ok                     bool
	keyPairName            string
	primaryNodeSubnetId    string
	secondary0NodeSubnetId string
	secondary1NodeSubnetId string
	mongoDBAdminPassword   string
	vpcId                  string
	bastionSecurityGroupId string
)

func TestMongodb(t *testing.T) {
	BeforeSuite(func() {
		CreateEC2Service()
		CreateVPC()
		primaryNodeSubnetId = CreateSubnet("10.0.0.0/24")
		secondary0NodeSubnetId = CreateSubnet("10.0.1.0/24")
		secondary1NodeSubnetId = CreateSubnet("10.0.2.0/24")
		bastionSecurityGroupId = CreateSecurityGroup("bastion-mongodb-test")
		mongoDBAdminPassword = "volunteer pilot"

		mongoDBService, err = mongodb.NewService(region)
		Expect(err).NotTo(HaveOccurred())
		assertEnvVar(&keyPairName, "ASB_KEY_PAIR")
	})

	AfterSuite(func() {
		DeleteSecurityGroup(bastionSecurityGroupId)
		DeleteSubnet(secondary0NodeSubnetId)
		DeleteSubnet(secondary1NodeSubnetId)
		DeleteSubnet(primaryNodeSubnetId)
		DeleteVPC()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongodb Suite")
}

func CreateEC2Service() {
	awsSession, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	Expect(err).NotTo(HaveOccurred())
	ec2Service = ec2.New(awsSession)
}

func CreateVPC() {
	vpcOutput, err := ec2Service.CreateVpc(&ec2.CreateVpcInput{
		CidrBlock: aws.String("10.0.0.0/16"),
	})
	Expect(err).NotTo(HaveOccurred())
	vpcId = *vpcOutput.Vpc.VpcId
}

func DeleteVPC() {
	_, err := ec2Service.DeleteVpc(&ec2.DeleteVpcInput{
		VpcId: aws.String(vpcId),
	})
	Expect(err).NotTo(HaveOccurred())
}

func CreateSubnet(cidrBlock string) string {
	createSubnetOutput, err := ec2Service.CreateSubnet(&ec2.CreateSubnetInput{
		CidrBlock: aws.String(cidrBlock),
		VpcId:     aws.String(vpcId),
	})
	Expect(err).NotTo(HaveOccurred())
	return *createSubnetOutput.Subnet.SubnetId
}

func DeleteSubnet(subnetId string) {
	_, err := ec2Service.DeleteSubnet(&ec2.DeleteSubnetInput{
		SubnetId: aws.String(subnetId),
	})
	Expect(err).NotTo(HaveOccurred())
}

func CreateSecurityGroup(groupName string) string {
	createSecurityGroupOutput, err := ec2Service.CreateSecurityGroup(
		&ec2.CreateSecurityGroupInput{
			GroupName: aws.String(groupName),
			Description: aws.String("Created for integration testing. " +
				"If you are seeing this either a test is running or " +
				"this has not been cleaned up properly"),
			VpcId: aws.String(vpcId),
		},
	)
	Expect(err).NotTo(HaveOccurred())
	return *createSecurityGroupOutput.GroupId
}

func DeleteSecurityGroup(groupId string) {
	_, err := ec2Service.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(groupId),
	})
	Expect(err).NotTo(HaveOccurred())
}

func assertEnvVar(parameter *string, key string) {
	*parameter, ok = os.LookupEnv(key)
	Expect(ok).To(BeTrue(), "key "+key+" not set")
	Expect(*parameter).NotTo(BeEmpty(), "for key "+key)
}
