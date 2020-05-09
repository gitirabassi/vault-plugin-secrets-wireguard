exit_after_auth = false

vault {
  address = "${vault_address}"
}

auto_auth {
  method "aws" {
    mount_path = "auth/aws"
    config = {
      type = "iam"
      role = "${vault_role_name}"
    }
  }
}

cache {
  use_auto_auth_token = true
}

listener "unix" {
  address     = "/opt/vault.sock"
  tls_disable = true
}
