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
	Type       string `json:"type"`
	Service    string `json:"service"`
	StackId    string `json:"stack_id,omitempty"`
	InstanceID string `json:"instance_id,omitempty"`
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
			Service: service.Name,
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

func (ap *AWSProvider) Deprovision(ctx context.Context, deprovisionData usbProvider.DeprovisionData) (
	operationData string, err error,
) {
	service, err := findServiceById(deprovisionData.Service.ID, &ap.Config.Catalog)
	if err != nil {
		return "", errors.New("could not find service ID: " + deprovisionData.Service.ID)
	}

	switch service.Name {
	case "mongodb":
		err := ap.MongoDBService.DeleteStack(deprovisionData.InstanceID)
		if err != nil {
			return "", err
		}
		operationDataJSON, err := json.Marshal(OperationData{
			Type:       "deprovision",
			Service:    service.Name,
			InstanceID: deprovisionData.InstanceID,
		})
		if err != nil {
			return "", err
		}
		return string(operationDataJSON), nil
	default:
		return "", errors.New("no provider for service name " + service.Name)
	}
}

func (ap *AWSProvider) Bind(context.Context, usbProvider.BindData) (
	binding brokerapi.Binding, err error,
) {
	return brokerapi.Binding{}, errors.New("Error: not implemented")
}

func (ap *AWSProvider) Unbind(context.Context, usbProvider.UnbindData) (err error) {
	return errors.New("Error: not implemented")
}

func (ap *AWSProvider) Update(ctx context.Context, updateData usbProvider.UpdateData) (operationData string, err error) {
	if len(updateData.Details.RawParameters) > 0 {
		return "", errors.New("update parameters are not supported")
	}

	service, err := findServiceById(updateData.Service.ID, &ap.Config.Catalog)
	if err != nil {
		return "", errors.New("could not find service ID: " + updateData.Service.ID)
	}

	newPlan, err := findPlanById(updateData.Plan.ID, service)
	if err != nil {
		return "", errors.New("could not find plan ID: " + updateData.Plan.ID)
	}

	currentPlan, err := findPlanById(updateData.Details.PreviousValues.PlanID, service)
	if err != nil {
		return "", errors.New("could not find plan ID: " + updateData.Details.PreviousValues.PlanID)
	}

	if err := validPlanUpdate(currentPlan, newPlan); err != nil {
		return "", err
	}

	switch service.Name {
	case "mongodb":
		updateParameters := buildMongoDBUpdateParameters(currentPlan, newPlan)
		updateStackOutput, err := ap.MongoDBService.UpdateStack(ctx, updateData.InstanceID, updateParameters)
		if err != nil {
			return "", err
		}
		operationDataJSON, err := json.Marshal(OperationData{
			Type:    "update",
			Service: service.Name,
			StackId: *updateStackOutput.StackId,
		})
		if err != nil {
			return "", err
		}
		return string(operationDataJSON), nil
	default:
		return "", errors.New("no provider for service name " + service.Name)
	}
}

func (ap *AWSProvider) LastOperation(ctx context.Context, lastOperationData usbProvider.LastOperationData) (
	state brokerapi.LastOperationState, description string, err error,
) {
	var operationData OperationData
	err = json.Unmarshal([]byte(lastOperationData.OperationData), &operationData)
	if err != nil {
		return "", "", err
	}

	switch operationData.Service {
	case "mongodb":
		switch operationData.Type {
		case "provision":
			completed, err := ap.MongoDBService.CreateStackCompleted(lastOperationData.InstanceID)
			if completed {
				if err == nil {
					return brokerapi.Succeeded, "provision succeeded", nil
				} else {
					return brokerapi.Failed, err.Error(), nil
				}
			}
			return brokerapi.InProgress, "provision in progress", nil
		case "deprovision":
			completed, err := ap.MongoDBService.DeleteStackCompleted(lastOperationData.InstanceID)
			if completed {
				if err == nil {
					return brokerapi.Succeeded, "deprovision succeeded", nil
				} else {
					return brokerapi.Failed, err.Error(), nil
				}
			}
			return brokerapi.InProgress, "deprovision in progress", nil
		case "update":
			completed, err := ap.MongoDBService.UpdateStackCompleted(lastOperationData.InstanceID)
			if completed {
				if err == nil {
					return brokerapi.Succeeded, "update succeeded", nil
				} else {
					return brokerapi.Failed, err.Error(), nil
				}
			}
			return brokerapi.InProgress, "update in progress", nil
		default:
			return "", "", errors.New("unknown operation type '" + operationData.Type + "'")
		}
	default:
		return "", "", errors.New("unknown service '" + operationData.Service + "'")
	}
}

func validPlanUpdate(currentPlan, newPlan Plan) error {
	if currentPlan.MongoDBAdminUsername != newPlan.MongoDBAdminUsername {
		return errors.New("updating MongoDB admin username is not supported")
	}
	if currentPlan.MongoDBVersion != newPlan.MongoDBVersion {
		return errors.New("updating MongoDB version is not supported")
	}
	if currentPlan.ClusterReplicaSetCount != newPlan.ClusterReplicaSetCount {
		return errors.New("updating cluster replica set count is not supported")
	}
	if currentPlan.ReplicaShardIndex != newPlan.ReplicaShardIndex {
		return errors.New("updating replica shard index is not supported")
	}
	if currentPlan.VolumeSize != newPlan.VolumeSize {
		return errors.New("updating volume size is not supported")
	}
	if currentPlan.VolumeType != newPlan.VolumeType {
		return errors.New("updating volume type is not supported")
	}
	if currentPlan.Iops != newPlan.Iops {
		return errors.New("updating IOPS is not supported")
	}
	return nil
}

func buildMongoDBUpdateParameters(currentPlan, newPlan Plan) mongodb.InputParameters {
	updateParameters := mongodb.InputParameters{}
	if currentPlan.MongoDBVersion != newPlan.MongoDBVersion {
		updateParameters.MongoDBVersion = newPlan.MongoDBVersion
	}
	if currentPlan.MongoDBAdminUsername != newPlan.MongoDBAdminUsername {
		updateParameters.MongoDBAdminUsername = newPlan.MongoDBAdminUsername
	}
	if currentPlan.ClusterReplicaSetCount != newPlan.ClusterReplicaSetCount {
		updateParameters.ClusterReplicaSetCount = newPlan.ClusterReplicaSetCount
	}
	if currentPlan.ReplicaShardIndex != newPlan.ReplicaShardIndex {
		updateParameters.ReplicaShardIndex = newPlan.ReplicaShardIndex
	}
	if currentPlan.VolumeSize != newPlan.VolumeSize {
		updateParameters.VolumeSize = newPlan.VolumeSize
	}
	if currentPlan.VolumeType != newPlan.VolumeType {
		updateParameters.VolumeType = newPlan.VolumeType
	}
	if currentPlan.Iops != newPlan.Iops {
		updateParameters.Iops = newPlan.Iops
	}
	if currentPlan.NodeInstanceType != newPlan.NodeInstanceType {
		updateParameters.NodeInstanceType = newPlan.NodeInstanceType
	}
	return updateParameters
}
