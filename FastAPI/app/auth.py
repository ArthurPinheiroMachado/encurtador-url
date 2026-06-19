import secrets

from fastapi import Depends, HTTPException, status
from fastapi.security import HTTPBasic, HTTPBasicCredentials

from .config import settings

security = HTTPBasic()


def verify_auth(credentials: HTTPBasicCredentials = Depends(security)):
    is_user = secrets.compare_digest(credentials.username, settings.USER)
    is_pass = secrets.compare_digest(credentials.password, settings.PASS)

    if not (is_user and is_pass):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Unauthorized",
            headers={"WWW-Authenticate": "Basic"},
        )

    return credentials
