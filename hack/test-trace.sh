#!/bin/bash

# Ensure bash runs this
GATEWAY_HOST="localhost:8087"
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.dtmo7EtANIsh4V84XavB30DR4KCk2rINmRdyNZoLcG8"

curl -v -X POST http://$GATEWAY_HOST/api/v1/payments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_id": "99999999-9999-9999-9999-999999999999",
    "from_account": "11111111-1111-1111-1111-111111111111",
    "to_account": "22222222-2222-2222-2222-222222222222",
    "amount": 10.50,
    "currency": "TRY",
    "idempotency_key": "trace-test-4"
  }'
