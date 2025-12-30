# ============================================
# Terraform Configuration for AWS Infrastructure
# ============================================
# This configuration creates:
# - VPC with public subnet
# - Security Group for the application
# - EC2 instance (t2.micro - Free Tier)
# - Elastic IP (optional)

terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  # Uncomment for remote state (recommended for teams)
  # backend "s3" {
  #   bucket = "your-terraform-state-bucket"
  #   key    = "url-shortener/terraform.tfstate"
  #   region = "us-east-1"
  # }
}

# ============================================
# Provider Configuration
# ============================================
provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "url-shortener"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  }
}

# ============================================
# Data Sources
# ============================================
# Get latest Amazon Linux 2023 AMI
data "aws_ami" "amazon_linux" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

# Get default VPC (Free Tier friendly)
data "aws_vpc" "default" {
  default = true
}

# Get available AZs
data "aws_availability_zones" "available" {
  state = "available"
}

# Get default subnets
data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# ============================================
# Security Group
# ============================================
resource "aws_security_group" "url_shortener" {
  name        = "${var.project_name}-sg"
  description = "Security group for URL Shortener application"
  vpc_id      = data.aws_vpc.default.id

  # SSH access (restrict to your IP in production)
  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = var.allowed_ssh_cidr
  }

  # HTTP access
  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HTTPS access
  ingress {
    description = "HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # All outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_name}-sg"
  }

  lifecycle {
    create_before_destroy = true
  }
}

# ============================================
# EC2 Instance
# ============================================
resource "aws_instance" "url_shortener" {
  ami           = data.aws_ami.amazon_linux.id
  instance_type = var.instance_type

  subnet_id                   = data.aws_subnets.default.ids[0]
  vpc_security_group_ids      = [aws_security_group.url_shortener.id]
  associate_public_ip_address = true

  # SSH key pair (create one in AWS Console first)
  key_name = var.key_pair_name

  # Root volume
  root_block_device {
    volume_type           = "gp3"
    volume_size           = 8  # GB - Free Tier allows up to 30GB
    delete_on_termination = true
    encrypted             = true
  }

  # User data script to install Docker
  user_data = base64encode(templatefile("${path.module}/user_data.sh", {
    docker_compose_version = "2.24.0"
  }))

  # Instance metadata options (security best practice)
  metadata_options {
    http_endpoint               = "enabled"
    http_tokens                 = "required"  # IMDSv2 only
    http_put_response_hop_limit = 1
  }

  tags = {
    Name = "${var.project_name}-server"
  }

  lifecycle {
    create_before_destroy = true
  }
}

# ============================================
# Elastic IP (Optional - comment out to save costs)
# ============================================
resource "aws_eip" "url_shortener" {
  count    = var.create_elastic_ip ? 1 : 0
  instance = aws_instance.url_shortener.id
  domain   = "vpc"

  tags = {
    Name = "${var.project_name}-eip"
  }

  depends_on = [aws_instance.url_shortener]
}
