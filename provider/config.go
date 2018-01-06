package provider

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/pivotal-cf/brokerapi"
)

type Config struct {
	Catalog Catalog `json:"catalog"`
}

type Catalog struct {
	Services []Service `json:"services"`
}

type Service struct {
	brokerapi.Service
	MongoDBServiceParameters
	Plans []Plan `json:"plans"`
}

type Plan struct {
	brokerapi.ServicePlan
	MongoDBPlanParameters
}

type MongoDBServiceParameters struct {
	BastionSecurityGroupId string `json:"bastion_security_group_id"`
	KeyPairName            string `json:"key_pair_name"`
	VPC                    string `json:"vpc_id"`
	PrimaryNodeSubnet      string `json:"primary_node_subnet_id"`
	Secondary0NodeSubnet   string `json:"secondary_0_node_subnet_id"`
	Secondary1NodeSubnet   string `json:"secondary_1_node_subnet_id"`
}

type MongoDBPlanParameters struct {
	ClusterReplicaSetCount string `json:"cluster_replica_set_count"`
	ReplicaShardIndex      string `json:"replica_shard_index"`
	VolumeSize             string `json:"volume_size"`
	VolumeType             string `json:"volume_type"`
	Iops                   string `json:"iops"`
	NodeInstanceType       string `json:"node_instance_type"`
}

func DecodeConfig(rawMessage json.RawMessage) (Config, error) {
	var config Config
	err := json.Unmarshal([]byte(rawMessage), &config)
	if err != nil {
		return config, err
	}
	if reflect.DeepEqual(config.Catalog, Catalog{}) {
		return config, errors.New("Error decoding config: no catalog found")
	}
	if len(config.Catalog.Services) == 0 {
		return config, errors.New("Error decoding config: at least one service must be configured")
	}

	for _, service := range config.Catalog.Services {
		switch service.Name {
		case "mongodb":
		default:
			return config, errors.New("Error decoding config: service name " + service.Name + " not recognised")
		}

		if len(service.Plans) == 0 {
			return config, errors.New("Error decoding config: at least one plan must be configured for service " + service.Name)
		}
	}

	return config, nil
}
