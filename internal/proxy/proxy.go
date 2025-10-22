package proxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// New создаёт http.Handler, который проксирует на target.
// Работает с chi.Mount: chi срежет mount-префикс и передаст оставшийся путь сюда.
// Если target.URL содержит basePath (например, http://svc:8080/api/v1),
// стандартный Director аккуратно склеит его с путём запроса.
func New(target *url.URL) http.Handler {
	rp := httputil.NewSingleHostReverseProxy(target)

	// Транспорт с таймаутами (безопасные дефолты)
	rp.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	}

	// Сохраняем исходные заголовки клиента и корректно выставляем Host апстрима
	origDirector := rp.Director
	rp.Director = func(r *http.Request) {
		origDirector(r)

		// Host апстрима (важно для виртуальных хостов)
		r.Host = target.Host

		// Проксируема схемa/хост установлены Director-ом, путь/квери уже склеены:
		//   finalPath = singleJoiningSlash(target.Path, r.URL.Path)
		//   r.URL.RawQuery остаётся нетронутой

		// X-Forwarded-*
		if xf := r.Header.Get("X-Forwarded-For"); xf == "" {
			if h, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
				r.Header.Set("X-Forwarded-For", h)
			} else {
				r.Header.Set("X-Forwarded-For", r.RemoteAddr)
			}
		}
		if r.Header.Get("X-Forwarded-Proto") == "" {
			if r.TLS != nil {
				r.Header.Set("X-Forwarded-Proto", "https")
			} else {
				r.Header.Set("X-Forwarded-Proto", "http")
			}
		}
		if r.Header.Get("X-Forwarded-Host") == "" {
			r.Header.Set("X-Forwarded-Host", r.Host)
		}
	}

	// Акуратная обработка ошибок апстрима
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "upstream error: "+err.Error(), http.StatusBadGateway)
	}

	return rp
}
