import hmac, hashlib, base64

msg = b'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXIiLCJleHAiOjE5MjIzOTEyNjR9'
key = b'securepay-secret-key'

sig = hmac.new(key, msg, hashlib.sha256).digest()
sig_b64 = base64.urlsafe_b64encode(sig).decode().rstrip("=")

print(msg.decode() + "." + sig_b64)
