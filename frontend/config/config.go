package config

import (
	"os"

	"go.uber.org/zap"
)

type Config struct {
	BackendURL string
}

func SetBackendURL() string {
	var config Config
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		zap.L().Error("BACKEND_URL environment variable is not set. Using localhost address for backend target")
		backendURL = "http://localhost:3000"
	} else {
		zap.L().Info("BACKEND_URL is set. Using BACKEND_URL from environment variable", zap.String("backendURL", backendURL))
	}

	config.BackendURL = backendURL
	return config.BackendURL
}
