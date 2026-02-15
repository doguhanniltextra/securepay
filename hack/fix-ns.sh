#!/bin/bash
kubectl get ns spire -o json | python3 -c "import sys, json; data=json.load(sys.stdin); data['spec']['finalizers']=[]; print(json.dumps(data))" > spire.json
kubectl replace --raw "/api/v1/namespaces/spire/finalize" -f ./spire.json
rm spire.json
