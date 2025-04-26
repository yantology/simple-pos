package product

import "time"

// ProductResponse represents the product data returned in API responses
// @Description Product model
type Product struct {
	ID          int       `json:"id" example:"1"` // Changed from string to int
	Name        string    `json:"name" example:"Laptop Pro"`
	Price       float64   `json:"price" example:"15000000"`
	IsAvailable bool      `json:"is_available" example:"true"`
	CategoryID  int       `json:"category_id" example:"1"` // Changed from string to int
	UserID      int       `json:"user_id" example:"1"`     // Changed from string to int
	CreatedAt   time.Time `json:"created_at" example:"2025-04-25T15:04:05Z07:00"`
	UpdatedAt   time.Time `json:"updated_at" example:"2025-04-25T15:04:05Z07:00"`
}

// ProductListResponse represents the response for listing products
// @Description Product list response model
type ProductList struct {
	Products []Product `json:"products"`
}

// UpdateProduct defines the structure for updating a product
// @Description Update product request model
type UpdateProduct struct {
	Name        string  `json:"name" binding:"required" example:"Laptop Pro X"`
	Price       float64 `json:"price" binding:"required,gt=0" example:"16500000"`
	IsAvailable bool    `json:"is_available" example:"false"`
	CategoryID  int     `json:"category_id" binding:"required" example:"2"` // Changed from string to int
}

// CreateProduct defines the structure for creating a new product
// @Description Create product request model
type CreateProduct struct {
	Name        string  `json:"name" binding:"required" example:"Wireless Mouse"`
	Price       float64 `json:"price" binding:"required,gt=0" example:"250000"`
	IsAvailable bool    `json:"is_available" example:"true"`
	CategoryID  int     `json:"category_id" binding:"required" example:"1"` // Changed from string to int
}
