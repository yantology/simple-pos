# Go Router Architecture Guide (Updated)

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
}
```

### 4. Service Layer (Optional)
- Should NOT call repositories
- Only handlers can call repositories
- Services can contain business logic that doesn't need data access
- May be omitted if there's no need for repository-independent logic
- It will canby handler
- Example:

```go
// service.go
type UserService interface {
    ValidateUserInput(req UserRequest) *customerror.CustomError
    FormatUserProfile(user *User) (*UserProfile, *customerror.CustomError)
}

type userService struct {
    // NO repository dependencies
}

func (s *userService) ValidateUserInput(req UserRequest) *customerror.CustomError {
    // Validation logic that doesn't need database access
    return nil
}
```

### 5. Handler Layer
- Coordinates between services and repositories
- ONLY component allowed to call repositories
- Handles HTTP requests and responses
- Using dto.MessageResponse or dto.DataResponse<model>
- Example:

```go
// handler.go
func (h *userHandler) GetUser(c *gin.Context) {
    id := c.Param("id")
    
    // Use service for validation or business logic that doesn't need DB access
    if h.userService != nil {
        if customErr := h.userService.ValidateID(id); customErr != nil {
            c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
            return
        }
    }
    
    // ONLY the handler calls the repository
    user, customErr := h.userRepository.GetUserByID(id)
    if customErr != nil {
        c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
        return
    }
    
    // Use service for non-DB operations like formatting
    var formattedUser interface{} = user
    if h.userService != nil {
        formattedUser, customErr = h.userService.FormatUserProfile(user)
        if customErr != nil {
            c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
            return
        }
    }
    
    c.JSON(http.StatusOK, formattedUser)
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
4. **Create Service Interface** (Optional) - Only for non-DB related operations
5. **Implement Service** (Optional) - Only if needed for complex non-DB logic
6. **Create Handler** - HTTP request handling, calling repositories directly
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
    
    // Create services (optional, for non-DB logic only)
    var userService service.UserService
    if complexLogicNeeded {
        userService = service.NewUserService() // No repo dependency
    }
    
    // Create handlers - repositories injected directly
    userHandler := handler.NewUserHandler(userRepo, userService)
    
    // Register routes
    api := router.Group("/api")
    userHandler.RegisterRoutes(api)
    
    return router
}
```

This updated architecture ensures that only handlers can directly access repositories, maintaining a clean separation of concerns while simplifying the service layer or making it optional when no complex business logic is needed apart from data access.
