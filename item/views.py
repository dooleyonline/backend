# items/views.py
from rest_framework import viewsets, status
from .models import Item
from django.core.files.storage import default_storage
from .serializers import ItemSerializer
from rest_framework.response import Response
from rest_framework.parsers import MultiPartParser, FormParser
from django.shortcuts import render
from django.conf import settings
from django.core.files.storage import FileSystemStorage
import uuid
from drf_yasg.utils import swagger_auto_schema
from drf_yasg import openapi
from django_filters.rest_framework import DjangoFilterBackend
from rest_framework import filters
from rest_framework.decorators import action
from django.db.models import Q



class ItemViewSet(viewsets.ModelViewSet):
    queryset = Item.objects.all().order_by("-id")
    serializer_class = ItemSerializer 
    parser_classes = [MultiPartParser, FormParser]
    # filter_backends = [filters.SearchFilter, DjangoFilterBackend]
    # search_fields = ['name'] # add description if needed
    # filterset_fields = ['category__category_name']
    
    def get_queryset(self):
        queryset = super().get_queryset()
        q = self.request.query_params.get('q')
        category = self.request.query_params.get('category')

        if q:
            queryset = queryset.filter(Q(name__icontains=q)) # | Q(description__icontains=q))
        if category:
            queryset = queryset.filter(category__category_name__iexact=category)

        return queryset
    
    @action(detail=False, methods=['get'], url_path='ids')
    def get_ids(self, request):
        ids = list(Item.objects.values_list('id', flat=True))
        return Response({'ids': ids})
    
    def perform_file_upload(self, file):
        """
        Handles the file upload and returns the path where the file is stored.
        """
        filename = f"{uuid.uuid4()}_{file.name}"
        if settings.USE_S3:
            return default_storage.save(f"uploads/{filename}", file)
        else:
            fs = FileSystemStorage()
            return fs.save(filename, file)
        
    @swagger_auto_schema(
        manual_parameters=[
            openapi.Parameter('q', openapi.IN_QUERY, type=openapi.TYPE_STRING, description='Keyword search'),
            openapi.Parameter('category', openapi.IN_QUERY, type=openapi.TYPE_STRING, description='Filter by category name'),
        ]
    )
    def list(self, request, *args, **kwargs):
        return super().list(request, *args, **kwargs)

    @swagger_auto_schema(
        manual_parameters=[
            openapi.Parameter('name', openapi.IN_FORM, type=openapi.TYPE_STRING),
            openapi.Parameter('description', openapi.IN_FORM, type=openapi.TYPE_STRING),
            openapi.Parameter(
                'images',
                openapi.IN_FORM,
                type=openapi.TYPE_FILE,
                description="Multiple image files",
                required=True,
                collection_format='multi',
            ),
        ]
    )
    def create(self, request, *args, **kwargs):
        serializer = ItemSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)

        images = request.FILES.getlist("images")  
        image_urls = []
        self.validate_images(images)
        for img in images:
            path = self.perform_file_upload(img)
            # path = default_storage.save(f"uploads/{img.name}", f=img)
            url = default_storage.url(path) if settings.USE_S3 else FileSystemStorage().url(path)
            image_urls.append(default_storage.url(path))

        item = Item.objects.create(
            name=serializer.validated_data["name"],
            description=serializer.validated_data["description"],
            image_urls=image_urls
        )
        output = ItemSerializer(item)
        return Response(output.data, status=status.HTTP_201_CREATED)
    
    def validate_images(self, images):
        """
        Validates the uploaded images.
        """
        if not images:
            raise serializers.ValidationError("No images provided.")
        
        for img in images:
            if not img.name.lower().endswith(('.png', '.jpg', '.jpeg')):
                raise serializers.ValidationError(f"Invalid image format: {img.name}")
            if img.size > 5 * 1024 * 1024:
                raise ValidationError(f"{img.name} exceeds 5MB limit.")

def upload_item_page(request):
    return render(request, "upload_item.html")
