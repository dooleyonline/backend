package storagesvc

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"

	"github.com/stretchr/testify/assert"
)

func TestPresignUploadText(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal("failed to initialize config:", err)
	}

	db, err := db.New(t.Context(), cfg)
	if err != nil {
		t.Fatal("failed to initialize db:", err)
	}

	contentType := "text/plain"
	bucket := "item"

	storage := New(cfg, db)

	presign, err := storage.PresignUpload(t.Context(), PresignParams{
		ContentType: contentType,
		Bucket:      bucket,
	})
	if err != nil {
		t.Fatal("failed to presign:", err)
	}

	payload := "hello!"
	req, err := http.NewRequest(http.MethodPut, presign.URL, bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal("failed to create request:", err)
	}
	req.Header.Set("Content-Type", contentType)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("failed to make request:", err)
	}
	defer res.Body.Close()

	assert := assert.New(t)
	assert.Equal(http.StatusOK, res.StatusCode)
}

func TestPresignUploadImage(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal("failed to initialize config:", err)
	}

	payload, err := os.Open("../../../test.png")
	if err != nil {
		t.Fatal("failed to open test image:", err)
	}
	defer payload.Close()

	contentType := "image/png"
	bucket := "item"


	db, err := db.New(t.Context(), cfg)
	if err != nil {
		t.Fatal("failed to initialize db:", err)
	}

	storage := New(cfg, db)
	presign, err := storage.PresignUpload(t.Context(), PresignParams{
		ContentType: contentType,
		Bucket:      bucket,
	})
	if err != nil {
		t.Fatal("failed to presign:", err)
	}

	req, err := http.NewRequest(http.MethodPut, presign.URL, payload)
	if err != nil {
		t.Fatal("failed to create request:", err)
	}
	req.Header.Set("Content-Type", contentType)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("failed to make request:", err)
	}
	defer res.Body.Close()

	assert := assert.New(t)
	assert.Equal(http.StatusOK, res.StatusCode)
}
