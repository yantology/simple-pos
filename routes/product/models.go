package product

import (
	"time"
)

// Product represents a product entity
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	IsAvailable bool      `json:"is_available"`
	CategoryID  int       `json:"category_id"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductRequest represents the data needed to create a new product
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	IsAvailable bool    `json:"is_available"`
	CategoryID  int     `json:"category_id" binding:"required"`
	UserID      int     `json:"user_id" binding:"required"`
}

// UpdateProductRequest represents the data needed to update an existing product
type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	IsAvailable bool    `json:"is_available"`
	CategoryID  int     `json:"category_id"`
}
