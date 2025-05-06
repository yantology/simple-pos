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
	port := "3000" // Default port

	// Prioritize APP_PORT environment variable (set in cloudrun.tf)
	env_APP_PORT := os.Getenv("APP_PORT")
	if env_APP_PORT != "" {
		log.Println("Using APP_PORT from environment => ", env_APP_PORT)
		port = env_APP_PORT
	} else {
		// Fallback to PORT if APP_PORT is not set (standard for Cloud Run, Heroku, etc.)
		env_PORT := os.Getenv("PORT")
		if env_PORT != "" {
			log.Println("Using PORT from environment => ", env_PORT)
			port = env_PORT
		} else {
			log.Println("Using default port => ", port)
		}
	}

	return &AppConfig{
		Port:            port,
		PublicRoute:     "/public",
		PublicAssetsDir: "./public",
	}
}
