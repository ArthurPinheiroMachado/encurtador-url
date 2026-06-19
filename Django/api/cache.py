import threading

from .models import Url


class UrlCache:
    def __init__(self):
        self._lock = threading.Lock()
        self._urls: dict[str, dict] = {}

    def load(self):
        with self._lock:
            self._urls.clear()
            for u in Url.objects.values("id", "original", "accesses"):
                self._urls[u["id"]] = {
                    "original": u["original"],
                    "accesses": u["accesses"],
                }

    def get_all(self) -> dict[str, dict]:
        with self._lock:
            return dict(self._urls)

    def get(self, id: str) -> dict | None:
        with self._lock:
            return self._urls.get(id)

    def exists(self, id: str) -> bool:
        with self._lock:
            return id in self._urls

    def set(self, id: str, info: dict):
        with self._lock:
            self._urls[id] = info

    def delete(self, id: str):
        with self._lock:
            self._urls.pop(id, None)

    def increment_accesses(self, id: str) -> int:
        with self._lock:
            if id in self._urls:
                self._urls[id]["accesses"] += 1
                return self._urls[id]["accesses"]
            return 0


url_cache = UrlCache()
