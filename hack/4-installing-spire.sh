#!/bin/bash
set -e

echo "Starting SPIRE Installation (Step 4)"
echo "===================================="

echo "Step 4: Installing SPIRE..."
helm repo add spiffe https://spiffe.github.io/helm-charts-hardened/ || true
helm repo update

# Cleanup existing installation forcefully
echo "Cleaning up any existing SPIRE installation..."
helm uninstall spire -n spire || true
# Try to delete namespace, ignore if not found
kubectl delete ns spire --force --grace-period=0 || true

# Wait a bit for cleanup
sleep 5

# Clean up Minikube persistent data for SPIRE server 
# This is critical because Minikube's hostpath provisioner might not clean up data immediately, causing 'trust domain mismatch' errors on reinstall
echo "Cleaning up Minikube persistent data..."
minikube ssh 'sudo rm -rf /tmp/hostpath-provisioner/spire/spire-data-spire-server-0' || true


# Create namespace manually to ensure it exists and is fresh
echo "Creating spire namespace..."
kubectl create ns spire

# Install SPIRE using 'install' instead of 'upgrade' to avoid issues with missing releases
echo "Installing SPIRE with Helm..."
helm install spire spiffe/spire \
  --namespace spire \
  --set global.spire.trustDomain=securepay.dev \
  --set spire-server.controllerManager.enabled=false

echo "Waiting for SPIRE to be ready..."
# Updated selectors based on actual pod labels
kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance=spire,app.kubernetes.io/name=server -n spire --timeout=300s
kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance=spire,app.kubernetes.io/name=agent -n spire --timeout=300s
echo "SPIRE installed and ready."
