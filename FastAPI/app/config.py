import os


class Settings:
    DB_TYPE: str = os.environ.get("DB_TYPE", "postgres")
    DB_NAME: str = os.environ.get("DB_NAME", "encurtador")
    DB_PASS: str = os.environ.get("DB_PASS", "postgres")
    DB_PORT: str = os.environ.get("DB_PORT", "5432")
    DB_HOST: str = os.environ.get("DB_HOST", "0.0.0.0")
    DB_USER: str = os.environ.get("DB_USER", "postgres")
    HTTP_PORT: int = int(os.environ.get("HTTP_PORT", "6060"))
    HTTP_BASE: str = os.environ.get("HTTP_BASE", "/api/")
    TIMEOUT_TIME: int = int(os.environ.get("TIMEOUT_TIME", "3"))
    USER: str = os.environ.get("USER", "user")
    PASS: str = os.environ.get("PASS", "pass123")


settings = Settings()
