package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin" // Import middleware package
	"github.com/yantology/simple-pos/pkg/dto"
)

// Handler holds the dependencies for the product handlers
type Handler struct {
	repository Repository
}

// NewHandler creates a new Handler instance
func NewHandler(repository Repository) *Handler {
	return &Handler{
		repository: repository,
	}
}

// Routes sets up all the routes for product management
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {

	router.POST("", h.CreateProduct)
	router.GET("", h.GetAllProducts)
	router.PUT("/:id", h.UpdateProduct)
	router.DELETE("/:id", h.DeleteProduct)
	router.GET("/category/:categoryID", h.GetProductsByCategoryID)
}

// @Summary Create a new product
// @Description Creates a new product associated with the authenticated user.
// @Tags products
// @Accept json
// @Produce json
// @Param product body CreateProduct true "Product details"
// @Security ApiKeyAuth
// @Success 201 {object} Product "Product created successfully"
// @Failure 400 {object} dto.MessageResponse "Invalid request data"
// @Failure 401 {object} dto.MessageResponse "Unauthorized: User ID not found in context"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /products [post]
func (h *Handler) CreateProduct(c *gin.Context) {
	var request CreateProduct // Use CreateProduct from models.go
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

	fmt.Println("userIDVal:", userIDVal)

	userID, err := strconv.Atoi(userIDVal.(string)) // Assert userID as int
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MessageResponse{Message: "Internal Server Error: User ID in context is not an integer"})
		return
	}

	// Create the product entity, associating it with the authenticated user
	// Pass the request directly as it matches CreateProduct struct now
	createdProduct, customErr := h.repository.Create(&request, userID) // Pass int userID
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{
			Message: customErr.Message(),
		})
		return
	}

	// No need to format, createdProduct is already *Product
	c.JSON(http.StatusCreated, dto.DataResponse[*Product]{ // Use *Product
		Data: createdProduct,
	})
}

// @Summary Get all products
// @Description Retrieves a list of all products available in the system. (No user filtering currently)
// @Tags products
// @Produce json
// @Success 200 {object} []Product "Successfully retrieved products"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /products [get]
func (h *Handler) GetAllProducts(c *gin.Context) {
	// Note: This endpoint currently retrieves ALL products, regardless of user.
	// If you need to restrict this to a specific user's products,
	// you should add the userID extraction logic here and call GetByUserID instead.
	// For now, keeping it as a general "get all" endpoint.
	products, customErr := h.repository.GetAll()
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
		return
	}

	// Products are already []*Product, no need for formatting
	c.JSON(http.StatusOK, dto.DataResponse[[]*Product]{ // Use []*Product
		Data: products,
	})
}

// @Summary Update an existing product
// @Description Updates an existing product by its ID. User must own the product.
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID" // Changed param type to int
// @Param product body UpdateProduct true "Updated product details"
// @Security ApiKeyAuth
// @Success 200 {object} Product "Product updated successfully"
// @Failure 400 {object} dto.MessageResponse "Invalid request data or ID format"
// @Failure 401 {object} dto.MessageResponse "Unauthorized: User ID not found in context or not owner"
// @Failure 404 {object} dto.MessageResponse "Product not found"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /products/{id} [put]
func (h *Handler) UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam) // Convert idParam to int
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "Invalid product ID format"})
		return
	}

	var request UpdateProduct // Use UpdateProduct from models.go
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

	// Pass idParam, userID and the request directly
	updatedProduct, customErr := h.repository.Update(id, userID, &request) // Pass int id and userID
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{
			Message: customErr.Message(),
		})
		return
	}

	// No need to format, updatedProduct is already *Product
	c.JSON(http.StatusOK, dto.DataResponse[*Product]{ // Use *Product
		Data: updatedProduct,
	})
}

// @Summary Delete a product
// @Description Deletes a product by its ID. User must own the product.
// @Tags products
// @Produce json
// @Param id path int true "Product ID" // Changed param type to int
// @Security ApiKeyAuth
// @Success 200 {object} dto.MessageResponse "Product deleted successfully"
// @Failure 400 {object} dto.MessageResponse "Invalid product ID format"
// @Failure 401 {object} dto.MessageResponse "Unauthorized: User ID not found in context or not owner"
// @Failure 404 {object} dto.MessageResponse "Product not found"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /products/{id} [delete]
func (h *Handler) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam) // Convert idParam to int
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "Invalid product ID format"})
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

	// Pass idParam and userID for authorization check in repository
	customErr := h.repository.Delete(id, userID) // Pass int id and userID
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{
			Message: customErr.Message(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Product deleted successfully"})
}

// @Summary Get products by Category ID
// @Description Retrieves all products belonging to a specific Category ID.
// @Tags products
// @Produce json
// @Param categoryID path int true "Category ID" // Changed param type to int
// @Success 200 {object} []Product "Successfully retrieved products"
// @Failure 400 {object} dto.MessageResponse "Invalid Category ID format"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /products/category/{categoryID} [get]
func (h *Handler) GetProductsByCategoryID(c *gin.Context) {
	categoryIDParam := c.Param("categoryID")
	categoryID, err := strconv.Atoi(categoryIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "Invalid Category ID format"})
		return
	}

	products, customErr := h.repository.GetByCategoryID(categoryID)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{
			Message: customErr.Message(),
		})
		return
	}

	// Products are already []*Product, no need for formatting
	c.JSON(http.StatusOK, dto.DataResponse[[]*Product]{ // Use []*Product
		Data: products,
	})
}
