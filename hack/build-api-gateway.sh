#!/bin/bash
eval $(minikube -p minikube docker-env)
docker build -t api-gateway:latest -f docker/api-gateway.dockerfile .
