#!/bin/bash
# hack/helm-jaeger.sh — Deploy Jaeger via Helm

# Add Jaeger Helm repository
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts

# Update Helm repositories
helm repo update

# Install or Upgrade Jaeger (all-in-one, in-memory storage)
helm upgrade --install secure-pay-jaeger jaegertracing/jaeger \
  --set provisionDataStore.cassandra=false \
  --set allInOne.enabled=true \
  --set storage.type=memory \
  --set agent.enabled=false \
  --set collector.enabled=false \
  --set query.enabled=false
