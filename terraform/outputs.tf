output "bastion_ip" {
  description = "ipv4 address of bastion instance"
  value       = aws_instance.bastion.public_ip
}

output "cluster_endpoint" {
  description = "endpoint of EKS control plane"
  value       = module.eks.cluster_endpoint
}

output "db_endpoint" {
  description = "endpoint of mysql instance"
  value       = aws_db_instance.db.endpoint
}

output "ecr_repository" {
  description = "endpoint of ecr repository"
  value       = aws_ecr_repository.repository.repository_url
}
