resource "aws_vpc" "vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
}

resource "aws_subnet" "public_subnet_1" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "ap-northeast-1a"
}

resource "aws_subnet" "public_subnet_2" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = "10.0.2.0/24"
  availability_zone = "ap-northeast-1c"
}

resource "aws_subnet" "private_subnet_1" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = "10.0.3.0/24"
  availability_zone = "ap-northeast-1a"
}

resource "aws_subnet" "private_subnet_2" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = "10.0.4.0/24"
  availability_zone = "ap-northeast-1c"
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.vpc.id
}

resource "aws_eip" "eip" {
  vpc                       = true
  associate_with_private_ip = "10.0.1.12"
  depends_on                = [aws_internet_gateway.igw]
}

resource "aws_nat_gateway" "ngw" {
  allocation_id = aws_eip.eip.id
  subnet_id     = aws_subnet.public_subnet_1.id
  depends_on    = [aws_internet_gateway.igw]
}

resource "aws_route_table" "public_rt" {
  vpc_id = aws_vpc.vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }
}

resource "aws_route_table" "private_rt" {
  vpc_id = aws_vpc.vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_nat_gateway.ngw.id
  }
}

resource "aws_route_table_association" "public_rta1" {
  subnet_id      = aws_subnet.public_subnet_1.id
  route_table_id = aws_route_table.public_rt.id
}

resource "aws_route_table_association" "public_rta2" {
  subnet_id      = aws_subnet.public_subnet_2.id
  route_table_id = aws_route_table.public_rt.id
}

resource "aws_route_table_association" "private_rta1" {
  subnet_id      = aws_subnet.private_subnet_1.id
  route_table_id = aws_route_table.private_rt.id
}

resource "aws_route_table_association" "private_rta2" {
  subnet_id      = aws_subnet.private_subnet_2.id
  route_table_id = aws_route_table.private_rt.id
}
