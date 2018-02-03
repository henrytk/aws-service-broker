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
	internetGatewayId      *string
	routeTableId           *string
	association1           *string
	association2           *string
	association3           *string
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
		CreateRouteTable()
		mongoDBAdminPassword = "volunteer-pilot"

		mongoDBService, err = mongodb.NewService(region)
		Expect(err).NotTo(HaveOccurred())
		assertEnvVar(&keyPairName, "ASB_KEY_PAIR")
	})

	AfterSuite(func() {
		DeleteRouteTable()
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
	awsVpcId := vpcOutput.Vpc.VpcId

	createInternetGatewayOutput, err := ec2Service.CreateInternetGateway(
		&ec2.CreateInternetGatewayInput{},
	)
	Expect(err).NotTo(HaveOccurred())

	internetGatewayId = createInternetGatewayOutput.InternetGateway.InternetGatewayId
	_, err = ec2Service.AttachInternetGateway(
		&ec2.AttachInternetGatewayInput{
			InternetGatewayId: internetGatewayId,
			VpcId:             awsVpcId,
		},
	)
	Expect(err).NotTo(HaveOccurred())

	vpcId = *awsVpcId
}

func DeleteVPC() {
	_, err := ec2Service.DetachInternetGateway(
		&ec2.DetachInternetGatewayInput{
			InternetGatewayId: internetGatewayId,
			VpcId:             aws.String(vpcId),
		},
	)
	Expect(err).NotTo(HaveOccurred())

	_, err = ec2Service.DeleteInternetGateway(&ec2.DeleteInternetGatewayInput{
		InternetGatewayId: internetGatewayId,
	})
	Expect(err).NotTo(HaveOccurred())

	_, err = ec2Service.DeleteVpc(&ec2.DeleteVpcInput{
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

func CreateRouteTable() {
	createRouteTableOutput, err := ec2Service.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId: aws.String(vpcId),
	})
	Expect(err).NotTo(HaveOccurred())

	routeTableId = createRouteTableOutput.RouteTable.RouteTableId

	_, err = ec2Service.CreateRoute(&ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		GatewayId:            internetGatewayId,
		RouteTableId:         routeTableId,
	})
	Expect(err).NotTo(HaveOccurred())

	associateRouteTableOutput1, err := ec2Service.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		RouteTableId: routeTableId,
		SubnetId:     aws.String(primaryNodeSubnetId),
	})
	Expect(err).NotTo(HaveOccurred())
	association1 = associateRouteTableOutput1.AssociationId

	associateRouteTableOutput2, err := ec2Service.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		RouteTableId: routeTableId,
		SubnetId:     aws.String(secondary0NodeSubnetId),
	})
	Expect(err).NotTo(HaveOccurred())
	association2 = associateRouteTableOutput2.AssociationId

	associateRouteTableOutput3, err := ec2Service.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		RouteTableId: routeTableId,
		SubnetId:     aws.String(secondary1NodeSubnetId),
	})
	Expect(err).NotTo(HaveOccurred())
	association3 = associateRouteTableOutput3.AssociationId
}

func DeleteRouteTable() {
	_, err := ec2Service.DisassociateRouteTable(&ec2.DisassociateRouteTableInput{
		AssociationId: association3,
	})
	Expect(err).NotTo(HaveOccurred())

	_, err = ec2Service.DisassociateRouteTable(&ec2.DisassociateRouteTableInput{
		AssociationId: association2,
	})
	Expect(err).NotTo(HaveOccurred())

	_, err = ec2Service.DisassociateRouteTable(&ec2.DisassociateRouteTableInput{
		AssociationId: association1,
	})
	Expect(err).NotTo(HaveOccurred())

	Expect(err).NotTo(HaveOccurred())
	_, err = ec2Service.DeleteRouteTable(&ec2.DeleteRouteTableInput{
		RouteTableId: routeTableId,
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
