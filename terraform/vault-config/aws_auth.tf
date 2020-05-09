resource "vault_auth_backend" "aws" {
  type = "aws"
}

resource "vault_aws_auth_backend_role" "main" {
  for_each             = var.servers
  backend              = vault_auth_backend.aws.path
  role                 = each.key
  auth_type            = "iam"
  bound_vpc_ids        = [each.value.vpc_id]
  bound_iam_role_arns  = [each.value.role_arn]
  inferred_entity_type = "ec2_instance"
  inferred_aws_region  = each.value.region
  token_ttl            = 24 * 3600
  token_max_ttl        = 72 * 3600
  token_policies       = ["wireguard-server-${each.key}"]
}

resource "vault_policy" "main" {
  for_each = var.servers
  name     = "wireguard-server-${each.key}"
  policy   = templatefile("${path.module}/policies/server.hcl", { server_name = each.key })
}
