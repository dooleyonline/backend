from django.shortcuts import render
from .serializers import RegisterSerializer, CookieTokenRefreshSerializer
from rest_framework.permissions import AllowAny, IsAuthenticated
from rest_framework import generics
from django.contrib.auth.models import User
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework import status
from rest_framework_simplejwt.tokens import RefreshToken
from rest_framework_simplejwt.views import TokenObtainPairView, TokenRefreshView
from django.conf import settings

class RegisterView(generics.CreateAPIView): # for create-only endpoints
    queryset = User.objects.all()
    permission_classes = (AllowAny,)
    serializer_class = RegisterSerializer

class LogoutView(APIView):
    permission_classes = (IsAuthenticated,)

    def post(self, request):
        try:
            refresh_token = request.COOKIES.get('refresh_token')
            token = RefreshToken(refresh_token)
            token.blacklist()  # Blacklist the refresh token

            response = Response(status=status.HTTP_205_RESET_CONTENT)
            response.delete_cookie('refresh_token')
            return response
        except Exception as e:
            return Response(status=status.HTTP_400_BAD_REQUEST)

class CookieTokenObtainPairView(TokenObtainPairView):
    def post(self, request, *args, **kwargs):
        response = super().post(request, *args, **kwargs)

        if response.status_code == 200:
            refresh_token = response.data.get("refresh")
            if refresh_token:
                secure = not settings.DEBUG
                samesite = "None" if not settings.DEBUG else "Lax"

                response.set_cookie(
                    "refresh_token",
                    refresh_token,
                    httponly=True,
                    samesite=samesite,
                    secure=secure,
                    path="/api/",
                )
                del response.data["refresh"]

        return response

class CookieTokenRefreshView(TokenRefreshView):
    serializer_class = CookieTokenRefreshSerializer

    def post(self, request, *args, **kwargs):
        response = super().post(request, *args, **kwargs)

        if response.status_code == 200:
            access_token = response.data.get("access")
            refresh_token = response.data.get("refresh")

            if access_token and refresh_token:
                secure = not settings.DEBUG
                samesite = "None" if not settings.DEBUG else "Lax"

                response.set_cookie(
                    "refresh_token",
                    refresh_token,
                    httponly=True,
                    samesite=samesite,
                    secure=secure,
                    path="/api/",
                )

        return response
