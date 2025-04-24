package category

import (
	"net/http"
	"strconv"

	"github.com/yantology/simple-ecommerce/pkg/dto"

	"github.com/gin-gonic/gin"
)

// CategoryHandler handles HTTP requests for categories
type CategoryHandler struct {
	repository Repository
	service    Service
}

// NewCategoryHandler creates a new handler instance
func NewCategoryHandler(repository Repository, service Service) *CategoryHandler {
	return &CategoryHandler{
		repository: repository,
		service:    service,
	}
}

// RegisterRoutes registers category routes to the router
func (h *CategoryHandler) RegisterRoutes(router *gin.RouterGroup) {
	categoryGroup := router.Group("/categories")
	{
		categoryGroup.GET("/:id", h.GetCategoryByID)
		categoryGroup.GET("/name/:name", h.GetCategoryByName)
		categoryGroup.GET("/user/:userId", h.GetCategoriesByUserID)
		categoryGroup.POST("/", h.CreateCategory)
		categoryGroup.PUT("/:id", h.UpdateCategory)
		categoryGroup.DELETE("/:id", h.DeleteCategory)
	}
}

// formatCategoryResponse converts a Category model to a CategoryResponse DTO
func formatCategoryResponse(category *Category) CategoryResponse {
	return CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		UserID:    category.UserID,
		CreatedAt: category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// GetCategoryByID handles GET /categories/:id
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{
			Message: "Invalid category ID format",
		})
		return
	}

	// Use service for input validation if available
	if h.service != nil {
		if customErr := h.service.ValidateCategoryID(id); customErr != nil {
			c.JSON(customErr.Code(), dto.MessageResponse{
				Message: customErr.Message(),
			})
			return
		}
	}

	// Direct repository call instead of service call for data access
	category, customErr := h.repository.GetCategoryByID(id)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{
			Message: customErr.Message(),
		})
		return
	}

	response := formatCategoryResponse(category)
	// Use the generic DataResponse from pkg/dto
	c.JSON(http.StatusOK, dto.DataResponse[CategoryResponse]{
		Data:    response,
		Message: "Category retrieved successfully",
	})
}

// GetCategoryByName handles GET /categories/name/:name
func (h *CategoryHandler) GetCategoryByName(c *gin.Context) {
	name := c.Param("name")

	// Use service for input validation if available
	if h.service != nil {
		if customErr := h.service.ValidateCategoryName(name); customErr != nil {
			c.JSON(customErr.Code(), dto.MessageResponse{
				Message: customErr.Message(),
			})
			return
		}
	}

	// Direct repository call
	category, customErr := h.repository.GetCategoryByName(name)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{ // Use MessageResponse for consistency
			Message: customErr.Message(),
		})
		return
	}

	response := formatCategoryResponse(category)
	// Use the generic DataResponse from pkg/dto
	c.JSON(http.StatusOK, dto.DataResponse[CategoryResponse]{
		Data:    response,
		Message: "Category retrieved successfully",
	})
}

// GetCategoriesByUserID handles GET /categories/user/:userId
func (h *CategoryHandler) GetCategoriesByUserID(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.Atoi(userIDStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{ // Use MessageResponse
			Message: "Invalid user ID format",
		})
		return
	}

	// Use service for input validation if available
	if h.service != nil {
		if customErr := h.service.ValidateUserID(userID); customErr != nil {
			c.JSON(customErr.Code(), dto.MessageResponse{ // Use MessageResponse
				Message: customErr.Message(),
			})
			return
		}
	}

	// Direct repository call
	categories, customErr := h.repository.GetCategoriesByUserID(userID)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{ // Use MessageResponse
			Message: customErr.Message(),
		})
		return
	}

	var responseCategories []CategoryResponse
	for _, category := range categories {
		categoryCopy := category // Create a copy to avoid pointer issues
		responseCategories = append(responseCategories, formatCategoryResponse(&categoryCopy))
	}

	// Use the generic DataResponse with a slice of CategoryResponse
	c.JSON(http.StatusOK, dto.DataResponse[[]CategoryResponse]{
		Data:    responseCategories,
		Message: "Categories retrieved successfully",
	})
}

