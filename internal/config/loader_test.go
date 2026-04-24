package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/driftwatch/internal/config"
)

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return path
}

func TestLoad_YAML(t *testing.T) {
	path := writeTempFile(t, "svc.yaml", `
name: api-server
version: "1.2.3"
image: myrepo/api:1.2.3
replicas: 3
environment:
  LOG_LEVEL: info
  PORT: "8080"
ports:
  - "8080:8080"
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Name != "api-server" {
		t.Errorf("expected name %q, got %q", "api-server", cfg.Name)
	}
	if cfg.Replicas != 3 {
		t.Errorf("expected replicas 3, got %d", cfg.Replicas)
	}
	if cfg.Environment["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %q", cfg.Environment["LOG_LEVEL"])
	}
}

func TestLoad_TOML(t *testing.T) {
	path := writeTempFile(t, "svc.toml", `
name = "worker"
version = "0.9.1"
image = "myrepo/worker:0.9.1"
replicas = 1

[environment]
  QUEUE = "tasks"

ports = []
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Name != "worker" {
		t.Errorf("expected name %q, got %q", "worker", cfg.Name)
	}
}

func TestLoad_MissingName(t *testing.T) {
	path := writeTempFile(t, "bad.yaml", `version: "1.0.0"`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing name, got nil")
	}
}

func TestLoad_UnsupportedExtension(t *testing.T) {
	path := writeTempFile(t, "svc.json", `{"name":"test"}`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for unsupported extension, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/svc.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
