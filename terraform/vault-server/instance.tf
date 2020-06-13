resource "aws_eip" "main" {
  vpc        = true
  depends_on = [aws_internet_gateway.main]

  tags = {
    Name = var.name
  }
}

resource "aws_instance" "main" {
  ami                    = var.instance_ami != "" ? var.instance_ami : data.aws_ami.flatcar.id
  instance_type          = var.instance_type
  user_data              = data.ignition_config.main.rendered
  subnet_id              = aws_subnet.main.id
  iam_instance_profile   = aws_iam_instance_profile.main.name
  vpc_security_group_ids = [aws_security_group.main.id]
  key_name               = var.ssh_key_name
  ebs_optimized          = true
  tags = {
    Name = var.name
  }
}

resource "aws_eip_association" "eip_assoc" {
  instance_id   = aws_instance.main.id
  allocation_id = aws_eip.main.id
}
