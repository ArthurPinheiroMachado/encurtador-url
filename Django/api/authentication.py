import os
import secrets

from rest_framework import authentication, exceptions


class EnvBasicAuthentication(authentication.BaseAuthentication):
    def authenticate(self, request):
        auth_header = request.META.get("HTTP_AUTHORIZATION", "")

        if not auth_header.startswith("Basic "):
            return None

        import base64

        try:
            token = auth_header[6:]
            decoded = base64.b64decode(token).decode("utf-8")
        except Exception:
            raise exceptions.AuthenticationFailed("Invalid Authorization header")

        parts = decoded.split(":", 1)
        if len(parts) != 2:
            raise exceptions.AuthenticationFailed("Invalid Authorization header")

        username, password = parts
        expected_user = os.environ.get("USER", "user")
        expected_pass = os.environ.get("PASS", "pass123")

        user_ok = secrets.compare_digest(username, expected_user)
        pass_ok = secrets.compare_digest(password, expected_pass)

        if not (user_ok and pass_ok):
            raise exceptions.AuthenticationFailed("Unauthorized")

        return (username, None)
