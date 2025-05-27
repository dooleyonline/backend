# items/views.py
from rest_framework import viewsets, status
from .models import Item
from django.core.files.storage import default_storage
from .serializers import ItemSerializer
from rest_framework.response import Response
from rest_framework.parsers import MultiPartParser, FormParser
from django.shortcuts import render

class ItemViewSet(viewsets.ModelViewSet):
    queryset = Item.objects.all()
    serializer_class = ItemSerializer 
    parser_classes = [MultiPartParser, FormParser]

    def create(self, request, *args, **kwargs):
        serializer = self.get_serializer(data=request.data)
        serializer.is_valid(raise_exception=True)

        images = serializer.validated_data.pop("images", [])
        image_urls = []
        for f in images:
            path = default_storage.save(f"uploads/{f.name}", f)
            image_urls.append(default_storage.url(path))

        item = Item.objects.create(
            name=serializer.validated_data["name"],
            image_urls=image_urls
        )
        out = ItemSerializer(item)
        return Response(out.data, status=status.HTTP_201_CREATED)

def upload_item_page(request):
    return render(request, "upload_item.html")
