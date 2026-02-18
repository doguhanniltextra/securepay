# SecurePay Makefile
# Replaces various hack scripts for easier management.

# Variables
MINIKUBE_PROFILE := minikube
KUBECTL := kubectl
HELM := helm
DOCKER_BUILD := eval $$(minikube -p $(MINIKUBE_PROFILE) docker-env) && docker build

.PHONY: help
help:
	@echo "SecurePay Development Commands:"
	@echo "  make ping          - Ping the cluster"
	@echo "  make start         - Start Minikube cluster"
	@echo "  make stop          - Stop Minikube cluster"
	@echo "  make delete        - Delete Minikube cluster"
	@echo "  make infrastructure - Install all infrastructure (Spire, Kafka, Redis, Jaeger)"
	@echo "  make spire         - Install and Configure SPIRE"
	@echo "  make spire-register - Register workloads with SPIRE"
	@echo "  make kafka         - Install Kafka"
	@echo "  make redis         - Install Redis"
	@echo "  make jaeger        - Install Jaeger"
	@echo "  make build         - Build all service Docker images"
	@echo "  make deploy        - Deploy application manifests to K8s"
	@echo "  make all           - Start, Infra, Build, Deploy"
	@echo "  make test          - Run end-to-end payment test"
	@echo "  make clean         - Remove all resources"


# --- Ping ---
.PHONY: ping 
ping:
	@echo "Pong"

# --- Cluster Management ---
.PHONY: start
start:
	@echo "Starting Minikube..."
	minikube start --driver=docker
	@echo "Cluster started."

.PHONY: stop
stop:
	minikube stop

.PHONY: delete
delete:
	minikube delete

# --- Infrastructure ---
.PHONY: infrastructure
infrastructure: spire spire-register kafka redis jaeger

.PHONY: spire
spire:
	@echo "Installing SPIRE..."
	$(HELM) repo add spiffe https://spiffe.github.io/helm-charts-hardened/
	$(HELM) repo update
	-$(HELM) uninstall spire -n spire
	-$(KUBECTL) delete ns spire --force --grace-period=0
	# Cleanup Minikube hostpath (critical for SPIRE re-install)
	-minikube ssh 'sudo rm -rf /tmp/hostpath-provisioner/spire/spire-data-spire-server-0'
	$(KUBECTL) create ns spire
	$(HELM) install spire spiffe/spire --namespace spire \
		--set global.spire.trustDomain=securepay.dev \
		--set spire-server.controllerManager.enabled=false
	@echo "Waiting for SPIRE to be ready..."
	$(KUBECTL) wait --for=condition=ready pod -l app.kubernetes.io/instance=spire,app.kubernetes.io/name=server -n spire --timeout=300s
	$(KUBECTL) wait --for=condition=ready pod -l app.kubernetes.io/instance=spire,app.kubernetes.io/name=agent -n spire --timeout=300s

.PHONY: spire-register
spire-register:
	@echo "Registering Workloads..."
	./hack/5-register-spire.sh

.PHONY: kafka
kafka:
	@echo "Installing Kafka..."
	-$(HELM) uninstall secure-pay-kafka
	$(KUBECTL) apply -f k8s/kafka/kafka.yaml
	$(KUBECTL) wait --for=condition=available deployment/secure-pay-kafka --timeout=300s

.PHONY: redis
redis:
	@echo "Installing Redis..."
	$(HELM) repo add bitnami https://charts.bitnami.com/bitnami
	$(HELM) repo update
	$(HELM) upgrade --install secure-pay-redis bitnami/redis -f k8s/redis/values.yaml

.PHONY: jaeger
jaeger:
	@echo "Installing Jaeger..."
	$(HELM) repo add jaegertracing https://jaegertracing.github.io/helm-charts
	$(HELM) repo update
	$(HELM) upgrade --install secure-pay-jaeger jaegertracing/jaeger \
		--set provisionDataStore.cassandra=false \
		--set allInOne.enabled=true \
		--set storage.type=memory \
		--set agent.enabled=false \
		--set collector.enabled=false \
		--set query.enabled=false

# --- Application ---
.PHONY: build
build:
	@echo "Building Docker images..."
	$(DOCKER_BUILD) -t securepay/api-gateway:v1.0.2 api-gateway/
	$(DOCKER_BUILD) -t securepay/payment-service:v0.0.4 payment-service/
	$(DOCKER_BUILD) -t securepay/account-service:v0.0.3 account-service/
	$(DOCKER_BUILD) -t securepay/notification-service:v0.0.5 notification-service/

.PHONY: deploy
deploy:
	@echo "Deploying application..."
	$(KUBECTL) apply -f k8s/api-gateway/
	$(KUBECTL) apply -f k8s/payment-service/
	$(KUBECTL) apply -f k8s/account-service/
	$(KUBECTL) apply -f k8s/notification-service/

.PHONY: all
all: start infrastructure build deploy

# --- Testing ---
.PHONY: test
test:
	@echo "Running E2E tests..."
	./scripts/test_payment.sh

.PHONY: verify
verify:
	$(KUBECTL) get pods -A

.PHONY: clean
clean: delete
