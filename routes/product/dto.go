package product

import (
	"github.com/gin-gonic/gin"
)

// ProductResponse represents the product data returned in API responses
type ProductResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	IsAvailable bool    `json:"is_available"`
	CategoryID  int     `json:"category_id"`
	UserID      int     `json:"user_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// ProductListResponse represents the response for listing products
type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
}

// CreateProductDTO converts CreateProductRequest to Product
func CreateProductDTO(c *gin.Context) (CreateProductRequest, error) {
	var productRequest CreateProductRequest
	if err := c.ShouldBindJSON(&productRequest); err != nil {
		return productRequest, err
	}
	return productRequest, nil
}

// UpdateProductDTO converts UpdateProductRequest to Product
func UpdateProductDTO(c *gin.Context) (UpdateProductRequest, error) {
	var productRequest UpdateProductRequest
	if err := c.ShouldBindJSON(&productRequest); err != nil {
		return productRequest, err
	}
	return productRequest, nil
}

// ToProductResponse converts a Product to ProductResponse
func ToProductResponse(product Product) ProductResponse {
	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		IsAvailable: product.IsAvailable,
		CategoryID:  product.CategoryID,
		UserID:      product.UserID,
		CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToProductListResponse converts a slice of Products to ProductListResponse
func ToProductListResponse(products []Product) ProductListResponse {
	var productResponses []ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, ToProductResponse(product))
	}
	return ProductListResponse{
		Products: productResponses,
	}
}
