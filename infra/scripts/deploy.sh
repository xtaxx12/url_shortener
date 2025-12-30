#!/bin/bash
# ============================================
# Deployment Script
# ============================================
# Usage: ./deploy.sh [environment]
# Example: ./deploy.sh production

# Exit on error, undefined vars, and pipe failures
set -euo pipefail

ENVIRONMENT=${1:-production}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

echo "🚀 Deploying URL Shortener to $ENVIRONMENT"
echo "============================================"

cd "$PROJECT_ROOT"

echo "📦 Pulling latest images..."
docker-compose pull

echo "🔄 Stopping existing containers..."
docker-compose down --remove-orphans

echo "🏗️ Building and starting services..."
docker-compose up -d --build

echo "⏳ Waiting for services to be healthy..."
sleep 10

echo "🔍 Checking service health..."
if curl -s http://localhost/health | grep -q "healthy"; then
    echo "✅ API is healthy"
else
    echo "❌ API health check failed"
    docker-compose logs api-1 api-2
    exit 1
fi

echo "🧹 Cleaning up old images..."
docker image prune -f

echo ""
echo "✅ Deployment complete!"
echo "============================================"
echo "📍 Application URL: http://localhost"
echo "📊 API Health: http://localhost/health"
echo ""
echo "📋 Running containers:"
docker-compose ps
