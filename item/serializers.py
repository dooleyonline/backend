# serializers.py
from rest_framework import serializers
from .models import Item

class ItemSerializer(serializers.ModelSerializer):
    images = serializers.ListField(
        child=serializers.ImageField(), # validates iamge files in the the images list
        write_only=True,
        required=False  
    ) # custom image field, not item_model 

    image_urls = serializers.ListField(read_only=True, required=False) 
    class Meta:
        model = Item
        fields = ["id", "name", "image_urls", 'images']