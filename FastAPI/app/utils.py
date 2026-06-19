import secrets


def generate_short_id(length: int = 8, exists=None) -> str:
    charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    max_attempts = 100

    for _ in range(max_attempts):
        id = "".join(secrets.choice(charset) for _ in range(length))
        if exists is None or not exists(id):
            return id

    raise RuntimeError(f"failed to generate unique ID after {max_attempts} attempts")
