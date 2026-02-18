#!/bin/bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXIiLCJleHAiOjE5MjIzOTEyNjR9.EzptHaNipsPmZMrIX60Q2XrLPfAU57C_DfKHyDp4FxQ"
GATEWAY_URL="http://localhost:8087"

echo "--- Request 1 (Expected: Cache Miss) ---"
curl -s -X GET "$GATEWAY_URL/api/v1/accounts/11111111-1111-1111-1111-111111111111/balance" \
  -H "Authorization: Bearer $TOKEN"
echo ""

echo "--- Request 2 (Expected: Cache Hit) ---"
curl -s -X GET "$GATEWAY_URL/api/v1/accounts/11111111-1111-1111-1111-111111111111/balance" \
  -H "Authorization: Bearer $TOKEN"
echo ""

echo "--- Account Service Logs ---"
kubectl logs -l app=account-service --tail 10
