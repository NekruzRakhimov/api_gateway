package proxy

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestProxy_ForwardsRequestAndSetsXForwardedHost(t *testing.T) {
	// Поднимаем фейковый апстрим
	var gotReq *http.Request

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotReq = r
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer upstream.Close()

	u, err := url.Parse(upstream.URL)
	if err != nil {
		t.Fatalf("failed to parse upstream url: %v", err)
	}

	h := New(u)

	// Запрос в прокси
	req := httptest.NewRequest(http.MethodGet, "/api/test?x=1", nil)
	req.Host = "gateway.local"

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 from proxy, got %d", rr.Code)
	}

	if gotReq == nil {
		t.Fatalf("upstream did not receive request")
	}

	if gotReq.URL.Path != "/api/test" {
		t.Fatalf("unexpected upstream path: got %q, want %q", gotReq.URL.Path, "/api/test")
	}

	// ВАЖНО: сейчас реализация ставит X-Forwarded-Host = target.Host (u.Host)
	if gotReq.Header.Get("X-Forwarded-Host") != u.Host {
		t.Fatalf("expected X-Forwarded-Host=%s, got %q",
			u.Host, gotReq.Header.Get("X-Forwarded-Host"))
	}
}
