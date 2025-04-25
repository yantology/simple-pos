package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yantology/simple-ecommerce/config"
	_ "github.com/yantology/simple-ecommerce/docs"
	"github.com/yantology/simple-ecommerce/middleware"
	"github.com/yantology/simple-ecommerce/pkg/jwt"
	"github.com/yantology/simple-ecommerce/pkg/resendutils"
	"github.com/yantology/simple-ecommerce/routes/auth"
	"github.com/yantology/simple-ecommerce/routes/category"
	"github.com/yantology/simple-ecommerce/routes/order"
	"github.com/yantology/simple-ecommerce/routes/product"
)

// initMigrations initializes and runs database migrations
func initMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// @title           Retail Pro API
// @version         1.0
// @description     This is a retail management system server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api/v1
// @schemes   http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	log.Println("Starting Retail Pro Backend Service...")

	// Initialize configurations
	appConfig := config.InitAppConfig()
	dbConfig := config.InitDatabaseConfig()
	jwtConfig, err := config.InitJWTConfig()
	tokenConfig := config.InitTokenConfig()
	if err != nil {
		log.Fatal("Failed to initialize JWT config:", err)
	}
	resendConfig, err := config.InitResendConfig()
	if err != nil {
		log.Fatal("Failed to initialize Resend config:", err)
	}

	db := config.ConnectDatabase(dbConfig, sql.Open)
	defer db.Close()

	// Run database migrations
	if err := initMigrations(db); err != nil {
		log.Fatal(err)
	}

	jwtService := jwt.NewJWTService(
		jwtConfig.AccessSecret,
		jwtConfig.RefreshSecret,
		jwtConfig.AccessDuration,
		jwtConfig.RefreshDuration,
		jwtConfig.Issuer,
	)
	if err != nil {
		log.Fatal("Failed to initialize JWT service:", err)
	}
	emailSender := resendutils.NewResendUtils(resendConfig.ApiKey, resendConfig.ResendDomain)

	// Initialize Auth middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService, tokenConfig)

	// Initialize Gin router with CORS configuration
	router := gin.Default()
	router.Use(config.CorsConfig())
	router.Use(gin.Recovery())

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		emailTemplate := auth.NewEmailTemplate()
		authPostgres := auth.NewAuthPostgres(db)
		authRepo := auth.NewAuthRepository(authPostgres)
		authService := auth.NewAuthService(jwtService, tokenConfig)
		authHandler := auth.NewAuthHandler(authService, authRepo, emailSender, emailTemplate, tokenConfig)
		authHandler.RegisterRoutes(v1)

		// Category routes (protected by auth middleware)
		categoryPostgres := category.NewPostgresRepository(db)           // Corrected: NewPostgresRepository
		categoryRepo := category.NewCategoryRepository(categoryPostgres) // Corrected: NewCategoryRepository
		categoryHandler := category.NewCategoryHandler(categoryRepo)
		categoryGroup := v1.Group("/categories")
		categoryGroup.Use(authMiddleware.AuthRequired()) // Apply auth middleware to all category routes
		{
			categoryGroup.GET("/", categoryHandler.GetAllCategories)
			categoryGroup.GET("/:id", categoryHandler.GetCategoryByID)
			categoryGroup.GET("/name/:name", categoryHandler.GetCategoryByName)
			categoryGroup.POST("/", categoryHandler.CreateCategory)
			categoryGroup.PUT("/:id", categoryHandler.UpdateCategory)
			categoryGroup.DELETE("/:id", categoryHandler.DeleteCategory)
		}

		// Product routes (protected by auth middleware)
		productPostgres := product.NewPostgresRepository(db)  // Corrected: NewPostgresRepository
		productRepo := product.NewRepository(productPostgres) // Corrected: NewRepository
		productHandler := product.NewHandler(productRepo)
		productGroup := v1.Group("/products")
		productGroup.Use(authMiddleware.AuthRequired()) // Apply auth middleware to all product routes
		{
			productGroup.POST("", productHandler.CreateProduct)
			productGroup.GET("", productHandler.GetAllProducts) // Consider if this should be public or user-specific
			productGroup.PUT("/:id", productHandler.UpdateProduct)
			productGroup.DELETE("/:id", productHandler.DeleteProduct)
			productGroup.GET("/user/:userID", productHandler.GetProductsByUserID) // Keep for admin/specific use, ensure proper authorization
			productGroup.GET("/category/:categoryID", productHandler.GetProductsByCategoryID)
		}

		// Order routes (protected by auth middleware)
		orderPostgres := order.NewPostgresRepository(db)     // Corrected: NewPostgresRepository
		orderRepo := order.NewOrderRepository(orderPostgres) // Corrected: NewOrderRepository
		orderHandler := order.NewOrderHandler(orderRepo)
		orderGroup := v1.Group("/orders")
		orderGroup.Use(authMiddleware.AuthRequired()) // Apply auth middleware to all order routes
		{
			orderGroup.GET("/", orderHandler.GetOrders)
			orderGroup.GET("/:id", orderHandler.GetOrderByID)
			orderGroup.POST("/", orderHandler.CreateOrder)
			orderGroup.DELETE("/:id", orderHandler.DeleteOrder)
		}

	}

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Retail Pro Backend Service",
		})
	})

	// Start the server with configured port
	serverAddr := fmt.Sprintf(":%s", appConfig.Port)
	log.Printf("Server is running on %s...\n", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
