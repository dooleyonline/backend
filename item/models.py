from django.db import models

class Category(models.Model):
    category_name = models.CharField(max_length=100, primary_key=True)
    category_desc = models.TextField(blank=True)

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

