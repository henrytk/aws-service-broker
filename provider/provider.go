package provider

import (
	"context"
	"errors"

	"github.com/henrytk/aws-service-broker/aws/cloudformation/mongodb"
	usbProvider "github.com/henrytk/universal-service-broker/provider"
	"github.com/pivotal-cf/brokerapi"
)

type AWSProvider struct {
	Config         *Config
	MongoDBService mongodb.Service
}

func NewAWSProvider(rawConfig []byte) (AWSProvider, error) {
	config, err := DecodeConfig(rawConfig)
	if err != nil {
		return AWSProvider{}, err
	}
	mongoDBService, err := mongodb.NewService(config.AWSConfig.Region)
	if err != nil {
		return AWSProvider{}, err
	}
	return AWSProvider{
		Config:         config,
		MongoDBService: mongoDBService,
	}, nil
}

type OperationData struct {
	Type string
}

func (ap AWSProvider) Provision(ctx context.Context, provisionData usbProvider.ProvisionData) (
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
		_, err = ap.MongoDBService.CreateStack(
			provisionData.InstanceID,
			mongodb.InputParameters{
				ap.MongoDBService.GenerateAdminPassword("irrelevant"),
				service.BastionSecurityGroupId,
				service.KeyPairName,
				service.VpcId,
				service.PrimaryNodeSubnetId,
				service.Secondary0NodeSubnetId,
				service.Secondary1NodeSubnetId,
				plan.MongoDBVersion,
				plan.ClusterReplicaSetCount,
				plan.ReplicaShardIndex,
				plan.VolumeSize,
				plan.VolumeType,
				plan.Iops,
				plan.NodeInstanceType,
			},
		)
		if err != nil {
			return "", "", err
		}
		return "", "", errors.New("not implemented")
	default:
		return "", "", errors.New("no provider for service name " + service.Name)
	}
}

func (ap AWSProvider) Deprovision(context.Context, usbProvider.DeprovisionData) (
	operationData string, err error,
) {
	return "", errors.New("Error: not implemented")
}

func (ap AWSProvider) Bind(context.Context, usbProvider.BindData) (
	binding brokerapi.Binding, err error,
) {
	return brokerapi.Binding{}, errors.New("Error: not implemented")
}

func (ap AWSProvider) Unbind(context.Context, usbProvider.UnbindData) (err error) {
	return errors.New("Error: not implemented")
}

func (ap AWSProvider) Update(context.Context, usbProvider.UpdateData) (operationData string, err error) {
	return "", errors.New("Error: not implemented")
}

func (ap AWSProvider) LastOperation(context.Context, usbProvider.LastOperationData) (
	state brokerapi.LastOperationState, description string, err error,
) {
	return brokerapi.LastOperationState(""), "", errors.New("Error: not implemented")
}
