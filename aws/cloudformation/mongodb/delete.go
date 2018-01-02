package mongodb

import (
	"github.com/aws/aws-sdk-go/aws"
	awscf "github.com/aws/aws-sdk-go/service/cloudformation"
)

func (m MongoDBService) DeleteStack(id string) error {
	_, err := m.Client.DeleteStack(&awscf.DeleteStackInput{
		ClientRequestToken: aws.String("delete-" + id),
		StackName:          aws.String(id),
	})
	return err
}
