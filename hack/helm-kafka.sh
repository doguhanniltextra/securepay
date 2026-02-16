#!/bin/bash
set -e

echo "Uninstalling any previous Helm Kafka releases..."
helm uninstall secure-pay-kafka 2>/dev/null || true

echo "Applying standard Kafka deployment (lensesio/fast-data-dev)..."
# Using single-container Kafka for development avoid persistent Bitnami image issues.
# Service name will be: secure-pay-kafka
# Access: secure-pay-kafka.default.svc.cluster.local:9092
kubectl apply -f hack/kafka.yaml

echo "Waiting for Kafka to be ready..."
kubectl wait --for=condition=available deployment/secure-pay-kafka --timeout=300s
