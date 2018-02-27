package mongodb

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awscf "github.com/aws/aws-sdk-go/service/cloudformation"
)

func (s *Service) GetStackState(stackName string) (string, string, error) {
	describeStacksOutput, err := s.Client.DescribeStacks(&awscf.DescribeStacksInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return "", "", err
	}

	if len(describeStacksOutput.Stacks) != 1 {
		return "", "", errors.New("Error checking stack state: number of stacks was not 1")
	}

	stack := describeStacksOutput.Stacks[0]
	reason := "no reason returned via the API"
	if stack.StackStatusReason != nil {
		reason = *stack.StackStatusReason
	}
	return *stack.StackStatus, reason, nil
}

func (s *Service) CreateStackCompleted(id string) (bool, error) {
	stackName := s.GenerateStackName(id)
	state, reason, err := s.GetStackState(stackName)
	if err != nil {
		return false, err
	}

	if state == awscf.StackStatusCreateComplete {
		return true, nil
	} else if stackStateIsFinal(state) {
		return true, errors.New("Final state of stack was not " + awscf.StackStatusCreateComplete + ". Got: " + state + ". Reason: " + reason)
	}
	return false, nil
}

func (s *Service) DeleteStackCompleted(id string) (bool, error) {
	stackName := s.GenerateStackName(id)
	state, reason, err := s.GetStackState(stackName)
	if err != nil {
		if strings.Contains(err.Error(), "Stack with id "+stackName+" does not exist") {
			return true, nil
		}
		return false, err
	}

	if state == awscf.StackStatusDeleteComplete {
		return true, nil
	} else if stackStateIsFinal(state) {
		return true, errors.New("Final state of stack was not " + awscf.StackStatusDeleteComplete + ". Got: " + state + ". Reason: " + reason)
	}
	return false, nil
}

func (s *Service) UpdateStackCompleted(id string) (bool, error) {
	stackName := s.GenerateStackName(id)
	state, reason, err := s.GetStackState(stackName)
	if err != nil {
		return false, err
	}

	if state == awscf.StackStatusUpdateComplete {
		return true, nil
	} else if stackStateIsFinal(state) {
		return true, errors.New("Final state of stack was not " + awscf.StackStatusUpdateComplete + ". Got: " + state + ". Reason: " + reason)
	}
	return false, nil
}

func stackStateIsFinal(state string) bool {
	switch state {
	case awscf.StackStatusCreateInProgress:
		return false
	case awscf.StackStatusCreateFailed:
		return true
	case awscf.StackStatusCreateComplete:
		return true
	case awscf.StackStatusRollbackInProgress:
		return false
	case awscf.StackStatusRollbackFailed:
		return true
	case awscf.StackStatusRollbackComplete:
		return true
	case awscf.StackStatusDeleteInProgress:
		return false
	case awscf.StackStatusDeleteFailed:
		return true
	case awscf.StackStatusDeleteComplete:
		return true
	case awscf.StackStatusUpdateInProgress:
		return false
	case awscf.StackStatusUpdateCompleteCleanupInProgress:
		return false
	case awscf.StackStatusUpdateComplete:
		return true
	case awscf.StackStatusUpdateRollbackInProgress:
		return false
	case awscf.StackStatusUpdateRollbackFailed:
		return true
	case awscf.StackStatusUpdateRollbackCompleteCleanupInProgress:
		return false
	case awscf.StackStatusUpdateRollbackComplete:
		return true
	case awscf.StackStatusReviewInProgress:
		return false
	default:
		return true
	}
}
