package config

import (
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var defaultOrigins = []string{"*"}

func CorsConfig() gin.HandlerFunc {
	origins := GetEnvAsSlice("CORS_ALLOW_ORIGINS", defaultOrigins)
	log.Println("CORS_ALLOW_ORIGINS => ", strings.Join(origins, ", "))

	config := cors.DefaultConfig()
	config.AllowOrigins = origins
	return cors.New(config)
}

// GetEnvAsSlice returns environment variable as a slice of strings
func GetEnvAsSlice(name string, defaultVal []string) []string {
	valStr := os.Getenv(name)
	if valStr == "" {
		return defaultVal
	}
	return strings.Split(valStr, ",")
}
