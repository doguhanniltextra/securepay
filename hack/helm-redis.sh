#!/bin/bash
# hack/helm-redis.sh — Deploy Redis via Helm

# Add Bitnami Helm repository
helm repo add bitnami https://charts.bitnami.com/bitnami

# Update Helm repositories
helm repo update

# Install Redis with values file
helm install secure-pay-redis bitnami/redis -f k8s/redis/values.yaml
