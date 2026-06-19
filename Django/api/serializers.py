from rest_framework import serializers

from .models import Url


class UrlCreateInput(serializers.Serializer):
    url = serializers.URLField()


class UrlCreateOutput(serializers.Serializer):
    id = serializers.CharField()
    url = serializers.CharField()


class InfoSerializer(serializers.Serializer):
    original = serializers.CharField()
    accesses = serializers.IntegerField()


class UrlModelSerializer(serializers.ModelSerializer):
    class Meta:
        model = Url
        fields = ["id", "original", "accesses"]
