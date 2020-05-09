data "ignition_config" "main" {
  systemd = [
    data.ignition_systemd_unit.caddy.rendered,
    data.ignition_systemd_unit.vault.rendered,
  ]
  files = [
    data.ignition_file.vault.rendered,
    data.ignition_file.caddy.rendered,
  ]
}

locals {
  server = {
    vault_address  = var.vault_address
    kms_arn        = aws_kms_key.main.arn
    dynamodb_table = aws_dynamodb_table.main.name
    region         = var.region
  }
}

data "ignition_file" "caddy" {
  filesystem = "root"
  path       = "/opt/caddy-conf/Caddyfile"
  mode       = "0664"
  content {
    content = var.auto_tls ? templatefile("${path.module}/ignition/Caddyfile", local.server) : templatefile("${path.module}/ignition/Caddyfile_notls", local.server)
  }
}

data "ignition_file" "vault" {
  filesystem = "root"
  path       = "/opt/vault/server.hcl"
  mode       = 0755
  content {
    content = templatefile("${path.module}/ignition/server.hcl", local.server)
  }
}

data "ignition_systemd_unit" "vault" {
  name    = "vault.service"
  content = file("${path.module}/ignition/vault.service")
}

data "ignition_systemd_unit" "caddy" {
  name    = "caddy.service"
  content = file("${path.module}/ignition/caddy.service")
}
