from django.apps import AppConfig


class ApiConfig(AppConfig):
    name = "api"

    def ready(self):
        from .cache import url_cache
        try:
            url_cache.load()
        except Exception:
            pass
