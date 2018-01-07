package provider

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/pivotal-cf/brokerapi"
)

type Config struct {
	AWSConfig AWSConfig `json:"aws_config"`
	Catalog   Catalog   `json:"catalog"`
}

type AWSConfig struct {
	Region string `json:"region"`
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
	VpcId                  string `json:"vpc_id"`
	PrimaryNodeSubnetId    string `json:"primary_node_subnet_id"`
	Secondary0NodeSubnetId string `json:"secondary_0_node_subnet_id"`
	Secondary1NodeSubnetId string `json:"secondary_1_node_subnet_id"`
}

type MongoDBPlanParameters struct {
	MongoDBVersion         string `json:"mongodb_version"`
	ClusterReplicaSetCount string `json:"cluster_replica_set_count"`
	ReplicaShardIndex      string `json:"replica_shard_index"`
	VolumeSize             string `json:"volume_size"`
	VolumeType             string `json:"volume_type"`
	Iops                   string `json:"iops"`
	NodeInstanceType       string `json:"node_instance_type"`
}

func DecodeConfig(b []byte) (*Config, error) {
	var config *Config
	err := json.Unmarshal(b, &config)
	if err != nil {
		return config, err
	}
	if config.AWSConfig.Region == "" {
		return config, errors.New("Config error: must provide AWS region")
	}
	if reflect.DeepEqual(config.Catalog, Catalog{}) {
		return config, errors.New("Config error: no catalog found")
	}
	if len(config.Catalog.Services) == 0 {
		return config, errors.New("Config error: at least one service must be configured")
	}

	for _, service := range config.Catalog.Services {
		switch service.Name {
		case "mongodb":
			if service.BastionSecurityGroupId == "" {
				return config, errors.New("Config error: must provide bastion security group ID")
			}
			if service.KeyPairName == "" {
				return config, errors.New("Config error: must provide key pair name")
			}
			if service.VpcId == "" {
				return config, errors.New("Config error: must provide VPC ID")
			}
			if service.PrimaryNodeSubnetId == "" {
				return config, errors.New("Config error: must provide primary node subnet ID")
			}
			if service.Secondary0NodeSubnetId == "" {
				return config, errors.New("Config error: must provide secondary 0 node subnet ID")
			}
			if service.Secondary1NodeSubnetId == "" {
				return config, errors.New("Config error: must provide secondary 1 node subnet ID")
			}
		default:
			return config, errors.New("Config error: service name " + service.Name + " not recognised")
		}

		if len(service.Plans) == 0 {
			return config, errors.New("Config error: at least one plan must be configured for service " + service.Name)
		}
	}

	return config, nil
}

func findServiceById(id string, catalog *Catalog) (Service, error) {
	for _, service := range catalog.Services {
		if service.ID == id {
			return service, nil
		}
	}
	return Service{}, errors.New("could not find service with id " + id)
}

func findPlanById(id string, service Service) (Plan, error) {
	for _, plan := range service.Plans {
		if plan.ID == id {
			return plan, nil
		}
	}
	return Plan{}, errors.New("could not find plan with id " + id)
}
