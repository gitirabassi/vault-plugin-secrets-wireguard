

provider "aws" {
  region = "eu-central-1"
}

data "http" "ssh_key" {
  url = "https://github.com/someone.keys"
}

resource "aws_key_pair" "main" {
  key_name   = "my_super_key"
  public_key = data.http.ssh_key.body
}

variable vault_address {
  default = "https://vault.example.com"
}

module "first-server" {
  source              = "./server"
  name                = "wireguard-server"
  vault_address       = var.vault_addres
  aws_role_name       = "wireguard-server"
  webhook_source_cidr = "0.0.0.0/0"
  instance_type       = "t3.small"
  ssh_key_name        = aws_key_pair.main.name
  enable_ssh          = true
}

provider "vault" {
  address = var.vault_address
}

module "vault-configuration" {
  source             = "./vault-config"
  backend_mount_path = "wireguard"
  wireguard_cidr     = "172.12.0.0/16"
  servers = {
    "default" = {
      address  = module.first-server.public_ip
      port     = module.first-server.public_port
      role_arn = module.first-server.role_arn
      vpc_id   = module.first-server.vpc_id
      region   = "eu-central-1"
    },
  }
}
