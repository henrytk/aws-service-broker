package cloudformation

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awscf "github.com/aws/aws-sdk-go/service/cloudformation"
)

func NewCloudFormationClient(region string) (*awscf.CloudFormation, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, err
	}
	return awscf.New(sess), nil
}
