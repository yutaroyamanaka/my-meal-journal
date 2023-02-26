module "eks" {
  source                         = "terraform-aws-modules/eks/aws"
  cluster_name                   = "eks-cluster"
  cluster_version                = "1.24"
  subnet_ids                     = [aws_subnet.private_subnet_1.id, aws_subnet.private_subnet_2.id]
  vpc_id                         = aws_vpc.vpc.id
  cluster_endpoint_public_access = true

  eks_managed_node_group_defaults = {
    ami_type = "AL2_x86_64"
  }
  eks_managed_node_groups = {
    one = {
      name           = "node-group-1"
      instance_types = ["t3.nano"]
      min_size       = 1
      max_size       = 1
    }
  }
    depends_on = [
    aws_db_instance.db
  ]
}
