#!/bin/bash
set -e

echo "Registering SPIRE Workload for Payment Service"
echo "============================================="

echo "Checking Helm connectivity..."
if ! command -v helm &> /dev/null; then
    echo "Helm is not installed or not in PATH."
    exit 1
fi
echo "Helm is working."

echo "Fetching the attested SPIRE Agent ID..."
# Fetching the first attested agent ID. In a single-node Minikube setup with the default chart, this is the dynamic ID.
# Format: SPIFFE ID         : spiffe://securepay.dev/spire/agent/k8s_psat/example-cluster/<UUID>
AGENT_ID=$(kubectl exec -n spire spire-server-0 -- /opt/spire/bin/spire-server agent list | grep "SPIFFE ID" | awk '{print $NF}' | head -n 1 | tr -d '\r')

if [ -z "$AGENT_ID" ]; then
    echo "Error: No attested SPIRE Agent found. Please ensure SPIRE is installed, running, and attested."
    exit 1
fi
echo "Using Agent ID: $AGENT_ID"

echo "Checking for existing Payment Service registration..."
# Check if entry exists to avoid duplicates. Identify by SPIFFE ID.
EXISTING_ENTRY=$(kubectl exec -n spire spire-server-0 -- /opt/spire/bin/spire-server entry show -spiffeID spiffe://securepay.dev/payment-service 2>/dev/null | grep "Entry ID" | awk '{print $NF}' | head -n 1 | tr -d '\r' || true)

if [ ! -z "$EXISTING_ENTRY" ]; then
    echo "Entry for payment-service already exists (ID: $EXISTING_ENTRY). Deleting to recreate..."
    kubectl exec -n spire spire-server-0 -- /opt/spire/bin/spire-server entry delete -entryID $EXISTING_ENTRY
fi

echo "Creating registration entry for Payment Service..."
# Registering the workload
# SPIFFE ID: spiffe://securepay.dev/payment-service
# Parent ID: The attested Agent ID
# Selectors: Kubernetes Namespace (default) and ServiceAccount (payment-service)
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://securepay.dev/payment-service \
    -parentID $AGENT_ID \
    -selector k8s:ns:default \
    -selector k8s:sa:payment-service

echo "Payment Service workload registered successfully."
