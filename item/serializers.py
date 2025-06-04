# serializers.py
from rest_framework import serializers
from .models import Item

class ItemSerializer(serializers.ModelSerializer):
    # images = serializers.ListField(
    #     child=serializers.ImageField(), # validates iamge files in the the images list
    #     write_only=True,
    #     required=False  
    # ) 
    image_urls = serializers.ListField(read_only=True) 
    class Meta:
        model = Item
        fields = ["id", "description", "name", "image_urls"]