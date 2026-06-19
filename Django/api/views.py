from django.conf import settings
from django.db import transaction
from django.shortcuts import redirect
from rest_framework import status
from rest_framework.decorators import api_view, authentication_classes, permission_classes
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response

from .authentication import EnvBasicAuthentication
from .cache import url_cache
from .models import Url
from .serializers import UrlCreateInput
from .utils import generate_short_id

@api_view(["GET", "POST"])
@authentication_classes([EnvBasicAuthentication])
@permission_classes([IsAuthenticated])
def handle_urls(request):
    if request.method == "GET":
        return Response(url_cache.get_all())

    elif request.method == "POST":
        serializer = UrlCreateInput(data=request.data)
        if not serializer.is_valid():
            return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)

        original = serializer.validated_data["url"]

        existing = Url.objects.filter(original=original).first()
        if existing:
            return Response(
                {"id": existing.id, "url": original},
                status=status.HTTP_200_OK,
            )

        short_id = generate_short_id(8, lambda id: url_cache.exists(id))

        with transaction.atomic():
            Url.objects.create(id=short_id, original=original, accesses=0)

        url_cache.set(short_id, {"original": original, "accesses": 0})

        return Response(
            {"id": short_id, "url": original},
            status=status.HTTP_201_CREATED,
        )


@api_view(["GET"])
@authentication_classes([EnvBasicAuthentication])
@permission_classes([IsAuthenticated])
def get_url_info(request, id):
    info = url_cache.get(id)
    if info is None:
        return Response({"detail": "URL not found"}, status=status.HTTP_400_BAD_REQUEST)
    return Response(info)


@api_view(["GET", "DELETE"])
@authentication_classes([EnvBasicAuthentication])
@permission_classes([IsAuthenticated])
def handle_id(request, id):
    if request.method == "GET":
        info = url_cache.get(id)
        if info is None:
            return Response({"detail": "URL not found"}, status=status.HTTP_404_NOT_FOUND)

        new_accesses = url_cache.increment_accesses(id)
        Url.objects.filter(id=id).update(accesses=new_accesses)

        return redirect(info["original"], status=302)

    elif request.method == "DELETE":
        if not url_cache.exists(id):
            return Response({"detail": "URL not found"}, status=status.HTTP_400_BAD_REQUEST)

        with transaction.atomic():
            Url.objects.filter(id=id).delete()

        url_cache.delete(id)
        return Response(status=status.HTTP_200_OK)
