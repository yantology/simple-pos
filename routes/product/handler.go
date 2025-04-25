package product

import (
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
func (h *Handler) Routes(router *gin.RouterGroup) {
	productGroup := router.Group("/products")
	{
		productGroup.POST("", h.CreateProduct)
		productGroup.GET("", h.GetAllProducts)
		productGroup.PUT("/:id", h.UpdateProduct)
		productGroup.DELETE("/:id", h.DeleteProduct)
		productGroup.GET("/user/:userID", h.GetProductsByUserID) // Keep this route for admin/specific use cases if needed, but primary user access should rely on context userID
		productGroup.GET("/category/:categoryID", h.GetProductsByCategoryID)
	}
}

// @Summary Create a new product
// @Description Creates a new product associated with the authenticated user.
// @Tags products
// @Accept json
// @Produce json
// @Param product body CreateProduct true "Product details"
// @Security ApiKeyAuth
// @Success 201 {object} dto.DataResponse[*Product] "Product created successfully"
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

	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.MessageResponse{Message: "Internal Server Error: User ID in context is not a string"})
		return
	}

	// Create the product entity, associating it with the authenticated user
	// Pass the request directly as it matches CreateProduct struct now
	createdProduct, customErr := h.repository.Create(&request, userID) // Pass the request and userID
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
		return
	}

	// No need to format, createdProduct is already *Product
	c.JSON(http.StatusCreated, dto.DataResponse[*Product]{ // Use *Product
		Data:    createdProduct,
		Message: "Product created successfully",
	})
}

// @Summary Get all products
// @Description Retrieves a list of all products available in the system. (No user filtering currently)
// @Tags products
// @Produce json
// @Success 200 {object} dto.DataResponse[[]*Product] "Successfully retrieved products"
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
		Data:    products,
		Message: "Products retrieved successfully",
	})
}

// @Summary Update an existing product
// @Description Updates an existing product by its ID. User must own the product.
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body UpdateProduct true "Updated product details"
// @Security ApiKeyAuth
// @Success 200 {object} dto.DataResponse[*Product] "Product updated successfully"
// @Failure 400 {object} dto.MessageResponse "Invalid request data or ID format"
// @Failure 401 {object} dto.MessageResponse "Unauthorized: User ID not found in context or not owner"
// @Failure 404 {object} dto.MessageResponse "Product not found"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /products/{id} [put]
func (h *Handler) UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")

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

	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.MessageResponse{Message: "Internal Server Error: User ID in context is not a string"})
		return
	}

	// Pass idParam, userID and the request directly
	updatedProduct, customErr := h.repository.Update(idParam, userID, &request)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
		return
	}

	// No need to format, updatedProduct is already *Product
	c.JSON(http.StatusOK, dto.DataResponse[*Product]{ // Use *Product
		Data:    updatedProduct,
		Message: "Product updated successfully",
	})
}

// @Summary Delete a product
// @Description Deletes a product by its ID. User must own the product.
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Security ApiKeyAuth
// @Success 200 {object} dto.MessageResponse "Product deleted successfully"
// @Failure 400 {object} dto.MessageResponse "Invalid product ID format"
// @Failure 401 {object} dto.MessageResponse "Unauthorized: User ID not found in context or not owner"
// @Failure 404 {object} dto.MessageResponse "Product not found"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /products/{id} [delete]
func (h *Handler) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")

	// Get userID from middleware context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{Message: "Unauthorized: User ID not found in context"})
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.MessageResponse{Message: "Internal Server Error: User ID in context is not a string"})
		return
	}

	// Pass idParam and userID for authorization check in repository
	customErr := h.repository.Delete(idParam, userID)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Product deleted successfully"})
}

// @Summary Get products by User ID
// @Description Retrieves all products associated with a specific User ID. (Requires appropriate authorization)
// @Tags products
// @Produce json
// @Param userID path string true "User ID"
// @Security ApiKeyAuth // Assuming admin or specific permissions needed
// @Success 200 {object} dto.DataResponse[[]*Product] "Successfully retrieved products"
// @Failure 400 {object} dto.MessageResponse "Invalid User ID format"
// @Failure 401 {object} dto.MessageResponse "Unauthorized"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /products/user/{userID} [get]
func (h *Handler) GetProductsByUserID(c *gin.Context) {
	// This route uses the userID from the URL parameter.
	// If you want the currently logged-in user's products, use the context userID.
	targetUserID := c.Param("userID") // Get target userID from URL

	// Optional: Add logic here to check if the requesting user (from context)
	// has permission to view products of the targetUserID.
	// For example:
	// requestingUserIDVal, exists := c.Get("user_id")
	// if !exists || requestingUserIDVal.(string) != targetUserID { // Basic check: only allow users to see their own products via this route
	//     // Or check if requesting user is an admin
	//     c.JSON(http.StatusForbidden, dto.MessageResponse{Message: "Forbidden"})
	// 	   return
	// }

	products, customErr := h.repository.GetByUserID(targetUserID) // Use targetUserID from URL
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
		return
	}

	// Products are already []*Product, no need for formatting
	c.JSON(http.StatusOK, dto.DataResponse[[]*Product]{ // Use []*Product
		Data:    products,
		Message: "Products retrieved successfully for user",
	})
}

// @Summary Get products by Category ID
// @Description Retrieves all products belonging to a specific Category ID.
// @Tags products
// @Produce json
// @Param categoryID path string true "Category ID"
// @Success 200 {object} dto.DataResponse[[]*Product] "Successfully retrieved products"
// @Failure 400 {object} dto.MessageResponse "Invalid Category ID format"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /products/category/{categoryID} [get]
func (h *Handler) GetProductsByCategoryID(c *gin.Context) {
	categoryIDParam := c.Param("categoryID")
	categoryID, err := strconv.Atoi(categoryIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "Invalid category ID format"})
		return
	}

	products, customErr := h.repository.GetByCategoryID(categoryID)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
		return
	}

	// Products are already []*Product, no need for formatting
	c.JSON(http.StatusOK, dto.DataResponse[[]*Product]{ // Use []*Product
		Data:    products,
		Message: "Products retrieved successfully for category",
	})
}
