import asyncpg

from .config import settings

_pool: asyncpg.Pool | None = None


async def get_pool() -> asyncpg.Pool:
    global _pool
    if _pool is None:
        _pool = await asyncpg.create_pool(
            user=settings.DB_USER,
            password=settings.DB_PASS,
            database=settings.DB_NAME,
            host=settings.DB_HOST,
            port=settings.DB_PORT,
        )
    return _pool


async def close_pool():
    global _pool
    if _pool:
        await _pool.close()
        _pool = None


QUERIES = {
    "create_migrations_table": """
        CREATE TABLE IF NOT EXISTS migrations(
            id INT NOT NULL,
            content TEXT NOT NULL,
            PRIMARY KEY(id)
        );
    """,
    "create_url_table": """
        CREATE TABLE IF NOT EXISTS url(
            id TEXT NOT NULL,
            original TEXT NOT NULL,
            accesses BIGINT DEFAULT 0,
            UNIQUE(original),
            PRIMARY KEY(id)
        );
    """,
    "last_position": "SELECT COALESCE(MAX(id), -1) FROM migrations",
    "insert_stage": "INSERT INTO migrations(id, content) VALUES($1, $2)",
    "get_urls": "SELECT id, original, accesses FROM url",
    "get_url_by_url": "SELECT id, original, accesses FROM url WHERE original = $1",
    "save_url": "INSERT INTO url(id, original, accesses) VALUES($1, $2, 0)",
    "delete_url": "DELETE FROM url WHERE id = $1",
    "update_accesses": "UPDATE url SET accesses = $1 WHERE id = $2",
}


async def migrate():
    pool = await get_pool()
    async with pool.acquire() as conn:
        last_position = await conn.fetchval(QUERIES["last_position"])
        statements = [QUERIES["create_migrations_table"], QUERIES["create_url_table"]]

        for idx, query in enumerate(statements):
            if idx <= last_position:
                continue

            await conn.execute(query)
            await conn.execute(QUERIES["insert_stage"], idx, query)


async def get_all_urls() -> list[dict]:
    pool = await get_pool()
    async with pool.acquire() as conn:
        rows = await conn.fetch(QUERIES["get_urls"])
        return [dict(row) for row in rows]


async def get_url_by_original(original: str) -> dict | None:
    pool = await get_pool()
    async with pool.acquire() as conn:
        row = await conn.fetchrow(QUERIES["get_url_by_url"], original)
        return dict(row) if row else None


async def save_url(id: str, original: str):
    pool = await get_pool()
    async with pool.acquire() as conn:
        async with conn.transaction():
            await conn.execute(QUERIES["save_url"], id, original)


async def delete_url(id: str):
    pool = await get_pool()
    async with pool.acquire() as conn:
        async with conn.transaction():
            await conn.execute(QUERIES["delete_url"], id)


async def update_accesses(id: str, accesses: int):
    pool = await get_pool()
    async with pool.acquire() as conn:
        await conn.execute(QUERIES["update_accesses"], accesses, id)
