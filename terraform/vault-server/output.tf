output role_arn {
  value = aws_iam_role.main.arn
}

output public_ip {
  value = aws_eip.main.public_ip
}
