variable "env" {
  description = "The environment name"
}

variable "zones" {
  description = "AWS availability zones"

  default = {
    zone0 = "eu-west-1a"
    zone1 = "eu-west-1b"
    zone2 = "eu-west-1c"
  }
}

variable "public_subnet_cidrs" {
  description = "CIDR for public subnets"

  default = {
    zone0 = "10.0.0.0/24"
    zone1 = "10.0.1.0/24"
    zone2 = "10.0.2.0/24"
  }
}

variable "private_subnet_cidrs" {
  description = "CIDR for private subnets"

  default = {
    zone0 = "10.0.150.0/24"
    zone1 = "10.0.151.0/24"
    zone2 = "10.0.152.0/24"
  }
}
