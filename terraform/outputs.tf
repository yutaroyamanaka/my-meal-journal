output "bastion_ip" {
  description = "ipv4 address of bastion instance"
  value       = aws_instance.bastion.public_ip
}

output "cluster_endpoint" {
  description = "endpoint of EKS control plane"
  value       = module.eks.cluster_endpoint
}

output "cluster_name" {
  description = "eks cluster name"
  value       = module.eks.cluster_name
}

output "db_endpoint" {
  description = "endpoint of mysql instance"
  value       = split(":", aws_db_instance.db.endpoint)[0]
}

output "db_name" {
  description = "name of database schema"
  value       = aws_db_instance.db.db_name
}

output "ecr_repository" {
  description = "endpoint of ecr repository"
  value       = aws_ecr_repository.repository.repository_url
}
