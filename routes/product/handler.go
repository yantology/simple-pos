package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service    Service
	repository Repository
}

// NewHandler creates a new handler instance
func NewHandler(repository Repository, service Service) *Handler {
	return &Handler{
		repository: repository,
		service:    service,
	}
}

// Routes sets up all the routes for product management
func (h *Handler) Routes(router *gin.RouterGroup) {
	productGroup := router.Group("/products")
	{
		productGroup.POST("", h.CreateProduct)
		productGroup.GET("", h.GetAllProducts)
		productGroup.GET("/:id", h.GetProductByID)
		productGroup.PUT("/:id", h.UpdateProduct)
		productGroup.DELETE("/:id", h.DeleteProduct)
		productGroup.GET("/user/:userID", h.GetProductsByUserID)
		productGroup.GET("/category/:categoryID", h.GetProductsByCategoryID)
	}
}

// CreateProduct handles the creation of a new product
func (h *Handler) CreateProduct(c *gin.Context) {
	productRequest, err := CreateProductDTO(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use service for validation (business logic)
	if err := h.service.ValidateProductInput(productRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use service to prepare the product for creation
	product := h.service.PrepareProductForCreation(productRequest)

	// Use repository for database operations
	createdProduct, err := h.repository.Create(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Use service to format the response
	response := h.service.FormatProductResponse(createdProduct)
	c.JSON(http.StatusCreated, response)
}

// GetProductByID handles retrieval of a product by its ID
func (h *Handler) GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	// Use repository for database operations
	product, err := h.repository.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Use service to format the response
	response := h.service.FormatProductResponse(product)
	c.JSON(http.StatusOK, response)
}

// GetAllProducts handles retrieval of all products
func (h *Handler) GetAllProducts(c *gin.Context) {
	// Use repository for database operations
	products, err := h.repository.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Use service to format the response
	response := h.service.FormatProductsResponse(products)
	c.JSON(http.StatusOK, response)
}

// UpdateProduct handles updating an existing product
func (h *Handler) UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	productRequest, err := UpdateProductDTO(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use service for validation (business logic)
	if err := h.service.ValidateUpdateRequest(productRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// First check if the product exists
	_, err = h.repository.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Use repository for database operations
	updatedProduct, err := h.repository.Update(id, productRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Use service to format the response
	response := h.service.FormatProductResponse(updatedProduct)
	c.JSON(http.StatusOK, response)
}

// DeleteProduct handles deletion of a product
func (h *Handler) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	// Use repository for database operations
	if err := h.repository.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

// GetProductsByUserID handles retrieval of products by user ID
func (h *Handler) GetProductsByUserID(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Use repository for database operations
	products, err := h.repository.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Use service to format the response
	response := h.service.FormatProductsResponse(products)
	c.JSON(http.StatusOK, response)
}

// GetProductsByCategoryID handles retrieval of products by category ID
func (h *Handler) GetProductsByCategoryID(c *gin.Context) {
	categoryIDParam := c.Param("categoryID")
	categoryID, err := strconv.Atoi(categoryIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	// Use repository for database operations
	products, err := h.repository.GetByCategoryID(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Use service to format the response
	response := h.service.FormatProductsResponse(products)
	c.JSON(http.StatusOK, response)
}
