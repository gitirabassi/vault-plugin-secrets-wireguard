resource "vault_auth_backend" "aws" {
  type = "aws"
}

resource "vault_aws_auth_backend_role" "main" {
  for_each = var.servers

  backend              = vault_auth_backend.aws.path
  role                 = "wireguard-server-${each.key}"
  auth_type            = "iam"
  bound_account_ids    = [var.aws_account]
  bound_vpc_ids        = [var.vpc_id]
  bound_iam_role_arns  = [var.master_role_arn != "" ? var.master_role_arn : "arn:aws:iam::${var.aws_account}:role/k8s-master-role-${var.cluster_name}.${var.cluster_domain}"]
  inferred_entity_type = "ec2_instance"
  inferred_aws_region  = var.region
  token_ttl            = 6 * 3600
  token_max_ttl        = 24 * 3600
  token_policies       = ["master-node-${var.cloud_provider}-${var.cluster_name}"]
}

resource "vault_policy" "main" {
  for_each = toset(var.policies)
  name     = each.value
  policy   = file("${path.module}/policies/${each.value}.hcl")
}
