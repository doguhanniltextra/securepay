#!/bin/bash

# Kubernetes Deployment Script for E-Commerce SPIFFE Project
# This script automates the entire deployment process

set -e

echo "🚀 Starting Kubernetes Deployment"
echo "=================================="

# Step 1: Install Minikube
echo ""
echo "📦 Step 1: Installing Minikube..."
if ! command -v minikube &> /dev/null; then
    curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
    sudo install minikube-linux-amd64 /usr/local/bin/minikube
    rm minikube-linux-amd64
    echo "✅ Minikube installed"
else
    echo "✅ Minikube already installed"
fi

# Step 2: Install Helm
echo ""
echo "📦 Step 2: Installing Helm..."
if ! command -v helm &> /dev/null; then
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
    echo "✅ Helm installed"
else
    echo "✅ Helm already installed"
fi

# Step 3: Start Minikube
echo ""
echo "🎯 Step 3: Starting Minikube..."
minikube start --driver=docker
echo "✅ Minikube started"

# Step 4: Install SPIRE (simplified - without controller manager)
echo ""
echo "🔐 Step 4: Installing SPIRE..."
helm repo add spiffe https://spiffe.github.io/helm-charts-hardened/ || true
helm repo update

# Install SPIRE without controller manager to avoid CRD issues
helm install spire spiffe/spire \
  --create-namespace \
  --namespace spire \
  --set global.spire.trustDomain=ecommerce.local \
  --set spire-server.controllerManager.enabled=false

echo "⏳ Waiting for SPIRE to be ready..."
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=spire-server -n spire --timeout=300s
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=spire-agent -n spire --timeout=300s
echo "✅ SPIRE installed and ready"

# Step 4.5: Register workloads manually
echo ""
echo "📝 Step 4.5: Registering workloads..."

# Register each service with SPIRE
kubectl exec -n spire spire-server-0 -- \
  /opt/spire/bin/spire-server entry create \
  -spiffeID spiffe://ecommerce.local/ns/ecommerce/sa/api-gateway \
  -parentID spiffe://ecommerce.local/spire/agent/k8s_psat/minikube/default \
  -selector k8s:ns:ecommerce \
  -selector k8s:sa:api-gateway

kubectl exec -n spire spire-server-0 -- \
  /opt/spire/bin/spire-server entry create \
  -spiffeID spiffe://ecommerce.local/ns/ecommerce/sa/order-service \
  -parentID spiffe://ecommerce.local/spire/agent/k8s_psat/minikube/default \
  -selector k8s:ns:ecommerce \
  -selector k8s:sa:order-service

kubectl exec -n spire spire-server-0 -- \
  /opt/spire/bin/spire-server entry create \
  -spiffeID spiffe://ecommerce.local/ns/ecommerce/sa/inventory-service \
  -parentID spiffe://ecommerce.local/spire/agent/k8s_psat/minikube/default \
  -selector k8s:ns:ecommerce \
  -selector k8s:sa:inventory-service

kubectl exec -n spire spire-server-0 -- \
  /opt/spire/bin/spire-server entry create \
  -spiffeID spiffe://ecommerce.local/ns/ecommerce/sa/payment-service \
  -parentID spiffe://ecommerce.local/spire/agent/k8s_psat/minikube/default \
  -selector k8s:ns:ecommerce \
  -selector k8s:sa:payment-service

kubectl exec -n spire spire-server-0 -- \
  /opt/spire/bin/spire-server entry create \
  -spiffeID spiffe://ecommerce.local/ns/ecommerce/sa/notification-service \
  -parentID spiffe://ecommerce.local/spire/agent/k8s_psat/minikube/default \
  -selector k8s:ns:ecommerce \
  -selector k8s:sa:notification-service

echo "✅ Workloads registered"

# Step 5: Build Docker Images
echo ""
echo "🐳 Step 5: Building Docker images..."
eval $(minikube docker-env)

docker build -t ecommerce/api-gateway:latest ./services/api-gateway
docker build -t ecommerce/order-service:latest ./services/order-service
docker build -t ecommerce/inventory-service:latest ./services/inventory-service
docker build -t ecommerce/payment-service:latest ./services/payment-service
docker build -t ecommerce/notification-service:latest ./services/notification-service

echo "✅ All images built"

# Step 6: Deploy Services
echo ""
echo "☸️  Step 6: Deploying services to Kubernetes..."

kubectl apply -f k8s/namespace.yaml
kubectl label namespace ecommerce spiffe=enabled --overwrite
kubectl apply -f k8s/serviceaccounts.yaml
kubectl apply -f k8s/notification-service.yaml
kubectl apply -f k8s/payment-service.yaml
kubectl apply -f k8s/inventory-service.yaml
kubectl apply -f k8s/order-service.yaml
kubectl apply -f k8s/api-gateway.yaml

echo "⏳ Waiting for all pods to be ready..."
kubectl wait --for=condition=ready pod --all -n ecommerce --timeout=300s
echo "✅ All services deployed"

# Step 7: Verify SPIRE Registration
echo ""
echo "🔍 Step 7: Verifying SPIRE workload registration..."
kubectl exec -n spire spire-server-0 -- \
  /opt/spire/bin/spire-server entry show

# Step 8: Display Access Information
echo ""
echo "=================================="
echo "✅ Deployment Complete!"
echo "=================================="
echo ""
echo "📍 API Gateway URL: http://$(minikube ip):30080"
echo ""
echo "🧪 Test command:"
echo "curl -X POST http://$(minikube ip):30080/api/orders \\"
echo "  -H \"Authorization: Bearer demo_token\" \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"user_id\":\"user_123\",\"product_id\":\"prod_456\",\"quantity\":2,\"amount\":199.98}'"
echo ""
echo "📊 View logs:"
echo "kubectl logs -f -n ecommerce deployment/api-gateway"
echo ""
echo "🔍 Check pods:"
echo "kubectl get pods -n ecommerce"
