variable vault_address {
  type        = string
  description = "Address of vault server (e.g https://vault.example.com)"
}

variable aws_role_name {
  type = string

}

variable vpc_cidr {
  type    = string
  default = "10.240.0.0/16"
}

variable name {
  type    = string
  default = "wireguard-server"
}

variable wireguard_server_port {
  type    = number
  default = 51820
}

variable wireguard_webook_port {
  type    = number
  default = 51821
}

variable webhook_source_cidr {
  type    = string
  default = "0.0.0.0/0"
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
    values = ["Flatcar-edge-*"]
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

data "aws_availability_zones" "available" {
}
