#!/bin/bash
set -e

echo "Registering SPIRE Workload for Account Service"
echo "============================================="

echo "Checking Helm connectivity..."
if ! command -v helm &> /dev/null; then
    echo "Helm is not installed or not in PATH."
    exit 1
fi
echo "Helm is working."

echo "Fetching the attested SPIRE Agent ID..."
AGENT_ID=$(kubectl exec -n spire spire-server-0 -- /opt/spire/bin/spire-server agent list | grep "SPIFFE ID" | awk '{print $NF}' | head -n 1 | tr -d '\r')

if [ -z "$AGENT_ID" ]; then
    echo "Error: No attested SPIRE Agent found. Please ensure SPIRE is installed, running, and attested."
    exit 1
fi
echo "Using Agent ID: $AGENT_ID"

echo "Checking for existing Account Service registration..."
EXISTING_ENTRY=$(kubectl exec -n spire spire-server-0 -- /opt/spire/bin/spire-server entry show -spiffeID spiffe://securepay.dev/account-service 2>/dev/null | grep "Entry ID" | awk '{print $NF}' | head -n 1 | tr -d '\r' || true)

if [ ! -z "$EXISTING_ENTRY" ]; then
    echo "Entry for account-service already exists (ID: $EXISTING_ENTRY). Deleting to recreate..."
    kubectl exec -n spire spire-server-0 -- /opt/spire/bin/spire-server entry delete -entryID $EXISTING_ENTRY
fi

echo "Creating registration entry for Account Service..."
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://securepay.dev/account-service \
    -parentID $AGENT_ID \
    -selector k8s:ns:default \
    -selector k8s:sa:account-service

echo "Account Service workload registered successfully."
