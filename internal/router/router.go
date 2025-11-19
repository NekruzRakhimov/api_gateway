package router

import (
	"fmt"
	// std
	"net/http"
	"net/url"

	// internal
	"github.com/NekruzRakhimov/api_gateway/internal/config"
	"github.com/NekruzRakhimov/api_gateway/internal/proxy"
	"github.com/go-chi/chi/v5"
)

func mount(r chi.Router, prefix string, target string) {
	u, _ := url.Parse(target)
	h := proxy.New(u)  // ваш http.Handler на базе httputil.NewSingleHostReverseProxy
	r.Mount(prefix, h) // chi сам срежет prefix для h
}

func Setup(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			fmt.Println(err)
			return
		}
	})

	mount(r, "/auth", cfg.AuthServiceURL)
	mount(r, "/api/products", cfg.ProductServiceURL)

	return r
}
