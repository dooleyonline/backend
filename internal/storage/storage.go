package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"

	"github.com/dooleyonline/backend/internal/config"
)

func Presign(ctx context.Context, cfg *config.Config, sreq *PresignParams) (string, http.Header, error) {
	if !sreq.validate() {
		return "", nil, fmt.Errorf("invalid storage request")
	}

	req, err := http.NewRequest(sreq.Method, sreq.url(cfg.StorageUrl), nil)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", sreq.ContentType)

	signer := v4.NewSigner(func(signer *v4.SignerOptions) {
		signer.DisableHeaderHoisting = true
	})

	return signer.PresignHTTP(ctx, aws.Credentials{
		AccessKeyID:     cfg.StorageAccessId,
		SecretAccessKey: cfg.StorageAccessSecret,
	},
		req,
		hashPayload(""),
		"s3",
		cfg.StorageRegion,
		time.Now(),
	)
}

type PresignParams struct {
	Method      string
	Bucket      string
	Key         string
	ContentType string
}

func (sreq *PresignParams) validate() bool {
	// only allow get and put methods
	switch sreq.Method {
	case http.MethodGet, http.MethodPut:
		// do nothing
	default:
		return false
	}

	return sreq.Bucket != "" && sreq.Key != "" && sreq.ContentType != ""
}

func (sreq *PresignParams) url(baseUrl string) string {
	return strings.Join([]string{baseUrl, sreq.Bucket, sreq.Key}, "/")
}

func hashPayload(payload string) string {
	ha := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(ha[:])
}
