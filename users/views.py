from django.shortcuts import render
from rest_framework.generics import RetrieveUpdateAPIView, CreateAPIView
from rest_framework.views import APIView
from rest_framework.permissions import IsAuthenticated
from .serializers import CustomUserSerializer, RegisterUserSerializer, LoginUserSerializer
from rest_framework.response import Response

# Create your views here.

class UserInfoView(RetrieveUpdateAPIView):
    permission_classes = [IsAuthenticated]
    serializer_class = CustomUserSerializer

    def get_object(self):
        return self.request.user
    
class UserRegistrationView(CreateAPIView):
    serializer_class = RegisterUserSerializer

    def perform_create(self, serializer):
        user = serializer.save()
        # Optionally, you can send a welcome email or perform other actions here
        return user
    
class UserLoginView(APIView):
    permission_classes = []

    def post(self, request):
        serializer = LoginUserSerializer(data=request.data)
        if serializer.is_valid():
            user = serializer.validated_data
            # Here you would typically generate a token or session for the user
            return Response({"message": "Login successful", "user_id": user.id}, status=200)
        return Response(serializer.errors, status=400)