from django.db import models
from django.contrib.postgres.fields import ArrayField
from django.core.validators import MinValueValidator, MaxValueValidator


class Category(models.Model):
    name = models.CharField(max_length=20, primary_key=True)
    subcategories = ArrayField(models.CharField(max_length=20), default=list, blank=True) # List of subcategories
    # subcategories = models.CharField(max_length=20, default='', blank=True)
    icon = models.CharField(max_length=100)  # URL or path to the icon image
    
    def __str__(self):
        return self.name
    
    
    
class Item(models.Model):
    id = models.AutoField(primary_key=True)
    name = models.CharField(max_length=100)
    description = models.CharField(max_length=300)
    # image_urls = models.JSONField(default=list, blank=True)
    images = ArrayField(models.CharField(max_length=255), default=list, blank=True)
    price = models.DecimalField(max_digits=10, decimal_places=2, default=0.00)
    postedAt = models.DateTimeField(auto_now_add=True)
    isSold = models.BooleanField(default=False)
    views = models.PositiveIntegerField(default=0)
    condition = models.PositiveSmallIntegerField(default=1, validators=[MinValueValidator(1), MaxValueValidator(5)])  # 1 to 5 
    isNegotiable = models.BooleanField(default=True)
    category = models.ForeignKey(Category, on_delete=models.CASCADE, related_name='items')
    subcategory = models.CharField(max_length=20) 
    seller = models.CharField(max_length=100)  # temporary field for seller's name

    class Meta:
        ordering = ['-postedAt']    

    def __str__(self): 
        return self.name 

