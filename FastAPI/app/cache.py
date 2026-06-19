import asyncio

from .models import Info


class UrlCache:
    def __init__(self):
        self._lock = asyncio.Lock()
        self._urls: dict[str, Info] = {}

    async def load(self, urls: list[dict]):
        async with self._lock:
            self._urls.clear()
            for u in urls:
                self._urls[u["id"]] = Info(original=u["original"], accesses=u["accesses"])

    async def get_all(self) -> dict[str, Info]:
        async with self._lock:
            return dict(self._urls)

    async def get(self, id: str) -> Info | None:
        async with self._lock:
            return self._urls.get(id)

    async def exists(self, id: str) -> bool:
        async with self._lock:
            return id in self._urls

    async def set(self, id: str, info: Info):
        async with self._lock:
            self._urls[id] = info

    async def delete(self, id: str):
        async with self._lock:
            self._urls.pop(id, None)

    async def increment_accesses(self, id: str) -> int:
        async with self._lock:
            if id in self._urls:
                self._urls[id].accesses += 1
                return self._urls[id].accesses
            return 0

    async def id_exists(self, id: str) -> bool:
        async with self._lock:
            return id in self._urls


url_cache = UrlCache()
