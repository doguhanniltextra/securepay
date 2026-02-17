#!/bin/bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXIiLCJleHAiOjE5MjIzOTEyNjR9.EzptHaNipsPmZMrIX60Q2XrLPfAU57C_DfKHyDp4FxQ"
curl -X POST http://localhost:8092/api/v1/payments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"payment_id": "4736f563-8756-425d-a602-0e3926868661", "from_account": "ce983b63-d14d-4e92-bc1a-697669d290fb", "to_account": "614be79f-67f7-4340-9a3d-368297753173", "amount": 250.0, "currency": "TRY", "idempotency_key": "key-987"}'
