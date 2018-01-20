package provider

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/henrytk/aws-service-broker/aws/cloudformation/mongodb"
	usbProvider "github.com/henrytk/universal-service-broker/provider"
	"github.com/pivotal-cf/brokerapi"
)

type AWSProvider struct {
	Config         *Config
	MongoDBService *mongodb.Service
}

func NewAWSProvider(rawConfig []byte) (*AWSProvider, error) {
	config, err := DecodeConfig(rawConfig)
	if err != nil {
		return &AWSProvider{}, err
	}
	mongoDBService, err := mongodb.NewService(config.AWSConfig.Region)
	if err != nil {
		return &AWSProvider{}, err
	}
	return &AWSProvider{
		Config:         config,
		MongoDBService: mongoDBService,
	}, nil
}

type OperationData struct {
	Type    string `json:"type"`
	StackId string `json:"stack_id,omitempty"`
}

func (ap *AWSProvider) Provision(ctx context.Context, provisionData usbProvider.ProvisionData) (
	dashboardURL string, operationData string, err error,
) {
	service, err := findServiceById(provisionData.Service.ID, &ap.Config.Catalog)
	if err != nil {
		return "", "", errors.New("could not find service ID: " + provisionData.Service.ID)
	}

	plan, err := findPlanById(provisionData.Plan.ID, service)
	if err != nil {
		return "", "", errors.New("could not find plan ID: " + provisionData.Plan.ID)
	}

	switch service.Name {
	case "mongodb":
		createStackOutput, err := ap.MongoDBService.CreateStack(
			provisionData.InstanceID,
			mongodb.InputParameters{
				BastionSecurityGroupId: service.BastionSecurityGroupId,
				KeyPairName:            service.KeyPairName,
				VpcId:                  service.VpcId,
				PrimaryNodeSubnetId:    service.PrimaryNodeSubnetId,
				Secondary0NodeSubnetId: service.Secondary0NodeSubnetId,
				Secondary1NodeSubnetId: service.Secondary1NodeSubnetId,
				MongoDBVersion:         plan.MongoDBVersion,
				MongoDBAdminUsername:   plan.MongoDBAdminUsername,
				MongoDBAdminPassword: ap.MongoDBService.GenerateAdminPassword(
					ap.Config.Secret + provisionData.InstanceID,
				),
				ClusterReplicaSetCount: plan.ClusterReplicaSetCount,
				ReplicaShardIndex:      plan.ReplicaShardIndex,
				VolumeSize:             plan.VolumeSize,
				VolumeType:             plan.VolumeType,
				Iops:                   plan.Iops,
				NodeInstanceType:       plan.NodeInstanceType,
			},
		)
		if err != nil {
			return "", "", err
		}
		operationDataJSON, err := json.Marshal(OperationData{
			Type:    "provision",
			StackId: *createStackOutput.StackId,
		})
		if err != nil {
			return "", "", err
		}
		return "", string(operationDataJSON), nil
	default:
		return "", "", errors.New("no provider for service name " + service.Name)
	}
}

func (ap *AWSProvider) Deprovision(context.Context, usbProvider.DeprovisionData) (
	operationData string, err error,
) {
	return "", errors.New("Error: not implemented")
}

func (ap *AWSProvider) Bind(context.Context, usbProvider.BindData) (
	binding brokerapi.Binding, err error,
) {
	return brokerapi.Binding{}, errors.New("Error: not implemented")
}

func (ap *AWSProvider) Unbind(context.Context, usbProvider.UnbindData) (err error) {
	return errors.New("Error: not implemented")
}

func (ap *AWSProvider) Update(context.Context, usbProvider.UpdateData) (operationData string, err error) {
	return "", errors.New("Error: not implemented")
}

func (ap *AWSProvider) LastOperation(context.Context, usbProvider.LastOperationData) (
	state brokerapi.LastOperationState, description string, err error,
) {
	return brokerapi.LastOperationState(""), "", errors.New("Error: not implemented")
}
