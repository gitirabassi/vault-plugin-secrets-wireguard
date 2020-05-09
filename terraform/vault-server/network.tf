resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "wireguard-server"
  }
}

resource "aws_subnet" "main" {
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  vpc_id                  = aws_vpc.main.id
  cidr_block              = cidrsubnet(var.vpc_cidr, 8, count.index + var.public_subnet_index)
  map_public_ip_on_launch = true

  tags = {
    Name              = "Public-${var.cluster_name}-${count.index}"
    ClusterName       = var.cluster_name
    KubernetesCluster = var.cluster_name
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name              = var.cluster_name
    Cluster           = var.cluster_name
    KubernetesCluster = var.cluster_name
  }
}


resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id
  tags = {
    Name              = "public-${var.cluster_name}"
    ClusterName       = var.cluster_name
    KubernetesCluster = var.cluster_name
  }
}

resource "aws_route_table_association" "public" {
  count          = var.az-count
  subnet_id      = element(aws_subnet.public.*.id, count.index)
  route_table_id = aws_route_table.public.id
}

resource "aws_route" "public-internet" {
  route_table_id         = aws_route_table.public.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.main.id
  depends_on = [
    aws_route_table.public
  ]
}
