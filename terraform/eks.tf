resource "aws_security_group" "eks_security_group" {
  name   = "eks-security-group"
  vpc_id = aws_vpc.vpc.id

  ingress {
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    security_groups = [aws_security_group.bastion_secruity_group.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "19.5.1"

  cluster_name    = "cluster"
  cluster_version = "1.24"

  vpc_id                         = aws_vpc.vpc.id
  subnet_ids                     = [aws_subnet.private_subnet_1.id, aws_subnet.private_subnet_2.id]
  cluster_endpoint_public_access = true

  eks_managed_node_group_defaults = {
    ami_type               = "AL2_x86_64"
    instance_types         = ["t3.small"]
    min_size               = 1
    max_size               = 2
    desired_size           = 1
    key_name               = aws_key_pair.admin.key_name
    vpc_security_group_ids = [aws_security_group.eks_security_group.id]
  }

  eks_managed_node_groups = {
    one = {
      name = "node-group-1"
    }

    two = {
      name = "node-group-2"
    }
  }
}
