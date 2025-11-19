package config

import (
	"os"
	"testing"
)

func Test_getEnv_DefaultValue(t *testing.T) {
	const key = "TEST_PORT_KEY"
	const def = ":1234"

	_ = os.Unsetenv(key)

	got := getEnv(key, def)
	if got != def {
		t.Fatalf("expected default %q, got %q", def, got)
	}
}

func Test_getEnv_EnvOverridesDefault(t *testing.T) {
	const key = "TEST_PORT_KEY"
	const def = ":1234"
	const val = ":9999"

	if err := os.Setenv(key, val); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer os.Unsetenv(key)

	got := getEnv(key, def)
	if got != val {
		t.Fatalf("expected %q from env, got %q", val, got)
	}
}

func TestLoad_UsesEnvValues(t *testing.T) {
	_ = os.Setenv("PORT", ":9000")
	_ = os.Setenv("AUTH_SERVICE_URL", "http://auth:8080")
	_ = os.Setenv("PRODUCT_SERVICE_URL", "http://products:8081")

	defer func() {
		_ = os.Unsetenv("PORT")
		_ = os.Unsetenv("AUTH_SERVICE_URL")
		_ = os.Unsetenv("PRODUCT_SERVICE_URL")
	}()

	cfg := Load()

	if cfg.Port != ":9000" {
		t.Fatalf("expected Port=:9000, got %q", cfg.Port)
	}
	if cfg.AuthServiceURL != "http://auth:8080" {
		t.Fatalf("expected AuthServiceURL=http://auth:8080, got %q", cfg.AuthServiceURL)
	}
	if cfg.ProductServiceURL != "http://products:8081" {
		t.Fatalf("expected ProductServiceURL=http://products:8081, got %q", cfg.ProductServiceURL)
	}
}
