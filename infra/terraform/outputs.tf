# ============================================
# Terraform Outputs
# ============================================

output "instance_id" {
  description = "ID of the EC2 instance"
  value       = aws_instance.url_shortener.id
}

output "instance_public_ip" {
  description = "Public IP of the EC2 instance"
  value       = aws_instance.url_shortener.public_ip
}

output "instance_public_dns" {
  description = "Public DNS of the EC2 instance"
  value       = aws_instance.url_shortener.public_dns
}

output "elastic_ip" {
  description = "Elastic IP address (if created)"
  value       = var.create_elastic_ip ? aws_eip.url_shortener[0].public_ip : null
}

output "security_group_id" {
  description = "ID of the security group"
  value       = aws_security_group.url_shortener.id
}

output "application_url" {
  description = "URL to access the application"
  value       = var.create_elastic_ip ? "http://${aws_eip.url_shortener[0].public_ip}" : "http://${aws_instance.url_shortener.public_ip}"
}

output "ssh_command" {
  description = "SSH command to connect to the instance"
  value       = "ssh -i ${var.key_pair_name}.pem ec2-user@${aws_instance.url_shortener.public_ip}"
}
