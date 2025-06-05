from django.db import models
from django.contrib.postgres.fields import ArrayField


class Category(models.Model):
    category_name = models.CharField(max_length=20, primary_key=True)
    # subcategories = ArrayField(models.CharField(max_length=20), default=list, blank=True)
    subcategories = models.CharField(max_length=20, default='', blank=True)

    def __str__(self):
        return self.category_name
    
    
    
class Item(models.Model):
    id = models.AutoField(primary_key=True)
    name = models.CharField(max_length=100)
    description = models.CharField(max_length=300)
    image_urls = models.JSONField(default=list, blank=True)
    
    category = models.ForeignKey(Category, on_delete=models.CASCADE)
    
    def __str__(self): 
        return self.name 

