from django.contrib import admin
from .models import Item, Category
from django.contrib.auth.models import Group

# Register your models here.

@admin.register(Item)
class ItemAdmin(admin.ModelAdmin):
    list_display = ('pk', 'name', 'description')
    list_filter = ('category',)
    
@admin.register(Category)
class CategoryAdmin(admin.ModelAdmin):
    list_display = ('pk', 'name', 'subcategories')

admin.site.unregister(Group)  # Unregister the Group model 

admin.site.site_header = "Dooley Online Admin"