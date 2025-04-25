package order

// CreateOrderRequest represents the request to create a new order
// @Description Create order request model
type CreateOrderRequest struct {
	Total   float64     `json:"total" binding:"required" example:"1299.99"`
	Product ProductJSON `json:"product" binding:"required"`
	UserID  int         `json:"user_id" binding:"required" example:"2"`
}

// OrderResponse represents the response for order endpoints
// @Description Order response model
type OrderResponse struct {
	ID        string      `json:"id" example:"1"`
	Total     float64     `json:"total" example:"1299.99"`
	Product   ProductJSON `json:"product"`
	UserID    string      `json:"user_id" example:"2"`
	CreatedAt string      `json:"created_at" example:"2025-04-25T15:04:05Z07:00"`
	UpdatedAt string      `json:"updated_at" example:"2025-04-25T15:04:05Z07:00"`
}
