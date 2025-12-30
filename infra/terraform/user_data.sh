#!/bin/bash
# ============================================
# EC2 User Data Script
# ============================================
# This script runs on first boot to configure the instance

set -e

# Update system
dnf update -y

# Install Docker
dnf install -y docker
systemctl enable docker
systemctl start docker

# Add ec2-user to docker group
usermod -aG docker ec2-user

# Install Docker Compose
DOCKER_COMPOSE_VERSION="${docker_compose_version}"
curl -L "https://github.com/docker/compose/releases/download/v$${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Install Git
dnf install -y git

# Create app directory
mkdir -p /opt/url-shortener
chown ec2-user:ec2-user /opt/url-shortener

# Install AWS CLI (for ECR access if needed)
dnf install -y aws-cli

# Configure Docker to start on boot
systemctl enable docker

# Log completion
echo "User data script completed at $(date)" >> /var/log/user-data.log
