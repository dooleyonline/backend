# serializers.py
from rest_framework import serializers
from .models import Item, Category

class ItemSerializer(serializers.ModelSerializer):
    # images = serializers.ListField(
    #     child=serializers.ImageField(), # validates iamge files in the the images list
    #     write_only=True,
    #     required=False  
    # ) 
    category = serializers.SlugRelatedField(
        queryset=Category.objects.all(),
        slug_field='category_name'
    )
    image_urls = serializers.ListField(read_only=True) 
    class Meta:
        model = Item
        fields = ["id", "description", "name", "image_urls", 'category']
        
class CategorySerializer(serializers.ModelSerializer):
    class Meta:
        model = Category
        fields = '__all__'