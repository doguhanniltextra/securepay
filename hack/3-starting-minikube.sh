#!/bin/bash
set -e

echo "Starting Minikube Setup (Step 3)"
echo "=================================="

# Check if minikube is running via status command
echo "Checking Minikube status..."
if minikube status | grep -q "Running"; then
    echo "Minikube is already running."
else
    echo "Starting Minikube with docker driver..."
    minikube start --driver=docker
    echo "Minikube started successfully."
fi

# Verify Helm connectivity as requested in the task description
echo "Verifying Helm connectivity..."
if command -v helm &> /dev/null; then
    helm list -A
    echo "Helm is working and connected to the cluster."
else
    echo "Helm command not found."
    exit 1
fi
