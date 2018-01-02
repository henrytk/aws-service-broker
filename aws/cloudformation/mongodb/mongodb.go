package mongodb

import (
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/henrytk/aws-service-broker/aws/cloudformation"
)

type MongoDBService struct {
	Client cloudformationiface.CloudFormationAPI
}

func NewMongoDBService(region string) (MongoDBService, error) {
	client, err := cloudformation.NewCloudFormationClient(region)
	if err != nil {
		return MongoDBService{}, err
	}
	return MongoDBService{
		Client: client,
	}, nil
}
