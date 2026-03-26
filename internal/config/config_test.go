package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	os.Setenv("BITBUCKET_SERVER_SECRET", "server-secret")
	defer os.Unsetenv("BITBUCKET_SERVER_SECRET")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != "8080" {
		t.Errorf("expected port 8080, got %s", cfg.Port)
	}
	if cfg.ArgocdWebhookURL != "http://argocd-applicationset-controller.argocd.svc.cluster.local:7000/api/webhook" {
		t.Errorf("unexpected default ArgoCD URL: %s", cfg.ArgocdWebhookURL)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected log level info, got %s", cfg.LogLevel)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("ARGOCD_WEBHOOK_URL", "http://custom:7000/api/webhook")
	os.Setenv("BITBUCKET_SERVER_SECRET", "my-secret")
	os.Setenv("BITBUCKET_CLOUD_SECRET", "cloud-secret")
	os.Setenv("LOG_LEVEL", "debug")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("ARGOCD_WEBHOOK_URL")
		os.Unsetenv("BITBUCKET_SERVER_SECRET")
		os.Unsetenv("BITBUCKET_CLOUD_SECRET")
		os.Unsetenv("LOG_LEVEL")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != "9090" {
		t.Errorf("expected port 9090, got %s", cfg.Port)
	}
	if cfg.ArgocdWebhookURL != "http://custom:7000/api/webhook" {
		t.Errorf("unexpected ArgoCD URL: %s", cfg.ArgocdWebhookURL)
	}
	if cfg.BitbucketServerSecret != "my-secret" {
		t.Errorf("unexpected server secret: %s", cfg.BitbucketServerSecret)
	}
	if cfg.BitbucketCloudSecret != "cloud-secret" {
		t.Errorf("unexpected cloud secret: %s", cfg.BitbucketCloudSecret)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("unexpected log level: %s", cfg.LogLevel)
	}
}

func TestLoad_NoSecrets_Error(t *testing.T) {
	os.Unsetenv("BITBUCKET_SERVER_SECRET")
	os.Unsetenv("BITBUCKET_CLOUD_SECRET")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when no secrets configured")
	}
}

func TestLoad_ServerSecretOnly_OK(t *testing.T) {
	os.Setenv("BITBUCKET_SERVER_SECRET", "s")
	os.Unsetenv("BITBUCKET_CLOUD_SECRET")
	defer os.Unsetenv("BITBUCKET_SERVER_SECRET")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BitbucketServerSecret != "s" {
		t.Errorf("unexpected server secret: %s", cfg.BitbucketServerSecret)
	}
	if cfg.BitbucketCloudSecret != "" {
		t.Errorf("expected empty cloud secret, got %s", cfg.BitbucketCloudSecret)
	}
}
