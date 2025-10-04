package storage

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"github.com/dooleyonline/backend/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestPresignUploadText(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal("failed to initialize config:", err)
	}

	sreq := &StorageRequest{
		Method:      http.MethodPut,
		Bucket:      "image",
		Key:         "test-plain.txt",
		ContentType: "text/plain",
	}

	url, header, err := Presign(t.Context(), cfg, sreq)
	if err != nil {
		t.Fatal("failed to presign:", err)
	}

	payload := "hello!"
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal("failed to create request:", err)
	}
	req.Header = header

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

	payload, err := os.Open("../../test.png")
	if err != nil {
		t.Fatal("failed to open test image:", err)
	}
	defer payload.Close()

	sreq := &StorageRequest{
		Method:      http.MethodPut,
		Bucket:      "image",
		Key:         "test-image.png",
		ContentType: "image/png",
	}

	url, header, err := Presign(t.Context(), cfg, sreq)
	if err != nil {
		t.Fatal("failed to presign:", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, payload)
	if err != nil {
		t.Fatal("failed to create request:", err)
	}
	req.Header = header

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("failed to make request:", err)
	}
	defer res.Body.Close()

	assert := assert.New(t)
	assert.Equal(http.StatusOK, res.StatusCode)
}
