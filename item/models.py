from django.db import models

class Item(models.Model):
    id = models.AutoField(primary_key=True)
    name = models.CharField(max_length=100)
    description = models.CharField(max_length=300)
    image_urls = models.JSONField(default=list, blank=True)
    
    def __str__(self): 
        return self.name 

