package provider

import (
	"encoding/json"

	"github.com/pivotal-cf/brokerapi"
)

type Provider struct {
	Config  json.RawMessage `json:"provider"`
	Catalog ProviderCatalog `json:"catalog"`
}

type ProviderCatalog struct {
	Services []ProviderService `json:"services"`
}

type ProviderService struct {
	ID             string          `json:"id"`
	ProviderConfig json.RawMessage `json:"provider"`
	Plans          []ProviderPlan  `json:"plans"`
}

type ProviderPlan struct {
	ID             string          `json:"id"`
	ProviderConfig json.RawMessage `json:"provider"`
}

type ProvisionData struct {
	InstanceID      string
	Details         brokerapi.ProvisionDetails
	Service         brokerapi.Service
	Plan            brokerapi.ServicePlan
	ProviderCatalog ProviderCatalog
}

type DeprovisionData struct {
	InstanceID      string
	Details         brokerapi.DeprovisionDetails
	Service         brokerapi.Service
	Plan            brokerapi.ServicePlan
	ProviderCatalog ProviderCatalog
}

type BindData struct {
	InstanceID      string
	BindingID       string
	Details         brokerapi.BindDetails
	ProviderCatalog ProviderCatalog
}

type UnbindData struct {
	InstanceID      string
	BindingID       string
	Details         brokerapi.UnbindDetails
	ProviderCatalog ProviderCatalog
}

type UpdateData struct {
	InstanceID      string
	Details         brokerapi.UpdateDetails
	Service         brokerapi.Service
	Plan            brokerapi.ServicePlan
	ProviderCatalog ProviderCatalog
}

type LastOperationData struct {
	InstanceID      string
	OperationData   string
	ProviderCatalog ProviderCatalog
}
