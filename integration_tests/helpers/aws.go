package helpers

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/gomega"
)

var (
	svc *ec2.EC2
)

type Vpc struct {
	VpcId           *string
	InternetGateway *ec2.InternetGateway
	KeyPairName     string
	Subnets         []*ec2.Subnet
	RouteTables     []*ec2.RouteTable
	NatGateways     []*ec2.NatGateway
	SecurityGroups  []*ec2.SecurityGroup
}

func SetupVpc(region, id string) *Vpc {
	createEC2Service(region)
	vpcId := createVpc()
	internetGatewayId := createInternetGateway(vpcId)

	keyPairName := "mongodb-test-" + id
	createKeyPair(keyPairName)

	publicSubnet1 := createSubnet("10.0.0.0/24", vpcId)
	natGateway1 := createNatGateway(publicSubnet1.SubnetId)
	publicSubnet2 := createSubnet("10.0.1.0/24", vpcId)
	natGateway2 := createNatGateway(publicSubnet2.SubnetId)
	publicSubnet3 := createSubnet("10.0.2.0/24", vpcId)
	natGateway3 := createNatGateway(publicSubnet3.SubnetId)

	routeTable := createRouteTable(vpcId, internetGatewayId, nil)
	associateRouteTable(routeTable, publicSubnet1)
	associateRouteTable(routeTable, publicSubnet2)
	associateRouteTable(routeTable, publicSubnet3)

	privateSubnet1 := createSubnet("10.0.3.0/24", vpcId)
	routeTable1 := createRouteTable(vpcId, nil, natGateway1.NatGatewayId)
	associateRouteTable(routeTable1, privateSubnet1)

	privateSubnet2 := createSubnet("10.0.4.0/24", vpcId)
	routeTable2 := createRouteTable(vpcId, nil, natGateway2.NatGatewayId)
	associateRouteTable(routeTable2, privateSubnet2)

	privateSubnet3 := createSubnet("10.0.5.0/24", vpcId)
	routeTable3 := createRouteTable(vpcId, nil, natGateway3.NatGatewayId)
	associateRouteTable(routeTable3, privateSubnet3)

	bastionSecurityGroup := createSecurityGroup("bastion-mongodb-test", vpcId)

	return &Vpc{
		VpcId:           vpcId,
		InternetGateway: &ec2.InternetGateway{InternetGatewayId: internetGatewayId},
		KeyPairName:     keyPairName,
		Subnets: []*ec2.Subnet{
			publicSubnet1,
			publicSubnet2,
			publicSubnet3,
			privateSubnet1,
			privateSubnet2,
			privateSubnet3,
		},
		RouteTables: []*ec2.RouteTable{
			routeTable,
			routeTable1,
			routeTable2,
			routeTable3,
		},
		NatGateways:    []*ec2.NatGateway{natGateway1, natGateway2, natGateway3},
		SecurityGroups: []*ec2.SecurityGroup{bastionSecurityGroup},
	}
}

func DestroyVpc(vpc *Vpc) {
	for _, sg := range vpc.SecurityGroups {
		deleteSecurityGroup(sg.GroupId)
	}
	for _, ng := range vpc.NatGateways {
		deleteNatGateway(ng)
	}
	for _, rt := range vpc.RouteTables {
		deleteRouteTable(rt)
	}
	for _, s := range vpc.Subnets {
		deleteSubnet(s)
	}

	deleteKeyPair(vpc.KeyPairName)
	deleteVpc(vpc)
}

func createEC2Service(region string) {
	awsSession, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	Expect(err).NotTo(HaveOccurred())
	svc = ec2.New(awsSession)
}

func createVpc() *string {
	vpcOutput, err := svc.CreateVpc(&ec2.CreateVpcInput{
		CidrBlock: aws.String("10.0.0.0/16"),
	})
	Expect(err).NotTo(HaveOccurred())

	return vpcOutput.Vpc.VpcId
}

