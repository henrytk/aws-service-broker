# MongoDB Integration Tests

## Requirements

* Your AWS credentials loaded
* Terraform installed
* An SSH Key Pair in the region defined in the Terraform provider

## Instructions

* Apply the Terraform configuration. It will output variables you need to set up environment variables.

  ```
  cd integration_tests/mongodb/terraform
  terraform apply .
  ```

* In the VPC created above, create a security group for a bastion host and note its ID. You don't need to add any inbound rules or apply it to anything.
* Set up the following environment variables:
  * `ASB_KEY_PAIR`: The name of the SSH key pair which will be used for accessing EC2 instances created by the tests.
  * `ASB_PRIMARY_NODE`: Set to `private_subnet_1` from Terraform outputs
  * `ASB_SECONDARY_0_NODE`: Set to `private_subnet_2` from Terraform outputs
  * `ASB_SECONDARY_1_NODE`: Set to `private_subnet_3` from Terraform outputs
  * `ASB_MONGODB_ADMIN_PASSWORD`: any non-empty value
  * `ASB_VPC_ID`: Set to `vpc_id` from Terraform outputs
  * `ASB_BASTION_SECURITY_GROUP`: The ID of the security group you created manually.
* Run the tests:

  ```
  cd integration_tests/mongodb
  ginkgo .
  ```
* Delete the bastion security group you created
* Run Terraform destroy
