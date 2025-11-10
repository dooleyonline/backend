package storagesvc

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
	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	"github.com/google/uuid"
)

type Service struct {
	cfg    *config.Config
	db     *db.DB
	signer *v4.Signer
}

func New(cfg *config.Config, db *db.DB) *Service {
	signer := v4.NewSigner(func(signer *v4.SignerOptions) {
		signer.DisableHeaderHoisting = true
	})
	return &Service{cfg, db, signer}
}

type PresignResult struct {
	URL     string      `json:"url"`
	Headers http.Header `json:"headers"`
	ImageID string      `json:"image_id"`
}

type PresignParams struct {
	ContentType string
	Bucket      string
}

func (s *Service) PresignUpload(ctx context.Context, params PresignParams) (*PresignResult, error) {
	key := generateImageID()

	url, err := buildURL(key, params.Bucket, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", params.ContentType)
	signedURI, signedHeaders, err := s.signer.PresignHTTP(ctx,
		credentials(s.cfg),
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

func generateImageID() string {
	return uuid.NewString()
}

func credentials(cfg *config.Config) aws.Credentials {
	return aws.Credentials{
		AccessKeyID:     cfg.StorageAccessId,
		SecretAccessKey: cfg.StorageAccessSecret,
	}
}

func buildURL(key string, bucket string, cfg *config.Config) (string, error) {
	return url.JoinPath(cfg.StorageS3Url, bucket, key)
}

func hashPayload(payload string) string {
	ha := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(ha[:])
}
