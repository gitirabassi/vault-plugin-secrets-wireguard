resource "aws_eip" "main" {
  vpc        = true
  depends_on = [aws_internet_gateway.main]

  tags = {
    Name = var.name
  }
}

resource "aws_instance" "main" {
  ami                    = data.aws_ami.ubuntu.id
  instance_type          = var.instance_type
  user_data              = templatefile("${path.module}/files/init.sh", local.server)
  subnet_id              = aws_subnet.main.id
  iam_instance_profile   = aws_iam_instance_profile.main.name
  vpc_security_group_ids = [aws_security_group.main.id]
  key_name               = var.ssh_key_name
  ebs_optimized          = true
  tags = {
    Name = var.name
  }
  # lifecycle {
  #   ignore_changes = [user_data]
  # }
}

resource "aws_eip_association" "eip_assoc" {
  instance_id   = aws_instance.main.id
  allocation_id = aws_eip.main.id
}
