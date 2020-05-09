data "ignition_config" "main" {
  systemd = [
    data.ignition_systemd_unit.wg-server-agent.rendered,
    data.ignition_systemd_unit.vault.rendered,
  ]
  files = [
    data.ignition_file.vault-agent.rendered,
  ]
}

locals {
  server = {
    vault_address = var.vault_address
    aws_role_name = var.aws_role_name
  }
}

data "ignition_file" "vault-agent" {
  filesystem = "root"
  path       = "/opt/conf/vault-agent.hcl"
  mode       = "0664"
  content {
    content = templatefile("${path.module}/ignition/vault-agent.hcl", local.server)
  }
}

data "ignition_systemd_unit" "vault" {
  name    = "vault-agent.service"
  content = file("${path.module}/ignition/vault-agent.service")
}

data "ignition_systemd_unit" "wg-server-agent" {
  name    = "wg-server-agent.service"
  content = file("${path.module}/ignition/wg-server-agent.service")
}
