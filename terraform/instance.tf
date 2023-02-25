# Create a security group that allows SSH access
resource "aws_security_group" "bastion_secruity_group" {
  vpc_id = aws_vpc.vpc.id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_key_pair" "admin" {
  key_name   = "admin"
  public_key = var.public_key
}

resource "aws_instance" "bastion" {
  ami                         = "ami-0822295a729d2a28e"
  instance_type               = "t2.nano"
  associate_public_ip_address = true
  key_name                    = aws_key_pair.admin.key_name
  subnet_id                   = aws_subnet.public_subnet.id
  vpc_security_group_ids      = [aws_security_group.bastion_secruity_group.id]

  connection {
    type        = "ssh"
    user        = "ubuntu"
    private_key = file(var.private_key_path)
    host        = self.public_ip
  }

  depends_on = [
    aws_db_instance.db
  ]

  provisioner "file" {
    source      = "../deploy/mysql/schema.sql"
    destination = "/home/ubuntu/schema.sql"
  }

  provisioner "remote-exec" {
    inline = [
      "nslookup ${local.db_host}  && echo 'DNS lookup successful'",
      "sudo apt update",
      "sudo apt -y install mysql-client",
      "mysql -h ${local.db_host} -u ${var.db_username} -p${var.db_password} ${aws_db_instance.db.db_name} < schema.sql"
    ]
  }
}
