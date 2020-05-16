data "external" "server" {
  program = ["sh", "${path.module}/files/genkeys.sh"]
}

data "external" "client" {
  count   = var.num_clients
  program = ["bash", "${path.module}/files/genkeys.sh"]
}
