resource "aws_iam_role" "main" {
  name               = var.name
  assume_role_policy = data.aws_iam_policy_document.main.json
}

resource "aws_iam_instance_profile" "main" {
  name = var.name
  role = aws_iam_role.main.name
}

data "aws_iam_policy_document" "main" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy" "storage" {
  name   = "${var.name}-storage"
  role   = aws_iam_role.main.id
  policy = data.aws_iam_policy_document.storage.json
}

data "aws_iam_policy_document" "storage" {
  statement {
    sid = "VaultStorage"
    actions = [
      "dynamodb:DescribeLimits",
      "dynamodb:DescribeTimeToLive",
      "dynamodb:ListTagsOfResource",
      "dynamodb:DescribeReservedCapacityOfferings",
      "dynamodb:DescribeReservedCapacity",
      "dynamodb:ListTables",
      "dynamodb:BatchGetItem",
      "dynamodb:BatchWriteItem",
      "dynamodb:CreateTable",
      "dynamodb:DeleteItem",
      "dynamodb:GetItem",
      "dynamodb:GetRecords",
      "dynamodb:PutItem",
      "dynamodb:Query",
      "dynamodb:UpdateItem",
      "dynamodb:Scan",
      "dynamodb:DescribeTable",
    ]
    resources = [
      aws_dynamodb_table.main.arn,
    ]
  }
}

resource "aws_iam_role_policy" "kms" {
  name   = "${var.name}-kms"
  role   = aws_iam_role.main.id
  policy = data.aws_iam_policy_document.kms.json
}

data "aws_iam_policy_document" "kms" {
  statement {
    sid = "VaultKMSSeal"
    actions = [
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:DescribeKey",
    ]
    resources = [
      aws_kms_key.main.arn,
    ]
  }
}

resource "aws_iam_role_policy" "aws-auth" {
  name   = "${var.name}-aws-auth"
  role   = aws_iam_role.main.id
  policy = data.aws_iam_policy_document.aws-auth.json
}

data "aws_iam_policy_document" "aws-auth" {
  statement {
    sid = ""
    actions = [
      "ec2:DescribeInstances",
      "iam:GetInstanceProfile",
      "iam:GetUser",
      "iam:GetRole",
    ]
    resources = [
      "*",
    ]
  }

}
