#!/bin/bash

# Clean WSL Deployment Script for E-Commerce SPIFFE Project

set -e

echo "Starting clean deployment..."

# 1. Remove old directory
echo "Removing old directory..."
rm -rf ~/ecommerce-spiffe

# 2. Create fresh directory
echo "Creating fresh directory..."
mkdir -p ~/ecommerce-spiffe

# 3. Copy entire project
echo "Copying project files..."
cp -r /mnt/c/Users/user/Desktop/spiffe-spike-example-folder/e-commerce/* ~/ecommerce-spiffe/

# 4. Navigate to project
cd ~/ecommerce-spiffe

# 5. Verify structure
echo ""
echo "Verifying project structure..."
echo "Services:"
ls -la services/
echo ""
echo "SPIRE directories:"
ls -la spire/
echo ""
echo "SPIRE scripts:"
ls -la spire/scripts/
echo ""
echo "SPIRE configs:"
ls -la spire/server/
ls -la spire/agent/

# 6. Make scripts executable
echo ""
echo "Making scripts executable..."
chmod +x spire/scripts/*.sh

# 7. Create necessary directories
echo "Creating SPIRE directories..."
mkdir -p spire/bin
mkdir -p spire/server/data
mkdir -p spire/agent/data

echo ""
echo "✅ Deployment complete!"
echo ""
echo "Next steps:"
echo "1. cd ~/ecommerce-spiffe"
echo "2. ./spire/scripts/setup.sh"
echo "3. ./spire/scripts/start-server.sh (Terminal 1)"
echo "4. ./spire/scripts/start-agent.sh (Terminal 2)"
echo "5. ./spire/scripts/register-workloads.sh (Terminal 3)"
