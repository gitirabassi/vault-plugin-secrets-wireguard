resource "vault_mount" "wireguard" {
  path = var.backend_mount_path
  type = "wireguard"
}

locals {
  config = {
    cidr = var.wireguard_cidr
  }
}

resource "vault_generic_endpoint" "config" {
  depends_on           = [vault_mount.wireguard]
  path                 = "${vault_mount.wireguard.path}/config"
  ignore_absent_fields = true
  data_json            = jsonencode(local.config)
}

resource "vault_generic_endpoint" "servers" {
  depends_on = [vault_mount.wireguard]
  for_each   = var.servers

  path                 = "${vault_mount.wireguard.path}/servers/${each.key}"
  ignore_absent_fields = true
  data_json            = <<EOT
{
  "address": "${each.value.address}",
  "port": "${each.value.port}",
}
EOT
}
