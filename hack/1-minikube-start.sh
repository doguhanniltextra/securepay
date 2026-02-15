#!/bin/bash
set -e

echo "Starting Kubernetes Setup (Minikube)"
echo "======================================"

# Step 1: Install Minikube
echo ""
echo "Step 1: Installing Minikube..."
if ! command -v minikube &> /dev/null; then
    echo "Minikube not found. Downloading..."
    curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
    sudo install minikube-linux-amd64 /usr/local/bin/minikube
    rm minikube-linux-amd64
    echo "Minikube installed successfully."
else
    echo "Minikube is already installed."
fi

# Step 2: Install Helm (Optional but recommended)
echo ""
echo "Step 2: Installing Helm..."
if ! command -v helm &> /dev/null; then
    echo "Helm not found. Downloading..."
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
    echo "Helm installed successfully."
else
    echo "Helm is already installed."
fi

# Step 3: Start Minikube
echo ""
echo "Step 3: Starting Minikube..."

# Check if minikube is running via status command (check exit code)
if minikube status | grep -q "Running"; then
    echo "Minikube is already running."
else
    echo "Starting Minikube with docker driver..."
    # The --driver=docker is crucial for WSL 2 environments
    minikube start --driver=docker
    echo "Minikube started successfully."
fi

# Step 4: Verify Status
echo ""
echo "Step 4: Verifying Cluster Status..."
kubectl cluster-info
echo ""
echo "Kubernetes Setup Complete!"
