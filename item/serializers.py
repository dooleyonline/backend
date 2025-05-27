# serializers.py
from rest_framework import serializers
from .models import Item

class ItemSerializer(serializers.ModelSerializer):
    images = serializers.ListField(
        child=serializers.ImageField(), 
        write_only=True,
        required=False  
    ) 
    # this is not an item field, the attribute is only f
    class Meta:
        model = Item
        fields = ["id", "name", "images", "image_urls"]