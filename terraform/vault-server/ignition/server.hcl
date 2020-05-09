api_addr     = "${vault_address}"
cluster_addr = "http://127.0.0.1:8201"

listener "tcp" {
  address     = "127.0.0.1:8200"
  tls_disable = "true"
}

ui = true
storage "dynamodb" {
  ha_enabled = "true"
  table      = "${dynamodb_table}"
  region     = "${region}"
}
seal "awskms" {
  kms_key_id = "${kms_id}"
}
cluster_name = "${name}"
