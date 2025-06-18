from django.shortcuts import render
from .serializers import CustomTokenObtainPairSerializer, RegisterSerializer
from rest_framework.permissions import AllowAny
from rest_framework_simplejwt.views import TokenObtainPairView
from rest_framework import generics
from django.contrib.auth.models import User


# Create your views here.

class CustomTokenObtainPairView(TokenObtainPairView):
    """
    Custom view to handle token generation with additional user information.
    """
    serializer_class = CustomTokenObtainPairSerializer
    permission_classes = (AllowAny,)


class RegisterView(generics.CreateAPIView): # for create-only endpoints
    queryset = User.objects.all()
    permission_classes = (AllowAny,)
    serializer_class = RegisterSerializer