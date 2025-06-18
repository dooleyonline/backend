from django.shortcuts import render
from .serializers import CustomTokenObtainPairSerializer
from rest_framework.permissions import AllowAny
from rest_framework_simplejwt.views import TokenObtainPairView


# Create your views here.

class CustomTokenObtainPairView(TokenObtainPairView):
    """
    Custom view to handle token generation with additional user information.
    """
    serializer_class = CustomTokenObtainPairSerializer
    permission_classes = (AllowAny,)
