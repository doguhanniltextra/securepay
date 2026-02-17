#!/bin/bash
# Add Prometheus Community Helm repo
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Install kube-prometheus-stack (includes Prometheus, Grafana, Alertmanager)
# We disable validation hooks if they cause issues in some Minikube envs, but usually fine.
helm install secure-pay-monitoring prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  --set grafana.adminPassword=admin \
  --set prometheus.service.type=NodePort \
  --set grafana.service.type=NodePort

echo "Waiting for monitoring stack to be ready..."
kubectl rollout status deployment/secure-pay-monitoring-grafana -n monitoring
