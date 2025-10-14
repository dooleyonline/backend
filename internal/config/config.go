package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerAddr          string
	DatabaseUrl         string
	StorageUrl          string
	PublicUrl           string
	StorageRegion       string
	StorageAccessId     string
	StorageAccessSecret string
}

const (
	envPort                = "PORT"
	envDatabaseUrl         = "DATABASE_URL"
	envStorageUrl          = "STORAGE_URL"
	envPublicUrl           = "PUBLIC_URL"
	envStorageRegion       = "STORAGE_REGION"
	envStorageAccessId     = "STORAGE_ACCESS_ID"
	envStorageAccessSecret = "STORAGE_ACCESS_SECRET"
)

func New() (*Config, error) {
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

	publicUrl, ok := os.LookupEnv(envPublicUrl)
	if !ok {
		return nil, fmt.Errorf("environment variable %s is required", envPublicUrl)
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

	cfg := &Config{
		ServerAddr:          ":" + port,
		DatabaseUrl:         databaseUrl,
		StorageUrl:          storageUrl,
		StorageRegion:       storageRegion,
		StorageAccessId:     storageAccessId,
		StorageAccessSecret: storageAccessSecret,
		PublicUrl:          publicUrl,
	}

	return cfg, nil
}
