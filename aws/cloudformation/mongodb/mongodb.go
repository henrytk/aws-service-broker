package mongodb

import (
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/henrytk/aws-service-broker/aws/cloudformation"
)

type Service struct {
	Client cloudformationiface.CloudFormationAPI
}

func NewService(region string) (Service, error) {
	client, err := cloudformation.NewCloudFormationClient(region)
	if err != nil {
		return Service{}, err
	}
	return Service{
		Client: client,
	}, nil
}
