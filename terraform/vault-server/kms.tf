resource "aws_kms_key" "main" {
  description             = var.name
  deletion_window_in_days = 30
  enable_key_rotation     = true
}

