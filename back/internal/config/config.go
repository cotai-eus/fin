package config

import (
	"fmt"
	"os"
)

// Config holds application configuration
type Config struct {
	Environment string
	Port        string
	DatabaseURL string

	// APISIX Gateway
	TrustedProxyIP string

	// Encryption
	EncryptionKey string // 32 bytes for AES-256
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Environment:    getEnv("ENVIRONMENT", "development"),
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		TrustedProxyIP: getEnv("TRUSTED_PROXY_IP", "127.0.0.1"),
		EncryptionKey:  getEnv("ENCRYPTION_KEY", ""),
	}

	// Validate required fields
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	if cfg.EncryptionKey == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY environment variable is required (32 bytes)")
	}

	if len(cfg.EncryptionKey) != 32 {
		return nil, fmt.Errorf("ENCRYPTION_KEY must be exactly 32 bytes, got %d", len(cfg.EncryptionKey))
	}

	return cfg, nil
}

// getEnv retrieves environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
