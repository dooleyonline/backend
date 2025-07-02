# dooleyonline/settings/prod.py
from .base import *
import os

DEBUG = False
SECURE_SSL_REDIRECT = True
USE_S3 = False

# AWS_ACCESS_KEY_ID = os.getenv("AWS_ACCESS_KEY_ID")
# AWS_SECRET_ACCESS_KEY = os.getenv("AWS_SECRET_ACCESS_KEY")
# AWS_STORAGE_BUCKET_NAME = "dooleybay"
# AWS_S3_REGION_NAME = "ap-southeast-2"
# AWS_S3_SIGNATURE_VERSION = 's3v4'
# AWS_S3_FILE_OVERWRITE = False
# AWS_DEFAULT_ACL = None

STORAGES = {
    "default": {
        "BACKEND": "storages.backends.s3boto3.S3Boto3Storage",
    },
    "staticfiles": {
        "BACKEND": "django.contrib.staticfiles.storage.StaticFilesStorage",
    },
}


SECURE_PROXY_SSL_HEADER = ('HTTP_X_FORWARDED_PROTO', 'https')
USE_X_FORWARDED_HOST = True
SECURE_SSL_REDIRECT = True

CORS_ALLOW_HEADERS = [
    "accept",
    "accept-encoding",
    "authorization",
    "content-type",
    "dnt",
    "origin",
    "user-agent",
    "x-csrftoken",
    "x-requested-with",
]
CORS_ALLOWED_ORIGINS = ["http://localhost:8000","http://0.0.0.0:8080", "https://api.dooleyonline.net", "https://backend-production-9918.up.railway.app", " https://dooleyonline.net/"]
CSRF_TRUSTED_ORIGINS = ["http://0.0.0.0:8080", "https://api.dooleyonline.net", "https://backend-production-9918.up.railway.app", "https://dooleyonline.net/"]