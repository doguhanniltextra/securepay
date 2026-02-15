#!/bin/bash
set -e

echo "Starting Helm Installation"
echo "======================================"

# Step 2: Install Helm
echo ""
echo "Step 2: Installing Helm..."
if ! command -v helm &> /dev/null; then
    echo "Helm not found. Downloading..."
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
    echo "Helm installed successfully."
else
    echo "Helm is already installed."
fi

# Step 3: Verify Helm
echo ""
echo "Step 3: Verifying Helm..."
helm version
echo "Helm Setup Complete!"
