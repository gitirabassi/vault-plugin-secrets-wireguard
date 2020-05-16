# Vault-secret-plugin-wireguard

## Paths

- GET/POST/DELETE `/config`
  - `cidr` --> `10.20.0.0/24`
- GET/POST/DELETE `/servers/:server_name`
  - `port` --> defaults to `51820`
  - `public_endpoint` -->
  - `post_up_script` 
  - `post_down_script`
  - `private_webhook_address` --> defaults to `public_endpoint` in case it's not specified
  - `allowed_ips` --> list of subnets to route traffic trhu this server can be `["0.0.0.0/0"]` or `["10.0.0.0/24","192.68.0.0/24"]`
- GET/POST/DELETE `/roles/:role_name`
  - `servers` --> list of servers to connect to `["default", "aws-europe", "gcp-testing"]` --> must exist `:server_name`
  - `dns` 
  - `client_persistent_keepalive` --> defaults to 25 seconds
  - `client_subnet_mask` --> defaults to 32
- GET `/creds/:role_name`
  - `conf` --> complete wireguard configuration to be used with wg-quick for a client
- GET `/server-creds/:server_name`
  - `conf` --> complete wireguard configuration to be used with wg-quick for a client
  - `webhook_secret` --> webhook secret that vault will use to POST updates to wireguard servers

## Apply policy
```shell
vault policy write wireguard_client_develper contrib/client.hcl
```

## Client User Experience

```shell
export VAULT_ADDR=https://vault.example.com
vault status
vault login -method=oidc
vault read -field=conf wireguard/crets/default |clipcopy
```

## Terraform deployment

> IMPORTANT: this cannot be applied all at once as it will breack. there is a order:

- First create Vault server without TLS
- Configure DNS such that will resolve your domain to the host
- enable TLS
- Create Wireguard server disabling agent
- Configure vault with right servers and configurations

2 main changes need to happen to this example:
- ssh key: if you use github to distribute your public ssh key, please add your account name instead of `someone` in the github link
- change the `vault_address` and `module.dns.domain` to match your domain

The modules are opinionated:
- both Vault and Wireguard server create and live in their own VPC
- these VPCs are dedidcated to running them, and them only. 
- You shuould enable Aws Ec2 Transit Gateways or VPC to connect the Wireguard Server with your VPC.
- THis way you'll get much more control of what goes where and you may have different VPCs connect to the Wireguard VPC

> A video recording on how to do all this will come very soon

```terraform
provider "aws" {
  region = "eu-central-1"
}

data "http" "ssh_key" {
  url = "https://github.com/someone.keys"
}

resource "aws_key_pair" "main" {
  key_name   = "wireguard_infra_key"
  public_key = data.http.ssh_key.body
}

variable vault_address {
  default = "vault.example.com"
}

module "dns" {
  source = "github.com/gitirabassi/vault-plugin-secrets-wireguard//terraform/route53"
  domain = "example.com"
  a_records = {
    "vault"   = module.vault-server.public_ip
    "first-wireguard-server" = module.wireguard-server.public_ip
  }
}

module "vault-server" {
  source        = "github.com/gitirabassi/vault-plugin-secrets-wireguard//terraform/vault-server"
  name          = "vault-server"
  vpc_cidr      = "10.210.0.0/16"
  vault_address = var.vault_address
  instance_type = "t3.small"
  region        = "eu-central-1"
  ssh_key_name  = aws_key_pair.main.key_name
  enable_ssh    = false
  auto_tls      = true
}

module "wireguard-server" {
  source              = "github.com/gitirabassi/vault-plugin-secrets-wireguard//terraform/wg-server"
  name                = "wireguard-server"
  vpc_cidr            = "10.220.0.0/16"
  vault_address       = "https://${var.vault_address}"
  vault_role_name     = "wireguard-server"
  webhook_source_cidr = "${module.vault-server.public_ip}/32"
  instance_type       = "t3.small"
  ssh_key_name        = aws_key_pair.main.key_name
  enable_ssh          = true
  disable_agent       = true
}

provider "vault" {
  address = "https://${var.vault_address}"
}

module "vault-configuration" {
  source             = "github.com/gitirabassi/vault-plugin-secrets-wireguard//terraform/vault-config"
  backend_mount_path = "wireguard"
  wireguard_cidr     = "172.16.0.0/16"
  servers = {
    "wireguard-server" = {
      address  = module.wireguard-server.public_ip
      port     = module.wireguard-server.public_port
      role_arn = module.wireguard-server.role_arn
      vpc_id   = module.wireguard-server.vpc_id
      region   = "eu-central-1"
    },
  }
}
```


## Webhook

The webook is a trick to not make the wireguard poll every X secods but to reload it's configuration only when a user gets added or deleted

To simulate the hook that Vault will send to the `wg-server-agent` curl can be used

```shell
curl -XPOST -H 'Content-Type: application/json' -d '{"token":"example-super-secret-token"}' http://dev.aws.example.com:51821/webhook
```


## Future features
- rotate server keys (find a ways to use multiple keys to make migration smoother)
