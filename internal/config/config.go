package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	IsProd     bool
	Url        string
	ServerAddr string

	DatabaseUrl string

	StorageUrl          string
	StorageS3Url        string
	StorageRegion       string
	StorageAccessId     string
	StorageAccessSecret string

	AuthTokenName   string
	AuthTokenExp    time.Duration
	AuthTokenSecret string
}

const (
	envEnv  = "ENV"
	envPort = "PORT"

	envDatabaseUrl = "DATABASE_URL"

	envStorageUrl          = "STORAGE_URL"
	envStorageS3Url        = "STORAGE_S3_URL"
	envStorageRegion       = "STORAGE_REGION"
	envStorageAccessId     = "STORAGE_ACCESS_ID"
	envStorageAccessSecret = "STORAGE_ACCESS_SECRET"

	envAuthTokenSecret = "AUTH_TOKEN_SECRET"
)

func New() (*Config, error) {
	env, ok := os.LookupEnv(envEnv)
	if !ok {
		return nil, fmt.Errorf("environment variable %s is required", envEnv)
	}

	port, ok := os.LookupEnv(envPort)
	if !ok {
		return nil, fmt.Errorf("environment variable %s is required", envPort)
	}

	databaseUrl, ok := os.LookupEnv(envDatabaseUrl)
	if !ok {
		return nil, fmt.Errorf("environment variable %s is required", envDatabaseUrl)
	}

	storageUrl, ok := os.LookupEnv(envStorageUrl)
	if !ok {
		return nil, fmt.Errorf("environment variable %s is required", envStorageUrl)
	}

	storageS3Url, ok := os.LookupEnv(envStorageS3Url)
	if !ok {
		return nil, fmt.Errorf("environment variable %s is required", envStorageS3Url)
	}

	storageRegion, ok := os.LookupEnv(envStorageRegion)
	if !ok {
		return nil, fmt.Errorf("environment variable %s is required", envStorageRegion)
	}

	storageAccessId, ok := os.LookupEnv(envStorageAccessId)
	if !ok {
		return nil, fmt.Errorf("environment variable %s is required", envStorageAccessId)
	}

	storageAccessSecret, ok := os.LookupEnv(envStorageAccessSecret)
	if !ok {
		return nil, fmt.Errorf("environment variable %s is required", envStorageAccessSecret)
	}

	authTokenSecret, ok := os.LookupEnv(envAuthTokenSecret)
	if !ok {
		return nil, fmt.Errorf("environment variable %s is required", envAuthTokenSecret)
	}

	cfg := &Config{
		IsProd:     env == "prod",
		Url:        "https://api.dooleyonline.net",
		ServerAddr: ":" + port,

		DatabaseUrl: databaseUrl,

		StorageS3Url:        storageS3Url,
		StorageRegion:       storageRegion,
		StorageAccessId:     storageAccessId,
		StorageAccessSecret: storageAccessSecret,

		AuthTokenName:   "dooleyonline_jwt",
		AuthTokenExp:    time.Hour * 240,
		AuthTokenSecret: authTokenSecret,
		StorageUrl:      storageUrl,
	}

	return cfg, nil
}
