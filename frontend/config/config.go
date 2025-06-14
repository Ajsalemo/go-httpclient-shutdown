package config

import (
	"os"
)

type Config struct {
	BackendURL string
}

func GetBackendURL() string {
	var config Config
	url := os.Getenv("BACKEND_URL")
	// If BACKEND_URL is not set, default to localhost address
	if url == "" {
		config.BackendURL = "http://localhost:3000"
		return config.BackendURL
	} else {
		config.BackendURL = url
		return config.BackendURL
	}
}

