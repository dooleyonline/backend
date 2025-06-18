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
    image_urls = serializers.ListField(child=serializers.CharField(), required=False) 
   
    class Meta:
        model = Item
        fields = [
            "id", "name", "description", "image_urls", "price",
            "postedAt", "isSold", "views", "condition", "isNegotiable",
            "category", "subcategory", "seller"
        ]
        read_only_fields = ["id", "postedAt", "views"]
        
class CategorySerializer(serializers.ModelSerializer):
    subcategories = serializers.ListField(
        child=serializers.CharField(max_length=20),
    )
    
    class Meta:
        model = Category
        fields = '__all__'