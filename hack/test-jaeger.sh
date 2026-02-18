#!/bin/bash

JAEGER_HOST="secure-pay-jaeger.default.svc.cluster.local"
JAEGER_PORT="16686"
OTLP_PORT="4317"

echo "Using Jaeger Host: $JAEGER_HOST"

# 1. Test UI Port via curl pod
echo "Testing Jaeger Query UI ($JAEGER_PORT)..."
kubectl run check-jaeger-ui --rm -i --restart='Never' --image=curlimages/curl -- curl -s -o /dev/null -w "%{http_code}" http://$JAEGER_HOST:$JAEGER_PORT
echo ""

# 2. Test OTLP Port (TCP) via nc (netcat) pod
echo "Testing OTLP Port ($OTLP_PORT)..."
kubectl run check-jaeger-otlp --rm -i --restart='Never' --image=busybox -- sh -c "nc -zv $JAEGER_HOST $OTLP_PORT"
echo ""