func createInternetGateway(vpcId *string) *string {
	createInternetGatewayOutput, err := svc.CreateInternetGateway(
		&ec2.CreateInternetGatewayInput{},
	)
	Expect(err).NotTo(HaveOccurred())

	internetGatewayId := createInternetGatewayOutput.InternetGateway.InternetGatewayId
	_, err = svc.AttachInternetGateway(
		&ec2.AttachInternetGatewayInput{
			InternetGatewayId: internetGatewayId,
			VpcId:             vpcId,
		},
	)
	Expect(err).NotTo(HaveOccurred())

	return internetGatewayId
}

func createKeyPair(keyPairName string) {
	createKeyPairOutput, err := svc.CreateKeyPair(&ec2.CreateKeyPairInput{
		KeyName: aws.String(keyPairName),
	})
	Expect(err).NotTo(HaveOccurred())
	err = ioutil.WriteFile("/tmp/"+keyPairName, []byte(*createKeyPairOutput.KeyMaterial), 0400)
	Expect(err).NotTo(HaveOccurred())
}

func createSubnet(cidrBlock string, vpcId *string) *ec2.Subnet {
	createSubnetOutput, err := svc.CreateSubnet(&ec2.CreateSubnetInput{
		CidrBlock: aws.String(cidrBlock),
		VpcId:     vpcId,
	})
	Expect(err).NotTo(HaveOccurred())

	return createSubnetOutput.Subnet
}

func createRouteTable(vpcId, gatewayId, natGatewayId *string) *ec2.RouteTable {
	createRouteTableOutput, err := svc.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId: vpcId,
	})
	Expect(err).NotTo(HaveOccurred())

	routeTable := createRouteTableOutput.RouteTable

	_, err = svc.CreateRoute(&ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		GatewayId:            gatewayId,
		NatGatewayId:         natGatewayId,
		RouteTableId:         routeTable.RouteTableId,
	})
	Expect(err).NotTo(HaveOccurred())

	routeTable.Routes = append(routeTable.Routes, &ec2.Route{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		GatewayId:            gatewayId,
		NatGatewayId:         natGatewayId,
	})

	return createRouteTableOutput.RouteTable
}

func associateRouteTable(routeTable *ec2.RouteTable, subnet *ec2.Subnet) {
	associateRouteTableOutput, err := svc.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		RouteTableId: routeTable.RouteTableId,
		SubnetId:     subnet.SubnetId,
	})
	Expect(err).NotTo(HaveOccurred())

	routeTable.Associations = append(routeTable.Associations, &ec2.RouteTableAssociation{
		RouteTableAssociationId: associateRouteTableOutput.AssociationId,
		RouteTableId:            routeTable.RouteTableId,
		SubnetId:                subnet.SubnetId,
	})
}

func createNatGateway(subnetId *string) *ec2.NatGateway {
	allocateAddressOutput, err := svc.AllocateAddress(&ec2.AllocateAddressInput{})
	Expect(err).NotTo(HaveOccurred())

	createNatGatewayOutput, err := svc.CreateNatGateway(&ec2.CreateNatGatewayInput{
		AllocationId: allocateAddressOutput.AllocationId,
		SubnetId:     subnetId,
	})
	Expect(err).NotTo(HaveOccurred())

	natGateway := createNatGatewayOutput.NatGateway

	CREATE_NAT_GATEWAY_TIMEOUT := 5 * time.Minute
	POLLING_INTERVAL := 10 * time.Second
	Eventually(
		func() string {
			describeNatGatewaysOutput, err := svc.DescribeNatGateways(
				&ec2.DescribeNatGatewaysInput{
					Filter: []*ec2.Filter{
						{
							Name:   aws.String("nat-gateway-id"),
							Values: []*string{natGateway.NatGatewayId},
						},
					},
				},
			)
			if err != nil {
				return ""
			}
			if len(describeNatGatewaysOutput.NatGateways) != 1 {
				return ""
			}
			return *describeNatGatewaysOutput.NatGateways[0].State
		},
		CREATE_NAT_GATEWAY_TIMEOUT,
		POLLING_INTERVAL,
	).Should(Equal("available"))

	return natGateway
}

