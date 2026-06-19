from django.db import models


class Url(models.Model):
    id = models.CharField(max_length=8, primary_key=True)
    original = models.TextField(unique=True)
    accesses = models.BigIntegerField(default=0)

    class Meta:
        db_table = "url"
