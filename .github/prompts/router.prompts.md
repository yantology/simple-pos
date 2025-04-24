# Go Router Architecture Guide

Below is a comprehensive guide for creating a structured router in a Go application following clean architecture principles.

## Component Structure

### 1. DTO (Data Transfer Objects)
- Contains request/response structures
- Uses the common dto package for standardized responses
- Example usage:

```go
// dto/user_dto.go
type UserRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

// For responses, utilize the common dto package
// c.JSON(http.StatusOK, dto.MessageResponse{Message: "Success"})
```

### 2. Database Layer
- Implements database-specific queries
- Follows interfaces defined in interface.go
- Example:

```go
// postgres.go
type postgresRepository struct {
    db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *postgresRepository {
    return &postgresRepository{db: db}
}

func (r *postgresRepository) GetUserByID(id string) (*User, *customerror.CustomError) {
    // SQL query implementation
    // Return model and error handling
}
```

### 3. Repository Layer
- Acts as intermediary between handlers and database
- Implements interfaces for testability
- Example:

```go
// repository.go
type UserRepository interface {
    GetUserByID(id string) (*User, *customerror.CustomError)
    CreateUser(user *CreateUserRequest) *customerror.CustomError
}

// Implementation in handler
type userHandler struct {
    userRepository UserRepository
    userService    UserService
}
```

### 4. Service Layer 
- Contains business logic and validations
- Processes data transformations
- Example:

```go
// service.go
type UserService interface {
    ValidateUserInput(req UserRequest) *customerror.CustomError
    GenerateUserProfile(userID string) (*UserProfile, *customerror.CustomError)
}

type userService struct {
    // Dependencies
}

func (s *userService) ValidateUserInput(req UserRequest) *customerror.CustomError {
    // Validation logic
    return nil
}
```

### 5. Handler Layer
- Coordinates between services and repositories
- Handles HTTP requests and responses
- Example:

```go
// handler.go
func (h *userHandler) GetUser(c *gin.Context) {
    id := c.Param("id")
    
    // Use service for validation or business logic
    if customErr := h.userService.ValidateID(id); customErr != nil {
        c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
        return
    }
    
    // Use repository for data access
    user, customErr := h.userRepository.GetUserByID(id)
    if customErr != nil {
        c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
        return
    }
    
    c.JSON(http.StatusOK, user)
}

// Register routes
func (h *userHandler) RegisterRoutes(router *gin.RouterGroup) {
    userGroup := router.Group("/users")
    {
        userGroup.GET("/:id", h.GetUser)
        userGroup.POST("/", h.CreateUser)
        // More routes...
    }
}
```

### 6. Model Layer
- Contains data structures used across layers
- Example:

```go
// model.go
type User struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    FullName  string    `json:"full_name"`
    CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
    Email     string `json:"email"`
    FullName  string `json:"full_name"`
    Password  string `json:"password"`
}
```

## Setting Up a New Router

1. **Define Models** - Create data structures first
2. **Create Repository Interface** - Define data access methods
3. **Implement Database Layer** - Create concrete repository
4. **Create Service Interface** - Define business logic operations
5. **Implement Service** - Business logic implementation
6. **Create Handler** - HTTP request handling using services and repositories
7. **Register Routes** - Connect handlers to routing system

## Dependency Injection

```go
// main.go or setup.go
func SetupRouter() *gin.Engine {
    router := gin.Default()
    
    // Setup database connection
    db := database.Connect()
    
    // Create repositories
    userRepo := postgres.NewUserRepository(db)
    
    // Create services
    userService := service.NewUserService()
    
    // Create handlers
    userHandler := handler.NewUserHandler(userService, userRepo)
    
    // Register routes
    api := router.Group("/api")
    userHandler.RegisterRoutes(api)
    
    return router
}
```

This architecture promotes clean separation of concerns, testability, and maintainability by ensuring each component has a single responsibility and clear dependencies.