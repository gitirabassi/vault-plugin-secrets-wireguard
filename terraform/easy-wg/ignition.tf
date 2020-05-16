data "ignition_config" "main" {
  systemd = [
    data.ignition_systemd_unit.wireguard.rendered,
    data.ignition_systemd_unit.coredns.rendered,
  ]
  directories = [
    data.ignition_directory.wireguard.rendered,
  ]
  files = [
    data.ignition_file.wireguard.rendered,
    data.ignition_file.coredns.rendered,
    data.ignition_file.clients.rendered,
  ]
}

data "ignition_directory" "wireguard" {
  filesystem = "root"
  path       = "/etc/wireguard"
  mode       = "0755"
}

data "ignition_file" "wireguard" {
  filesystem = "root"
  path       = "/etc/wireguard/wg0.conf"
  mode       = 644
  content {
    content = templatefile("${path.module}/ignition/wg0.conf", local.server)
  }
}

data "ignition_file" "coredns" {
  filesystem = "root"
  path       = "/opt/Corefile"
  mode       = 644
  content {
    content = templatefile("${path.module}/ignition/Corefile", local.server)
  }
}

data "ignition_file" "clients" {
  filesystem = "root"
  path       = "/opt/clients.conf"
  mode       = 644
  content {
    content = join("----\n", local.client_configs)
  }
}

data "ignition_systemd_unit" "wireguard" {
  name    = "wireguard.service"
  content = file("${path.module}/ignition/wireguard.service")
}

data "ignition_systemd_unit" "coredns" {
  name    = "coredns.service"
  content = file("${path.module}/ignition/coredns.service")
}
