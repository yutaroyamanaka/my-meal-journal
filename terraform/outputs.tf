output "bastion_ip" {
  description = "ipv4 address of bastion instance"
  value       = aws_instance.bastion.public_ip
}
