package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/google/uuid"

	"github.com/dooleyonline/backend/internal/config"
)

type Storage struct {
	cfg    *config.Config
	signer *v4.Signer
}

func New(cfg *config.Config) *Storage {
	signer := v4.NewSigner(func(signer *v4.SignerOptions) {
		signer.DisableHeaderHoisting = true
	})
	return &Storage{cfg, signer}
}

type PresignResult struct {
	URL     string      `json:"url"`
	Headers http.Header `json:"headers"`
	ImageID string      `json:"image_id"`
}

func (s *Storage) PresignUpload(ctx context.Context, contentType string) (*PresignResult, error) {
	key := s.GenerateImageID()

	url, err := s.buildURL(key)
	if err != nil {
		return nil, fmt.Errorf("failed to build url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)
	signedURI, signedHeaders, err := s.signer.PresignHTTP(ctx,
		s.credentials(),
		req,
		hashPayload(""),
		"s3",
		s.cfg.StorageRegion,
		time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to presign request: %w", err)
	}

	res := &PresignResult{
		URL:     signedURI,
		Headers: signedHeaders,
		ImageID: key,
	}
	return res, nil
}

func (s *Storage) GenerateImageID() string {
	return uuid.NewString()
}

func (s *Storage) credentials() aws.Credentials {
	return aws.Credentials{
		AccessKeyID:     s.cfg.StorageAccessId,
		SecretAccessKey: s.cfg.StorageAccessSecret,
	}
}

func (s *Storage) buildURL(key string) (string, error) {
	return url.JoinPath(s.cfg.StorageS3Url, s.cfg.StorageBucket, key)
}

func hashPayload(payload string) string {
	ha := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(ha[:])
}