// CreateCategory handles POST /categories/
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var requestDTO CreateCategoryDTO

	if err := c.ShouldBindJSON(&requestDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{ // Use MessageResponse
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	// Extract user ID from authenticated user (assuming middleware sets this)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{ // Use MessageResponse
			Message: "User not authenticated",
		})
		return
	}

	// Create the category request with the authenticated user's ID
	createRequest := &CreateCategoryRequest{
		Name:   requestDTO.Name,
		UserID: userID.(int),
	}

	// Use service for input validation if available
	if h.service != nil {
		if customErr := h.service.ValidateCreateCategoryRequest(createRequest); customErr != nil {
			c.JSON(customErr.Code(), dto.MessageResponse{ // Use MessageResponse
				Message: customErr.Message(),
			})
			return
		}
	}

	// Direct repository call
	category, customErr := h.repository.CreateCategory(createRequest)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{ // Use MessageResponse
			Message: customErr.Message(),
		})
		return
	}

	response := formatCategoryResponse(category)
	// Use the generic DataResponse
	c.JSON(http.StatusCreated, dto.DataResponse[CategoryResponse]{
		Data:    response,
		Message: "Category created successfully",
	})
}

// UpdateCategory handles PUT /categories/:id
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{ // Use MessageResponse
			Message: "Invalid category ID format",
		})
		return
	}

	var requestDTO UpdateCategoryDTO
	if err := c.ShouldBindJSON(&requestDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{ // Use MessageResponse
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	// Extract user ID from authenticated user to check ownership
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{ // Use MessageResponse
			Message: "User not authenticated",
		})
		return
	}

	// Check if the category exists and belongs to the user
	existingCategory, customErr := h.repository.GetCategoryByID(id)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{ // Use MessageResponse
			Message: customErr.Message(),
		})
		return
	}

	// Use service for ownership validation if available
	if h.service != nil {
		if customErr := h.service.CheckOwnership(existingCategory.UserID, userID.(int)); customErr != nil {
			c.JSON(customErr.Code(), dto.MessageResponse{ // Use MessageResponse
				Message: customErr.Message(),
			})
			return
		}
	} else {
		// Basic ownership check if no service available
		if existingCategory.UserID != userID.(int) {
			c.JSON(http.StatusForbidden, dto.MessageResponse{ // Use MessageResponse
				Message: "You don't have permission to edit this category",
			})
			return
		}
	}

	updateRequest := &UpdateCategoryRequest{
		Name: requestDTO.Name,
	}

	// Use service for input validation if available
	if h.service != nil {
		if customErr := h.service.ValidateUpdateCategoryRequest(id, updateRequest); customErr != nil {
			c.JSON(customErr.Code(), dto.MessageResponse{ // Use MessageResponse
				Message: customErr.Message(),
			})
			return
		}
	}

	// Direct repository call
	category, customErr := h.repository.UpdateCategory(id, updateRequest)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{ // Use MessageResponse
			Message: customErr.Message(),
		})
		return
	}

	response := formatCategoryResponse(category)
	// Use the generic DataResponse
	c.JSON(http.StatusOK, dto.DataResponse[CategoryResponse]{
		Data:    response,
		Message: "Category updated successfully",
	})
}

// DeleteCategory handles DELETE /categories/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{ // Use MessageResponse
			Message: "Invalid category ID format",
		})
		return
	}

	// Extract user ID from authenticated user to check ownership
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{ // Use MessageResponse
			Message: "User not authenticated",
		})
		return
	}

	// Check if the category exists and belongs to the user
	existingCategory, customErr := h.repository.GetCategoryByID(id)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{ // Use MessageResponse
			Message: customErr.Message(),
		})
		return
	}

	// Use service for ownership validation if available
	if h.service != nil {
		if customErr := h.service.CheckOwnership(existingCategory.UserID, userID.(int)); customErr != nil {
			c.JSON(customErr.Code(), dto.MessageResponse{ // Use MessageResponse
				Message: customErr.Message(),
			})
			return
		}
	} else {
		// Basic ownership check if no service available
		if existingCategory.UserID != userID.(int) {
			c.JSON(http.StatusForbidden, dto.MessageResponse{ // Use MessageResponse
				Message: "You don't have permission to delete this category",
			})
			return
		}
	}

	// Direct repository call
	customErr = h.repository.DeleteCategory(id)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{ // Use MessageResponse
			Message: customErr.Message(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{ // Use MessageResponse
		Message: "Category deleted successfully",
	})
}
