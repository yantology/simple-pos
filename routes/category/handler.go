package category

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yantology/simple-pos/pkg/dto"
)

// CategoryHandler handles HTTP requests for categories
type CategoryHandler struct {
	repository Repository
}

// NewCategoryHandler creates a new handler instance
// S
func NewCategoryHandler(repository Repository) *CategoryHandler {
	return &CategoryHandler{
		repository: repository,
	}
}

// RegisterRoutes registers category routes to the router
// @Summary Register routes with authentication
func (h *CategoryHandler) RegisterRoutes(router *gin.RouterGroup) {

	router.GET("/", h.GetAllCategories) // Add route for getting all categories
	router.GET("/:id", h.GetCategoryByID)
	router.GET("/name/:name", h.GetCategoryByName)
	router.POST("/", h.CreateCategory)
	router.PUT("/:id", h.UpdateCategory)
	router.DELETE("/:id", h.DeleteCategory)

}

// @Summary Get category by ID
// @Description Retrieves a specific category by its ID for the authenticated user.
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} Category "Successfully retrieved categories"
// @Failure 400 {object} dto.MessageResponse "Invalid category ID format (if applicable)"
// @Failure 401 {object} dto.MessageResponse "Unauthorized: User ID not found in context"
// @Failure 404 {object} dto.MessageResponse "Category not found"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "Invalid category ID format"})
		return
	}

	// Get userID from middleware context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{Message: "Unauthorized: User ID not found in context"})
		return
	}

	userID, err := strconv.Atoi(userIDVal.(string)) // Assert userID as int
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MessageResponse{Message: "Internal Server Error: User ID in context is not an integer"})
		return
	}

	// Direct repository call, passing userID for authorization
	category, customErr := h.repository.GetCategoryByID(id, userID)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{
			Message: customErr.Message(),
		})
		return
	}

	// Use the Category struct directly from models.go
	c.JSON(http.StatusOK, dto.DataResponse[*Category]{ // Use *Category
		Data: category,
	})
}

// @Summary Get category by name
// @Description Retrieves a specific category by its name for the authenticated user.
// @Tags categories
// @Produce json
// @Param name path string true "Category Name"
// @Success 200 {object} Category "Successfully retrieved categories"
// @Failure 401 {object} dto.MessageResponse "Unauthorized: User ID not found in context"
// @Failure 404 {object} dto.MessageResponse "Category not found"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /categories/name/{name} [get]
func (h *CategoryHandler) GetCategoryByName(c *gin.Context) {
	name := c.Param("name")

	// Get userID from middleware context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{Message: "Unauthorized: User ID not found in context"})
		return
	}

	userID, err := strconv.Atoi(userIDVal.(string)) // Assert userID as int
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MessageResponse{Message: "Internal Server Error: User ID in context is not an integer"})
		return
	}

	// Direct repository call, passing userID for authorization
	category, customErr := h.repository.GetCategoryByName(name, userID)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{
			Message: customErr.Message(),
		})
		return
	}

	// Use the Category struct directly from models.go
	c.JSON(http.StatusOK, dto.DataResponse[*Category]{ // Use *Category
		Data: category,
	})
}

// @Summary Get all categories for the authenticated user
// @Description Retrieves a list of all categories associated with the logged-in user.
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} []Category
// @Failure 400 {object} dto.MessageResponse "Invalid request data"
// @Failure 401 {object} dto.MessageResponse "Unauthorized: User ID not found in context"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /categories [get]
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	// Get userID from middleware context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{Message: "Unauthorized: User ID not found in context"})
		return
	}

	userID, err := strconv.Atoi(userIDVal.(string)) // Assert userID as int
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MessageResponse{Message: "Internal Server Error: User ID in context is not an integer"})
		return
	}

	// Direct repository call to get all categories for the user
	categories, customErr := h.repository.GetAllCategoriesByUserID(userID)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{
			Message: customErr.Message(),
		})
		return
	}
	// Use the Category struct directly from models.go
	c.JSON(http.StatusOK, dto.DataResponse[[]Category]{
		Data: categories,
	})
}

// @Summary Create a new category
// @Description Creates a new category for the authenticated user.
// @Tags categories
// @Accept json
// @Produce json
// @Param category body category.CreateCategory true "Category details"
// @Success 201 {object} Category "Successfully retrieved categories"
// @Failure 400 {object} dto.MessageResponse "Invalid request data"
// @Failure 401 {object} dto.MessageResponse "Unauthorized: User ID not found in context"
// @Failure 409 {object} dto.MessageResponse "Category with this name already exists"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var request CreateCategory // Use CreateCategory from models.go
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "Invalid request data: " + err.Error()})
		return
	}

	// Get userID from middleware context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{Message: "Unauthorized: User ID not found in context"})
		return
	}

	userID, err := strconv.Atoi(userIDVal.(string)) // Assert userID as int
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MessageResponse{Message: "Internal Server Error: User ID in context is not an integer"})
		return
	}

	// Direct repository call, passing the request struct and userID
	category, customErr := h.repository.CreateCategory(&request, userID)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{
			Message: customErr.Message(),
		})
		return
	}

	// Use the Category struct directly from models.go
	c.JSON(http.StatusCreated, dto.DataResponse[*Category]{ // Use *Category
		Data: category,
	})
}

// @Summary Update an existing category
// @Description Updates an existing category by its ID for the authenticated user.
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body category.UpdateCategoryRequest true "Updated category details"
// @Success 200 {object} Category
// @Failure 400 {object} dto.MessageResponse "Invalid request data or ID format"
// @Failure 401 {object} dto.MessageResponse "Unauthorized: User ID not found in context or not owner"
// @Failure 404 {object} dto.MessageResponse "Category not found"
// @Failure 409 {object} dto.MessageResponse "Category with this name already exists"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "Invalid category ID format"})
		return
	}

	var request UpdateCategoryRequest // Use UpdateCategoryRequest from models.go
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "Invalid request data: " + err.Error()})
		return
	}

	// Get userID from middleware context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{Message: "Unauthorized: User ID not found in context"})
		return
	}

	userID, err := strconv.Atoi(userIDVal.(string)) // Assert userID as int
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MessageResponse{Message: "Internal Server Error: User ID in context is not an integer"})
		return
	}

	// Direct repository call, passing userID for authorization and the request struct
	category, customErr := h.repository.UpdateCategory(id, userID, &request)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{
			Message: customErr.Message(),
		})
		return
	}

	// Use the Category struct directly from models.go
	c.JSON(http.StatusOK, dto.DataResponse[*Category]{ // Use *Category
		Data: category,
	})
}

// @Summary Delete a category
// @Description Deletes a category by its ID for the authenticated user.
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} Category
// @Failure 400 {object} dto.MessageResponse "Invalid category ID format"
// @Failure 401 {object} dto.MessageResponse "Unauthorized: User ID not found in context or not owner"
// @Failure 404 {object} dto.MessageResponse "Category not found"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "Invalid category ID format"})
		return
	}

	// Get userID from middleware context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{Message: "Unauthorized: User ID not found in context"})
		return
	}

	userID, err := strconv.Atoi(userIDVal.(string)) // Assert userID as int
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MessageResponse{Message: "Internal Server Error: User ID in context is not an integer"})
		return
	}

	// Direct repository call, passing userID for authorization
	customErr := h.repository.DeleteCategory(id, userID)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{
			Message: customErr.Message(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{ // Use MessageResponse
		Message: "Category deleted successfully",
	})
}
