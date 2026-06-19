from django.conf import settings
from django.urls import path

from . import views

prefix = settings.HTTP_BASE

urlpatterns = [
    path(f"{prefix}/urls", views.handle_urls, name="handle-urls"),
    path(f"{prefix}/urls/<str:id>", views.get_url_info, name="get-url-info"),
    path(f"{prefix}/<str:id>", views.handle_id, name="handle-id"),
]
