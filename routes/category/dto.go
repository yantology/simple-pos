package category

// CategoryResponse represents the response structure for a category
type CategoryResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	UserID    int    `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetCategoriesResponse wraps multiple categories response
type GetCategoriesResponse struct {
	Categories []CategoryResponse `json:"categories"`
}

// CreateCategoryDTO represents the request structure for creating a category
type CreateCategoryDTO struct {
	Name string `json:"name" binding:"required"`
}

// UpdateCategoryDTO represents the request structure for updating a category
type UpdateCategoryDTO struct {
	Name string `json:"name" binding:"required"`
}
