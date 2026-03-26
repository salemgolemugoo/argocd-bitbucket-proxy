package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port                  string
	ArgocdWebhookURL      string
	BitbucketServerSecret string
	BitbucketCloudSecret  string
	LogLevel              string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:                  envOrDefault("PORT", "8080"),
		ArgocdWebhookURL:      envOrDefault("ARGOCD_WEBHOOK_URL", "http://argocd-applicationset-controller.argocd.svc.cluster.local:7000/api/webhook"),
		BitbucketServerSecret: os.Getenv("BITBUCKET_SERVER_SECRET"),
		BitbucketCloudSecret:  os.Getenv("BITBUCKET_CLOUD_SECRET"),
		LogLevel:              envOrDefault("LOG_LEVEL", "info"),
	}

	if cfg.BitbucketServerSecret == "" && cfg.BitbucketCloudSecret == "" {
		return nil, fmt.Errorf("at least one of BITBUCKET_SERVER_SECRET or BITBUCKET_CLOUD_SECRET must be set")
	}

	return cfg, nil
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
