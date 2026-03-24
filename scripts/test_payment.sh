#!/bin/bash
# End-to-end payment test script for SecurePay.
# Generates a fresh JWT token (valid for 1 hour) before making the request.

set -euo pipefail

# Generate a fresh token (reads JWT_SECRET from env, defaults to dev secret)
TOKEN=$(python3 "$(dirname "$0")/gen_token.py" --ttl 3600 | tail -1)

echo "=== SecurePay E2E Payment Test ==="
echo "Token: ${TOKEN:0:20}..."

curl -s -X POST http://localhost:8092/api/v1/payments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_id": "'"$(uuidgen 2>/dev/null || cat /proc/sys/kernel/random/uuid)"'",
    "from_account": "11111111-1111-1111-1111-111111111111",
    "to_account": "22222222-2222-2222-2222-222222222222",
    "amount": 250.0,
    "currency": "TRY",
    "idempotency_key": "key-'"$RANDOM"'"
  }' | python3 -m json.tool 2>/dev/null || echo "(raw response above)"

echo ""
echo "=== Test Complete ==="
