package templates

var MongoDBStack = []byte(`
{
    "AWSTemplateFormatVersion": "2010-09-09",
    "Description": "(000F) Deploy MongoDB Replica Set on AWS (Existing VPC)",
    "Metadata": {
        "AWS::CloudFormation::Interface": {
            "ParameterGroups": [
                {
                    "Label": {
                        "default": "Network Configuration"
                    },
                    "Parameters": [
                        "VPC",
                        "PrimaryNodeSubnet",
                        "Secondary0NodeSubnet",
                        "Secondary1NodeSubnet",
                        "BastionSecurityGroupID"
                    ]
                },
                {
                    "Label": {
                        "default": "Security Configuration"
                    },
                    "Parameters": [
                        "KeyPairName"
                    ]
                },
                {
                    "Label": {
                        "default": "MongoDB Database Configuration"
                    },
                    "Parameters": [
                        "ClusterReplicaSetCount",
                        "Iops",
                        "MongoDBVersion",
                        "MongoDBAdminUsername",
                        "MongoDBAdminPassword",
                        "NodeInstanceType",
                        "ReplicaShardIndex",
                        "VolumeSize",
                        "VolumeType"
                    ]
                },
                {
                    "Label": {
                        "default": "AWS Quick Start Configuration"
                    },
                    "Parameters": [
                        "QSS3BucketName",
                        "QSS3KeyPrefix"
                    ]
                }
            ],
            "ParameterLabels": {
                "BastionSecurityGroupID": {
                    "default": "Bastion Security Group ID"
                },
                "ClusterReplicaSetCount": {
                    "default": "Cluster Replica Set Count"
                },
                "Iops": {
                    "default": "Iops"
                },
                "KeyPairName": {
                    "default": "Key Pair Name"
                },
                "MongoDBAdminPassword": {
                    "default": "MongoDB Admin Password"
                },
                "MongoDBAdminUsername": {
                    "default": "MongoDB Admin Username"
                },
                "MongoDBVersion": {
                    "default": "MongoDB Version"
                },
                "NodeInstanceType": {
                    "default": "Node Instance Type"
                },
                "PrimaryNodeSubnet": {
                    "default": "Primary Node Subnet"
                },
                "QSS3BucketName": {
                    "default": "Quick Start S3 Bucket Name"
                },
                "QSS3KeyPrefix": {
                    "default": "Quick Start S3 Key Prefix"
                },
                "ReplicaShardIndex": {
                    "default": "Replica Shard Index"
                },
                "Secondary0NodeSubnet": {
                    "default": "Secondary0 Node Subnet"
                },
                "Secondary1NodeSubnet": {
                    "default": "Secondary1 Node Subnet"
                },
                "VPC": {
                    "default": "VPC"
                },
                "VolumeSize": {
                    "default": "Volume Size"
                },
                "VolumeType": {
                    "default": "Volume Type"
                }
            }
        }
    },
    "Parameters": {
        "BastionSecurityGroupID": {
            "Description": "ID of the Bastion Security Group (e.g., sg-7f16e910)",
            "Type": "AWS::EC2::SecurityGroup::Id"
        },
        "ClusterReplicaSetCount": {
            "Description": "Number of Replica Set Members. Choose 1 or 3",
            "Type": "String",
            "Default": "1",
            "AllowedValues": [
                "1",
                "3"
            ]
        },
        "MongoDBVersion": {
            "Description": "MongoDB version",
            "Type": "String",
            "Default": "3.4",
            "AllowedValues": [
                "3.4",
                "3.2"
            ]
        },
        "MongoDBAdminUsername": {
            "Default": "admin",
            "NoEcho": "true",
            "Description": "MongoDB admin account username",
            "Type": "String",
            "MinLength": "1",
            "MaxLength": "16",
            "AllowedPattern": "[a-zA-Z][a-zA-Z0-9]*",
            "ConstraintDescription": "must begin with a letter and contain only alphanumeric characters."
        },
        "MongoDBAdminPassword": {
            "AllowedPattern": "([A-Za-z0-9_@-]{8,32})",
            "ConstraintDescription": "Input your MongoDB database password, Min 8, Maximum of 32 characters. . Allowed characters are: [A-Za-z0-9_@-]",
            "Description": "Enter your MongoDB Database Password, Min 8, maximum of 32 characters.",
            "NoEcho": "true",
            "Type": "String"
        },
        "ReplicaShardIndex": {
            "Description": "Shard Index of this replica set",
            "Type": "String",
            "Default": "0"
        },
        "QSS3BucketName": {
            "AllowedPattern": "^[0-9a-zA-Z]+([0-9a-zA-Z-]*[0-9a-zA-Z])*$",
            "Default": "quickstart-reference",
            "Type": "String",
            "ConstraintDescription": "Quick Start bucket name can include numbers, lowercase letters, uppercase letters, and hyphens (-). It cannot start or end with a hyphen (-).",
            "Description": "S3 bucket name for the Quick Start assets. Quick Start bucket name can include numbers, lowercase letters, uppercase letters, and hyphens (-). It cannot start or end with a hyphen (-)."
        },
        "QSS3KeyPrefix": {
            "AllowedPattern": "^[0-9a-zA-Z-/]*$",
            "Default": "mongodb/latest/",
            "Type": "String",
            "ConstraintDescription": "Quick Start key prefix can include numbers, lowercase letters, uppercase letters, hyphens (-), and forward slash (/).",
            "Description": "S3 key prefix for the Quick Start assets. Quick Start key prefix can include numbers, lowercase letters, uppercase letters, hyphens (-), and forward slash (/). It cannot start or end with a hyphen (-)."
        },
        "KeyPairName": {
            "Type": "AWS::EC2::KeyPair::KeyName",
            "Description": "Name of an existing EC2 KeyPair. MongoDB instances will launch with this KeyPair."
        },
        "VolumeSize": {
            "Type": "String",
            "Description": "EBS Volume Size (data) to be attached to node in GBs",
            "Default": "400"
        },
        "VolumeType": {
            "Type": "String",
            "Description": "EBS Volume Type (data) to be attached to node in GBs [io1,gp2]",
            "Default": "gp2",
            "AllowedValues": [
                "gp2",
                "io1"
            ]
        },
        "Iops": {
            "Type": "String",
            "Description": "Iops of EBS volume when io1 type is chosen. Otherwise ignored",
            "Default": "100"
        },
        "NodeInstanceType": {
            "Description": "Amazon EC2 instance type for the MongoDB nodes.",
            "Type": "String",
            "Default": "m4.large",
            "AllowedValues": [
                "m3.medium",
                "m3.large",
                "m3.xlarge",
                "m3.2xlarge",
                "m4.large",
                "m4.xlarge",
                "m4.2xlarge",
                "m4.4xlarge",
                "m4.10xlarge",
                "c3.large",
                "c3.xlarge",
                "c3.2xlarge",
                "c3.4xlarge",
                "c3.8xlarge",
                "r3.large",
                "r3.xlarge",
                "r3.2xlarge",
                "r3.4xlarge",
                "r3.8xlarge",
                "i2.xlarge",
                "i2.2xlarge",
                "i2.4xlarge",
                "i2.8xlarge"
            ]
        },
        "VPC": {
            "Type": "AWS::EC2::VPC::Id",
            "Description": "VPC-ID of your existing Virtual Private Cloud (VPC) where you want to depoy MongoDB cluster.",
            "AllowedPattern": "vpc-[0-9a-z]{8}"
        },
        "PrimaryNodeSubnet": {
            "Type": "AWS::EC2::Subnet::Id",
            "Description": "Subnet-ID the existing subnet in your VPC where you want to deploy Primary node(s).",
            "AllowedPattern": "subnet-[0-9a-z]{8}"
        },
        "Secondary0NodeSubnet": {
            "Type": "AWS::EC2::Subnet::Id",
            "Description": "Subnet-ID the existing subnet in your VPC where you want to deploy Primary node(s).",
            "AllowedPattern": "subnet-[0-9a-z]{8}"
        },
        "Secondary1NodeSubnet": {
            "Type": "AWS::EC2::Subnet::Id",
            "Description": "Subnet-ID the existing subnet in your VPC where you want to deploy Primary node(s).",
            "AllowedPattern": "subnet-[0-9a-z]{8}"
        }
    },
    "Conditions": {
        "CreateThreeReplicaSet": {
            "Fn::Equals": [
                {
                    "Ref": "ClusterReplicaSetCount"
                },
                "3"
            ]
        },
        "GovCloudCondition": {
            "Fn::Equals": [
                {
                    "Ref": "AWS::Region"
                },
                "us-gov-west-1"
            ]
        }
    },
    "Mappings": {
        "AWSAMIRegionMap": {
            "AMI": {
                "AMZNLINUX": "amzn-ami-hvm-2017.09.1.20171120-x86_64-gp2"
            },
            "ap-northeast-1": {
                "AMZNLINUX": "ami-da9e2cbc"
            },
            "ap-northeast-2": {
                "AMZNLINUX": "ami-1196317f"
            },
            "ap-south-1": {
                "AMZNLINUX": "ami-d5c18eba"
            },
            "ap-southeast-1": {
                "AMZNLINUX": "ami-c63d6aa5"
            },
            "ap-southeast-2": {
                "AMZNLINUX": "ami-ff4ea59d"
            },
            "ca-central-1": {
                "AMZNLINUX": "ami-d29e25b6"
            },
            "eu-central-1": {
                "AMZNLINUX": "ami-bf2ba8d0"
            },
            "eu-west-1": {
                "AMZNLINUX": "ami-1a962263"
            },
            "eu-west-2": {
                "AMZNLINUX": "ami-e7d6c983"
            },
            "sa-east-1": {
                "AMZNLINUX": "ami-286f2a44"
            },
            "us-east-1": {
                "AMZNLINUX": "ami-55ef662f"
            },
            "us-east-2": {
                "AMZNLINUX": "ami-15e9c770"
            },
            "us-west-1": {
                "AMZNLINUX": "ami-a51f27c5"
            },
            "us-west-2": {
                "AMZNLINUX": "ami-bf4193c7"
            }
        }
    },
    "Resources": {
        "MongoDBServerAccessSecurityGroup": {
            "Type": "AWS::EC2::SecurityGroup",
            "Properties": {
                "VpcId": {
                    "Ref": "VPC"
                },
                "GroupDescription": "Instances with access to MongoDB servers"
            }
        },
        "MongoDBServerSecurityGroup": {
            "Type": "AWS::EC2::SecurityGroup",
            "Properties": {
                "VpcId": {
                    "Ref": "VPC"
                },
                "GroupDescription": "MongoDB server management and access ports",
                "SecurityGroupIngress": [
                    {
                        "IpProtocol": "tcp",
                        "FromPort": "22",
                        "ToPort": "22",
                        "SourceSecurityGroupId": {
                            "Ref": "BastionSecurityGroupID"
                        }
                    },
                    {
                        "IpProtocol": "tcp",
                        "FromPort": "27017",
                        "ToPort": "27030",
                        "SourceSecurityGroupId": {
                            "Ref": "MongoDBServerAccessSecurityGroup"
                        }
                    },
                    {
                        "IpProtocol": "tcp",
                        "FromPort": "28017",
                        "ToPort": "28017",
                        "SourceSecurityGroupId": {
                            "Ref": "MongoDBServerAccessSecurityGroup"
                        }
                    }
                ]
            }
        },
        "MongoDBServersSecurityGroup": {
            "Type": "AWS::EC2::SecurityGroup",
            "Properties": {
                "VpcId": {
                    "Ref": "VPC"
                },
                "GroupDescription": "MongoDB inter-server communication and management ports",
                "SecurityGroupIngress": [
                    {
                        "IpProtocol": "tcp",
                        "FromPort": "22",
                        "ToPort": "22",
                        "SourceSecurityGroupId": {
                            "Ref": "MongoDBServerSecurityGroup"
                        }
                    },
                    {
                        "IpProtocol": "tcp",
                        "FromPort": "27017",
                        "ToPort": "27030",
                        "SourceSecurityGroupId": {
                            "Ref": "MongoDBServerSecurityGroup"
                        }
                    },
                    {
                        "IpProtocol": "tcp",
                        "FromPort": "28017",
                        "ToPort": "28017",
                        "SourceSecurityGroupId": {
                            "Ref": "MongoDBServerSecurityGroup"
                        }
                    }
                ]
            }
        },
        "MongoDBNodeIAMRole": {
            "Type": "AWS::IAM::Role",
            "Properties": {
                "AssumeRolePolicyDocument": {
                    "Statement": [
                        {
                            "Effect": "Allow",
                            "Principal": {
                                "Service": [
                                    "ec2.amazonaws.com"
                                ]
                            },
                            "Action": [
                                "sts:AssumeRole"
                            ]
                        }
                    ]
                },
                "Path": "/",
                "Policies": [
                    {
                        "PolicyName": "Backup",
                        "PolicyDocument": {
                            "Statement": [
                                {
                                    "Effect": "Allow",
                                    "Action": [
                                        "s3:*",
                                        "ec2:Describe*",
                                        "ec2:AttachNetworkInterface",
                                        "ec2:AttachVolume",
                                        "ec2:CreateTags",
                                        "ec2:CreateVolume",
                                        "ec2:RunInstances",
                                        "ec2:StartInstances",
                                        "ec2:DeleteVolume",
                                        "ec2:CreateSecurityGroup",
                                        "ec2:CreateSnapshot"
                                    ],
                                    "Resource": "*"
                                },
                                {
                                    "Effect": "Allow",
                                    "Action": [
                                        "dynamodb:*",
                                        "dynamodb:Scan",
                                        "dynamodb:Query",
                                        "dynamodb:GetItem",
                                        "dynamodb:BatchGetItem",
                                        "dynamodb:UpdateTable"
                                    ],
                                    "Resource": [
                                        "*"
                                    ]
                                }
                            ]
                        }
                    }
                ]
            }
        },
        "MongoDBNodeIAMProfile": {
            "Type": "AWS::IAM::InstanceProfile",
            "Properties": {
                "Path": "/",
                "Roles": [
                    {
                        "Ref": "MongoDBNodeIAMRole"
                    }
                ]
            }
        },
        "PrimaryReplicaNode0WaitForNodeInstallWaitHandle": {
            "Type": "AWS::CloudFormation::WaitConditionHandle",
            "Properties": {}
        },
        "PrimaryReplicaNode0": {
            "DependsOn": "PrimaryReplicaNode0WaitForNodeInstallWaitHandle",
            "Type": "AWS::CloudFormation::Stack",
            "Properties": {
                "TemplateURL": {
                    "Fn::Sub": [
                        "https://${QSS3BucketName}.${QSS3Region}.amazonaws.com/${QSS3KeyPrefix}templates/mongodb-node.template",
                        {
                            "QSS3Region": {
                                "Fn::If": [
                                    "GovCloudCondition",
                                    "s3-us-gov-west-1",
                                    "s3"
                                ]
                            }
                        }
                    ]
                },
                "Parameters": {
                    "QSS3BucketName": {
                        "Ref": "QSS3BucketName"
                    },
                    "QSS3KeyPrefix": {
                        "Ref": "QSS3KeyPrefix"
                    },
                    "ClusterReplicaSetCount": {
                        "Ref": "ClusterReplicaSetCount"
                    },
                    "Iops": {
                        "Ref": "Iops"
                    },
                    "KeyName": {
                        "Ref": "KeyPairName"
                    },
                    "MongoDBVersion": {
                        "Ref": "MongoDBVersion"
                    },
                    "MongoDBAdminUsername": {
                        "Ref": "MongoDBAdminUsername"
                    },
                    "MongoDBAdminPassword": {
                        "Ref": "MongoDBAdminPassword"
                    },
                    "NodeInstanceType": {
                        "Ref": "NodeInstanceType"
                    },
                    "NodeSubnet": {
                        "Ref": "PrimaryNodeSubnet"
                    },
                    "MongoDBServerSecurityGroupID": {
                        "Ref": "MongoDBServerSecurityGroup"
                    },
                    "MongoDBServersSecurityGroupID": {
                        "Ref": "MongoDBServersSecurityGroup"
                    },
                    "MongoDBNodeIAMProfileID": {
                        "Ref": "MongoDBNodeIAMProfile"
                    },
                    "VPC": {
                        "Ref": "VPC"
                    },
                    "VolumeSize": {
                        "Ref": "VolumeSize"
                    },
                    "VolumeType": {
                        "Ref": "VolumeType"
                    },
                    "StackName": {
                        "Ref": "AWS::StackName"
                    },
                    "ImageId": {
                        "Fn::FindInMap": [
                            "AWSAMIRegionMap",
                            {
                                "Ref": "AWS::Region"
                            },
                            "AMZNLINUX"
                        ]
                    },
                    "ReplicaNodeNameTag": "PrimaryReplicaNode0",
                    "NodeReplicaSetIndex": "0",
                    "ReplicaShardIndex": {
                        "Ref": "ReplicaShardIndex"
                    },
                    "ReplicaNodeWaitForNodeInstallWaitHandle": {
                        "Ref": "PrimaryReplicaNode0WaitForNodeInstallWaitHandle"
                    }
                }
            }
        },
        "PrimaryReplicaNode0WaitForNodeInstall": {
            "Type": "AWS::CloudFormation::WaitCondition",
            "DependsOn": "PrimaryReplicaNode0",
            "Properties": {
                "Handle": {
                    "Ref": "PrimaryReplicaNode0WaitForNodeInstallWaitHandle"
                },
                "Timeout": "3600"
            }
        },
        "SecondaryReplicaNode0WaitForNodeInstallWaitHandle": {
            "Type": "AWS::CloudFormation::WaitConditionHandle",
            "Properties": {},
            "Condition": "CreateThreeReplicaSet"
        },
        "SecondaryReplicaNode0": {
            "DependsOn": "SecondaryReplicaNode0WaitForNodeInstallWaitHandle",
            "Condition": "CreateThreeReplicaSet",
            "Type": "AWS::CloudFormation::Stack",
            "Properties": {
                "TemplateURL": {
                    "Fn::Sub": [
                        "https://${QSS3BucketName}.${QSS3Region}.amazonaws.com/${QSS3KeyPrefix}templates/mongodb-node.template",
                        {
                            "QSS3Region": {
                                "Fn::If": [
                                    "GovCloudCondition",
                                    "s3-us-gov-west-1",
                                    "s3"
                                ]
                            }
                        }
                    ]
                },
                "Parameters": {
                    "QSS3BucketName": {
                        "Ref": "QSS3BucketName"
                    },
                    "QSS3KeyPrefix": {
                        "Ref": "QSS3KeyPrefix"
                    },
                    "ClusterReplicaSetCount": {
                        "Ref": "ClusterReplicaSetCount"
                    },
                    "Iops": {
                        "Ref": "Iops"
                    },
                    "KeyName": {
                        "Ref": "KeyPairName"
                    },
                    "MongoDBVersion": {
                        "Ref": "MongoDBVersion"
                    },
                    "MongoDBAdminUsername": {
                        "Ref": "MongoDBAdminUsername"
                    },
                    "MongoDBAdminPassword": {
                        "Ref": "MongoDBAdminPassword"
                    },
                    "NodeInstanceType": {
                        "Ref": "NodeInstanceType"
                    },
                    "NodeSubnet": {
                        "Ref": "Secondary0NodeSubnet"
                    },
                    "MongoDBServerSecurityGroupID": {
                        "Ref": "MongoDBServerSecurityGroup"
                    },
                    "MongoDBServersSecurityGroupID": {
                        "Ref": "MongoDBServersSecurityGroup"
                    },
                    "MongoDBNodeIAMProfileID": {
                        "Ref": "MongoDBNodeIAMProfile"
                    },
                    "VPC": {
                        "Ref": "VPC"
                    },
                    "VolumeSize": {
                        "Ref": "VolumeSize"
                    },
                    "VolumeType": {
                        "Ref": "VolumeType"
                    },
                    "StackName": {
                        "Ref": "AWS::StackName"
                    },
                    "ImageId": {
                        "Fn::FindInMap": [
                            "AWSAMIRegionMap",
                            {
                                "Ref": "AWS::Region"
                            },
                            "AMZNLINUX"
                        ]
                    },
                    "ReplicaNodeNameTag": "SecondaryReplicaNode0",
                    "NodeReplicaSetIndex": "1",
                    "ReplicaShardIndex": {
                        "Ref": "ReplicaShardIndex"
                    },
                    "ReplicaNodeWaitForNodeInstallWaitHandle": {
                        "Ref": "SecondaryReplicaNode0WaitForNodeInstallWaitHandle"
                    }
                }
            }
        },
        "SecondaryReplicaNode0WaitForNodeInstall": {
            "Type": "AWS::CloudFormation::WaitCondition",
            "Condition": "CreateThreeReplicaSet",
            "DependsOn": "SecondaryReplicaNode0",
            "Properties": {
                "Handle": {
                    "Ref": "SecondaryReplicaNode0WaitForNodeInstallWaitHandle"
                },
                "Timeout": "3600"
            }
        },
        "SecondaryReplicaNode1WaitForNodeInstallWaitHandle": {
            "Type": "AWS::CloudFormation::WaitConditionHandle",
            "Properties": {},
            "Condition": "CreateThreeReplicaSet"
        },
        "SecondaryReplicaNode1": {
            "DependsOn": "SecondaryReplicaNode1WaitForNodeInstallWaitHandle",
            "Condition": "CreateThreeReplicaSet",
            "Type": "AWS::CloudFormation::Stack",
            "Properties": {
                "TemplateURL": {
                    "Fn::Sub": [
                        "https://${QSS3BucketName}.${QSS3Region}.amazonaws.com/${QSS3KeyPrefix}templates/mongodb-node.template",
                        {
                            "QSS3Region": {
                                "Fn::If": [
                                    "GovCloudCondition",
                                    "s3-us-gov-west-1",
                                    "s3"
                                ]
                            }
                        }
                    ]
                },
                "Parameters": {
                    "QSS3BucketName": {
                        "Ref": "QSS3BucketName"
                    },
                    "QSS3KeyPrefix": {
                        "Ref": "QSS3KeyPrefix"
                    },
                    "ClusterReplicaSetCount": {
                        "Ref": "ClusterReplicaSetCount"
                    },
                    "Iops": {
                        "Ref": "Iops"
                    },
                    "KeyName": {
                        "Ref": "KeyPairName"
                    },
                    "MongoDBVersion": {
                        "Ref": "MongoDBVersion"
                    },
                    "MongoDBAdminUsername": {
                        "Ref": "MongoDBAdminUsername"
                    },
                    "MongoDBAdminPassword": {
                        "Ref": "MongoDBAdminPassword"
                    },
                    "NodeInstanceType": {
                        "Ref": "NodeInstanceType"
                    },
                    "NodeSubnet": {
                        "Ref": "Secondary1NodeSubnet"
                    },
                    "MongoDBServerSecurityGroupID": {
                        "Ref": "MongoDBServerSecurityGroup"
                    },
                    "MongoDBServersSecurityGroupID": {
                        "Ref": "MongoDBServersSecurityGroup"
                    },
                    "MongoDBNodeIAMProfileID": {
                        "Ref": "MongoDBNodeIAMProfile"
                    },
                    "VPC": {
                        "Ref": "VPC"
                    },
                    "VolumeSize": {
                        "Ref": "VolumeSize"
                    },
                    "VolumeType": {
                        "Ref": "VolumeType"
                    },
                    "StackName": {
                        "Ref": "AWS::StackName"
                    },
                    "ImageId": {
                        "Fn::FindInMap": [
                            "AWSAMIRegionMap",
                            {
                                "Ref": "AWS::Region"
                            },
                            "AMZNLINUX"
                        ]
                    },
                    "ReplicaNodeNameTag": "SecondaryReplicaNode1",
                    "NodeReplicaSetIndex": "2",
                    "ReplicaShardIndex": {
                        "Ref": "ReplicaShardIndex"
                    },
                    "ReplicaNodeWaitForNodeInstallWaitHandle": {
                        "Ref": "SecondaryReplicaNode1WaitForNodeInstallWaitHandle"
                    }
                }
            }
        },
        "SecondaryReplicaNode1WaitForNodeInstall": {
            "Type": "AWS::CloudFormation::WaitCondition",
            "Condition": "CreateThreeReplicaSet",
            "DependsOn": "SecondaryReplicaNode1",
            "Properties": {
                "Handle": {
                    "Ref": "SecondaryReplicaNode1WaitForNodeInstallWaitHandle"
                },
                "Timeout": "3600"
            }
        }
    },
    "Outputs": {
        "PrimaryReplicaNodeIp": {
            "Value": {
                "Fn::GetAtt": [
                    "PrimaryReplicaNode0",
                    "Outputs.NodePrivateIp"
                ]
            },
            "Description": "Private IP Address of Primary Replica Node"
        },
        "SecondaryReplicaNode0Ip": {
            "Value": {
                "Fn::GetAtt": [
                    "SecondaryReplicaNode0",
                    "Outputs.NodePrivateIp"
                ]
            },
            "Description": "Private IP Address of Secondary Replica 0 Node",
            "Condition": "CreateThreeReplicaSet"
        },
        "SecondaryReplicaNode1Ip": {
            "Value": {
                "Fn::GetAtt": [
                    "SecondaryReplicaNode1",
                    "Outputs.NodePrivateIp"
                ]
            },
            "Description": "Private IP Address of Secondary Replica 1 Node",
            "Condition": "CreateThreeReplicaSet"
        },
        "MongoDBServerAccessSecurityGroup": {
            "Value": {
                "Ref": "MongoDBServerAccessSecurityGroup"
            },
            "Description": "MongoDB Access Security Group"
        }
    }
}
`)
