package mongodb

import (
	"github.com/aws/aws-sdk-go/aws"
	awscf "github.com/aws/aws-sdk-go/service/cloudformation"
)

func (s Service) DeleteStack(id string) error {
	_, err := s.Client.DeleteStack(&awscf.DeleteStackInput{
		ClientRequestToken: aws.String("delete-" + id),
		StackName:          aws.String(id),
	})
	return err
}
