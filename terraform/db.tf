resource "aws_security_group" "db_security_group" {
  name   = "db-security-group"
  vpc_id = aws_vpc.vpc.id
  ingress {
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    security_groups = [aws_security_group.bastion_secruity_group.id, aws_security_group.eks_security_group.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_db_subnet_group" "private_db_subnet_group" {
  name       = "private_db_subnet_group"
  subnet_ids = [aws_subnet.private_subnet_1.id, aws_subnet.private_subnet_2.id]
}

resource "aws_db_instance" "db" {
  allocated_storage    = 5
  availability_zone    = "ap-northeast-1a"
  db_name              = "APP"
  engine               = "mysql"
  engine_version       = "8.0"
  instance_class       = "db.t2.micro"
  username             = var.db_username
  password             = var.db_password
  parameter_group_name = "default.mysql8.0"
  storage_type         = "standard"
  skip_final_snapshot  = true

  db_subnet_group_name   = aws_db_subnet_group.private_db_subnet_group.name
  vpc_security_group_ids = [aws_security_group.db_security_group.id]
}

locals {
  db_host = split(":", aws_db_instance.db.endpoint)[0]
}
