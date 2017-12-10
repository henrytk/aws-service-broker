package provider

import (
	"context"

	"github.com/pivotal-cf/brokerapi"
)

type ServiceProvider interface {
	Provision(context.Context, ProvisionData) (dashboardURL, operationData string, err error)
	Deprovision(context.Context, DeprovisionData) (operationData string, err error)
	Bind(context.Context, BindData) (binding brokerapi.Binding, err error)
	Unbind(context.Context, UnbindData) (err error)
	Update(context.Context, UpdateData) (operationData string, err error)
	LastOperation(context.Context, LastOperationData) (state brokerapi.LastOperationState, description string, err error)
}
