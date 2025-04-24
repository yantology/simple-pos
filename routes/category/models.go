package category

import (
	"time"
)

// Category represents a category entity
type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateCategoryRequest represents the data needed to create a new category
type CreateCategoryRequest struct {
	Name   string `json:"name" binding:"required"`
	UserID int    `json:"user_id" binding:"required"`
}

// UpdateCategoryRequest represents the data needed to update an existing category
type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}
