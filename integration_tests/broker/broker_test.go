package broker_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/henrytk/aws-service-broker/broker"
	"github.com/henrytk/aws-service-broker/provider"
	"github.com/pivotal-cf/brokerapi"
	uuid "github.com/satori/go.uuid"

	usb "github.com/henrytk/universal-service-broker/broker"
	broker_tester "github.com/henrytk/universal-service-broker/broker/testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	DEFAULT_TIMEOUT time.Duration = 15 * time.Minute
)

var _ = Describe("Broker", func() {

	const ASYNC bool = true

	var (
		err              error
		config           usb.Config
		awsProvider      *provider.AWSProvider
		awsServiceBroker http.Handler
		brokerTester     broker_tester.BrokerTester
		instanceID       string

		serviceID        string = "serviceID"
		plan1ID          string = "plan1ID"
		plan2ID          string = "plan2ID"
		organizationGUID string = "orgGUID"
		spaceGUID        string = "spaceGUID"
	)

	BeforeEach(func() {
		configFile := bytes.NewReader([]byte(`
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
						"id": "` + serviceID + `",
						"name": "mongodb",
						"description": "MongoDB clusters via AWS CloudFormation",
						"bindable": true,
						"plan_updateable": true,
						"requires": [],
						"metadata": {},
						"bastion_security_group_id": "` + bastionSecurityGroupId + `",
						"key_pair_name": "` + keyPairName + `",
						"vpc_id": "` + vpcId + `",
						"primary_node_subnet_id": "` + primaryNodeSubnetId + `",
						"secondary_0_node_subnet_id": "` + secondary0NodeSubnetId + `",
						"secondary_1_node_subnet_id": "` + secondary1NodeSubnetId + `",
						"plans": [{
							"id": "` + plan1ID + `",
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
							"id": "` + plan2ID + `",
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
		`))

		config, err = usb.NewConfig(configFile)
		Expect(err).NotTo(HaveOccurred())

		awsProvider, err = provider.NewAWSProvider(config.Provider)
		Expect(err).NotTo(HaveOccurred())

		awsServiceBroker = broker.NewAWSServiceBroker(config, awsProvider)

		brokerTester = broker_tester.New(brokerapi.BrokerCredentials{
			Username: config.API.BasicAuthUsername,
			Password: config.API.BasicAuthPassword,
		}, awsServiceBroker)

		instanceID = uuid.NewV4().String()
	})

	Describe("MongoDB", func() {
		It("should manage the MongoDB cluster lifecycle", func() {
			By("provisioning an instance")
			res := brokerTester.Provision(
				instanceID,
				broker_tester.RequestBody{
					ServiceID:        serviceID,
					PlanID:           plan1ID,
					OrganizationGUID: organizationGUID,
					SpaceGUID:        spaceGUID,
				},
				ASYNC,
			)
			Expect(res.Code).To(Equal(http.StatusAccepted))

			provisioningResponse := brokerapi.ProvisioningResponse{}
			err = json.Unmarshal(res.Body.Bytes(), &provisioningResponse)
			Expect(err).NotTo(HaveOccurred())

			By("reporting status of last operation: provision")
			pollForCompletion(brokerTester, instanceID, provisioningResponse.OperationData, brokerapi.LastOperationResponse{
				State:       brokerapi.Succeeded,
				Description: "provision succeeded",
			})

			By("updating the instance")
			res = brokerTester.Update(
				instanceID,
				broker_tester.RequestBody{
					ServiceID: serviceID,
					PlanID:    plan1ID,
					PreviousValues: &broker_tester.RequestBody{
						PlanID: plan2ID,
					},
				},
				ASYNC,
			)
			Expect(res.Code).To(Equal(http.StatusAccepted))

			updateResponse := brokerapi.UpdateResponse{}
			err = json.Unmarshal(res.Body.Bytes(), &updateResponse)
			Expect(err).NotTo(HaveOccurred())

			By("reporting status of last operation: update")
			pollForCompletion(brokerTester, instanceID, updateResponse.OperationData, brokerapi.LastOperationResponse{
				State:       brokerapi.Succeeded,
				Description: "update succeeded",
			})

			By("deprovisioning the instance")
			res = brokerTester.Deprovision(instanceID, serviceID, plan1ID, ASYNC)
			Expect(res.Code).To(Equal(http.StatusAccepted))

			deprovisionResponse := brokerapi.DeprovisionResponse{}
			err = json.Unmarshal(res.Body.Bytes(), &deprovisionResponse)
			Expect(err).NotTo(HaveOccurred())

			By("reporting the status of last operation: deprovision")
			pollForCompletion(brokerTester, instanceID, deprovisionResponse.OperationData, brokerapi.LastOperationResponse{
				State:       brokerapi.Succeeded,
				Description: "deprovision succeeded",
			})
		})
	})
})

func pollForCompletion(bt broker_tester.BrokerTester, instanceID, operationData string, expectedResponse brokerapi.LastOperationResponse) {
	Eventually(
		func() brokerapi.LastOperationResponse {
			lastOperationResponse := brokerapi.LastOperationResponse{}
			res := bt.LastOperation(instanceID, "", "", operationData)
			if res.Code != http.StatusOK {
				return lastOperationResponse
			}
			_ = json.Unmarshal(res.Body.Bytes(), &lastOperationResponse)
			return lastOperationResponse
		},
		DEFAULT_TIMEOUT,
		30*time.Second,
	).Should(Equal(expectedResponse))
}
