package mongodb

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/henrytk/aws-service-broker/aws/cloudformation"
	"github.com/henrytk/aws-service-broker/utils"
)

const adminPasswordMaxLength = 64

type Service struct {
	Client cloudformationiface.CloudFormationAPI
}

func NewService(region string) (*Service, error) {
	client, err := cloudformation.NewCloudFormationClient(region)
	if err != nil {
		return &Service{}, err
	}
	return &Service{
		Client: client,
	}, nil
}

func (s *Service) GenerateAdminPassword(input string) string {
	return utils.GetMD5Hex(input, adminPasswordMaxLength)
}

func (s *Service) GenerateStackName(input string) string {
	return "mongodb" + strings.Replace(input, "-", "", -1)
}
