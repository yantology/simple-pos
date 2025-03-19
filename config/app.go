package config

import (
	"log"
	"os"
)

// AppConfig holds application configuration
type AppConfig struct {
	Port            string
	PublicRoute     string
	PublicAssetsDir string
}

// InitAppConfig initializes and returns a new AppConfig
func InitAppConfig() *AppConfig {
	port := "3000"

	env_APP_PORT := os.Getenv("APP_PORT")
	if env_APP_PORT != "" {
		log.Println("APP_PORT => ", env_APP_PORT)
		port = env_APP_PORT
	}

	return &AppConfig{
		Port:            port,
		PublicRoute:     "/public",
		PublicAssetsDir: "./public",
	}
}
