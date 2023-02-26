output "bastion_ip" {
  description = "ipv4 address of bastion instance"
  value       = aws_instance.bastion.public_ip
}

output "cluster_endpoint" {
  description = "Endpoint for EKS control plane"
  value       = module.eks.cluster_endpoint
}

output "cluster_name" {
  description = "Kubernetes Cluster Name"
  value       = module.eks.cluster_name
}
