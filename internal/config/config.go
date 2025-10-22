package config

import (
	"os"
)

type Config struct {
	Port              string
	AuthServiceURL    string
	ProductServiceURL string
}

func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", ":8283"),
		AuthServiceURL:    getEnv("AUTH_SERVICE_URL", "http://localhost:8284"),
		ProductServiceURL: getEnv("PRODUCT_SERVICE_URL", "http://localhost:7780"),
	}
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}
