package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/yantology/retail-pro-be/cmd/docs" // Updated path to use cmd/docs instead of just docs
)

// @title           Retail Pro API
// @version         1.0
// @description     This is a retail management system server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
// @schemes   http
func main() {
	log.Println("Starting Retail Pro Backend Service...")

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS if needed
	router.Use(gin.Recovery())

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Welcome to Retail Pro Backend Service API v1",
			})
		})
	}

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Retail Pro Backend Service",
		})
	})

	// Start the server
	log.Println("Server is running on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