func createSecurityGroup(groupName string, vpcId *string) *ec2.SecurityGroup {
	createSecurityGroupOutput, err := svc.CreateSecurityGroup(
		&ec2.CreateSecurityGroupInput{
			GroupName: aws.String(groupName),
			Description: aws.String("Created for integration testing. " +
				"If you are seeing this either a test is running or " +
				"this has not been cleaned up properly"),
			VpcId: vpcId,
		},
	)
	Expect(err).NotTo(HaveOccurred())

	return &ec2.SecurityGroup{
		GroupName: aws.String(groupName),
		GroupId:   createSecurityGroupOutput.GroupId,
	}
}

func deleteSecurityGroup(groupId *string) {
	_, err := svc.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
		GroupId: groupId,
	})
	Expect(err).NotTo(HaveOccurred())
}

func deleteNatGateway(natGateway *ec2.NatGateway) {
	_, err := svc.DeleteNatGateway(&ec2.DeleteNatGatewayInput{
		NatGatewayId: natGateway.NatGatewayId,
	})
	Expect(err).NotTo(HaveOccurred())

	DELETE_NAT_GATEWAY_TIMEOUT := 5 * time.Minute
	POLLING_INTERVAL := 10 * time.Second
	Eventually(
		func() string {
			describeNatGatewaysOutput, err := svc.DescribeNatGateways(
				&ec2.DescribeNatGatewaysInput{
					Filter: []*ec2.Filter{
						{
							Name:   aws.String("nat-gateway-id"),
							Values: []*string{natGateway.NatGatewayId},
						},
					},
				},
			)
			if err != nil {
				return ""
			}
			if len(describeNatGatewaysOutput.NatGateways) != 1 {
				return ""
			}
			return *describeNatGatewaysOutput.NatGateways[0].State
		},
		DELETE_NAT_GATEWAY_TIMEOUT,
		POLLING_INTERVAL,
	).Should(Equal("deleted"))

	for _, address := range natGateway.NatGatewayAddresses {
		_, err = svc.ReleaseAddress(&ec2.ReleaseAddressInput{
			AllocationId: address.AllocationId,
		})
		Expect(err).NotTo(HaveOccurred())
	}
}

func deleteRouteTable(routeTable *ec2.RouteTable) {

	for _, ass := range routeTable.Associations {
		_, err := svc.DisassociateRouteTable(&ec2.DisassociateRouteTableInput{
			AssociationId: ass.RouteTableAssociationId,
		})
		Expect(err).NotTo(HaveOccurred())
	}

	_, err := svc.DeleteRouteTable(&ec2.DeleteRouteTableInput{
		RouteTableId: routeTable.RouteTableId,
	})
	Expect(err).NotTo(HaveOccurred())
}

func deleteSubnet(subnet *ec2.Subnet) {
	_, err := svc.DeleteSubnet(&ec2.DeleteSubnetInput{
		SubnetId: subnet.SubnetId,
	})
	Expect(err).NotTo(HaveOccurred())
}

func deleteKeyPair(keyPairName string) {
	_, err := svc.DeleteKeyPair(&ec2.DeleteKeyPairInput{
		KeyName: aws.String(keyPairName),
	})
	Expect(err).NotTo(HaveOccurred())
	err = os.Remove("/tmp/" + keyPairName)
	Expect(err).NotTo(HaveOccurred())
}

func deleteVpc(vpc *Vpc) {
	_, err := svc.DetachInternetGateway(
		&ec2.DetachInternetGatewayInput{
			InternetGatewayId: vpc.InternetGateway.InternetGatewayId,
			VpcId:             vpc.VpcId,
		},
	)
	Expect(err).NotTo(HaveOccurred())

	_, err = svc.DeleteInternetGateway(&ec2.DeleteInternetGatewayInput{
		InternetGatewayId: vpc.InternetGateway.InternetGatewayId,
	})
	Expect(err).NotTo(HaveOccurred())

	_, err = svc.DeleteVpc(&ec2.DeleteVpcInput{
		VpcId: vpc.VpcId,
	})
	Expect(err).NotTo(HaveOccurred())
}
