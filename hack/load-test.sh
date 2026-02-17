#!/bin/bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXIiLCJleHAiOjE5MjIzOTEyNjR9.EzptHaNipsPmZMrIX60Q2XrLPfAU57C_DfKHyDp4FxQ"
GATEWAY_URL="http://localhost:8085"

echo "Starting load test on $GATEWAY_URL..."

for i in {1..50}
do
  # 1. Success attempt (Valid JWT, although gRPC might fail downstream)
  curl -s -X POST "$GATEWAY_URL/api/v1/payments" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"payment_id\": \"$(cat /proc/sys/kernel/random/uuid)\", \"from_account\": \"ce983b63-d14d-4e92-bc1a-697669d290fb\", \"to_account\": \"614be79f-67f7-4340-9a3d-368297753173\", \"amount\": 10.0, \"currency\": \"TRY\", \"idempotency_key\": \"key-$i-$(date +%s)\"}" > /dev/null

  # 2. Unauthorized attempt (Failing status code)
  curl -s -X GET "$GATEWAY_URL/api/v1/accounts/ce983b63-d14d-4e92-bc1a-697669d290fb/balance" \
    -H "Authorization: Bearer invalid-token" > /dev/null

  # 3. Valid Check Balance (gRPC call)
  curl -s -X GET "$GATEWAY_URL/api/v1/accounts/ce983b63-d14d-4e92-bc1a-697669d290fb/balance" \
    -H "Authorization: Bearer $TOKEN" > /dev/null

  # 4. Health check (No middleware)
  curl -s -X GET "$GATEWAY_URL/health" > /dev/null

  if (( $i % 10 == 0 )); then
    echo "Sent $i batches of requests..."
  fi
  
  sleep 0.1
done

echo "Load test completed."
