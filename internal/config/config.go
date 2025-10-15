package config

import (
	"fmt"
	"os"
)

type Config struct {
	Env                 string
	ServerAddr          string
	DatabaseUrl         string
	StorageUrl          string
	StorageRegion       string
	StorageAccessId     string
	StorageAccessSecret string
	HmacSecretKey       string
}

const (
	envEnv                 = "ENV"
	envPort                = "PORT"
	envDatabaseUrl         = "DATABASE_URL"
	envStorageUrl          = "STORAGE_URL"
	envStorageRegion       = "STORAGE_REGION"
	envStorageAccessId     = "STORAGE_ACCESS_ID"
	envStorageAccessSecret = "STORAGE_ACCESS_SECRET"
	envHmacSecretKey       = "HMAC_SECRET_KEY"
)

func New() (*Config, error) {
	env, ok := os.LookupEnv(envEnv)
	if !ok {
		return nil, fmt.Errorf("environment variable %s is required", envPort)
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

	hmacSecretKey, ok := os.LookupEnv(envHmacSecretKey)
	if !ok {
		return nil, fmt.Errorf("environment variable %s is required", envHmacSecretKey)
	}

	cfg := &Config{
		Env:                 env,
		ServerAddr:          ":" + port,
		DatabaseUrl:         databaseUrl,
		StorageUrl:          storageUrl,
		StorageRegion:       storageRegion,
		StorageAccessId:     storageAccessId,
		StorageAccessSecret: storageAccessSecret,
		HmacSecretKey:       hmacSecretKey,
	}

	return cfg, nil
}
