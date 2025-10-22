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

type Client struct {
	cfg *config.Config
	signer *v4.Signer
}

func NewClient(cfg *config.Config) *Client {
	signer := v4.NewSigner(func(signer *v4.SignerOptions) {
		signer.DisableHeaderHoisting = true
	})
	return &Client{cfg, signer}
}


func (c *Client) Presign(ctx context.Context, sreq *PresignRequest) (string, http.Header, error) {
	url := c.buildURL(sreq.Bucket, sreq.Key)
	req, err := http.NewRequestWithContext(ctx, sreq.Method, url, nil)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", sreq.ContentType)

	return c.signer.PresignHTTP(ctx, 
		c.credentials(),
		req,
		hashPayload(""),
		"s3",
		c.cfg.StorageRegion,
		time.Now(),
	)
}

func (c *Client) credentials() aws.Credentials {
	return aws.Credentials{
		AccessKeyID:     c.cfg.StorageAccessId,
		SecretAccessKey: c.cfg.StorageAccessSecret,
	}
}

func (c *Client) buildURL(bucket, key string) string {
	return strings.Join([]string{c.cfg.StorageS3Url, bucket, key}, "/")
}

func hashPayload(payload string) string {
	ha := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(ha[:])
}
