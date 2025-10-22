package main

import (
	"log"
	"net/http"

	"github.com/NekruzRakhimov/api_gateway/internal/config"

	"github.com/NekruzRakhimov/api_gateway/internal/router"
)

func main() {
	cfg := config.Load()

	r := router.Setup(cfg)
	log.Printf("API Gateway listening on %s", cfg.Port)
	err := http.ListenAndServe(cfg.Port, r)
	if err != nil {
		log.Fatal(err)
		return
	}
}
