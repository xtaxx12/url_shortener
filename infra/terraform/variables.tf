# ============================================
# Terraform Variables
# ============================================

variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "url-shortener"
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"  # Free Tier eligible
}

variable "key_pair_name" {
  description = "Name of the SSH key pair"
  type        = string
  default     = "url-shortener-key"
}

variable "allowed_ssh_cidr" {
  description = "CIDR blocks allowed to SSH"
  type        = list(string)
  default     = ["0.0.0.0/0"]  # Restrict in production!
}

variable "create_elastic_ip" {
  description = "Whether to create an Elastic IP"
  type        = bool
  default     = false  # Set to true if you need a static IP
}
