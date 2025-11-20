#!/bin/bash

# =================================================================
# Deep Research Agent - Docker Quick Start Script
# =================================================================
# This script helps you get started with the Docker setup quickly
# =================================================================

set -e  # Exit on error

echo "🔍 Deep Research Agent - Docker Setup"
echo "====================================="
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed!"
    echo "Please install Docker Desktop from: https://www.docker.com/products/docker-desktop"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed!"
    echo "Please install Docker Compose"
    exit 1
fi

echo "✅ Docker is installed"
echo "✅ Docker Compose is installed"
echo ""

# Check if .env file exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env

    # Generate JWT secret
    if command -v openssl &> /dev/null; then
        JWT_SECRET=$(openssl rand -base64 32)
        # Replace JWT_SECRET in .env file (works on Linux and macOS)
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS
            sed -i '' "s|JWT_SECRET=.*|JWT_SECRET=$JWT_SECRET|g" .env
        else
            # Linux
            sed -i "s|JWT_SECRET=.*|JWT_SECRET=$JWT_SECRET|g" .env
        fi
        echo "✅ Generated secure JWT_SECRET"
    else
        echo "⚠️  Please manually set JWT_SECRET in .env file"
    fi
else
    echo "✅ .env file already exists"
fi

echo ""

# Create docker-data directory structure
echo "📁 Creating data directories..."
mkdir -p docker-data/{database,uploads,exports,config}
echo "✅ Data directories created"
echo ""

# Check if containers are already running
if docker-compose ps | grep -q "Up"; then
    echo "⚠️  Containers are already running"
    echo ""
    read -p "Do you want to restart them? (y/n) " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "🔄 Stopping containers..."
        docker-compose down
    else
        echo "👋 Exiting without changes"
        exit 0
    fi
fi

# Build and start containers
echo "🏗️  Building Docker images..."
echo "This may take a few minutes on first run..."
echo ""

docker-compose build

echo ""
echo "🚀 Starting containers..."
docker-compose up -d

echo ""
echo "⏳ Waiting for services to start..."
sleep 5

# Check if containers are running
echo ""
echo "📊 Container Status:"
docker-compose ps

echo ""
echo "============================================"
echo "✅ Setup Complete!"
echo "============================================"
echo ""
echo "Your services are now available at:"
echo ""
echo "  🔬 Research Agent:  http://localhost:8080"
echo "  ❤️  Health Check:    http://localhost:8080/health"
echo "  📊 Prometheus:       http://localhost:9090"
echo "  📈 Grafana:          http://localhost:3000 (admin/admin)"
echo ""
echo "============================================"
echo "Quick Commands:"
echo "============================================"
echo ""
echo "  View logs:           docker-compose logs -f research-agent"
echo "  Stop services:       docker-compose down"
echo "  Restart:             docker-compose restart"
echo "  Rebuild:             docker-compose up -d --build"
echo ""
echo "See DOCKER_SETUP.md for complete documentation"
echo ""
echo "Testing the API..."
echo ""

# Wait a bit more for the API to be fully ready
sleep 3

# Test health endpoint
if curl -f -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ API is responding!"
    echo ""
    echo "Response:"
    curl -s http://localhost:8080/health | python3 -m json.tool 2>/dev/null || curl -s http://localhost:8080/health
else
    echo "⚠️  API not responding yet. Check logs with:"
    echo "   docker-compose logs -f research-agent"
fi

echo ""
echo "🎉 Happy researching!"
