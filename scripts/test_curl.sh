#!/bin/bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXIiLCJleHAiOjE5MjIzOTEyNjR9.EzptHaNipsPmZMrIX60Q2XrLPfAU57C_DfKHyDp4FxQ"
curl -v -H "Authorization: Bearer $TOKEN" http://localhost:8091/api/v1/accounts/11111111-1111-1111-1111-111111111111/balance
