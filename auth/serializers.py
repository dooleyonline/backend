from rest_framework_simplejwt.serializers import TokenObtainPairSerializer

class CustomTokenObtainPairSerializer(TokenObtainPairSerializer):
    
    @classmethod
    def get_token(cls, user):
        token = super(CustomTokenObtainPairSerializer, cls).get_token(user)

        # Add custom claims
        token['username'] = user.username
        token['email'] = user.email

        return token
    
