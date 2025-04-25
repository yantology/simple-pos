package category

import (
	"time"
)

// Category represents a category entity
// @Description Category model
type Category struct {
	ID        string    `json:"id" binding:"required" example:"uuid-string-1"`
	Name      string    `json:"name" binding:"required" example:"Electronics"`
	UserID    string    `json:"user_id" binding:"required" example:"uuid-user-1"`
	CreatedAt time.Time `json:"created_at" binding:"required" example:"2025-04-25T15:04:05Z07:00"`
	UpdatedAt time.Time `json:"updated_at" binding:"required" example:"2025-04-25T15:04:05Z07:00"`
}

// CreateCategoryRequest represents the data needed to create a new category
// @Description Create category request model
type CreateCategory struct {
	Name string `json:"name" binding:"required" example:"Groceries"`
}

// UpdateCategoryRequest represents the data needed to update an existing category
// @Description Update category request model
type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required" example:"Home Goods"`
}

// DataResponse represents a generic data response
// @Description Generic data response model
type DataResponse[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message" example:"Operation completed successfully"`
}
