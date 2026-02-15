#!/bin/bash
eval $(minikube -p minikube docker-env)
docker build -t api-gateway:latest -f api-gateway/apigateway.dockerfile .
