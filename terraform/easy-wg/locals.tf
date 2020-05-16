locals {
  client = {
    server_public_endpoint = aws_eip.main.public_ip
    server_port            = var.wireguard_server_port
    server_public_key      = data.external.server.result.public
    allowed_ips            = var.allowed_ips
    dns                    = var.wireguard_gateway_address
  }

  temp_conf      = templatefile("${path.module}/ignition/clients.conf", local.client)
  server_cidr    = "${var.wireguard_gateway_address}/${var.wireguard_gateway_netmask}"
  client_ips     = [for count in range(var.num_clients) : cidrhost(local.server_cidr, tonumber(count) + 2)]
  client_configs = formatlist(local.temp_conf, local.client_ips, flatten(data.external.client.*.result.private))
  peers          = formatlist(file("${path.module}/ignition/peer.conf"), local.client_ips, flatten(data.external.client.*.result.public))

  server = {
    wireguard_gateway_address = var.wireguard_gateway_address
    wireguard_gateway_netmask = var.wireguard_gateway_netmask
    server_private_key        = data.external.server.result.private
    server_port               = var.wireguard_server_port
    main_interface_name       = var.main_interface
    peers                     = join("", local.peers)
  }
}
