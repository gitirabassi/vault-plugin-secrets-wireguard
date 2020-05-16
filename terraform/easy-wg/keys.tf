data "external" "server" {
  program = ["sh", "${path.module}/genkeys.sh"]
}

data "external" "client" {
  count   = var.num_clients
  program = ["bash", "${path.module}/genkeys.sh"]
}
