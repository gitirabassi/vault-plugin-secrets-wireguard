variable vpc_cidr {
  type    = string
  default = "10.240.0.0/16"
}

variable name {
  type    = string
  default = "simple-wireguard-server"
}

variable wireguard_server_port {
  type    = number
  default = 51820
}

variable enable_ssh {
  type    = bool
  default = true
}

variable "ssh_key_name" {
  type        = string
  description = "The AWS ssh key name to install in the instance"
}

data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"]

  filter {
    name   = "name"
    values = ["ubuntu-minimal/images/hvm-ssd/ubuntu-focal-20.04-amd64-minimal-*"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }
}

data "aws_availability_zones" "available" {}

variable instance_type {
  type    = string
  default = "t3.micro"
}

variable num_clients {
  type    = number
  default = 1
}

variable allowed_ips {
  type    = string
  default = "0.0.0.0/0"
}

variable wireguard_gateway_address {
  type    = string
  default = "172.16.0.1"
}

variable wireguard_gateway_netmask {
  type    = string
  default = "24"
}

variable main_interface {
  type    = string
  default = "ens5"
}
