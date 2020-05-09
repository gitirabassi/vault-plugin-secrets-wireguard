
variable vpc_cidr {
  type    = string
  default = "10.240.0.0/16"
}

variable name {
  type    = string
  default = "vault-server"
}

variable enable_ssh {
  type    = bool
  default = true
}

variable "ssh_key_name" {
  type        = string
  description = "The AWS ssh key name to install in the instance"
}

data "aws_ami" "flatcar" {
  most_recent = true
  owners      = ["075585003325"]

  filter {
    name   = "name"
    values = ["Flatcar-stable-*"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

data "aws_availability_zones" "available" {}