package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NekruzRakhimov/api_gateway/internal/config"
)

func TestPingRoute(t *testing.T) {
	cfg := &config.Config{
		Port:              ":8080",
		AuthServiceURL:    "http://auth:8080",
		ProductServiceURL: "http://products:8081",
	}

	r := Setup(cfg)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	wantBody := `{"status":"ok"}`
	if rr.Body.String() != wantBody {
		t.Fatalf("unexpected body: got %q, want %q", rr.Body.String(), wantBody)
	}
}
