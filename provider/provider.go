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
		Config:         &config,
		MongoDBService: mongoDBService,
	}, nil
}

func (ap AWSProvider) Provision(context.Context, usbProvider.ProvisionData) (
	dashboardURL, operationData string, err error,
) {
	return "", "", errors.New("Error: not implemented")
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
