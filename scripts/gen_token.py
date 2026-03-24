#!/usr/bin/env python3
"""Generate a short-lived JWT token for SecurePay API testing.

Usage:
    python3 gen_token.py                           # 1-hour token (default)
    python3 gen_token.py --ttl 3600                # custom TTL in seconds
    JWT_SECRET=my-secret python3 gen_token.py      # custom secret via env
"""

import hmac
import hashlib
import base64
import json
import os
import sys
import time


def b64url_encode(data: bytes) -> str:
    """Base64url encode without padding."""
    return base64.urlsafe_b64encode(data).decode().rstrip("=")


def generate_token(secret: str, subject: str = "test-user", ttl_seconds: int = 3600) -> str:
    """Generate a HS256 JWT token with the given secret and TTL."""
    # Header
    header = b64url_encode(json.dumps(
        {"alg": "HS256", "typ": "JWT"}, separators=(",", ":")
    ).encode())

    # Payload
    now = int(time.time())
    payload = b64url_encode(json.dumps(
        {"sub": subject, "iat": now, "exp": now + ttl_seconds}, separators=(",", ":")
    ).encode())

    # Signature
    signing_input = f"{header}.{payload}".encode()
    signature = hmac.new(secret.encode(), signing_input, hashlib.sha256).digest()
    sig_b64 = b64url_encode(signature)

    return f"{header}.{payload}.{sig_b64}"


def main():
    secret = os.environ.get("JWT_SECRET", "securepay-secret-key")
    ttl = 3600  # default: 1 hour

    # Parse --ttl flag
    if "--ttl" in sys.argv:
        idx = sys.argv.index("--ttl")
        if idx + 1 < len(sys.argv):
            ttl = int(sys.argv[idx + 1])

    token = generate_token(secret=secret, ttl_seconds=ttl)

    exp_time = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(time.time() + ttl))
    print(f"# Token expires at: {exp_time} (TTL: {ttl}s)")
    print(token)


if __name__ == "__main__":
    main()
