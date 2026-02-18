#!/bin/bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXIiLCJleHAiOjE5MjIzOTEyNjR9.EzptHaNipsPmZMrIX60Q2XrLPfAU57C_DfKHyDp4FxQ"
GATEWAY_URL="http://localhost:8086"

ID_KEY="test-idempotency-$(date +%s)"
PAYMENT_ID=$(cat /proc/sys/kernel/random/uuid)

echo "--- Request 1 ---"
echo "Key: $ID_KEY"
curl -s -X POST "$GATEWAY_URL/api/v1/payments" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"payment_id\": \"$PAYMENT_ID\", \"from_account\": \"11111111-1111-1111-1111-111111111111\", \"to_account\": \"22222222-2222-2222-2222-222222222222\", \"amount\": 10.0, \"currency\": \"TRY\", \"idempotency_key\": \"$ID_KEY\"}"
echo -e "\n"

echo "--- Request 2 (Same Key) ---"
curl -s -X POST "$GATEWAY_URL/api/v1/payments" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"payment_id\": \"$PAYMENT_ID\", \"from_account\": \"11111111-1111-1111-1111-111111111111\", \"to_account\": \"22222222-2222-2222-2222-222222222222\", \"amount\": 10.0, \"currency\": \"TRY\", \"idempotency_key\": \"$ID_KEY\"}"
echo -e "\n"

echo "--- Logs from payment-service ---"
kubectl logs -l app=payment-service --tail 10
